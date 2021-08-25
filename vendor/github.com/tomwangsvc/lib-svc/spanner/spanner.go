package spanner

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_search "github.com/tomwangsvc/lib-svc/search"
	lib_validation "github.com/tomwangsvc/lib-svc/validation"
	"google.golang.org/grpc/codes"
)

const (
	errorDescriptionDatabaseNotFoundPrefix      = "Database not found"
	errorDescriptionSessionNotFoundPrefix       = "Session not found"
	errorDescriptionUniqueIndexViolationPostfix = " at index key"
	errorDescriptionUniqueIndexViolationPrefix  = "Unique index violation on index "

	BadGatewayEncounteredSpannerAbortedError = "ENCOUNTERED_SPANNER_ABORTED_ERROR"

	ConflictObjectAlreadyExists            = "OBJECT_ALREADY_EXISTS"
	ConflictObjectCannotBeCreatedOrUpdated = "OBJECT_CANNOT_BE_CREATED_OR_UPDATED"
	ConflictUniqueIndexViolation           = "UNIQUE_INDEX_VIOLATION"

	paramsKeyBase = "__PARAM__"
)

func WrapError(err error, defaultMessage string) error {
	return wrapError(err, defaultMessage)
}

func wrapError(err error, defaultMessage string) error {
	errorToCheckCodeOf := err
	if customError, ok := err.(lib_errors.Custom); ok && customError.Cause != nil {
		errorToCheckCodeOf = customError.Cause
	}

	switch spanner.ErrCode(errorToCheckCodeOf) {
	default:
		return lib_errors.Wrap(err, defaultMessage)

	case codes.Aborted:
		return lib_errors.NewCustomWithCause(http.StatusBadGateway, BadGatewayEncounteredSpannerAbortedError, err)

	case codes.AlreadyExists:
		errorDescription := spanner.ErrDesc(err)
		if start := strings.Index(errorDescription, errorDescriptionUniqueIndexViolationPrefix); start > -1 {
			errorDescription = errorDescription[start+len(errorDescriptionUniqueIndexViolationPrefix):]
			end := strings.Index(errorDescription, errorDescriptionUniqueIndexViolationPostfix)
			var indexViolated string
			if end == -1 {
				indexViolated = "UNKNOWN"
			} else {
				indexViolated = errorDescription[:end]
			}
			return lib_errors.NewCustomWithCauseAndMetadata(http.StatusConflict, ConflictObjectCannotBeCreatedOrUpdated, err, lib_validation.Metadata{
				indexViolated: ConflictUniqueIndexViolation,
			})
		}
		return lib_errors.NewCustomWithCause(http.StatusConflict, ConflictObjectAlreadyExists, err)

	case codes.DeadlineExceeded:
		return lib_errors.NewCustomWithCause(http.StatusGatewayTimeout, "Encountered Spanner deadline exceeded error", err)

	case codes.Internal:
		return lib_errors.NewCustomWithCause(http.StatusBadGateway, "Encountered Spanner internal error", err)

	case codes.NotFound:
		if strings.Contains(spanner.ErrDesc(err), errorDescriptionDatabaseNotFoundPrefix) {
			return lib_errors.NewCustomWithCause(http.StatusBadGateway, "Encountered Spanner database not found error", err)
		}
		if strings.Contains(spanner.ErrDesc(err), errorDescriptionSessionNotFoundPrefix) {
			return lib_errors.NewCustomWithCause(http.StatusBadGateway, "Encountered Spanner session not found error", err)
		}
		// TODO: Until NotFound error messages are swallowed (see errors.Render) this error message could be exposed to clients and hence it must not discuss internal infrastructure
		//return lib_errors.NewCustomWithCause(http.StatusNotFound, "Encountered Spanner not found error", err)
		return lib_errors.NewCustomWithCause(http.StatusNotFound, "Entity not found", err)
	}
}

func IsAborted(err error) bool {
	return lib_errors.IsCustomWithCode(err, http.StatusBadGateway) && strings.Contains(err.Error(), BadGatewayEncounteredSpannerAbortedError)
}

