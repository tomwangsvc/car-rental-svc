package search

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_validation "github.com/tomwangsvc/lib-svc/validation"
)

const (
	badRequestCannotParseQuery                                  = "CANNOT_PARSE_QUERY"
	badRequestCannotParseValueToValueType                       = "CANNOT_PARSE_VALUE_TO_VALUE_TYPE"
	badRequestCloseBracketPrefixedWithOperator                  = "CLOSE_BRACKET_PREFIXED_WITH_OPERATOR"
	badRequestCountOfOpenBracketsNotEqualToCountCloseBrackets   = "COUNT_OF_OPEN_BRACKETS_NOT_EQUAL_TO_COUNT_CLOSE_BRACKETS"
	badRequestDuplicateKeyInQueryPartExists                     = "DUPLICATE_KEY_IN_QUERY_PART_EXISTS"
	badRequestExclusiveOptionsUsedForSameFilter                 = "EXCLUSIVE_OPTIONS_USED_FOR_SAME_FILTER"
	badRequestFailedQueryUnescape                               = "FAILED_QUERY_UNESCAPE"
	badRequestFoundEmptyQueryPart                               = "FOUND_EMPTY_QUERY_PART"
	badRequestFoundOpenBracketExpectedClose                     = "FOUND_OPEN_BRACKET_EXPECTED_CLOSE"
	badRequestMalformedQueryNoEqualFoundToDetermineKeyValuePair = "MALFORMED_QUERY_NO_EQUAL_FOUND_TO_DETERMINE_KEY_VALUE_PAIR"
	badRequestMalformedQueryShouldBeginWithOpenBracket          = "MALFORMED_QUERY_SHOULD_BEGIN_WITH_OPEN_BRACKET"
	badRequestMissingValue                                      = "MISSING_VALUE"
	badRequestMissingValueOrValueType                           = "MISSING_VALUE_OR_VALUE_TYPE"
	badRequestNumberArrayValueCannotBeEmpty                     = "NUMBER_ARRAY_VALUE_CANNOT_BE_EMPTY"
	badRequestNumberArrayValueCannotMixIntegersWithDecimals     = "NUMBER_ARRAY_VALUE_CANNOT_MIX_INTEGERS_WITH_DECIMALS"
	badRequestNumberArrayValueCannotMixDecimalsWithIntegers     = "NUMBER_ARRAY_VALUE_CANNOT_MIX_DECIMALS_WITH_INTEGERS"

	badRequestOpenBracketPrefixedWithUnexpectedOperator                 = "OPEN_BRACKET_PREFIXED_WITH_UNEXPECTED_OPERATOR"
	badRequestOptionIsNullExistsCannotBeUsedWithValue                   = "OPTION_IS_NULL_EXISTS_CANNOT_BE_USED_WITH_VALUE"
	badRequestOptionIsNullExistsCannotBeUsedWithValueOrValueType        = "OPTION_IS_NULL_EXISTS_CANNOT_BE_USED_WITH_VALUE_OR_VALUE_TYPE"
	badRequestQueryDidNotContainAtLeastOneFilterWithTheRequiredBrackets = "QUERY_DID_NOT_CONTAIN_AT_LEAST_ONE_FILTER_WITH_THE_REQUIRED_BRACKETS"
	badRequestQueryDidNotEndWithCloseBracket                            = "QUERY_DID_NOT_END_WITH_CLOSE_BRACKET"
	badRequestQueryDidNotStartWithOpenBracket                           = "QUERY_DID_NOT_START_WITH_OPEN_BRACKET"
	badRequestQueryIsNotBase64Encoded                                   = "QUERY_IS_NOT_BASE64_ENCODED"
	badRequestQueryOperatorNotPrefixedByCloseBracket                    = "QUERY_OPERATOR_NOT_PREFIXED_BY_CLOSE_BRACKET"

	badRequestQueryFilterCannotAlsoBeAnOperator          = "QUERY_FILTER_CANNOT_ALSO_BE_AN_OPERATOR"
	badRequestQueryFilterDoesNotFollowBracketOrOperator  = "QUERY_FILTER_DOES_NOT_FOLLOW_BRACKET_OR_OPERATOR"
	badRequestQueryFilterDoesNotPreceedBracketOrOperator = "QUERY_FILTER_DOES_NOT_PRECEED_BRACKET_OR_OPERATOR"

	badRequestQueryOperatorAndDidNotFollowFilterOrCloseBracket = "QUERY_OPERATOR_AND_DID_NOT_FOLLOW_FILTER_OR_CLOSE_BRACKET"
	badRequestQueryOperatorAndDidNotPreceedFilterOrOpenBracket = "QUERY_OPERATOR_AND_DID_NOT_PRECEED_FILTER_OR_OPEN_BRACKET"
	badRequestQueryOperatorOrDidNotFollowFilterOrCloseBracket  = "QUERY_OPERATOR_OR_DID_NOT_FOLLOW_FILTER_OR_CLOSE_BRACKET"
	badRequestQueryOperatorOrDidNotPreceedFilterOrOpenBracket  = "QUERY_OPERATOR_OR_DID_NOT_PRECEED_FILTER_OR_OPEN_BRACKET"

	badRequestQueryPartDidNotStartWithOpenBracket                                    = "QUERY_PART_DID_NOT_START_WITH_OPEN_BRACKET"
	badRequestQueryValueIsNotBase64Encoded                                           = "QUERY_VALUE_IS_NOT_BASE64_ENCODED"
	badRequestStringArrayValueCannotBeEmpty                                          = "STRING_ARRAY_VALUE_CANNOT_BE_EMPTY"
	badRequestStringArrayValueCannotMixStringsAndDatetimes                           = "STRING_ARRAY_VALUE_CANNOT_MIX_DSTRINGS_AND_DATETIMES"
	badRequestUnexpectedAmountOfKeysWithinKeyAndValueAndOptionsAndValueTypeQueryPart = "UNEXPECTED_AMOUNT_OF_KEYS_WITHIN_KEY_AND_VALUE_AND_OPTIONS_AND_VALUE_TYPE_QUERY_PART"
	badRequestUnexpectedKeyInQueryPart                                               = "UNEXPECTED_KEY_IN_QUERY_PART"
	badRequestUnrecognizedDatetimeValue                                              = "UNRECOGNIZED_DATETIME_VALUE"
	badRequestUnrecognizedDatetimeValueArray                                         = "UNRECOGNIZED_DATETIME_VALUE_ARRAY"
	badRequestUnrecognizedFloat64Value                                               = "UNRECOGNIZED_FLOAT64_VALUE"
	badRequestUnrecognizedFloat64ValueArray                                          = "UNRECOGNIZED_FLOAT64_VALUE_ARRAY"
	badRequestUnrecognizedInt64Value                                                 = "UNRECOGNIZED_INT64_VALUE"
	badRequestUnrecognizedInt64ValueArray                                            = "UNRECOGNIZED_INT64_VALUE_ARRAY"
	badRequestUnrecognizedNumberValue                                                = "UNRECOGNIZED_NUMBER_VALUE"
	badRequestUnrecognizedNumberValueArray                                           = "UNRECOGNIZED_NUMBER_VALUE_ARRAY"
	badRequestUnrecognizedOption                                                     = "UNRECOGNIZED_OPTION"
	badRequestUnrecognizedValue                                                      = "UNRECOGNIZED_VALUE"
	badRequestUnrecognizedValueType                                                  = "UNRECOGNIZED_VALUE_TYPE"
	badRequestValueTypeIsNullButValueIsNot                                           = "VALUE_TYPE_IS_NULL_BUT_VALUE_IS_NOT"

	filterOptionArrayContains          = "ARRAY_CONTAINS"
	filterOptionCaseInsensitiveString  = "CASE_INSENSITIVE_STRING"
	filterOptionInArray                = "IN_ARRAY"
	filterOptionInRange                = "IN_RANGE"
	filterOptionIsGreaterThan          = "IS_GREATER_THAN"
	filterOptionIsGreaterThanOrEqualTo = "IS_GREATER_THAN_OR_EQUAL_TO"
	filterOptionIsLessThan             = "IS_LESS_THAN"
	filterOptionIsLessThanOrEqualTo    = "IS_LESS_THAN_OR_EQUAL_TO"
	filterOptionIsNull                 = "IS_NULL"
	filterOptionNotCondition           = "NOT_CONDITION"
	filterOptionPartialMatchString     = "PARTIAL_MATCH_STRING"

	LinkedFilterTypeAnd          = "AND"
	LinkedFilterTypeCloseBracket = "CLOSE_BRACKET"
	LinkedFilterTypeFilter       = "FILTER"
	LinkedFilterTypeOpenBracket  = "OPEN_BRACKET"
	LinkedFilterTypeOr           = "OR"

	literalLinkedFilterTypeAnd          = "+"
	literalLinkedFilterTypeCloseBracket = ")"
	literalLinkedFilterTypeOpenBracket  = "("
	literalLinkedFilterTypeOr           = "*"
	literalLinkedFilterTypeOrDeprecated = "|" // TODO: deprecate

	queryKeyQuery = "query"

	queryPartKeyKey       = "key"
	queryPartKeyOptions   = "options"
	queryPartKeyValue     = "value"
	queryPartKeyValueType = "value_type"

	queryPartKeyValueTypeBool          = "BOOL"
	queryPartKeyValueTypeDatetime      = "DATETIME"
	queryPartKeyValueTypeDatetimeArray = "DATETIME_ARRAY"
	queryPartKeyValueTypeFloat64       = "FLOAT64"
	queryPartKeyValueTypeFloat64Array  = "FLOAT64_ARRAY"
	queryPartKeyValueTypeInt64         = "INT64"
	queryPartKeyValueTypeInt64Array    = "INT64_ARRAY"
	queryPartKeyValueTypeString        = "STRING"
	queryPartKeyValueTypeStringArray   = "STRING_ARRAY"
)

