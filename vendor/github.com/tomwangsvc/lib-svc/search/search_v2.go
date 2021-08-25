package search

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_validation "github.com/tomwangsvc/lib-svc/validation"
)

func ParseQueryWithTestV2(query string, test bool) (filters []Filter, linkedFilters []LinkedFilter, err error) {
	filters, linkedFilters, err = ParseQueryV2(query)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed parsing query")
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

//revive:disable:cyclomatic
func ParseQueryV2(query string) (filters []Filter, linkedFilters []LinkedFilter, err error) {
	if query == "" {
		return
	}

	if i := strings.IndexAny(query, literalLinkedFilterTypeOpenBracket); i != 0 {
		err = lib_errors.NewCustom(http.StatusBadRequest, badRequestMalformedQueryShouldBeginWithOpenBracket)
		return
	}

	var countOfOpenBrackets, countOfCloseBrackets int
	var startIndexOfNextQueryPart int
	var previousFirstCharacterOfQuery string
	for query != "" {
		switch string(query[0]) {
		default:
			if previousFirstCharacterOfQuery != literalLinkedFilterTypeOpenBracket {
				err = lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryPartDidNotStartWithOpenBracket)
				return
			}

			startIndexOfNextQueryPart = strings.IndexAny(query, literalLinkedFilterTypeCloseBracket)
			if startIndexOfNextQueryPart < 0 {
				err = lib_errors.NewCustom(http.StatusBadRequest, badRequestFoundOpenBracketExpectedClose)
				return
			}

			var filter *Filter
			filter, err = generateFilterFromQueryPartV2(query[:startIndexOfNextQueryPart])
			if err != nil {
				err = lib_errors.Wrap(err, "Failed generating filter from query part")
				return
			}

			filters = append(filters, *filter)
			linkedFilters = append(linkedFilters, LinkedFilter{
				Filter: filter,
			})

		case literalLinkedFilterTypeOpenBracket:
			if len(linkedFilters) > 0 {
				switch previousFirstCharacterOfQuery {
				default:
					err = lib_errors.NewCustom(http.StatusBadRequest, badRequestOpenBracketPrefixedWithUnexpectedOperator)
					return
				case literalLinkedFilterTypeOpenBracket, literalLinkedFilterTypeAnd, literalLinkedFilterTypeOr, literalLinkedFilterTypeOrDeprecated:
				}
			}

			linkedFilterTypeOpenBracket := LinkedFilterTypeOpenBracket
			linkedFilters = append(linkedFilters, LinkedFilter{
				Type: &linkedFilterTypeOpenBracket,
			})
			startIndexOfNextQueryPart = 1
			countOfOpenBrackets++

		case literalLinkedFilterTypeCloseBracket:
			switch previousFirstCharacterOfQuery {
			default:
			case literalLinkedFilterTypeOpenBracket, literalLinkedFilterTypeAnd, literalLinkedFilterTypeOr, literalLinkedFilterTypeOrDeprecated:
				err = lib_errors.NewCustom(http.StatusBadRequest, badRequestCloseBracketPrefixedWithOperator)
				return
			}

			linkedFilterTypeCloseBracket := LinkedFilterTypeCloseBracket
			linkedFilters = append(linkedFilters, LinkedFilter{
				Type: &linkedFilterTypeCloseBracket,
			})
			startIndexOfNextQueryPart = 1
			countOfCloseBrackets++

		case literalLinkedFilterTypeAnd:
			if previousFirstCharacterOfQuery != literalLinkedFilterTypeCloseBracket {
				err = lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorNotPrefixedByCloseBracket)
				return
			}

			linkedFilterTypeAnd := LinkedFilterTypeAnd
			linkedFilters = append(linkedFilters, LinkedFilter{
				Type: &linkedFilterTypeAnd,
			})
			startIndexOfNextQueryPart = 1

		case literalLinkedFilterTypeOr, literalLinkedFilterTypeOrDeprecated:
			if previousFirstCharacterOfQuery != literalLinkedFilterTypeCloseBracket {
				err = lib_errors.NewCustom(http.StatusBadRequest, badRequestQueryOperatorNotPrefixedByCloseBracket)
				return
			}

			linkedFilterTypeOr := LinkedFilterTypeOr
			linkedFilters = append(linkedFilters, LinkedFilter{
				Type: &linkedFilterTypeOr,
			})
			startIndexOfNextQueryPart = 1
		}

		previousFirstCharacterOfQuery = string(query[0])
		query = query[startIndexOfNextQueryPart:]
	}

	if countOfCloseBrackets != countOfOpenBrackets {
		err = lib_errors.NewCustom(http.StatusBadRequest, badRequestCountOfOpenBracketsNotEqualToCountCloseBrackets)
		return
	}

	return
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func generateFilterFromQueryPartV2(queryPart string) (*Filter, error) {
	if queryPart == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestFoundEmptyQueryPart)
	}
	queryPartSections := strings.Split(queryPart, ":")
	if len(queryPartSections) > 4 {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestUnexpectedAmountOfKeysWithinKeyAndValueAndOptionsAndValueTypeQueryPart)
	}

	var filter Filter
	var queryPartKeyKeyExists, queryPartKeyOptionsExists, queryPartKeyValueExists, queryPartKeyValueTypeExists bool
	var value, valueType *string
	for _, queryPartSection := range queryPartSections {
		i := strings.Index(queryPartSection, "=")
		if i < 0 {
			return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestMalformedQueryNoEqualFoundToDetermineKeyValuePair)
		}

		var err error
		k, v := queryPartSection[:i], queryPartSection[i+1:]

		switch k {
		default:
			return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestUnexpectedKeyInQueryPart)

		case queryPartKeyKey:
			if queryPartKeyKeyExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			filter.Key = v
			queryPartKeyKeyExists = true

		case queryPartKeyOptions:
			if queryPartKeyOptionsExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			var options []string
			for _, val := range strings.Split(v, ",") {
				options = append(options, val)
			}
			filter, err = setOptions(filter, options)
			if err != nil {
				return nil, lib_errors.Wrap(err, "Failed setting options")
			}
			queryPartKeyOptionsExists = true

		case queryPartKeyValue:
			if queryPartKeyValueExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			value = &v
			queryPartKeyValueExists = true

		case queryPartKeyValueType:
			if queryPartKeyValueTypeExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			valueType = &v
			queryPartKeyValueTypeExists = true
		}
	}

	var err error
	filter.Value, err = convertValueToValueType(value, filter.IsNull, filter.Key, valueType)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed converting escaped value to value type")
	}

	return &filter, nil
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func convertValueToValueType(value *string, filterIsNull bool, key string, valueType *string) (interface{}, error) {
	if filterIsNull {
		if value != nil || valueType != nil {
			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				key: badRequestOptionIsNullExistsCannotBeUsedWithValueOrValueType,
			})
		}
		return nil, nil
	}

	if value == nil || valueType == nil {
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
			key: badRequestMissingValueOrValueType,
		})
	}

	v, err := base64.StdEncoding.DecodeString(*value)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestQueryValueIsNotBase64Encoded, lib_errors.Wrapf(err, "Failed decoding query value %q with type %q for key %q", *value, *valueType, key))
	}
	decodedValue := string(v)

	switch *valueType {
	default:
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{ // TODO: combine this to be top level metadata
			key: badRequestUnrecognizedValueType,
		})

	case queryPartKeyValueTypeBool:
		val, err := strconv.ParseBool(decodedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type BOOL"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return val, nil

	case queryPartKeyValueTypeDatetime: // TODO: consider indicating format or just change/add specific formats later as this is what we do in cases by default
		val, err := time.Parse(time.RFC3339Nano, decodedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type DATETIME"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return val, nil

	case queryPartKeyValueTypeDatetimeArray: // TODO: consider indicating format or just change/add specific formats later as this is what we do in cases by default
		var vals []time.Time
		for _, v := range strings.Split(decodedValue, ",") {
			val, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type DATETIME_ARRAY"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
					key: badRequestCannotParseValueToValueType,
				})
			}
			vals = append(vals, val)
		}
		return vals, nil

	case queryPartKeyValueTypeFloat64:
		val, err := strconv.ParseFloat(decodedValue, 64)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type FLOAT64"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return val, nil

	case queryPartKeyValueTypeFloat64Array:
		var vals []float64
		for _, v := range strings.Split(decodedValue, ",") {
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type FLOAT64_ARRAY"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
					key: badRequestCannotParseValueToValueType,
				})
			}
			vals = append(vals, val)
		}
		return vals, nil

	case queryPartKeyValueTypeInt64:
		val, err := strconv.ParseInt(decodedValue, 10, 64)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type INT64"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return val, nil

	case queryPartKeyValueTypeInt64Array:
		var vals []int64
		for _, v := range strings.Split(decodedValue, ",") {
			val, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", lib_errors.Wrap(err, "Failed parsing key value type INT64_ARRAY"), lib_validation.Metadata{ // TODO: combine this to be top level metadata
					key: badRequestCannotParseValueToValueType,
				})
			}
			vals = append(vals, val)
		}
		return vals, nil

	case queryPartKeyValueTypeString:
		return decodedValue, nil

	case queryPartKeyValueTypeStringArray:
		var vals []string
		for _, v := range strings.Split(decodedValue, ",") {
			vals = append(vals, v)
		}
		return vals, nil
	}
	//revive:enable:cyclomatic
}