func ClientConfigWithMinOpenedAndRestrictLocalhost(env lib_env.Env, minOpened int, restrictLocalhost bool) spanner.ClientConfig {
	if restrictLocalhost && env.Localhost {
		minOpened = 0
	}
	// If the protocol is not met (https://cloud.google.com/spanner/docs/reference/rpc/google.spanner.v1#session) the following error is generated
	// -> spanner: code = "InvalidArgument", desc = "Invalid CreateSession request."
	return spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: uint64(minOpened),
		},
		SessionLabels: map[string]string{
			"id":       env.RuntimeId,
			"location": env.RuntimeLabel,
		},
	}
}

func ClientConfigWithMinOpened(env lib_env.Env, minOpened int) spanner.ClientConfig {
	return ClientConfigWithMinOpenedAndRestrictLocalhost(env, minOpened, true)
}

func ClientConfig(env lib_env.Env) spanner.ClientConfig {
	return ClientConfigWithMinOpenedAndRestrictLocalhost(env, 3, true)
}

// Reader interface for row readers
type Reader interface {
	ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
	QueryWithOptions(ctx context.Context, statement spanner.Statement, options spanner.QueryOptions) *spanner.RowIterator
}

// ReadById reads an entity using its ID
func ReadById(ctx context.Context, r Reader, table string, columns []string, id string, dst interface{}) error {
	lib_log.Info(ctx, "Reading", lib_log.FmtString("id", id), lib_log.FmtString("table", table), lib_log.FmtStrings("columns", columns))

	row, err := r.ReadRow(ctx, table, spanner.Key{id}, columns)
	if err != nil {
		return wrapError(err, "Failed reading row")
	}

	if err := row.ToStruct(dst); err != nil {
		return lib_errors.Wrapf(err, "Failed unpacking row from table %q using ID %q", table, id)
	}

	lib_log.Info(ctx, "Read", lib_log.FmtAny("dst", dst))
	return nil
}

// ReadById reads an entity using its IDs
func ReadByIds(ctx context.Context, r Reader, table string, columns, ids []string, dst interface{}) error {
	lib_log.Info(ctx, "Reading", lib_log.FmtAny("ids", ids), lib_log.FmtString("table", table), lib_log.FmtStrings("columns", columns))

	var key spanner.Key
	for _, id := range ids {
		key = append(key, id)
	}
	row, err := r.ReadRow(ctx, table, key, columns)
	if err != nil {
		return wrapError(err, "Failed reading row")
	}

	if err := row.ToStruct(dst); err != nil {
		return lib_errors.Wrapf(err, "Failed unpacking row from table %q using IDs %q", table, ids)
	}

	lib_log.Info(ctx, "Read", lib_log.FmtAny("dst", dst))
	return nil
}

// HasDateUpdatedChanged checks whether an update is being attemped based on out-of-date data
// If the client does not send a date then data should be persisted i.e. this function should return false
func HasDateUpdatedChanged(clientViewOfDateUpdated *time.Time, databaseViewOfDateUpdated spanner.NullTime) bool {
	if !databaseViewOfDateUpdated.Valid {
		return false
	}
	if clientViewOfDateUpdated != nil && !(*clientViewOfDateUpdated).Equal(databaseViewOfDateUpdated.Time) {
		return true
	}
	return false
}

type SearchFiltersByColumn map[string]SearchFilter

type SearchFilter struct {
	ArrayContains          bool        // Matches where query value in spanner value (array), cannot be used in combination with InSlice, IsNull, CaseInsensitiveString, or PartialMatchString
	CaseInsensitiveString  bool        // Matches where spanner value (string) matches query value ignoring case
	InRange                bool        // Matches where spanner value between two query value(slice.len() == 2), cannot be used in combination with ArrayContains, IsNull, CaseInsensitiveString, PartialMatchString, IsGreaterThanOrEqualTo or IsLessThanOrEqualTo
	InSlice                bool        // Matches where spanner value in query value (slice), cannot be used in combination with ArrayContains, IsNull, CaseInsensitiveString, or PartialMatchString
	IsGreaterThanOrEqualTo bool        // Matches where spanner value is greater than query value
	IsLessThanOrEqualTo    bool        // Matches where spanner value is less than query value
	IsNull                 bool        // Matches where spanner value is null, cannot be used in combination with CaseInsensitiveString, ArrayContains, InSlice, or PartialMatchString
	NotCondition           bool        // Makes any condition into a "not" condition, e.g. "=" -> "!=", "LIKE" -> "NOT LIKE", "IS NULL" -> "IS NOT NULL"
	PartialMatchString     bool        // Matches where string spanner value (string) contains query value
	Value                  interface{} // HINT - if using filter InSlice, pass in the original query value of the slice.  JSON marshalling then unmarhsalling into map[string]interface{} causes slices to become []interface{}, which is not supported by spanner
}

