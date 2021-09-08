package spanner

import (
	"car-svc/internal/lib/dto"
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
	lib_spanner "github.com/tomwangsvc/lib-svc/spanner"
	"google.golang.org/api/option"
)

type Client interface {
	Close()
	TransformCarToJson(ctx context.Context, car Car) ([]byte, error)
	TransformCarsToJson(ctx context.Context, cars []Car) ([]byte, error)
	CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error)
	SearchCars(ctx context.Context, carsSearch dto.CarsSearch) ([]Car, *lib_pagination.Pagination, error)
	ReadCar(ctx context.Context, carRead dto.CarRead) (*Car, error)
	UpdateCar(ctx context.Context, carUpdate dto.CarUpdate) error
	DeleteCar(ctx context.Context, carDelete dto.CarDelete) error

	TransformBrandClassAssociationToJson(ctx context.Context, carCustomerAssociation CarCustomerAssociation) ([]byte, error)
	TransformBrandClassAssociationsToJson(ctx context.Context, carCustomerAssociations []CarCustomerAssociation) ([]byte, error)
}

type Config struct {
	spanner.ClientConfig
	DatabaseId string
	Env        lib_env.Env
	InstanceId string
	ProjectId  string
}

func NewClient(ctx context.Context, config Config) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))
	spannerClient, err := spanner.NewClientWithConfig(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", config.ProjectId, config.InstanceId, config.DatabaseId), config.ClientConfig)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating spanner client with config")
	}
	lib_log.Info(ctx, "Initialized")
	return client{
		config:        config,
		spannerClient: spannerClient,
	}, nil
}

type client struct {
	config        Config
	spannerClient *spanner.Client
}

func (c client) Close() {
	c.spannerClient.Close()
}

type ReadWriteTransaction func(context.Context, *spanner.ReadWriteTransaction) error

type searchFilters map[string]searchFilter

type searchFilter struct {
	CaseInsensitiveString bool
	PartialMatchString    bool
	Value                 interface{}
}

func (c client) ExecuteReadWriteTransaction(ctx context.Context, readWriteTransaction ReadWriteTransaction) error {
	lib_log.Info(ctx, "Executing")

	if _, err := c.spannerClient.ReadWriteTransaction(ctx, readWriteTransaction); err != nil {
		return lib_spanner.WrapError(err, "Failed executing spanner read write transaction")
	}

	lib_log.Info(ctx, "Executed")
	return nil
}