type LinkedFilter struct {
	Filter *Filter `json:"filter,omitempty"`
	Type   *string `json:"type,omitempty"`
}

type Filter struct {
	ArrayContains          bool        `json:"array_contains,omitempty"`              // Matches where value in resource value (array)
	CaseInsensitiveString  bool        `json:"case_insensitive_string,omitempty"`     // Matches where resource value (string) matches value ignoring case
	InArray                bool        `json:"in_array,omitempty"`                    // Matches where resource value in value (array)
	InRange                bool        `json:"in_range,omitempty"`                    // Matches where resource value between two value(slice.len() == 2)
	IsGreaterThan          bool        `json:"is_greater_than,omitempty"`             // Matches where resource value is greater than value
	IsGreaterThanOrEqualTo bool        `json:"is_greater_than_or_equal_to,omitempty"` // Matches where resource value is greater than or equal to value
	IsLessThan             bool        `json:"is_less_than,omitempty"`                // Matches where resource value is less than value
	IsLessThanOrEqualTo    bool        `json:"is_less_than_or_equal_to,omitempty"`    // Matches where resource value is less than or equal to value
	IsNull                 bool        `json:"is_null,omitempty"`                     // Matches where resource value is null
	Key                    string      `json:"key"`
	NotCondition           bool        `json:"not_condition,omitempty"`        // Makes any condition into a "not" condition, e.g. "=" -> "!=", "LIKE" -> "NOT LIKE", "IS NULL" -> "IS NOT NULL"
	PartialMatchString     bool        `json:"partial_match_string,omitempty"` // Matches where string resource value (string) contains value
	Value                  interface{} `json:"value,omitempty"`
	ValueType              *string     `json:"value_type,omitempty"`
}