func GenerateSQLWhereAndParamsForSearch(filters SearchFiltersByColumn) (string, map[string]interface{}, error) {
	return GenerateSQLWhereAndParamsForSearchWithInitialWhere(filters, "", nil)
}

//revive:disable:cyclomatic
func GenerateSQLWhereAndParamsForSearchWithInitialWhere(filters SearchFiltersByColumn, initialWhere string, initialWhereParams map[string]interface{}) (string, map[string]interface{}, error) {
	if len(filters) == 0 {
		return "", nil, nil
	}
	var sqlWhere string
	params := make(map[string]interface{})
	for k, v := range initialWhereParams {
		params[k] = v
	}
	for key, filter := range filters {
		if filter.IsNull && (filter.ArrayContains || filter.CaseInsensitiveString || filter.InRange || filter.InSlice || filter.IsGreaterThanOrEqualTo || filter.IsLessThanOrEqualTo || filter.PartialMatchString) ||
			filter.ArrayContains && (filter.CaseInsensitiveString || filter.InRange || filter.InSlice || filter.IsGreaterThanOrEqualTo || filter.IsLessThanOrEqualTo || filter.PartialMatchString) ||
			filter.InRange && (filter.CaseInsensitiveString || filter.InSlice || filter.IsGreaterThanOrEqualTo || filter.IsLessThanOrEqualTo || filter.PartialMatchString) ||
			filter.InSlice && (filter.CaseInsensitiveString || filter.IsGreaterThanOrEqualTo || filter.IsLessThanOrEqualTo || filter.PartialMatchString) ||
			filter.IsGreaterThanOrEqualTo && (filter.CaseInsensitiveString || filter.IsLessThanOrEqualTo || filter.PartialMatchString) ||
			filter.IsLessThanOrEqualTo && (filter.CaseInsensitiveString || filter.PartialMatchString) {
			return "", nil, lib_errors.New("Used exclusive filters with other filters")
		}

		if sqlWhere == "" {
			sqlWhere = "WHERE"
			if initialWhere != "" {
				sqlWhere = fmt.Sprintf("%s %s AND", sqlWhere, initialWhere)
			}
		} else {
			sqlWhere = fmt.Sprintf("%s AND", sqlWhere)
		}

		column := key
		var condition string
		value := filter.Value
		if filter.ArrayContains {
			if filter.NotCondition {
				condition = "NOT IN"
			} else {
				condition = "IN"
			}
			sqlWhere = fmt.Sprintf("%s @%s %s UNNEST(%s)", sqlWhere, key, condition, column)
			params[key] = value

		} else if filter.InRange {
			valueOf := reflect.ValueOf(filter.Value)
			if valueOf.Kind() == reflect.Slice && valueOf.Len() == 2 {
				if filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s >= @%s) AND NOT(%s <= @%s)", sqlWhere, column, fmt.Sprintf("%s_start", key), column, fmt.Sprintf("%s_end", key))
				} else {
					sqlWhere = fmt.Sprintf("%s %s >= @%s AND %s <= @%s", sqlWhere, column, fmt.Sprintf("%s_start", key), column, fmt.Sprintf("%s_end", key))
				}
				params[fmt.Sprintf("%s_start", key)] = valueOf.Index(0).Interface()
				params[fmt.Sprintf("%s_end", key)] = valueOf.Index(1).Interface()
			} else {
				return "", nil, lib_errors.Errorf("filters[%q].Value of type %T, expected slice of length 2 for range search", key, filter.Value)
			}

		} else if filter.InSlice {
			if filter.NotCondition {
				condition = "NOT IN"
			} else {
				condition = "IN"
			}
			sqlWhere = fmt.Sprintf("%s %s %s UNNEST(@%s)", sqlWhere, column, condition, key)
			params[key] = value

		} else if filter.IsGreaterThanOrEqualTo {
			if filter.NotCondition {
				sqlWhere = fmt.Sprintf("%s NOT(%s >= @%s)", sqlWhere, column, key)
			} else {
				sqlWhere = fmt.Sprintf("%s %s >= @%s", sqlWhere, column, key)
			}
			params[key] = value

		} else if filter.IsLessThanOrEqualTo {
			if filter.NotCondition {
				sqlWhere = fmt.Sprintf("%s NOT(%s <= @%s)", sqlWhere, column, key)
			} else {
				sqlWhere = fmt.Sprintf("%s %s <= @%s", sqlWhere, column, key)
			}
			params[key] = value

		} else if filter.IsNull {
			if filter.NotCondition {
				condition = "IS NOT NULL"
			} else {
				condition = "IS NULL"
			}
			sqlWhere = fmt.Sprintf("%s %s %s", sqlWhere, column, condition)

		} else {
			if filter.NotCondition {
				condition = "!="
			} else {
				condition = "="
			}
			if filter.CaseInsensitiveString {
				column = fmt.Sprintf("UPPER(%s)", column)
				strVal, ok := value.(string)
				if !ok {
					return "", nil, lib_errors.Errorf("filters[%q].Value of type %T, expected string for case insensitive string search", key, filter.Value)
				}
				value = strings.ToUpper(strVal)
			}
			if filter.PartialMatchString {
				if filter.NotCondition {
					condition = "NOT LIKE"
				} else {
					condition = "LIKE"
				}
				strVal, ok := value.(string)
				if !ok {
					return "", nil, lib_errors.Errorf("filters[%q].Value of type %T, expected string for partial match string search", key, filter.Value)
				}
				value = fmt.Sprintf("%%%s%%", strVal)
			}
			sqlWhere = fmt.Sprintf("%s %s %s @%s", sqlWhere, column, condition, key)
			params[key] = value
		}
	}
	return sqlWhere, params, nil
}

