package integration

import (
	lib_integration_mock "github.com/tomwangsvc/lib-svc/integration/mock"
)

var (
	clientErrorIntegration = client{
		integrationClient: lib_integration_mock.ClientError,
	}
	clientSuccessIntegrationStatusInternalServerError = client{
		integrationClient: lib_integration_mock.ClientSuccessStatusInternalServerError,
	}
	clientSuccessIntegrationStatusOkBodyList = client{
		integrationClient: lib_integration_mock.ClientSuccessStatusOkBodyList,
	}
	clientSuccessIntegrationStatusNoContent = client{
		integrationClient: lib_integration_mock.ClientSuccessStatusNoContent,
	}
	clientSuccessIntegrationStatusNotFound = client{
		integrationClient: lib_integration_mock.ClientSuccessStatusNotFound,
	}
	clientSuccessIntegrationStatusOkBodyObject = client{
		integrationClient: lib_integration_mock.ClientSuccessStatusOkBodyObject,
	}
)
