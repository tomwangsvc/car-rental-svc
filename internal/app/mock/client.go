package mock

import (
	"car-svc/internal/app"
	"car-svc/internal/lib/dto"
	"context"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
)

var (
	ClientError   app.Client = clientError{}
	ClientSuccess app.Client = clientSuccess{}

	ExpectedErrorClient = lib_errors.NewCustom(501, "Mock miscellaneous error")
)

type clientError struct{}

func (clientError) CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error) {
	return "", ExpectedErrorClient
}

func (clientError) SearchCars(ctx context.Context, carCreate dto.CarsSearch) ([]byte, *lib_pagination.Pagination, error) {
	return nil, nil, ExpectedErrorClient
}

func (clientError) ReadCar(ctx context.Context, carCreate dto.CarRead) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (clientError) UpdateCar(ctx context.Context, carCreate dto.CarUpdate) error {
	return ExpectedErrorClient
}

func (clientError) DeleteCar(ctx context.Context, carCreate dto.CarDelete) error {
	return ExpectedErrorClient
}

type clientSuccess struct{}

func (clientSuccess) CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error) {
	return lib_mock.ExpectedResultString, nil
}

func (clientSuccess) SearchCars(ctx context.Context, carCreate dto.CarsSearch) ([]byte, *lib_pagination.Pagination, error) {
	return lib_mock.ExpectedResultBytes, nil, nil
}

func (clientSuccess) ReadCar(ctx context.Context, carCreate dto.CarRead) ([]byte, error) {
	return lib_mock.ExpectedResultBytes, nil
}

func (clientSuccess) UpdateCar(ctx context.Context, carCreate dto.CarUpdate) error {
	return nil
}

func (clientSuccess) DeleteCar(ctx context.Context, carCreate dto.CarDelete) error {
	return nil
}