func GenerateSqlWhereAndParamsForSearchV2(linkedFilters []lib_search.LinkedFilter) (string, map[string]interface{}, error) {
	return GenerateSqlWhereAndParamsForSearchWithInitialWhereV2("", nil, linkedFilters)
}

//revive:disable:cyclomatic
func GenerateSqlWhereAndParamsForSearchWithInitialWhereV2(initialWhere string, initialWhereParams map[string]interface{}, linkedFilters []lib_search.LinkedFilter) (string, map[string]interface{}, error) {
	if initialWhere == "" && len(linkedFilters) == 0 {
		return "", nil, nil

	}
	sqlWhere := "WHERE"
	if initialWhere != "" {
		sqlWhere = fmt.Sprintf("%s %s", sqlWhere, initialWhere)
		if len(linkedFilters) > 0 {
			sqlWhere = fmt.Sprintf("%s AND", sqlWhere)
		}
	}
	params := make(map[string]interface{})
	for k, v := range initialWhereParams {
		params[k] = v
	}

	for i, linkedFilter := range linkedFilters {
		if linkedFilter.Type != nil && *linkedFilter.Type != lib_search.LinkedFilterTypeFilter {
			switch *linkedFilter.Type {
			default:
				return "", nil, lib_errors.Errorf("Linked filter type %q not recognized", *linkedFilter.Type)

			case lib_search.LinkedFilterTypeAnd:
				sqlWhere = fmt.Sprintf("%s AND", sqlWhere)

			case lib_search.LinkedFilterTypeCloseBracket:
				sqlWhere = fmt.Sprintf("%s )", sqlWhere)

			case lib_search.LinkedFilterTypeOpenBracket:
				sqlWhere = fmt.Sprintf("%s (", sqlWhere)

			case lib_search.LinkedFilterTypeOr:
				sqlWhere = fmt.Sprintf("%s OR", sqlWhere)
			}

		} else if linkedFilter.Filter != nil {

			if err := lib_search.CheckFilterOptionsForExclusivity(*linkedFilter.Filter); err != nil {
				return "", nil, lib_errors.Wrap(err, "Failed checking filter options for exclusivity")
			}

			column := linkedFilter.Filter.Key
			paramKey := fmt.Sprintf("%s%d_%s", paramsKeyBase, i, column)
			value := linkedFilter.Filter.Value

			if linkedFilter.Filter.ArrayContains {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s @%s NOT IN UNNEST(%s)", sqlWhere, paramKey, column)
				} else {
					sqlWhere = fmt.Sprintf("%s @%s IN UNNEST(%s)", sqlWhere, paramKey, column)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.InArray {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s %s NOT IN UNNEST(@%s)", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s IN UNNEST(@%s)", sqlWhere, column, paramKey)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.InRange {
				valueOf := reflect.ValueOf(linkedFilter.Filter.Value)
				if valueOf.Kind() != reflect.Slice || valueOf.Len() != 2 {
					return "", nil, lib_errors.Errorf("Linked filter %d with filter with key %q of type %T, expected slice of length 2 for range search", i, paramKey, linkedFilter.Filter.Value)
				}
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s >= @%s) AND NOT(%s <= @%s)", sqlWhere, column, fmt.Sprintf("%s_start", paramKey), column, fmt.Sprintf("%s_end", paramKey))
				} else {
					sqlWhere = fmt.Sprintf("%s %s >= @%s AND %s <= @%s", sqlWhere, column, fmt.Sprintf("%s_start", paramKey), column, fmt.Sprintf("%s_end", paramKey))
				}
				params[fmt.Sprintf("%s_start", paramKey)] = valueOf.Index(0).Interface()
				params[fmt.Sprintf("%s_end", paramKey)] = valueOf.Index(1).Interface()

			} else if linkedFilter.Filter.IsGreaterThan {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s > @%s)", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s > @%s", sqlWhere, column, paramKey)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.IsGreaterThanOrEqualTo {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s >= @%s)", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s >= @%s", sqlWhere, column, paramKey)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.IsLessThan {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s < @%s)", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s < @%s", sqlWhere, column, paramKey)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.IsLessThanOrEqualTo {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s NOT(%s <= @%s)", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s <= @%s", sqlWhere, column, paramKey)
				}
				params[paramKey] = value

			} else if linkedFilter.Filter.IsNull {
				if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s %s IS NOT NULL", sqlWhere, column)
				} else {
					sqlWhere = fmt.Sprintf("%s %s IS NULL", sqlWhere, column)
				}

			} else {
				if linkedFilter.Filter.CaseInsensitiveString {
					column = fmt.Sprintf("UPPER(%s)", column)
					strVal, ok := value.(string)
					if !ok {
						return "", nil, lib_errors.Errorf("Linked filter %d with filter with key %q of type %T, expected string for case insensitive string search", i, paramKey, linkedFilter.Filter.Value)
					}
					value = strings.ToUpper(strVal)
				}
				if linkedFilter.Filter.PartialMatchString {
					if linkedFilter.Filter.NotCondition {
						sqlWhere = fmt.Sprintf("%s %s NOT LIKE @%s", sqlWhere, column, paramKey)
					} else {
						sqlWhere = fmt.Sprintf("%s %s LIKE @%s", sqlWhere, column, paramKey)
					}
					strVal, ok := value.(string)
					if !ok {
						return "", nil, lib_errors.Errorf("Linked filter %d with filter with key %q of type %T, expected string for partial match string search", i, paramKey, linkedFilter.Filter.Value)
					}
					value = fmt.Sprintf("%%%s%%", strVal)
				} else if linkedFilter.Filter.NotCondition {
					sqlWhere = fmt.Sprintf("%s %s != @%s", sqlWhere, column, paramKey)
				} else {
					sqlWhere = fmt.Sprintf("%s %s = @%s", sqlWhere, column, paramKey)
				}
				params[paramKey] = value
			}

		} else {
			return "", nil, lib_errors.New("Linked filter exists without type or filter")
		}
	}

	return sqlWhere, params, nil
	//revive:enable:cyclomatic
}