func QueryEncodedQueryFromRawQuery(rawQuery string) (string, error) {
	var query string
	for rawQuery != "" {
		key := rawQuery
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, rawQuery = key[:i], key[i+1:]
		} else {
			rawQuery = ""
		}
		if key == "" {
			continue
		}
		var value string
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err := url.QueryUnescape(key)
		if err != nil {
			return "", lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
		}
		if key == queryKeyQuery {
			query = value
			break
		}
	}
	return query, nil
}

func ParseQueryWithTest(query string, test bool) (filters []Filter, linkedFilters []LinkedFilter, err error) {
	filters, linkedFilters, err = ParseQuery(query)
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
func ParseQuery(query string) (filters []Filter, linkedFilters []LinkedFilter, err error) {
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
			filter, err = generateFilterFromQueryPart(query[:startIndexOfNextQueryPart])
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
func generateFilterFromQueryPart(queryPart string) (*Filter, error) {
	if queryPart == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestFoundEmptyQueryPart)
	}
	queryPartSections := strings.Split(queryPart, ":")
	if len(queryPartSections) > 4 {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestUnexpectedAmountOfKeysWithinKeyAndValueAndOptionsAndValueTypeQueryPart)
	}

	var filter Filter
	var queryPartKeyKeyExists, queryPartKeyOptionsExists, queryPartKeyValueExists, queryPartKeyValueTypeExists bool
	var escapedValue, valueType *string
	for _, queryPartSection := range queryPartSections {
		i := strings.Index(queryPartSection, "=")
		if i < 0 {
			return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestMalformedQueryNoEqualFoundToDetermineKeyValuePair)
		}

		var err error
		k, v := queryPartSection[:i], queryPartSection[i+1:]
		k, err = url.QueryUnescape(k)
		if err != nil {
			return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
		}

		switch k {
		default:
			return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestUnexpectedKeyInQueryPart)

		case queryPartKeyKey:
			if queryPartKeyKeyExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			v, err = url.QueryUnescape(v)
			if err != nil {
				return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
			}
			filter.Key = v
			queryPartKeyKeyExists = true

		case queryPartKeyOptions:
			if queryPartKeyOptionsExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			var options []string
			for _, val := range strings.Split(v, ",") {
				val, err = url.QueryUnescape(val)
				if err != nil {
					return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
				}
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
			escapedValue = &v
			queryPartKeyValueExists = true

		case queryPartKeyValueType:
			if queryPartKeyValueTypeExists {
				return nil, lib_errors.NewCustom(http.StatusBadRequest, badRequestDuplicateKeyInQueryPartExists)
			}
			v, err = url.QueryUnescape(v)
			if err != nil {
				return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
			}
			valueType = &v
			queryPartKeyValueTypeExists = true
		}
	}

	var err error
	filter.Value, err = convertEscapedValueToValueType(escapedValue, filter.IsNull, filter.Key, valueType)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed converting escaped value to value type")
	}

	return &filter, nil
	//revive:enable:cyclomatic
}

