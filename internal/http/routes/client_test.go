package routes

import (
	app_mock "car-svc/internal/app/mock"
	parser_mock "car-svc/internal/http/routes/parser/mock"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_schema_mock "github.com/tomwangsvc/lib-svc/schema/mock"
)

var (
	clientErrorApp = client{
		appClient:    app_mock.ClientError,
		parserClient: parser_mock.ClientSuccess,
		config:       Config{Env: lib_env.Env{Id: lib_env.Dev}},
		schemaClient: lib_schema_mock.ClientSuccess,
	}
	clientErrorParser = client{
		appClient:    app_mock.ClientSuccess,
		parserClient: parser_mock.ClientError,
		config:       Config{Env: lib_env.Env{Id: lib_env.Dev}},
		schemaClient: lib_schema_mock.ClientSuccess,
	}
	clientErrorSchema = client{
		appClient:    app_mock.ClientSuccess,
		parserClient: parser_mock.ClientSuccess,
		config:       Config{Env: lib_env.Env{Id: lib_env.Dev}},
		schemaClient: lib_schema_mock.ClientError,
	}
	clientSuccess = client{
		appClient:    app_mock.ClientSuccess,
		parserClient: parser_mock.ClientSuccess,
		config:       Config{Env: lib_env.Env{Id: lib_env.Dev}},
		schemaClient: lib_schema_mock.ClientSuccess,
	}
)
