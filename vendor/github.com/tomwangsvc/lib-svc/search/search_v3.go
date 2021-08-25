package search

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_validation "github.com/tomwangsvc/lib-svc/validation"
)

func ParseQueryWithTestV3(encodedQueryValue string, test bool) (filters []Filter, linkedFilters []LinkedFilter, err error) {
	filters, linkedFilters, err = ParseQueryV3(encodedQueryValue)
	if err != nil {
		if lib_errors.IsCustomWithCode(err, http.StatusBadRequest) {
			return // Expose the real error
		}
		err = lib_errors.Wrapf(err, "Failed parsing encodedQueryValue %q", encodedQueryValue)
		return
	}
	if len(linkedFilters) > 0 {
		linkedFilterTypeAnd := LinkedFilterTypeAnd
		linkedFilters = append(linkedFilters, LinkedFilter{
			Type: &linkedFilterTypeAnd,
		})
	}
	linkedFilters = append(linkedFilters,
		LinkedFilter{
			Filter: &Filter{
				Key:   "test",
				Value: test,
			},
		},
	)
	return
}

func ParseQueryV3(encodedQueryValue string) (filters []Filter, linkedFilters []LinkedFilter, err error) {
	if encodedQueryValue == "" {
		return
	}

	decodedQueryValue, err := base64.StdEncoding.DecodeString(encodedQueryValue)
	if err != nil {
		err = lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestQueryIsNotBase64Encoded, lib_errors.Wrapf(err, "Failed decoding query value encodedQueryValue from Base64 %q", encodedQueryValue))
		return
	}

	d := json.NewDecoder(bytes.NewReader(decodedQueryValue))
	d.UseNumber()
	if err = d.Decode(&linkedFilters); err != nil {
		err = lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestCannotParseQuery, lib_errors.Wrapf(err, "Failed unmarshalling decoded query value decodedQueryValue into linkedFilters %q", decodedQueryValue))
		return
	}

	// if err = json.Unmarshal(decodedQueryValue, &linkedFilters); err != nil {
	// 	err = lib_errors.Wrap(err, "Failed json unmarshalling linked filters")
	// 	return
	// }

	err = checkFilters(linkedFilters)
	if err != nil {
		if lib_errors.IsCustomWithCode(err, http.StatusBadRequest) {
			return // Expose the real error
		}
		err = lib_errors.Wrap(err, "Failed checking filters")
		return
	}

	for _, v := range linkedFilters {
		if v.Filter != nil {
			var value interface{}
			value, err = convertValue(v.Filter)
			if err != nil {
				if lib_errors.IsCustomWithCode(err, http.StatusBadRequest) {
					return // Expose the real error
				}
				err = lib_errors.Wrapf(err, "Failed converting filter value to value with type: %v", v.Filter)
				return
			}
			v.Filter.Value = value

			filters = append(filters, *v.Filter)
		}
	}

	return
}

