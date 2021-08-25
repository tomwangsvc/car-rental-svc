package json

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lib_constants "github.com/tomwangsvc/lib-svc/constants"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_finance "github.com/tomwangsvc/lib-svc/finance"
)

func CopyFields(from, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return lib_errors.Wrap(err, "Failed marshalling from into bytes")
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		return lib_errors.Wrap(err, "Failed unmarshalling bytes into to")
	}

	return nil
}

type Redaction func(result string) string

func DefaultRedaction(_ string) string {
	return lib_constants.Redacted
}

func ApplyRedactions(data []byte, redactionsByPath map[string]Redaction) ([]byte, error) {
	for path, redaction := range redactionsByPath {
		result := gjson.GetBytes(data, path)
		if result.Exists() {
			var err error
			data, err = sjson.SetBytes(data, path, redaction(result.Str))
			if err != nil {
				return nil, lib_errors.Wrap(err, "Failed setting raw bytes")
			}
		}
	}
	return data, nil
}

func GenerateJson(spannerStruct interface{}, fieldMetadata []FieldMetadata, filterBy string) ([]byte, error) {
	bytes, err := json.Marshal(spannerStruct)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed marshalling struct to bytes")
	}
	blob := string(bytes)
	for _, f := range fieldMetadata {
		if blob, err = filter(blob, f.JsonTag, filterBy, f.Filters); err != nil {
			return nil, lib_errors.Wrap(err, "Failed filtering struct")
		}
		if blob, err = transform(blob, f.JsonTag, f.Transform); err != nil {
			return nil, lib_errors.Wrap(err, "Failed transforming struct")
		}
	}
	return []byte(blob), nil
}

func GenerateJsonList(spannerStructs []interface{}, fieldMetadata []FieldMetadata, filterBy string) ([]byte, error) {
	bytes := []byte{'['}
	for i, v := range spannerStructs {
		b, err := GenerateJson(v, fieldMetadata, filterBy)
		if err != nil {
			return nil, lib_errors.Wrapf(err, "Failed generating JSON at index %d", i)
		}
		if i > 0 {
			bytes = append(bytes, ',')
		}
		bytes = append(bytes, b...)
	}
	bytes = append(bytes, ']')
	return bytes, nil
}

type FieldMetadata struct {
	Filters   []string
	JsonTag   string
	Transform string
	Type      string
}

// StructFieldMetadata returns FieldMetadata for all fields in the passed struct type
//
// Filter tags
//   supported:
//     `filter:"xxx"` filters out field if filterBy != "xxx" during GenerateJson funcs
//     `filter:"-"` filters out field during GenerateJson funcs
//     `filter:""` or no filter tag never filters out field during GenerateJson funcs
//
// Transform tags
//   supported:
//     `transform:"money"` transforms float64 to money value during GenerateJson funcs
//     `transform:"raw"` uses the raw value during GenerateJson funcs, useful for extracting json from a string field
func StructFieldMetadata(t reflect.Type) []FieldMetadata {
	var metadata []FieldMetadata
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if jsonFieldName := field.Tag.Get("json"); jsonFieldName != "" && jsonFieldName != "-" {
			var filters []string
			if filter := field.Tag.Get("filter"); filter != "" {
				filters = strings.Split(filter, ",")
			}
			var jsonTag string
			for i, v := range strings.Split(jsonFieldName, ",") {
				if i == 0 {
					jsonTag = v
				}
			}
			metadata = append(metadata, FieldMetadata{
				Filters:   filters,
				JsonTag:   jsonTag,
				Transform: field.Tag.Get("transform"),
				Type:      field.Type.Name(),
			})
		} else if field.Type.Kind() == reflect.Struct {
			metadata = append(metadata, StructFieldMetadata(field.Type)...)
		}
	}
	return metadata
}

func filter(blob, path, filterBy string, filters []string) (string, error) {
	result := gjson.Get(blob, path)
	if !result.Exists() {
		return blob, nil
	}
	if shouldBeFilteredOut(filters, filterBy) {
		var err error
		if blob, err = sjson.Delete(blob, path); err != nil {
			return "", lib_errors.Wrapf(err, "Failed deleting field of struct at %q", path)
		}
	}
	return blob, nil
}

func shouldBeFilteredOut(filters []string, filterBy string) bool {
	if len(filters) == 0 {
		return false
	}
	for _, f := range filters {
		if f == "-" {
			return true
		}
	}
	if filterBy == "" {
		return false
	}
	gs := strings.Split(strings.TrimSpace(filterBy), ",")
	for _, f := range filters {
		for _, g := range gs {
			if strings.ToLower(strings.TrimSpace(f)) == strings.ToLower(g) {
				return false
			}
		}
	}
	return true
}

func transform(blob, path, transform string) (string, error) {
	result := gjson.Get(blob, path)
	if !result.Exists() {
		return blob, nil
	}

	if result.Raw == "null" {
		var err error
		if blob, err = sjson.Delete(blob, path); err != nil {
			return "", lib_errors.Wrapf(err, "Failed deleting field at path %q", path)
		}
		return blob, nil
	}

	result = gjson.Get(blob, path)
	if !result.Exists() {
		return blob, nil
	}

	var raw string
	switch transform {
	default:
		return blob, nil
	case "money":
		raw = transformMoney(result)
	case "raw":
		raw = transformRaw(result)
	}

	var err error
	if blob, err = sjson.SetRaw(blob, path, raw); err != nil {
		return "", lib_errors.Wrap(err, "Failed setting new raw value")
	}

	return blob, nil
}

func transformMoney(result gjson.Result) string {
	return lib_finance.DisplayMoney(result.Num)
}

func transformRaw(result gjson.Result) string {
	return result.Str
}
