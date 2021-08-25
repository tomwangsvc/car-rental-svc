package spanner

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

// DeprecatedFieldMetadata gathers certain meta-data for fields in a struct
type DeprecatedFieldMetadata struct {
	Filters    []string
	JsonTag    string
	NoCopyFrom bool
	NoCopyTo   bool
	Transform  string
	Type       string
}

// This function is deprecated since the copy to/from funcs are too magic. Please use a manual copy or, if required for generating json, use lib_json.
func DeprecatedStructFieldMetadata(t reflect.Type) []DeprecatedFieldMetadata {

	var metadata []DeprecatedFieldMetadata

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if jsonFieldName := field.Tag.Get("json"); jsonFieldName != "" && jsonFieldName != "-" {

			var filters []string
			if filter := field.Tag.Get("filter"); filter != "" {
				filters = strings.Split(filter, ",")
			}

			var jsonTag string
			var noCopyFrom, noCopyTo bool
			for i, v := range strings.Split(jsonFieldName, ",") {
				if i == 0 {
					jsonTag = v
				}
				if v == "nocopyfrom" {
					noCopyFrom = true
				}
				if v == "nocopyto" {
					noCopyTo = true
				}
			}

			metadata = append(metadata, DeprecatedFieldMetadata{
				Filters:    filters,
				JsonTag:    jsonTag,
				NoCopyFrom: noCopyFrom,
				NoCopyTo:   noCopyTo,
				Transform:  field.Tag.Get("transform"),
				Type:       field.Type.Name(),
			})

		} else if field.Type.Kind() == reflect.Struct {
			metadata = append(metadata, DeprecatedStructFieldMetadata(field.Type)...)
		}
	}

	return metadata
}

// This function is deprecated since it is too magic. Please use a manual copy.
func DeprecatedCopyFromWithConversion(dst interface{}, src interface{}, srcFs []DeprecatedFieldMetadata) error {

	o, err := json.Marshal(src)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to marshal src")
	}

	blob := string(o)
	for _, srcF := range srcFs {
		if blob, err = copyFromNullable(blob, srcF.JsonTag, srcF.Type, srcF.NoCopyFrom); err != nil {
			return lib_errors.Wrap(err, "Failed copying from nullables")
		}
	}

	err = json.Unmarshal([]byte(blob), dst)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to unmarshal into dst")
	}

	return nil
}

func copyFromNullable(blob, srcJsonTag, srcFType string, srcNoCopyFrom bool) (string, error) {

	result := gjson.Get(blob, srcJsonTag)
	if !result.Exists() {
		return blob, nil
	}

	if srcNoCopyFrom {
		var err error
		if blob, err = sjson.Delete(blob, srcJsonTag); err != nil {
			return "", lib_errors.Wrapf(err, "Failed deleting nocopyfrom at tag %q", srcJsonTag)
		}
		return blob, nil
	}

	nullable, valid, value, err := unmarshallIfNullType([]byte(result.Raw), srcFType)
	if err != nil {
		return "", lib_errors.Wrapf(err, "Failed unmarshalling if null type")
	}

	if nullable {
		var err error
		if valid {
			if blob, err = sjson.Set(blob, srcJsonTag, value); err != nil {
				return "", lib_errors.Wrapf(err, "Failed setting field at tag %q with value %v", srcJsonTag, value)
			}

		} else {
			if blob, err = sjson.Delete(blob, srcJsonTag); err != nil {
				return "", lib_errors.Wrapf(err, "Failed deleting field at tag %q", srcJsonTag)
			}
		}
	}

	return blob, nil
}

// This function is deprecated since it is too magic. Please use a manual copy.
func DeprecatedCopyToWithConversion(dst interface{}, dstFs []DeprecatedFieldMetadata, src interface{}) error {

	o, err := json.Marshal(src)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to marshal src")
	}

	blob := string(o)
	for _, dstF := range dstFs {
		if blob, err = copyToNullable(blob, dstF.JsonTag, dstF.Type, dstF.NoCopyTo); err != nil {
			return lib_errors.Wrap(err, "Failed copying into nullables")
		}
	}

	err = json.Unmarshal([]byte(blob), dst)
	if err != nil {
		return lib_errors.Wrap(err, "Failed to unmarshal into dst")
	}

	return nil
}

