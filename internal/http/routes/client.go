package routes

import (
	"car-svc/internal/app"
	"car-svc/internal/http/routes/parser"
	"net/http"

	lib_countries "github.com/tomwangsvc/lib-svc/countries"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_schema "github.com/tomwangsvc/lib-svc/schema"
)

type Client interface {
	Health() http.HandlerFunc
	CreateCar() http.HandlerFunc
	SearchCars() http.HandlerFunc
	ReadCar() http.HandlerFunc
	UpdateCar() http.HandlerFunc
	DeleteCar() http.HandlerFunc
}

type Config struct {
	Env lib_env.Env
}

func NewClient(config Config, appClient app.Client, schemaClient lib_schema.Client, countriesMetadata lib_countries.Metadata) Client {
	return client{
		config:            config,
		appClient:         appClient,
		countriesMetadata: countriesMetadata,
		parserClient:      parser.NewClient(parser.Config{Env: config.Env}, schemaClient, countriesMetadata),
		schemaClient:      schemaClient,
	}
}

type client struct {
	config            Config
	appClient         app.Client
	countriesMetadata lib_countries.Metadata
	parserClient      parser.Client
	schemaClient      lib_schema.Client
}
