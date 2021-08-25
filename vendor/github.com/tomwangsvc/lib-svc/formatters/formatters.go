package formatters

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

func IsLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func IsUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func FormatSentenceCase(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return ""
	}
	s = ReplaceWhiteSpaceWithSpace(s)

	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// FormatTitleCase is useful because 'strings.Title' converts 'Peter's socks' to 'Peter'S socks'
func FormatTitleCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = ReplaceWhiteSpaceWithSpace(s)

	var b strings.Builder
	var previous rune

	for i, v := range s {
		if i == 0 {
			b.WriteString(strings.ToUpper(string(v)))

		} else {
			if unicode.IsSpace(previous) {
				b.WriteString(strings.ToUpper(string(v)))

			} else {
				b.WriteString(strings.ToLower(string(v)))
			}
		}
		previous = v
	}
	return b.String()
}

// FormatStringGroupInFour inserts a spacer every forth character
func FormatStringGroupInFour(str, spacer string) string {
	str = StripSpace(str)

	paddedStr := ""

	count := 0
	for _, v := range str {
		if count > 0 && count%4 == 0 {
			paddedStr = fmt.Sprintf("%s%s%s", paddedStr, spacer, string(v))

		} else {
			paddedStr = fmt.Sprintf("%s%s", paddedStr, string(v))
		}
		count++
	}
	return strings.TrimSpace(strings.TrimPrefix(paddedStr, spacer))
}

// StripSpace strips spaces from the passed string
func StripSpace(str string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str))
}

// StripSpaceAndPunctuation strips spaces and punctuations from the passed string
func StripSpaceAndPunctuation(str string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, str))
}

// StripPunctuation strips punctuations from the passed string
func StripPunctuation(str string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, str))
}

// Replace punctuation with white space
func ReplacePunctuationWithWhiteSpace(str string) string {
	return ReplaceWhiteSpaceWithSpace(strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, str))
}

// Replace punctuation with underscore character
func ReplacePunctuationWithUnderscore(str string) string {
	return ReplaceWhiteSpaceWithSpace(strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, str))
}

//StripUspsZipPlus4CodeInPostalCode strips the USPS Zip+4 code from a postal code, see https://smartystreets.com/articles/zip-4-code. For example, 34561-9721 will be 34561 after applying the func. Only US postal codes should be used with this func.
func StripUspsZipPlus4CodeInPostalCode(str string) (string, error) {
	result := str
	r, err := regexp.Compile("^[0-9]{5}-[0-9]{4}$")
	if err != nil {
		return "", lib_errors.Errorf("Failed compiling regular expression")
	}
	if r.MatchString(str) {
		result = strings.Split(str, "-")[0]
	}

	return result, nil
}

// ReplaceWhiteSpaceWithSpace replaces new lines tabs etc from the passed string with a single space
// -> Useful for stripping newlines from address elments submitted in an api
func ReplaceWhiteSpaceWithSpace(str string) string {
	s := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, str)
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

// IsPhoneNumberSearchFormat checks if it is format of phone number
func IsPhoneNumberSearchFormat(str string) bool {
	if len(str) < 6 {
		return false
	}
	for _, v := range str {
		if !unicode.IsNumber(v) && !unicode.IsSpace(v) {
			return false
		}
	}
	return true
}

// IsPhoneNumberExtensionFormat checks if it is format of phone number extension
func IsPhoneNumberExtensionFormat(str string) bool {
	if len(str) > 10 {
		return false
	}
	for _, v := range str {
		if !unicode.IsNumber(v) {
			return false
		}
	}
	return true
}

// ExtractNameFromEmail returns a string with an email e.g. tom.wang@tomwang.com => Tom Wang
func ExtractNameFromEmail(email string) (string, error) {
	if email == "" {
		return "", lib_errors.NewCustom(http.StatusBadRequest, "Email should not be empty")
	}

	i := strings.Index(email, "@")
	if i == -1 {
		return "", lib_errors.NewCustom(http.StatusBadRequest, "Email should not right format")
	}
	emailPrefix := email[:i]

	var name []string
	a := strings.Split(emailPrefix, ".")
	for _, v := range a {
		name = append(name, strings.Title(strings.ToLower(v)))
	}

	return strings.Join(name, " "), nil
}