func copyToNullable(blob, dstJsonTag, dstType string, dstNoCopyTo bool) (string, error) {

	result := gjson.Get(blob, dstJsonTag)
	if !result.Exists() {
		return blob, nil
	}

	if dstNoCopyTo {
		var err error
		if blob, err = sjson.Delete(blob, dstJsonTag); err != nil {
			return "", lib_errors.Wrapf(err, "Failed deleting nocopyto at tag %q", dstJsonTag)
		}
		return blob, nil
	}

	nullable, _, value, err := unmarshallIfPointer([]byte(result.Raw), dstType)
	if err != nil {
		return "", lib_errors.Wrapf(err, "Failed unmarshalling if pointer")
	}

	if nullable {
		var err error
		if blob, err = sjson.Set(blob, dstJsonTag, value); err != nil {
			return "", lib_errors.Wrapf(err, "Failed setting field at tag %q with value %v", dstJsonTag, value)
		}
	}

	return blob, nil
}

func unmarshallIfNullType(raw []byte, typeOf string) (nullable, valid bool, value interface{}, err error) {

	nullable = true

	switch typeOf {
	default:
		nullable = false

	case "NullBool":
		var v spanner.NullBool
		err = json.Unmarshal(raw, &v)
		value = v.Bool
		valid = v.Valid

	case "NullFloat64":
		var v spanner.NullFloat64
		err = json.Unmarshal(raw, &v)
		value = v.Float64
		valid = v.Valid

	case "NullInt64":
		var v spanner.NullInt64
		err = json.Unmarshal(raw, &v)
		value = v.Int64
		valid = v.Valid

	case "NullString":
		var v spanner.NullString
		err = json.Unmarshal(raw, &v)
		value = v.StringVal
		valid = v.Valid

	case "NullTime":
		var v spanner.NullTime
		err = json.Unmarshal(raw, &v)
		value = v.Time
		valid = v.Valid

	case "NullDate":
		// If field NullDate in spanner is empty, its value is '0000-00-00' which cannot be unmarshalled into NullDate, but valid can be checked.
		var v struct {
			Valid bool
		}
		err = json.Unmarshal(raw, &v)
		if v.Valid {
			var v spanner.NullDate
			err = json.Unmarshal(raw, &v)
			value = v.Date
			valid = true
		}
	}

	if err != nil {
		err = lib_errors.Wrapf(err, "Failed unmarshalling field type %q from %v", typeOf, raw)
		return
	}

	return
}

//revive:disable:cyclomatic
func unmarshallIfPointer(raw []byte, typeOf string) (nullable, valid bool, value interface{}, err error) {

	nullable = true

	switch typeOf {
	default:
		nullable = false

	case "NullBool":
		var v *bool
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullBool
		if v != nil {
			nullVal = spanner.NullBool{Bool: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid

	case "NullFloat64":
		var v *float64
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullFloat64
		if v != nil {
			nullVal = spanner.NullFloat64{Float64: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid

	case "NullInt64":
		var v *int64
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullInt64
		if v != nil {
			nullVal = spanner.NullInt64{Int64: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid

	case "NullString":
		var v *string
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullString
		if v != nil {
			nullVal = spanner.NullString{StringVal: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid

	case "NullTime":
		var v *time.Time
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullTime
		if v != nil {
			nullVal = spanner.NullTime{Time: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid

	case "NullDate":
		var v *civil.Date
		err = json.Unmarshal(raw, &v)
		var nullVal spanner.NullDate
		if v != nil {
			nullVal = spanner.NullDate{Date: *v, Valid: true}
		}
		value = nullVal
		valid = nullVal.Valid
	}
	if err != nil {
		err = lib_errors.Wrapf(err, "Failed unmarshalling pointer %q into field type %s", raw, typeOf)
		return
	}

	return
	//revive:enable:cyclomatic
}