//revive:disable:cyclomatic
func checkFilters(linkedFilters []LinkedFilter) error {
	openBracketCount := 0
	closeBracketCount := 0

	if len(linkedFilters) < 3 {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryDidNotContainAtLeastOneFilterWithTheRequiredBrackets)
	}

	if linkedFilters[0].Type == nil || *linkedFilters[0].Type != LinkedFilterTypeOpenBracket {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryDidNotStartWithOpenBracket)
	}

	if linkedFilters[len(linkedFilters)-1].Type == nil || *linkedFilters[len(linkedFilters)-1].Type != LinkedFilterTypeCloseBracket {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryDidNotEndWithCloseBracket)
	}

	for i, v := range linkedFilters {
		if v.Type != nil && *v.Type == LinkedFilterTypeOpenBracket {
			openBracketCount++
		} else if v.Type != nil && *v.Type == LinkedFilterTypeCloseBracket {
			closeBracketCount++
		}

		if i > 0 && i < len(linkedFilters)-1 { // Not the first and not the last
			if v.Type != nil && v.Filter != nil {
				return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryFilterCannotAlsoBeAnOperator)
			}

			if v.Filter != nil && (linkedFilters[i-1].Type == nil || (*linkedFilters[i-1].Type != LinkedFilterTypeOpenBracket && *linkedFilters[i-1].Type != LinkedFilterTypeCloseBracket && *linkedFilters[i-1].Type != LinkedFilterTypeAnd && *linkedFilters[i-1].Type != LinkedFilterTypeOr)) {
				return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryFilterDoesNotFollowBracketOrOperator)
			}

			if v.Filter != nil && (linkedFilters[i+1].Type == nil || (*linkedFilters[i+1].Type != LinkedFilterTypeOpenBracket && *linkedFilters[i+1].Type != LinkedFilterTypeCloseBracket && *linkedFilters[i+1].Type != LinkedFilterTypeAnd && *linkedFilters[i+1].Type != LinkedFilterTypeOr)) {
				return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryFilterDoesNotPreceedBracketOrOperator)
			}

			if v.Type != nil && *v.Type == LinkedFilterTypeAnd {
				if linkedFilters[i-1].Filter == nil && (linkedFilters[i-1].Type != nil && *linkedFilters[i-1].Type != LinkedFilterTypeCloseBracket) {
					return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorAndDidNotFollowFilterOrCloseBracket)
				}
				if linkedFilters[i+1].Filter == nil && (linkedFilters[i+1].Type != nil && *linkedFilters[i+1].Type != LinkedFilterTypeOpenBracket) {
					return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorAndDidNotPreceedFilterOrOpenBracket)
				}
			}

			if v.Type != nil && *v.Type == LinkedFilterTypeOr {
				if linkedFilters[i-1].Filter == nil && (linkedFilters[i-1].Type != nil && *linkedFilters[i-1].Type != LinkedFilterTypeCloseBracket) {
					return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorOrDidNotFollowFilterOrCloseBracket)
				}
				if linkedFilters[i+1].Filter == nil && (linkedFilters[i+1].Type != nil && *linkedFilters[i+1].Type != LinkedFilterTypeOpenBracket) {
					return lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorOrDidNotPreceedFilterOrOpenBracket)
				}
			}

		}
	}

	if openBracketCount != closeBracketCount {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestCountOfOpenBracketsNotEqualToCountCloseBrackets)
	}

	return nil
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func convertValue(filter *Filter) (interface{}, error) {
	if filter.IsNull {
		if filter.Value != nil {
			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				filter.Key: badRequestOptionIsNullExistsCannotBeUsedWithValue,
			})
		}
		return nil, nil
	}

	if filter.Value == nil {
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
			filter.Key: badRequestMissingValue,
		})
	}

	switch v := filter.Value.(type) {
	default:
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
			filter.Key: badRequestUnrecognizedValue,
			"Value":    v,
		})

	case bool:
		return v, nil

	case float64:
		return v, nil

	case int:
		return v, nil

		// i, errInt64 := v.Int64()
		// if errInt64 == nil {
		// 	return i, nil
		// }
		// f, errFloat64 := v.Float64()
		// if errFloat64 == nil {
		// 	return f, nil
		// }

		// errMsg := badRequestUnrecognizedNumberValue
		// if errInt64 != nil && errFloat64 == nil {
		// 	errMsg = badRequestUnrecognizedInt64Value
		// }
		// if errInt64 == nil && errFloat64 != nil {
		// 	errMsg = badRequestUnrecognizedFloat64Value
		// }
		// return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
		// 	filter.Key: errMsg,
		// 	"Value":    v,
		// })

	case string:
		if filter.ValueType != nil && *filter.ValueType == queryPartKeyValueTypeDatetime {
			t, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
					filter.Key:  badRequestUnrecognizedDatetimeValue,
					"Value":     filter.Value,
					"ValueType": filter.ValueType,
				})
			}
			return t, nil
		}
		return v, nil

	case []interface{}:
		lib_log.Info(context.Background(), "777777")

		if len(v) == 0 {
			lib_log.Info(context.Background(), "66666")

			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				filter.Key: badRequestStringArrayValueCannotBeEmpty,
				"Value":    filter.Value,
			})
		}

		switch v[0].(type) {
		default:
			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				filter.Key: badRequestUnrecognizedValue,
				"Value":    filter.Value,
			})

		case bool:
			return valuesFromBools(filter.Key, v)

		case int:

			i, errInt64 := valuesFromJsonNumbersAsInt64(filter.Key, v)
			if errInt64 != nil {
				return nil, errInt64
			}
			return i, nil

		case float64:
			i, err := valuesFromJsonNumbersAsFloat64(filter.Key, v)
			if err != nil {
				return nil, err
			}
			return i, nil
			// f, errFloat64 := valuesFromJsonNumbersAsFloat64(filter.Key, v)
			// if errFloat64 == nil {
			// 	return f, nil
			// }

			// errMsg := badRequestUnrecognizedNumberValueArray
			// if errInt64 != nil && errFloat64 == nil {
			// 	errMsg = badRequestUnrecognizedInt64ValueArray
			// }
			// if errInt64 == nil && errFloat64 != nil {
			// 	errMsg = badRequestUnrecognizedFloat64ValueArray
			// }

			// return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
			// 	filter.Key: errMsg,
			// 	"Value":    v,
			// })

		case string:
			lib_log.Info(context.Background(), "555555")

			if filter.ValueType != nil && *filter.ValueType == queryPartKeyValueTypeDatetimeArray {
				return valuesFromStringsAsTime(filter.Key, v)
			}
			return valuesFromStrings(filter.Key, v)
		}
	}
	//revive:enable:cyclomatic
}