func setOptions(filter Filter, options []string) (updatedFilter Filter, err error) {
	updatedFilter = filter
	for _, option := range options {
		switch option {
		default:
			err = lib_errors.NewCustom(http.StatusBadRequest, badRequestUnrecognizedOption) // TODO: consider turning this into metadata and append, i.e. "query.1.option.2": "UNRECOGNIZED"
			return

		case filterOptionArrayContains:
			updatedFilter.ArrayContains = true

		case filterOptionCaseInsensitiveString:
			updatedFilter.CaseInsensitiveString = true

		case filterOptionInArray:
			updatedFilter.InArray = true

		case filterOptionInRange:
			updatedFilter.InRange = true

		case filterOptionIsGreaterThan:
			updatedFilter.IsGreaterThan = true

		case filterOptionIsGreaterThanOrEqualTo:
			updatedFilter.IsGreaterThanOrEqualTo = true

		case filterOptionIsLessThan:
			updatedFilter.IsLessThan = true

		case filterOptionIsLessThanOrEqualTo:
			updatedFilter.IsLessThanOrEqualTo = true

		case filterOptionIsNull:
			updatedFilter.IsNull = true

		case filterOptionNotCondition:
			updatedFilter.NotCondition = true

		case filterOptionPartialMatchString:
			updatedFilter.PartialMatchString = true
		}
	}

	if err = CheckFilterOptionsForExclusivity(updatedFilter); err != nil {
		err = lib_errors.Wrap(err, "Failed checking filter options for exclusivity")
		return
	}

	return
}

