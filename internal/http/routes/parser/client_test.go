package parser

import (
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
	lib_schema_mock "github.com/tomwangsvc/lib-svc/schema/mock"
)

var (
	clientErrorLibSchema = client{
		countriesMetadata: lib_mock.ExpectedResultLibCountryMetadata,
		schemaClient:      lib_schema_mock.ClientError,
	}

	clientSuccess = client{
		countriesMetadata: lib_mock.ExpectedResultLibCountryMetadata,
		schemaClient:      lib_schema_mock.ClientSuccess,
	}
)