func valuesFromStrings(key string, values []interface{}) ([]string, error) {
	var ss []string
	for _, v := range values {
		vv, ok := v.(string)
		if !ok {
			return nil, lib_errors.Errorf("Failed processing key %q value of type string from array: %v", key, v)
		}
		ss = append(ss, vv)
	}
	return ss, nil
}

func valuesFromStringsAsTime(key string, values []interface{}) ([]time.Time, error) {
	var ts []time.Time
	for _, v := range values {
		vv, ok := v.(string)
		if !ok {
			return nil, lib_errors.Errorf("Failed processing key %q value of type string from array: %v", key, v)
		}
		t, err := time.Parse(time.RFC3339Nano, vv)
		if err != nil {
			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				key:     badRequestUnrecognizedDatetimeValueArray,
				"Value": values,
			})
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func valuesFromBools(key string, values []interface{}) ([]bool, error) {
	var bb []bool
	for _, v := range values {
		vv, ok := v.(bool)
		if !ok {
			return nil, lib_errors.Errorf("Failed processing key %q value of type bool from array: %v", key, v)
		}
		bb = append(bb, vv)
	}
	return bb, nil
}

func valuesFromJsonNumbersAsInt64(key string, values []interface{}) ([]int, error) {
	lib_log.Info(context.Background(), "3333", lib_log.FmtAny("key", key), lib_log.FmtAny("values", values))

	_, ok := values[0].(int)
	if !ok {
		return nil, lib_errors.Errorf("Failed processing key %q value of type json.Number from array: %v", key, values[0])
	}

	var ii []int
	for _, v := range values {
		vv, ok := v.(int)
		if !ok {
			return nil, lib_errors.Errorf("Failed processing key %q value of type int64 from array: %v", key, v)
		}
		ii = append(ii, vv)
	}

	return ii, nil
}

func valuesFromJsonNumbersAsFloat64(key string, values []interface{}) ([]float64, error) {
	lib_log.Info(context.Background(), "44444", lib_log.FmtAny("key", key), lib_log.FmtAny("values", values))

	_, ok := values[0].(float64)
	if !ok {
		return nil, lib_errors.Errorf("Failed processing key %q value of type json.Number from array: %v", key, values[0])
	}

	var ii []float64
	for _, v := range values {
		vv, ok := v.(float64)
		if !ok {
			return nil, lib_errors.Errorf("Failed processing key %q value of type int64 from array: %v", key, v)
		}
		ii = append(ii, vv)
	}
	lib_log.Info(context.Background(), "888888", lib_log.FmtAny("ii", ii))

	return ii, nil
}
