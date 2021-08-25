package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"cloud.google.com/go/spanner"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/pkg/errors"
	lib_reflect "github.com/tomwangsvc/lib-svc/reflect"
	"google.golang.org/grpc/codes"
)

const (
	MetadataKeyIndex                             = "__INDEX__"
	MetadataKeyTruncate                          = "__TRUNCATE__"
	MetadataKeyValueConflictUniqueIndexViolation = "UNIQUE_INDEX_VIOLATION"
	MetadataTagRaw                               = "__RAW__"
	ConflictFromSpannerPrimaryKeyViolation       = "OBJECT_ALREADY_EXISTS"
	ConflictFromSpannerUniqueIndexViolation      = "OBJECT_CANNOT_BE_CREATED_OR_UPDATED"

	fieldLabel   = "field"
	messageLabel = "message"
)

type Item struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func UniqueItems(sample []Item) []Item {
	var unique []Item
	type key struct{ value1, value2 string }
	m := make(map[key]int)
	for _, v := range sample {
		k := key{v.Field, v.Message}
		if i, ok := m[k]; ok {
			// Overwrite previous value per requirement in question to keep last matching value.
			unique[i] = v
		} else {
			// Unique key found. Record position and collect in result.
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func IsDownstream(err error) bool {
	return spanner.ErrCode(err) == codes.Aborted || spanner.ErrCode(err) == codes.Internal
}

func IsTimeout(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, context.DeadlineExceeded.Error()) ||
		strings.Contains(msg, context.Canceled.Error()) ||
		spanner.ErrCode(err) == codes.DeadlineExceeded
}

func PublicString(err error) string {
	var msg string
	if cerr, ok := err.(Custom); ok {
		msg = cerr.PublicError()
	} else {
		msg = err.Error()
	}
	return msg
}

// FormattedMessageWithInlineMetadata creates a formatted message from the message and metadata of a custom error or runs the default stringify for any other error
// Metadata are sorted for elegance but this also helps with testing since map keys are not sorted so comparisons do not work
func FormattedMessageWithInlineMetadata(err error) string {

	if cerr, ok := err.(Custom); ok {
		// Convert enum into something more friendly
		msg := sentence(cerr.Message)
		if msg != "" && !strings.HasSuffix(msg, ".") {
			msg = msg + "."
		}

		if len(cerr.Metadata) == 0 {
			return msg
		}

		flattenedMetadata, _ := FlattenMetadata(cerr.Metadata)

		metadataSortedKeys := make([][]string, 0, len(flattenedMetadata))

		for _, v := range flattenedMetadata {
			var keys []string
			for key := range v {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			metadataSortedKeys = append(metadataSortedKeys, keys)
		}

		var metadataMsg strings.Builder
		for i, v := range metadataSortedKeys {
			metadata := flattenedMetadata[i]
			for _, key := range v {
				if key != fieldLabel {
					metadataMsg.WriteString(fmt.Sprintf(" | %s", metadata[key]))
				}
			}
		}

		s := strings.Replace(metadataMsg.String(), " | ", "", 1)

		if msg == "" {
			return s
		}

		return fmt.Sprintf("%s %s", msg, s)
	}

	return err.Error()
}

// FormattedMessageWithInlineItems creates a formatted message from the message and items of a custom error or runs the default stringify for any other error
func FormattedMessageWithInlineItems(err error) string {

	if cerr, ok := err.(Custom); ok {
		// Convert enum into something more friendly
		msg := sentence(cerr.Message)
		if msg != "" && !strings.HasSuffix(msg, ".") {
			msg = msg + "."
		}

		if len(cerr.Items) == 0 {
			return msg
		}

		var itemsMsg strings.Builder
		for _, v := range cerr.Items {
			if strings.HasPrefix(v.Message, MetadataTagRaw) {
				itemsMsg.WriteString(fmt.Sprintf(" | %s", v.Message[len(MetadataTagRaw):]))
			} else {
				itemsMsg.WriteString(fmt.Sprintf(" | %s", v.Message))
			}
		}

		s := strings.Replace(itemsMsg.String(), " | ", "", 1)

		if msg == "" {
			return s
		}

		return fmt.Sprintf("%s %s", msg, s)
	}

	return err.Error()
}

func sentence(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToLower(strings.Replace(s, "_", " ", -1))
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

type Custom struct {
	At              string                 `json:"at"`
	Cause           error                  `json:"cause,omitempty"`
	Code            int                    `json:"code"`
	Items           []Item                 `json:"items,omitempty"`
	Message         string                 `json:"message,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	WrappedMessages []string               `json:"wrapped_messages,omitempty"`
}

// Only used for proxies as we need additional information
// TODO: should be deprecated when we've moved entirely away from proxies
func NewCustomFromHttpResponse(code int, body []byte) (none Custom, err error) {
	if code == http.StatusBadRequest ||
		code == http.StatusConflict ||
		code == http.StatusTeapot ||
		code == http.StatusUnprocessableEntity {

		var items []Item
		if err = json.Unmarshal(body, &items); err != nil {
			return none, Wrap(err, "Failed unmarshalling HTTP response body to error items")
		}
		return NewCustomWithItems(code, "", items), nil
	}
	return NewCustomf(code, "HTTP response not recognized as Custom error: %v", string(body)), nil
}

func NewCustomf(code int, message string, a ...interface{}) Custom {
	return newCustom(code, fmt.Sprintf(message, a...), nil, nil, nil, nil)
}

func NewCustom(code int, message string) Custom {
	return newCustom(code, message, nil, nil, nil, nil)
}

func NewCustomWithCause(code int, message string, cause error) Custom {
	return newCustom(code, message, cause, nil, nil, nil)
}

func NewCustomWithCauseAndMetadata(code int, message string, cause error, metadata map[string]interface{}) Custom {
	return newCustom(code, message, cause, metadata, nil, nil)
}

func NewCustomWithMetadata(code int, message string, metadata map[string]interface{}) Custom {
	return newCustom(code, message, nil, metadata, nil, nil)
}

func NewCustomWithItems(code int, message string, items []Item) Custom {
	return newCustom(code, message, nil, nil, items, nil)
}

func NewCustomWithWrappedMessages(code int, message string, wrappedMessages []string) Custom {
	return newCustom(code, message, nil, nil, nil, wrappedMessages)
}

func newCustom(code int, message string, cause error, metadata map[string]interface{}, items []Item, wrappedmessages []string) Custom {
	switch code {
	default:
		status := http.StatusText(code)
		if message == "" {
			message = status
		} else {
			if _, err := strconv.ParseInt(message, 10, 64); err != nil { // Test for possible functional error codes
				if strings.HasPrefix(strings.ToLower(message), strings.ToLower(status)) {
					message = fmt.Sprintf("%s%s", status, message[len(status):]) // Cleanup the casing
				} else {
					message = fmt.Sprintf("%s: %s", status, message)
				}
			}
		}
	case http.StatusBadRequest, http.StatusConflict, http.StatusTeapot, http.StatusUnprocessableEntity:
		// http.StatusConflict, http.StatusTeapot and http.StatusUnprocessableEntity codes are not manipulated because the accompanying message is generally a functional error code
		// http.StatusBadRequest code is not manipulated because it would become a redundant message when rendered
	}
	if len(items) == 0 {
		items = nil
	}
	if len(metadata) == 0 {
		metadata = nil
	}
	return Custom{
		At:              lib_reflect.At(2),
		Cause:           cause,
		Code:            code,
		Items:           items,
		Message:         message,
		Metadata:        metadata,
		WrappedMessages: wrappedmessages,
	}
}

func (c Custom) Error() string {
	e := fmt.Sprintf("Code: %d, Message: %s, At: %s", c.Code, c.Message, c.At)
	if c.Cause != nil {
		e = fmt.Sprintf("%s, Cause: %+v", e, c.Cause)
	}
	if len(c.WrappedMessages) > 0 {
		e = fmt.Sprintf("%s, WrappedMessages: %v", e, strings.Join(c.WrappedMessages, ", "))
	}
	if len(c.Items) > 0 {
		e = fmt.Sprintf("%s, Items: %v", e, c.Items)
	}
	if c.Metadata != nil {
		e = fmt.Sprintf("%s, Metadata: %v", e, c.Metadata)
	}
	return e
}

func (c Custom) PublicError() string {
	e := fmt.Sprintf("Code: %d, Message: %s", c.Code, c.Message)
	if len(c.Items) > 0 {
		e = fmt.Sprintf("%s, Items: %v", e, c.Items)
	}
	if c.Metadata != nil {
		e = fmt.Sprintf("%s, Metadata: %v", e, c.Metadata)
	}
	return e
}

func ItemsFromRenderedCustomError(body []byte) ([]Item, error) {
	if lenBody := len(body); lenBody == 0 || lenBody == 1 && string(body) == "" {
		return nil, nil
	}
	var items []Item
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, Wrap(err, "Failed unmarshalling into []Item")
	}
	return items, nil
}

func (c Custom) Render() (body []byte, warnings []string) {
	var dataLostFromMisusage []map[string]string
	switch c.Code {
	default:
		return nil, nil
	case http.StatusBadRequest,
		http.StatusConflict,
		http.StatusNotFound,
		http.StatusPreconditionFailed,
		http.StatusTeapot,
		http.StatusUnprocessableEntity:

		var cerr []map[string]string
		if c.Message != "" {
			// Yes this is correct, an error message is a comma separated list of enums, plain text messages are not supported because the system is global
			for _, message := range strings.Split(c.Message, ",") {
				cerr = append(cerr, map[string]string{
					messageLabel: message,
				})
			}
		}
		if len(c.Items) > 0 {
			for _, v := range c.Items {
				m := map[string]string{
					messageLabel: v.Message,
				}
				if v.Field != "" {
					m[fieldLabel] = v.Field
				}
				cerr = append(cerr, m)
			}
		}
		if c.Metadata != nil {
			var flattenedMetadata []map[string]string
			flattenedMetadata, dataLostFromMisusage = FlattenMetadata(c.Metadata)
			cerr = append(cerr, flattenedMetadata...)
		}

		if len(dataLostFromMisusage) > 0 {
			for _, val := range dataLostFromMisusage {
				for k, v := range val {
					warnings = append(warnings, fmt.Sprintf("Lost metadata: %q = %q", k, v))
				}
			}
		}

		var err error
		body, err = json.Marshal(cerr)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("Failed json marshalling []map[string]string: %+v", err))
		}

		return
	}
}

func FlattenMetadata(metadata map[string]interface{}) (flattenedMetadata []map[string]string, dataLostFromMisusage []map[string]string) {
	var flatMap map[string]string
	flatMap, dataLostFromMisusage = replaceIndexInFlatMap(flatmap.Flatten(metadata))

	keys := make([]string, 0, len(flatMap))
	for key := range flatMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		if !strings.Contains(key, ".#") {
			if truncatePrefix := fmt.Sprintf("%s.", MetadataKeyTruncate); strings.HasPrefix(key, truncatePrefix) {
				key = strings.Replace(key, truncatePrefix, "", 1)
			}
			if strings.HasPrefix(flatMap[keys[i]], MetadataTagRaw) {
				flattenedMetadata = append(flattenedMetadata, map[string]string{
					fieldLabel:   key,
					messageLabel: flatMap[keys[i]][len(MetadataTagRaw):],
				})

			} else {
				for _, splitVal := range strings.Split(flatMap[keys[i]], ",") {
					flattenedMetadata = append(flattenedMetadata, map[string]string{
						fieldLabel:   key,
						messageLabel: splitVal,
					})
				}
			}
		}
	}

	return
}

//revive:disable:cyclomatic
func replaceIndexInFlatMap(flatMap map[string]string) (updatedFlatMap map[string]string, dataLostFromMisusage []map[string]string) {
	updatedFlatMap = make(map[string]string)
	for k, v := range flatMap {
		updatedFlatMap[k] = v
	}

	for {
		var keyOfIndexForReplacement string
		greatestIndexValueIfInt := int64(-1)
		for key, val := range updatedFlatMap {
			if strings.HasSuffix(key, MetadataKeyIndex) {
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					keyOfIndexForReplacement = key
					continue
				}
				if valInt > greatestIndexValueIfInt {
					greatestIndexValueIfInt = valInt
					keyOfIndexForReplacement = key
				}
			}
		}
		if keyOfIndexForReplacement == "" {
			break
		}

		tokensOfKeyOfIndexForReplacement := strings.Split(keyOfIndexForReplacement, ".")
		if len(tokensOfKeyOfIndexForReplacement) < 2 {
			dataLostFromMisusage = append(dataLostFromMisusage, map[string]string{
				keyOfIndexForReplacement: updatedFlatMap[keyOfIndexForReplacement],
			})
			delete(updatedFlatMap, keyOfIndexForReplacement)
			continue
		}
		if _, err := strconv.ParseInt(tokensOfKeyOfIndexForReplacement[len(tokensOfKeyOfIndexForReplacement)-2], 10, 64); err != nil {
			dataLostFromMisusage = append(dataLostFromMisusage, map[string]string{
				keyOfIndexForReplacement: updatedFlatMap[keyOfIndexForReplacement],
			})
			delete(updatedFlatMap, keyOfIndexForReplacement)
			continue
		}

		keyToReplace := strings.Join(tokensOfKeyOfIndexForReplacement[:len(tokensOfKeyOfIndexForReplacement)-1], ".")
		tokensOfKeyOfIndexForReplacement[len(tokensOfKeyOfIndexForReplacement)-2] = updatedFlatMap[keyOfIndexForReplacement]
		replacementKey := strings.Join(tokensOfKeyOfIndexForReplacement[:len(tokensOfKeyOfIndexForReplacement)-1], ".")
		delete(updatedFlatMap, keyOfIndexForReplacement)

		if keyToReplace != replacementKey {
			for k, v := range updatedFlatMap {
				if strings.HasPrefix(k, keyToReplace) {
					newKey := strings.Replace(k, keyToReplace, replacementKey, 1)
					if val, ok := updatedFlatMap[newKey]; ok {
						dataLostFromMisusage = append(dataLostFromMisusage, map[string]string{
							newKey: val,
						})
					}
					updatedFlatMap[newKey] = v
					delete(updatedFlatMap, k)
				}
			}
		}
	}

	for k, v := range updatedFlatMap {
		if strings.HasSuffix(k, MetadataKeyTruncate) {
			if split := strings.Split(k, "."); len(split) > 1 {
				delete(updatedFlatMap, k)
				updatedFlatMap[strings.Join(split[:len(split)-1], ".")] = v
			}
		}
	}
	return
	//revive:enable:cyclomatic
}

func Wrap(err error, message string) error {
	if err == nil {
		return New(message)
	}
	if customError, ok := err.(Custom); ok {
		customError.WrappedMessages = append(customError.WrappedMessages, message)
		return customError
	}
	if IsDownstream(err) {
		return transformDownstreamIntoCustomBadGateway(err)
	}
	if IsTimeout(err) {
		return transformTimeoutIntoCustomGatewayTimeout(err)
	}
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return Errorf(format, args...)
	}
	if customError, ok := err.(Custom); ok {
		customError.WrappedMessages = append(customError.WrappedMessages, fmt.Sprintf(format, args...))
		return customError
	}
	if IsDownstream(err) {
		return transformDownstreamIntoCustomBadGateway(err)
	}
	if IsTimeout(err) {
		return transformTimeoutIntoCustomGatewayTimeout(err)
	}
	return errors.Wrapf(err, format, args...)
}

func transformDownstreamIntoCustomBadGateway(err error) error {
	return Custom{
		At:      lib_reflect.At(2),
		Cause:   err,
		Code:    http.StatusBadGateway,
		Message: "Aborted",
	}
}

func transformTimeoutIntoCustomGatewayTimeout(err error) error {
	return Custom{
		At:      lib_reflect.At(2),
		Cause:   err,
		Code:    http.StatusGatewayTimeout,
		Message: "Timeout",
	}
}

func IsCustom(err error) bool {
	if _, ok := err.(Custom); ok {
		return true
	}
	return false
}

func IsCustomWithCode(err error, code int) bool {
	if cerr, ok := err.(Custom); ok {
		return cerr.Code == code
	}
	return false
}

func IsCustomUnprocessableEntityContainingMessage(err error, message string) bool {
	if cerr, ok := err.(Custom); ok {
		if cerr.Code == http.StatusUnprocessableEntity {
			if strings.Contains(cerr.Message, message) {
				return true
			}
			for _, v := range cerr.Items {
				if strings.Contains(v.Message, message) {
					return true
				}
			}
		}
		return false
	}
	return false
}
