package mock

import (
	"car-svc/internal/http/routes/parser"
	"car-svc/internal/lib/dto"
	"net/http"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

var (
	ClientError   parser.Client = clientError{}
	ClientSuccess parser.Client = clientSuccess{}

	ExpectedErrorClient = lib_errors.NewCustom(500, "Mock miscellaneous error")
)

type clientError struct{}

func (clientError) ParseCreateCar(_ *http.Request) (*dto.CarCreate, error) {
	return nil, ExpectedErrorClient
}

func (clientError) ParseSearchCars(_ *http.Request) (*dto.CarsSearch, error) {
	return nil, ExpectedErrorClient
}

func (clientError) ParseReadCar(_ *http.Request) (*dto.CarRead, error) {
	return nil, ExpectedErrorClient
}

func (clientError) ParseUpdateCar(_ *http.Request) (*dto.CarUpdate, error) {
	return nil, ExpectedErrorClient
}

func (clientError) ParseDeleteCar(_ *http.Request) (*dto.CarDelete, error) {
	return nil, ExpectedErrorClient
}

type clientSuccess struct{}

func (clientSuccess) ParseCreateCar(_ *http.Request) (*dto.CarCreate, error) {
	return &dto.CarCreate{}, nil
}

func (clientSuccess) ParseSearchCars(_ *http.Request) (*dto.CarsSearch, error) {
	return &dto.CarsSearch{}, nil
}

func (clientSuccess) ParseReadCar(_ *http.Request) (*dto.CarRead, error) {
	return &dto.CarRead{}, nil
}

func (clientSuccess) ParseUpdateCar(_ *http.Request) (*dto.CarUpdate, error) {
	return &dto.CarUpdate{}, nil
}

func (clientSuccess) ParseDeleteCar(_ *http.Request) (*dto.CarDelete, error) {
	return &dto.CarDelete{}, nil
}
