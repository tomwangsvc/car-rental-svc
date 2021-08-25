package http

import (
	"car-svc/internal/app"
	"car-svc/internal/http/routes"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	lib_countries "github.com/tomwangsvc/lib-svc/countries"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_schema "github.com/tomwangsvc/lib-svc/schema"
)

// @title car-svc

type Client interface {
	ListenAndServe() error
}

type Config struct {
	Env lib_env.Env
}

type client struct {
	config Config
	router *chi.Mux
}

func NewClient(
	config Config,
	appClient app.Client,
	schemaClient lib_schema.Client,
	countriesMetadata lib_countries.Metadata,
) (Client, error) {

	routesClient := routes.NewClient(routes.Config{Env: config.Env}, appClient, schemaClient, countriesMetadata)

	r := chi.NewRouter()
	lib_http.GeneralMiddleware(r, config.Env.Id, config.Env.MaintenanceMode, []string{"/car-svc?health=true"})

	r.Route("/car-svc", func(r chi.Router) {
		r.Get("/", routesClient.Health())
	})

	r.Route("/car-svc/v1", func(r chi.Router) {
		// TODO: authorize locally
		// r.Use(iamClient.Authorize)
		r.Route("/cars", func(r chi.Router) {
			r.Post("/", routesClient.CreateCar())
			r.Get("/", routesClient.SearchCars())

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", routesClient.ReadCar())
				r.Put("/", routesClient.UpdateCar())
				r.Delete("/", routesClient.DeleteCar())
			})
		})
	})

	return client{
		config: config,
		router: r,
	}, nil
}

func (c client) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", c.config.Env.Port), c.router)
}
