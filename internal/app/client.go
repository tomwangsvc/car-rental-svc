package app

import (
	"car-svc/internal/lib/dto"
	"car-svc/internal/lib/spanner"
	"context"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
)

type Client interface {
	CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error)
	SearchCars(ctx context.Context, categoriesSearch dto.CarsSearch) ([]byte, *lib_pagination.Pagination, error)
	ReadCar(ctx context.Context, carRead dto.CarRead) ([]byte, error)
	UpdateCar(ctx context.Context, carUpdate dto.CarUpdate) error
	DeleteCar(ctx context.Context, carDelete dto.CarDelete) error
}

type Config struct {
	Env lib_env.Env
}

func NewClient(config Config, spannerClient spanner.Client) Client {
	return client{
		config:        config,
		spannerClient: spannerClient,
	}
}

type client struct {
	config        Config
	spannerClient spanner.Client
}