//revive:disable:cyclomatic
func CheckFilterOptionsForExclusivity(filter Filter) error {
	if filter.ArrayContains && (filter.CaseInsensitiveString || filter.InArray || filter.InRange || filter.IsGreaterThan || filter.IsGreaterThanOrEqualTo || filter.IsLessThan || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.InArray && (filter.CaseInsensitiveString || filter.InRange || filter.IsGreaterThan || filter.IsGreaterThanOrEqualTo || filter.IsLessThan || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.InRange && (filter.CaseInsensitiveString || filter.IsGreaterThan || filter.IsGreaterThanOrEqualTo || filter.IsLessThan || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.IsGreaterThan && (filter.CaseInsensitiveString || filter.IsGreaterThanOrEqualTo || filter.IsLessThan || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.IsGreaterThanOrEqualTo && (filter.CaseInsensitiveString || filter.IsLessThan || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.IsLessThan && (filter.CaseInsensitiveString || filter.IsLessThanOrEqualTo || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.IsLessThanOrEqualTo && (filter.CaseInsensitiveString || filter.IsNull || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination

	} else if filter.IsNull && (filter.CaseInsensitiveString || filter.PartialMatchString) {
		return lib_errors.NewCustom(http.StatusBadRequest, badRequestExclusiveOptionsUsedForSameFilter) // TODO: consider specific errors based on the combination
	}

	return nil
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func convertEscapedValueToValueType(escapedValue *string, filterIsNull bool, key string, valueType *string) (interface{}, error) {
	if filterIsNull {
		if escapedValue != nil || valueType != nil {
			return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
				key: badRequestOptionIsNullExistsCannotBeUsedWithValueOrValueType,
			})
		}
		return nil, nil
	}

	if escapedValue == nil || valueType == nil {
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{
			key: badRequestMissingValueOrValueType,
		})
	}

	switch *valueType {
	default:
		return nil, lib_errors.NewCustomWithMetadata(http.StatusBadRequest, "", lib_validation.Metadata{ // TODO: combine this to be top level metadata
			key: badRequestUnrecognizedValueType,
		})

	case queryPartKeyValueTypeBool:
		value, err := parseQueryPartKeyValueTypeBool(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return *value, nil

	case queryPartKeyValueTypeDatetime: // TODO: consider indicating format or just change/add specific formats later as this is what we do in cases by default
		value, err := parseQueryPartKeyValueTypeDatetime(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return *value, nil

	case queryPartKeyValueTypeDatetimeArray: // TODO: consider indicating format or just change/add specific formats later as this is what we do in cases by default
		value, err := parseQueryPartKeyValueTypeDatetimeArray(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return value, nil

	case queryPartKeyValueTypeFloat64:
		value, err := parseQueryPartKeyValueTypeFloat64(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return *value, nil

	case queryPartKeyValueTypeFloat64Array:
		value, err := parseQueryPartKeyValueTypeFloat64Array(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return value, nil

	case queryPartKeyValueTypeInt64:
		value, err := parseQueryPartKeyValueTypeInt64(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return *value, nil

	case queryPartKeyValueTypeInt64Array:
		value, err := parseQueryPartKeyValueTypeInt64Array(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return value, nil

	case queryPartKeyValueTypeString:
		value, err := parseQueryPartKeyValueTypeString(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return *value, nil

	case queryPartKeyValueTypeStringArray:
		value, err := parseQueryPartKeyValueTypeStringArray(*escapedValue)
		if err != nil {
			return nil, lib_errors.NewCustomWithCauseAndMetadata(http.StatusBadRequest, "", err, lib_validation.Metadata{ // TODO: combine this to be top level metadata
				key: badRequestCannotParseValueToValueType,
			})
		}
		return value, nil
	}
	//revive:enable:cyclomatic
}

func parseQueryPartKeyValueTypeBool(escapedValue string) (*bool, error) {
	value, err := url.QueryUnescape(escapedValue)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
	}
	val, err := strconv.ParseBool(value)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing bool")
	}
	return &val, nil
}

func parseQueryPartKeyValueTypeDatetime(escapedValue string) (*time.Time, error) {
	value, err := url.QueryUnescape(escapedValue)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
	}
	val, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing datetime")
	}
	return &val, nil
}

func parseQueryPartKeyValueTypeDatetimeArray(escapedValue string) ([]time.Time, error) {
	var vals []time.Time
	for _, value := range strings.Split(escapedValue, ",") {
		val, err := parseQueryPartKeyValueTypeDatetime(value)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed parsing query part key value type datetime")
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func parseQueryPartKeyValueTypeFloat64(escapedValue string) (*float64, error) {
	value, err := url.QueryUnescape(escapedValue)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
	}
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing float")
	}
	return &val, nil
}

func parseQueryPartKeyValueTypeFloat64Array(escapedValue string) ([]float64, error) {
	var vals []float64
	for _, value := range strings.Split(escapedValue, ",") {
		val, err := parseQueryPartKeyValueTypeFloat64(value)
		if err != nil {
			return nil, lib_errors.Wrap(err, "parsing query part key value type float64")
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func parseQueryPartKeyValueTypeInt64(escapedValue string) (*int64, error) {
	value, err := url.QueryUnescape(escapedValue)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
	}
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing int")
	}
	return &val, nil
}

func parseQueryPartKeyValueTypeInt64Array(escapedValue string) ([]int64, error) {
	var vals []int64
	for _, value := range strings.Split(escapedValue, ",") {
		val, err := parseQueryPartKeyValueTypeInt64(value)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed parsing query part key value type int64")
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func parseQueryPartKeyValueTypeString(escapedValue string) (*string, error) {
	value, err := url.QueryUnescape(escapedValue)
	if err != nil {
		return nil, lib_errors.NewCustomWithCause(http.StatusBadRequest, badRequestFailedQueryUnescape, err)
	}
	return &value, nil
}

func parseQueryPartKeyValueTypeStringArray(escapedValue string) ([]string, error) {
	var vals []string
	for _, value := range strings.Split(escapedValue, ",") {
		val, err := parseQueryPartKeyValueTypeString(value)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed parsing query part key value type string")
		}
		vals = append(vals, *val)
	}
	return vals, nil
}
