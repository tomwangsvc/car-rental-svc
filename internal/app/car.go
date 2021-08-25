package app

import (
	"car-svc/internal/lib/dto"
	"context"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
)

func (c client) CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error) {
	lib_log.Info(ctx, "Creating", lib_log.FmtAny("carCreate", carCreate))

	carId, err := c.spannerClient.CreateCar(ctx, carCreate)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed creating car")
	}

	lib_log.Info(ctx, "Created", lib_log.FmtString("carId", carId))
	return carId, nil
}

func (c client) SearchCars(ctx context.Context, carsSearch dto.CarsSearch) ([]byte, *lib_pagination.Pagination, error) {
	lib_log.Info(ctx, "Searching", lib_log.FmtAny("carsSearch", carsSearch))

	cars, pagination, err := c.spannerClient.SearchCars(ctx, carsSearch)
	if err != nil {
		return nil, nil, lib_errors.Wrap(err, "Failed searching cars")
	}

	carsResponse, err := c.spannerClient.TransformCarsToJson(ctx, cars)
	if err != nil {
		return nil, nil, lib_errors.Wrap(err, "Failed transforming cars to response")
	}

	lib_log.Info(ctx, "Searched", lib_log.FmtInt("len(carsResponse)", len(carsResponse)))
	return carsResponse, pagination, nil
}

func (c client) ReadCar(ctx context.Context, carRead dto.CarRead) ([]byte, error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtAny("carRead", carRead))

	car, err := c.spannerClient.ReadCar(ctx, carRead)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed read car")
	}

	carResponse, err := c.spannerClient.TransformCarToJson(ctx, *car)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed transforming car to response")
	}

	lib_log.Info(ctx, "Read", lib_log.FmtInt("len(carResponse)", len(carResponse)))
	return carResponse, nil
}

func (c client) UpdateCar(ctx context.Context, carUpdate dto.CarUpdate) error {
	lib_log.Info(ctx, "Updating", lib_log.FmtAny("carUpdate", carUpdate))

	if err := c.spannerClient.UpdateCar(ctx, carUpdate); err != nil {
		return lib_errors.Wrap(err, "Failed updating car")
	}

	lib_log.Info(ctx, "Updated")
	return nil
}

func (c client) DeleteCar(ctx context.Context, carDelete dto.CarDelete) error {
	lib_log.Info(ctx, "Deleting", lib_log.FmtAny("carDelete", carDelete))

	if err := c.spannerClient.DeleteCar(ctx, carDelete); err != nil {
		return lib_errors.Wrap(err, "Failed deleting car")
	}

	lib_log.Info(ctx, "Deleted", lib_log.FmtAny("carDelete", carDelete))
	return nil
}
