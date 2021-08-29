package routes

import (
	"car-svc/internal/lib/schema"
	"net/http"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

// @Summary create car
// @Param Authorization header string true "IAM token"
// @Description create car
// @Description See schema file car_create.json for body requirements
// @Param Authorization header string not implementeed
// @Success 201
// @Header 201 {string} Location "id"
// @Router /v1/cars[post]
func (c client) CreateCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Creating")

		carCreate, err := c.parserClient.ParseCreateCar(r)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed parsing create car request"))
			return
		}

		carId, err := c.appClient.CreateCar(ctx, *carCreate)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed creating car"))
			return
		}

		lib_log.Info(ctx, "Created", lib_log.FmtString("carId", carId))
		lib_http.RenderCreated(ctx, w, carId)
	}
}

// @Summary search cars
// @Param Authorization header string true "IAM token"
// @Description search cars
// @Description See schema file cars_search.json for query params
// @Description See schema file cars.json for response
// @Param Authorization header string not implementeed
// @Success 200
// @Router /v1/cars[get]
func (c client) SearchCars() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Searching")

		carsSearch, err := c.parserClient.ParseSearchCars(r)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed parsing search cars request"))
			return
		}

		carsBytes, pagination, err := c.appClient.SearchCars(ctx, *carsSearch)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed search cars"))
			return
		}

		if len(carsBytes) == 0 {
			lib_http.RenderNoContent(ctx, w)
			return
		}

		if err := c.schemaClient.CheckContentAgainstSchema(ctx, schema.Cars, carsBytes); err != nil {
			if carsSearch.IntegrationTest {
				lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed checking request body against schema, in integration test, will return new error"))
				return
			}
			lib_log.Error(ctx, "Failed checking request body against schema, something is likely misconfigured, will return success", lib_log.FmtError(err))
		}

		lib_log.Info(ctx, "Searched", lib_log.FmtBytes("carsBytes", carsBytes), lib_log.FmtAny("pagination", pagination))
		lib_http.RenderJsonBytesWithPagination(ctx, w, carsBytes, *pagination)
	}
}

// @Summary read car
// @Param Authorization header string true "IAM token"
// @Description read car
// @Description See schema file car_read.json for response
// @Param Authorization header string not implementeed
// @Success 200
// @Router /v1/cars/{car_id}[get]
func (c client) ReadCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Reading")

		carRead, err := c.parserClient.ParseReadCar(r)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed parsing search cars request"))
			return
		}

		car, err := c.appClient.ReadCar(ctx, *carRead)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed search cars"))
			return
		}

		if err := c.schemaClient.CheckContentAgainstSchema(ctx, schema.Car, carRead); err != nil {
			if carRead.IntegrationTest {
				lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed checking request body against schema, in integration test, will return new error"))
				return
			}
			lib_log.Error(ctx, "Failed checking request body against schema, something is likely misconfigured, will return success", lib_log.FmtError(err))
		}

		lib_log.Info(ctx, "Read", lib_log.FmtInt("len(car)", len(car)))
		lib_http.RenderJsonBytes(ctx, w, car)
	}
}

// @Summary update car
// @Param Authorization header string true "IAM token"
// @Description update car
// @Description See schema file car_update.json for user input
// @Param Authorization header string not implementeed
// @Success 204
// @Router /v1/cars/{car_id}[put]
func (c client) UpdateCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Updating")

		carUpdate, err := c.parserClient.ParseUpdateCar(r)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed parsing update car request"))
			return
		}

		if err := c.appClient.UpdateCar(ctx, *carUpdate); err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed updating car"))
			return
		}

		lib_log.Info(ctx, "Updated")
		lib_http.RenderNoContent(ctx, w)
	}
}

// @Summary delete car
// @Param Authorization header string true "IAM token"
// @Description delete car
// @Description See schema file car_delete.json for user input
// @Param Authorization header string not implementeed
// @Success 204
// @Router /v1/cars/{car_id}[delete]
func (c client) DeleteCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Deleting")

		carDelete, err := c.parserClient.ParseDeleteCar(r)
		if err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed parsing delete car request"))
			return
		}

		if err := c.appClient.DeleteCar(ctx, *carDelete); err != nil {
			lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed deleting car"))
			return
		}

		lib_log.Info(ctx, "Deleted")
		lib_http.RenderNoContent(ctx, w)
	}
}
