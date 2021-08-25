package misc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_token "github.com/tomwangsvc/lib-svc/token"
)

// RecovererMaxPanics is the maximum number of panics 'Recoverer' will attempt recovery
const (
	RecovererMaxPanics   = 10
	RecovererPauseInMSec = 5000
)

// Recoverer generically recovers a goroutine that panics with a pause of RecovererPauseInMSec between attempts
// Usage: go util.Recoverer(util.RecovererMaxPanics, "someFuncIDForLogging", func() { someFunc() })
func Recoverer(maxPanics int, id string, f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("***** RECOVERED FROM PANIC *****", map[string]string{"id": id, "previousAttempts": strconv.Itoa(RecovererMaxPanics - maxPanics)})
			if maxPanics == 0 {
				panic(fmt.Sprintf("***** TOO MANY PANICS -> Attempted recovery '%d' times, '%s' has now failed *****", RecovererMaxPanics, id))
			} else {
				time.Sleep(RecovererPauseInMSec * time.Millisecond)
				go Recoverer(maxPanics-1, id, f)
			}
		}
	}()
	f()
}

// NewJsonEncoder returns json.NewEncoder with some defaults applied
func NewJsonEncoder(w io.Writer) *json.Encoder { // TODO: deprecate
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder
}

// GetFuncName returns the function name of the caller of this function
func GetFuncName() string {
	function, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(function).Name()
}

// Copy performs a deep copy of an object by converting to and from JSON format
func Copy(dst interface{}, src interface{}) error {
	o, err := json.Marshal(src)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to marshal src")
	}

	err = json.Unmarshal(o, dst)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to unmarshal into dst")
	}

	return nil
}

// GenerateIdFromFieldValue generates an ID from a field value by lower casing and stripping spaces and punctuations
// DO NOT change the implementation of this function because it is used by various services for database table constraints
func GenerateIdFromFieldValue(str string) string {
	return strings.ToLower(strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, str)))
}

// SaveFile saves a file in the current working directory
func SaveFile(fileName string, content []byte) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed to acquire current directory of process")
	}

	f := fmt.Sprintf("%s/%s", dir, fileName)

	if err := os.WriteFile(f, content, 0644); err != nil {
		return "", lib_errors.Wrapf(err, "Failed to save file '%s'", f)
	}

	return f, nil
}

// StructTaggedFieldNames returns the names of all struct fields in the passed type tagged with the passed tag
func StructTaggedFieldNames(t reflect.Type, tag string) []string {
	var fieldNames []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get(tag)
		fieldName = strings.Replace(fieldName, ",omitempty", "", -1)
		fieldName = strings.Replace(fieldName, ",noindex", "", -1)
		if fieldName != "" && fieldName != "-" {
			fieldNames = append(fieldNames, fieldName)

		} else if field.Type.Kind() == reflect.Struct {
			fieldNames = append(fieldNames, StructTaggedFieldNames(field.Type, tag)...)
		}
	}
	return fieldNames
}

// StructTaggedFieldNameTagMap returns a map of tag names indexed by struct field name of all struct fields in the passed struct tagged with the passed tag
func StructTaggedFieldNameTagMap(v interface{}, tag string) map[string]string {
	t := reflect.TypeOf(v)
	fieldNames := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := t.Field(i).Name
		fieldTagName := field.Tag.Get(tag)
		fieldTagName = strings.Replace(fieldTagName, ",omitempty", "", -1)
		fieldTagName = strings.Replace(fieldTagName, ",noindex", "", -1)
		fieldNames[fieldName] = fieldTagName
	}
	return fieldNames
}

func RedactTokenFromUrl(url string) string {
	tokenQuery := "token="
	lenTokenQuery := len(tokenQuery)
	if tokenQueryIndex := strings.Index(url, tokenQuery); tokenQueryIndex >= 0 && len(url) >= tokenQueryIndex+lenTokenQuery {
		token := url[tokenQueryIndex+lenTokenQuery:]
		var endOfUrl string
		if querySeparatorIndex := strings.Index(token, "&"); querySeparatorIndex >= 0 {
			endOfUrl = RedactTokenFromUrl(token[querySeparatorIndex:])
			token = token[:querySeparatorIndex]
		}
		return url[:tokenQueryIndex+lenTokenQuery] + lib_token.Redact(token) + endOfUrl
	}
	return url
}

func GenerateEtag(bytes []byte) string {
	return fmt.Sprintf("%x", (md5.Sum(bytes)))
}
