package integration

import (
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_integration "github.com/tomwangsvc/lib-svc/integration"
)

type Client interface {
}

type Config struct {
	Env        lib_env.Env
	UrlStorage string
}

func NewClient(config Config, integrationClient lib_integration.Client) Client {
	return client{
		config:            config,
		integrationClient: integrationClient,
	}
}

type client struct {
	config            Config
	integrationClient lib_integration.Client
}
