package parser

import (
	"car-svc/internal/lib/dto"
	"car-svc/internal/lib/schema"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
	lib_search "github.com/tomwangsvc/lib-svc/search"
)

func (c client) ParseCreateCar(r *http.Request) (*dto.CarCreate, error) {
	ctx := r.Context()
	lib_log.Info(ctx, "Parsing")

	body, err := lib_http.ReadRequestBody(r, true)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed decoding body for request")
	}
	if err := c.schemaClient.CheckBodyAgainstSchema(ctx, schema.CarCreate, body); err != nil {
		return nil, lib_errors.Wrap(err, "Failed checking body against schema")
	}

	carCreate := dto.CarCreate{
		Test: lib_context.Test(ctx),
	}
	if err := json.Unmarshal(body, &carCreate.UserInput); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling body into dto.CarCreate")
	}

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("carCreate", carCreate))
	return &carCreate, nil
}

func (c client) ParseSearchCars(r *http.Request) (*dto.CarsSearch, error) {
	ctx := r.Context()
	lib_log.Info(ctx, "Parsing")
	queryEncodedQuery, err := lib_search.QueryEncodedQueryFromRawQuery(r.URL.RawQuery)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed getting query encoded query from raw query")
	}
	test := lib_context.Test(ctx)
	filtersForSchemaCheck, linkedFilters, err := lib_search.ParseQueryWithTestV3(queryEncodedQuery, test)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing query with test")
	}
	if err := c.schemaClient.CheckContentAgainstSchema(ctx, schema.CarsSearch, struct {
		Query []lib_search.Filter `json:"query,omitempty"`
	}{
		Query: filtersForSchemaCheck,
	}); err != nil {
		return nil, lib_errors.Wrap(err, "Failed checking content against schema")
	}

	pagination, err := lib_pagination.NewPagination(r, nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating pagination")
	}

	carsSearch := dto.CarsSearch{
		Filters: dto.CarsSearchFilters{
			Test:          lib_context.Test(ctx),
			LinkedFilters: linkedFilters,
		},
		Pagination: *pagination,
	}

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("carsSearch", carsSearch))
	return &carsSearch, nil
}

func (c client) ParseReadCar(r *http.Request) (*dto.CarRead, error) {
	ctx := r.Context()
	lib_log.Info(ctx, "Parsing")

	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Missing id in url params")
	}

	carRead := dto.CarRead{
		Id: id,
	}

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("carRead", carRead))
	return &carRead, nil
}

func (c client) ParseUpdateCar(r *http.Request) (*dto.CarUpdate, error) {
	ctx := r.Context()
	lib_log.Info(ctx, "Parsing")

	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Missing id in url params")
	}

	carUpdate := dto.CarUpdate{
		Id:   id,
		Test: lib_context.Test(ctx),
	}

	body, err := lib_http.ReadRequestBody(r, true)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed decoding body for request")
	}
	if err := c.schemaClient.CheckBodyAgainstSchema(ctx, schema.CarUpdate, body); err != nil {
		return nil, lib_errors.Wrap(err, "Failed checking body against schema")
	}
	if err := json.Unmarshal(body, &carUpdate.UserInput); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling body into dto.CarUpdate")
	}

	// TODO: handle concurrency from mutiple users
	// ifUnmodifiedSinceHeader := r.Header.Get("If-Unmodified-Since")
	// if ifUnmodifiedSinceHeader != "" {
	// 	var err error
	// 	ifUnmodifiedSince, err := lib_time.ParseFormattedTimeWithFullPrecision(ifUnmodifiedSinceHeader)
	// 	if err != nil {
	// 		return nil, lib_errors.Wrap(err, "Failed to parse value in If-Unmodified-Since header")
	// 	}

	// 	carUpdate.IfUnmodifiedSince = ifUnmodifiedSince
	// }

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("carUpdate", carUpdate))
	return &carUpdate, nil
}

func (c client) ParseDeleteCar(r *http.Request) (*dto.CarDelete, error) {
	ctx := r.Context()
	lib_log.Info(ctx, "Parsing")

	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Missing id in url params")
	}

	carDelete := dto.CarDelete{
		Id:   id,
		Test: lib_context.Test(ctx),
	}

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("carDelete", carDelete))
	return &carDelete, nil
}
