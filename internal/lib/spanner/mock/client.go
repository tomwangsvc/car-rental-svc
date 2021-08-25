package mock

import (
	"car-svc/internal/lib/dto"
	"car-svc/internal/lib/spanner"
	"context"
	"encoding/binary"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
)

var (
	ClientError          spanner.Client = clientError{}
	ClientErrorTransform spanner.Client = clientErrorTransform{}
	ClientSuccess        spanner.Client = clientSuccess{}

	ExpectedErrorClient = lib_errors.NewCustom(int(binary.BigEndian.Uint64([]byte("INTEGRATION_CLIENT"))), "")
)

type clientError struct{}

func (c clientError) Close() {}

func (c clientError) TransformCarToJson(_ context.Context, _ spanner.Car) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientError) TransformCarsToJson(_ context.Context, _ []spanner.Car) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientError) CreateCar(_ context.Context, _ dto.CarCreate) (string, error) {
	return "", ExpectedErrorClient
}

func (c clientError) SearchCars(_ context.Context, _ dto.CarsSearch) ([]spanner.Car, *lib_pagination.Pagination, error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) ReadCar(_ context.Context, _ dto.CarRead) (*spanner.Car, error) {
	return nil, ExpectedErrorClient
}

func (c clientError) UpdateCar(_ context.Context, _ dto.CarUpdate) error {
	return ExpectedErrorClient
}

func (c clientError) DeleteCar(_ context.Context, _ dto.CarDelete) error {
	return ExpectedErrorClient
}

func (c clientError) TransformBrandClassAssociationToJson(_ context.Context, _ spanner.CarCustomerAssociation) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientError) TransformBrandClassAssociationsToJson(_ context.Context, _ []spanner.CarCustomerAssociation) ([]byte, error) {
	return nil, ExpectedErrorClient
}

type clientErrorTransform struct{}

func (c clientErrorTransform) Close() {}

func (c clientErrorTransform) TransformCarToJson(_ context.Context, _ spanner.Car) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientErrorTransform) TransformCarsToJson(_ context.Context, _ []spanner.Car) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientErrorTransform) CreateCar(_ context.Context, _ dto.CarCreate) (string, error) {
	return "", ExpectedErrorClient
}

func (c clientErrorTransform) SearchCars(_ context.Context, _ dto.CarsSearch) ([]spanner.Car, *lib_pagination.Pagination, error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientErrorTransform) ReadCar(_ context.Context, _ dto.CarRead) (*spanner.Car, error) {
	return nil, ExpectedErrorClient
}

func (c clientErrorTransform) UpdateCar(_ context.Context, _ dto.CarUpdate) error {
	return ExpectedErrorClient
}

func (c clientErrorTransform) DeleteCar(_ context.Context, _ dto.CarDelete) error {
	return ExpectedErrorClient
}

func (c clientErrorTransform) TransformBrandClassAssociationToJson(_ context.Context, _ spanner.CarCustomerAssociation) ([]byte, error) {
	return nil, ExpectedErrorClient
}

func (c clientErrorTransform) TransformBrandClassAssociationsToJson(_ context.Context, _ []spanner.CarCustomerAssociation) ([]byte, error) {
	return nil, ExpectedErrorClient
}

type clientSuccess struct{}

func (c clientSuccess) Close() {}

func (c clientSuccess) TransformCarToJson(_ context.Context, _ spanner.Car) ([]byte, error) {
	return lib_mock.ExpectedResultBytes, nil
}

func (c clientSuccess) TransformCarsToJson(_ context.Context, _ []spanner.Car) ([]byte, error) {
	return lib_mock.ExpectedResultBytes, nil
}

func (c clientSuccess) CreateCar(_ context.Context, _ dto.CarCreate) (string, error) {
	return lib_mock.ExpectedResultString, nil
}

func (c clientSuccess) SearchCars(_ context.Context, _ dto.CarsSearch) ([]spanner.Car, *lib_pagination.Pagination, error) {
	return []spanner.Car{{}}, nil, nil
}

func (c clientSuccess) ReadCar(_ context.Context, _ dto.CarRead) (*spanner.Car, error) {
	return &spanner.Car{}, nil
}

func (c clientSuccess) UpdateCar(_ context.Context, _ dto.CarUpdate) error {
	return nil
}

func (c clientSuccess) DeleteCar(_ context.Context, _ dto.CarDelete) error {
	return nil
}

func (c clientSuccess) TransformBrandClassAssociationToJson(_ context.Context, _ spanner.CarCustomerAssociation) ([]byte, error) {
	return lib_mock.ExpectedResultBytes, nil
}

func (c clientSuccess) TransformBrandClassAssociationsToJson(_ context.Context, _ []spanner.CarCustomerAssociation) ([]byte, error) {
	return lib_mock.ExpectedResultBytes, nil
}
