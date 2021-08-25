package parser

import (
	"car-svc/internal/lib/dto"
	"net/http"

	lib_countries "github.com/tomwangsvc/lib-svc/countries"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_schema "github.com/tomwangsvc/lib-svc/schema"
)

type Client interface {
	ParseCreateCar(r *http.Request) (*dto.CarCreate, error)
	ParseSearchCars(r *http.Request) (*dto.CarsSearch, error)
	ParseReadCar(r *http.Request) (*dto.CarRead, error)
	ParseUpdateCar(r *http.Request) (*dto.CarUpdate, error)
	ParseDeleteCar(r *http.Request) (*dto.CarDelete, error)
}

type Config struct {
	Env lib_env.Env
}

func NewClient(config Config, schemaClient lib_schema.Client, countriesMetadata lib_countries.Metadata) Client {
	return client{
		config:            config,
		schemaClient:      schemaClient,
		countriesMetadata: countriesMetadata,
	}
}

type client struct {
	config            Config
	schemaClient      lib_schema.Client
	countriesMetadata lib_countries.Metadata
}
