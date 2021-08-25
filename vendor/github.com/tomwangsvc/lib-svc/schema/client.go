package schema

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"cloud.google.com/go/civil"

	lib_domain "github.com/tomwangsvc/lib-svc/domain"
	lib_email "github.com/tomwangsvc/lib-svc/email"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_formatters "github.com/tomwangsvc/lib-svc/formatters"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_regexp "github.com/tomwangsvc/lib-svc/regexp"
	lib_time "github.com/tomwangsvc/lib-svc/time"
	"github.com/xeipuuv/gojsonschema"
)

const (
	requestBodyExpectedButFoundNone = "REQUEST_BODY_EXPECTED_BUT_FOUND_NONE"
)

var (
	constantCaseRegExpCompile   *regexp.Regexp
	hourRegExpCompile           *regexp.Regexp
	intRegExpCompile            *regexp.Regexp
	moneyStringRegExpCompile    *regexp.Regexp
	phoneNumberRegExpCompile    *regexp.Regexp
	timeZoneUtcOffsetExpCompile *regexp.Regexp
	websiteRegExpCompile        *regexp.Regexp
)

func init() {
	gojsonschema.FormatCheckers.Add("constant-case", constantCaseChecker{})
	gojsonschema.FormatCheckers.Add("datetime", datetimeChecker{})
	gojsonschema.FormatCheckers.Add("not-empty", notEmptyChecker{})
	gojsonschema.FormatCheckers.Add("locale", localeChecker{})
	gojsonschema.FormatCheckers.Add("local-time", localTimeChecker{})
	gojsonschema.FormatCheckers.Add("locale-or-empty", localeOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("language-code", languageCodeChecker{})
	gojsonschema.FormatCheckers.Add("language-code-or-empty", languageCodeOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("hour-not-empty", hourAndNotEmptyChecker{})
	gojsonschema.FormatCheckers.Add("hour-or-empty", hourOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("money-string", moneyStringChecker{})
	gojsonschema.FormatCheckers.Add("time", timeChecker{})
	gojsonschema.FormatCheckers.Add("time-or-empty", timeOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("country-code", countryCodeChecker{})
	gojsonschema.FormatCheckers.Add("currency", currencyChecker{})
	gojsonschema.FormatCheckers.Add("email", emailChecker{})
	gojsonschema.FormatCheckers.Add("email-or-empty", emailOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("email-with-optional-name", emailWithOptionalNameChecker{})
	gojsonschema.FormatCheckers.Add("website", websiteChecker{})
	gojsonschema.FormatCheckers.Add("website-or-empty", websiteOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("phone-number", phoneNumberChecker{})
	gojsonschema.FormatCheckers.Add("phone-number-or-empty", phoneNumberOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("timezone-utc-offset", timezoneUtcOffsetChecker{})
	gojsonschema.FormatCheckers.Add("timezone-name", timezoneNameChecker{})
	gojsonschema.FormatCheckers.Add("domain", domainChecker{})
	gojsonschema.FormatCheckers.Add("domain-from-or-to", domainFromOrToChecker{})
	gojsonschema.FormatCheckers.Add("date", dateChecker{})
	gojsonschema.FormatCheckers.Add("date-or-empty", dateOrEmptyChecker{})
	gojsonschema.FormatCheckers.Add("extension", extensionChecker{})
	gojsonschema.FormatCheckers.Add("extension-or-empty", extensionOrEmptyChecker{})

	constantCaseRegExpCompile = regexp.MustCompile(lib_regexp.ConstantCaseRegExp)
	hourRegExpCompile = regexp.MustCompile(lib_regexp.HourRegExp)
	moneyStringRegExpCompile = regexp.MustCompile(lib_regexp.MoneyRegExp)
	phoneNumberRegExpCompile = regexp.MustCompile(lib_regexp.PhoneNumberRegExp)
	timeZoneUtcOffsetExpCompile = regexp.MustCompile(lib_regexp.TimeZoneUtcOffsetExp)
	websiteRegExpCompile = regexp.MustCompile(lib_regexp.WebsiteRegExp)
}

type Client interface {
	CheckContentAgainstSchema(ctx context.Context, schemaFileName string, content interface{}) error
	CheckBodyAgainstSchema(ctx context.Context, schemaFileName string, body []byte) error
}

type client struct {
	schemaByFileName map[string]*gojsonschema.Schema
}

func NewClient(ctx context.Context, schemaFileNames []string) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtStrings("schemaFileNames", schemaFileNames))

	type schemaAndFileNameAndError struct {
		Err      error
		FileName string
		Schema   *gojsonschema.Schema
	}
	ch := make(chan schemaAndFileNameAndError)
	for _, schemaFileName := range schemaFileNames {
		go func(schemaFileName string) {
			lib_log.Info(ctx, "Loading", lib_log.FmtString("schemaFileName", schemaFileName))
			schema, err := loadSchemaFromFile(schemaFileName)
			if err != nil {
				ch <- schemaAndFileNameAndError{
					Err: lib_errors.Wrapf(err, "Failed to load schema %q", schemaFileName),
				}
				return
			}
			ch <- schemaAndFileNameAndError{
				FileName: schemaFileName,
				Schema:   schema,
			}
		}(schemaFileName)
	}

	var err error
	schemaByFileName := make(map[string]*gojsonschema.Schema)
	for range schemaFileNames {
		schemaAndFileNameAndError := <-ch
		if schemaAndFileNameAndError.Err != nil {
			lib_log.Error(ctx, "Failed loading schema in goroutine", lib_log.FmtError(schemaAndFileNameAndError.Err))
			err = lib_errors.Wrap(schemaAndFileNameAndError.Err, "Failed loading schema in goroutine")
		}
		schemaByFileName[schemaAndFileNameAndError.FileName] = schemaAndFileNameAndError.Schema
	}
	if err != nil {
		return nil, err
	}

	lib_log.Info(ctx, "Initialized")
	return &client{
		schemaByFileName: schemaByFileName,
	}, nil
}

func loadSchemaFromFile(fileName string) (*gojsonschema.Schema, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed to acquire current directory of process")
	}
	schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/%s", dir, fileName)))
	if err != nil {
		return nil, lib_errors.Wrapf(err, "Failed to initialize content schema for file %q", fileName)
	}
	return schema, nil
}

func (c client) CheckContentAgainstSchema(ctx context.Context, schemaFileName string, content interface{}) error {
	lib_log.Info(ctx, "Checking", lib_log.FmtString("schemaFileName", schemaFileName))

	if err := c.validateAgainstSchemaAndCheckResult(schemaFileName, gojsonschema.NewGoLoader(content)); err != nil {
		return err
	}

	lib_log.Info(ctx, "Checked")
	return nil
}

func (c client) CheckBodyAgainstSchema(ctx context.Context, schemaFileName string, body []byte) error {
	lib_log.Info(ctx, "Checking", lib_log.FmtString("schemaFileName", schemaFileName))

	if len(body) == 0 {
		return lib_errors.NewCustom(http.StatusBadRequest, requestBodyExpectedButFoundNone)
	}

	if err := c.validateAgainstSchemaAndCheckResult(schemaFileName, gojsonschema.NewBytesLoader(body)); err != nil {
		return err
	}

	lib_log.Info(ctx, "Checked")
	return nil
}

func (c client) validateAgainstSchemaAndCheckResult(schemaFileName string, loader gojsonschema.JSONLoader) error {
	v, ok := c.schemaByFileName[schemaFileName]
	if !ok {
		return lib_errors.Errorf("Schema %q not loaded", schemaFileName)
	}

	result, err := v.Validate(loader)
	if err != nil {
		return lib_errors.NewCustom(http.StatusBadRequest, err.Error())
	}

	if !result.Valid() {
		var structuralIssues []lib_errors.Item
		var validationIssues []lib_errors.Item

		issue := func(resultError gojsonschema.ResultError) lib_errors.Item {
			item := lib_errors.Item{Message: resultError.Description()}
			if resultError.Field() == "(root)" {
				if v, ok := resultError.Details()["property"]; ok { // There may not be a property depending on the type of error
					if s, ok := v.(string); ok {
						item.Field = s
					}
				}

			} else {
				item.Field = resultError.Field()
			}
			return item
		}

		for _, v := range result.Errors() {
			switch v.Type() { // See https://github.com/xeipuuv/gojsonschema
			case "condition_else", "condition_then", "internal", "missing_dependency", "number_all_of", "number_any_of", "number_one_of", "number_not":
				structuralIssues = append(structuralIssues, issue(v))
			default:
				validationIssues = append(validationIssues, issue(v))
			}
		}

		// structuralIssues duplicate validationIssues (with technical language) and so they are not useful to end users unless there are no validation items
		// E.g.
		// -> Must validate at least one schema (anyOf)
		// -> Must validate \"then\" as \"if\" was valid
		if len(validationIssues) > 0 {
			return lib_errors.NewCustomWithItems(http.StatusBadRequest, "", validationIssues)
		}
		return lib_errors.NewCustomWithItems(http.StatusBadRequest, "", structuralIssues)
	}
	return nil
}

type datetimeChecker struct{}

func (f datetimeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	asString = strings.TrimSpace(asString)
	if asString == "" {
		return false
	}

	_, err := time.Parse(time.RFC3339Nano, asString)
	return err == nil
}

type notEmptyChecker struct{}

func (f notEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(asString) != ""
}

type localeChecker struct{}

func (f localeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if asString == "" {
		return false
	}

	return checkLocale(asString)
}

type localeOrEmptyChecker struct{}

func (f localeOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if asString == "" {
		return true
	}

	return checkLocale(asString)
}

func checkLocale(l string) bool {
	if len(l) == 2 {
		return lib_formatters.IsLowerCase(l)
	}

	if len(l) != 5 {
		return false
	}

	if len(strings.Split(l, "-")) == 2 {
		c := strings.Split(l, "-")
		if !lib_formatters.IsLowerCase(c[0]) {
			return false
		}
		if !lib_formatters.IsUpperCase(c[1]) {
			return false
		}
		return true
	}

	if len(strings.Split(l, "_")) == 2 {
		c := strings.Split(l, "_")
		if !lib_formatters.IsLowerCase(c[0]) {
			return false
		}
		if !lib_formatters.IsUpperCase(c[1]) {
			return false
		}
		return true
	}

	return false
}

type timeChecker struct{}

func (f timeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	asString = strings.TrimSpace(asString)
	if asString == "" {
		return false
	}

	_, err := time.Parse(time.RFC3339Nano, asString)
	return err == nil
}

type localTimeChecker struct{}

func (f localTimeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	asString = strings.TrimSpace(asString)
	if asString == "" {
		return false
	}

	_, err := time.Parse(lib_time.LayoutLocal, asString)
	return err == nil
}

type timeOrEmptyChecker struct{}

func (f timeOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if asString == "" {
		return true
	}

	_, err := time.Parse(time.RFC3339, asString)
	return err == nil
}

type hourOrEmptyChecker struct{}

func (f hourOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if asString == "" {
		return true
	}
	return hourRegExpCompile.MatchString(asString)
}

type moneyStringChecker struct{}

func (f moneyStringChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return moneyStringRegExpCompile.MatchString(asString)
}

type hourAndNotEmptyChecker struct{}

func (f hourAndNotEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return hourRegExpCompile.MatchString(asString)
}

type languageCodeChecker struct{}

func (f languageCodeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if len(asString) != 2 {
		return false
	}
	for _, r := range asString {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

type languageCodeOrEmptyChecker struct{}

func (f languageCodeOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if asString == "" {
		return true
	}

	if len(asString) != 2 {
		return false
	}
	for _, r := range asString {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

type countryCodeChecker struct{}

func (f countryCodeChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	asString = strings.TrimSpace(asString)

	if len(asString) != 2 {
		return false
	}
	for _, r := range asString {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

type currencyChecker struct{}

func (f currencyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	if len(asString) != 3 {
		return false
	}
	for _, r := range asString {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

type constantCaseChecker struct{}

func (f constantCaseChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return constantCaseRegExpCompile.MatchString(asString)
}

type emailChecker struct{}

func (f emailChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	return lib_email.IsEmailLowercase(asString)
}

type emailOrEmptyChecker struct{}

func (f emailOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if asString == "" {
		return true
	}

	return lib_email.IsEmailLowercase(asString)
}

type emailWithOptionalNameChecker struct{}

func (f emailWithOptionalNameChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	return lib_email.IsEmailWithOptionalNameLowercase(asString)
}

type websiteChecker struct{}

func (f websiteChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return websiteRegExpCompile.MatchString(asString)
}

type websiteOrEmptyChecker struct{}

func (f websiteOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if asString == "" {
		return true
	}

	return websiteRegExpCompile.MatchString(asString)
}

type phoneNumberChecker struct{}

func (f phoneNumberChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return phoneNumberRegExpCompile.MatchString(asString)
}

type phoneNumberOrEmptyChecker struct{}

func (f phoneNumberOrEmptyChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if asString == "" {
		return true
	}
	return phoneNumberRegExpCompile.MatchString(asString)
}

type timezoneUtcOffsetChecker struct{}

func (f timezoneUtcOffsetChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	return timeZoneUtcOffsetExpCompile.MatchString(asString)
}

type timezoneNameChecker struct{}

func (f timezoneNameChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}

	if asString == "" {
		return false
	}

	if _, err := time.LoadLocation(asString); err != nil {
		return false
	}

	return true
}

type domainChecker struct{}

func (f domainChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}

	return lib_domain.Recognized(s)
}

type domainFromOrToChecker struct{}

func (f domainFromOrToChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}

	return lib_domain.RecognizedFromOrTo(s)
}

type dateChecker struct{}

func (f dateChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}
	if s == "" {
		return false
	}

	if _, err := civil.ParseDate(s); err != nil {
		return false
	}

	return true
}

type dateOrEmptyChecker struct{}

func (f dateOrEmptyChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}
	if s == "" {
		return true
	}

	if _, err := civil.ParseDate(s); err != nil {
		return false
	}

	return true
}

type extensionChecker struct{}

func (f extensionChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}
	if s == "" {
		return false
	}

	if len(s) > 10 {
		return false
	}
	for _, v := range s {
		if !unicode.IsNumber(v) {
			return false
		}
	}

	return true
}

type extensionOrEmptyChecker struct{}

func (f extensionOrEmptyChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}
	if s == "" {
		return true
	}

	if len(s) > 10 {
		return false
	}
	for _, v := range s {
		if !unicode.IsNumber(v) {
			return false
		}
	}

	return true
}
