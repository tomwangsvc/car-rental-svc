package mock

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_integration "github.com/tomwangsvc/lib-svc/integration"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
)

var (
	ClientError                                                    lib_integration.Client = clientError{}
	ClientSuccessStatusAccepted                                    lib_integration.Client = clientSuccessStatusAccepted{}
	ClientSuccessStatusAcceptedNoLocation                          lib_integration.Client = clientSuccessStatusAcceptedNoLocation{}
	ClientSuccessStatusBadGateway                                  lib_integration.Client = clientSuccessStatusBadGateway{}
	ClientSuccessStatusCreated                                     lib_integration.Client = clientSuccessStatusCreated{}
	ClientSuccessStatusCreatedNoLocation                           lib_integration.Client = clientSuccessStatusCreatedNoLocation{}
	ClientSuccessStatusCreatedBodyObject                           lib_integration.Client = clientSuccessStatusCreatedBodyObject{}
	ClientSuccessStatusCreatedBodyObjectNoLocation                 lib_integration.Client = clientSuccessStatusCreatedBodyObjectNoLocation{}
	ClientSuccessStatusConflict                                    lib_integration.Client = clientSuccessStatusConflict{}
	ClientSuccessStatusConflictWithLocation                        lib_integration.Client = clientSuccessStatusConflictWithLocation{}
	ClientSuccessStatusGatewayTimeout                              lib_integration.Client = clientSuccessStatusGatewayTimeout{}
	ClientSuccessStatusInternalServerError                         lib_integration.Client = clientSuccessStatusInternalServerError{}
	ClientSuccessStatusNoContent                                   lib_integration.Client = clientSuccessStatusNoContent{}
	ClientSuccessStatusNotFound                                    lib_integration.Client = clientSuccessStatusNotFound{}
	ClientSuccessStatusOkBodyList                                  lib_integration.Client = clientSuccessStatusOkBodyList{}
	ClientSuccessStatusOkBodyListWithObject                        lib_integration.Client = clientSuccessStatusOkBodyListWithObject{}
	ClientSuccessStatusOkBodyObject                                lib_integration.Client = clientSuccessStatusOkBodyObject{}
	ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOne  lib_integration.Client = ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject{}
	ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZero lib_integration.Client = ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject{}
	ClientSuccessStatusServiceUnavailable                          lib_integration.Client = clientSuccessStatusServiceUnavailable{}
	ClientSuccessStatusUnauthorized                                lib_integration.Client = clientSuccessStatusUnauthorized{}
	ClientSuccessStatusUnprocessableEntity                         lib_integration.Client = clientSuccessStatusUnprocessableEntity{}

	ExpectedErrorClient                                                         = lib_errors.NewCustom(int(binary.BigEndian.Uint64([]byte("LIB_INTEGRATION_CLIENT"))), "Mock miscellaneous error")
	ExpectedResultHttpResponseStatusAccepted                                    = http.Response{StatusCode: http.StatusAccepted, Header: http.Header{"Location": []string{lib_mock.ExpectedResultString}}}
	ExpectedResultHttpResponseStatusAcceptedNoLocation                          = http.Response{StatusCode: http.StatusAccepted, Header: http.Header{"Location": []string{lib_mock.ExpectedResultString}}}
	ExpectedResultHttpResponseStatusBadGateway                                  = http.Response{StatusCode: http.StatusBadGateway}
	ExpectedResultHttpResponseStatusCreated                                     = http.Response{StatusCode: http.StatusCreated, Header: http.Header{"Location": []string{lib_mock.ExpectedResultString}}}
	ExpectedResultHttpResponseStatusCreatedNoLocation                           = http.Response{StatusCode: http.StatusCreated}
	ExpectedResultHttpResponseStatusConflict                                    = http.Response{StatusCode: http.StatusConflict}
	ExpectedResultHttpResponseStatusConflictWithLocation                        = http.Response{StatusCode: http.StatusConflict, Header: http.Header{"Location": []string{lib_mock.ExpectedResultString}}}
	ExpectedResultHttpResponseStatusGatewayTimeout                              = http.Response{StatusCode: http.StatusGatewayTimeout}
	ExpectedResultHttpResponseStatusInternalServerError                         = http.Response{StatusCode: http.StatusInternalServerError}
	ExpectedResultHttpResponseStatusNoContent                                   = http.Response{StatusCode: http.StatusNoContent}
	ExpectedResultHttpResponseStatusNotFound                                    = http.Response{StatusCode: http.StatusNotFound}
	ExpectedResultHttpResponseStatusOk                                          = http.Response{StatusCode: http.StatusOK}
	ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne  = http.Response{StatusCode: http.StatusNoContent, Header: lib_mock.ExpectedHeaderForXLcPaginationTotalOfOne}
	ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero = http.Response{StatusCode: http.StatusNoContent, Header: lib_mock.ExpectedHeaderForXLcPaginationTotalOfZero}
	ExpectedResultHttpResponseStatusServiceUnavailable                          = http.Response{StatusCode: http.StatusServiceUnavailable}
	ExpectedResultHttpResponseStatusUnauthorized                                = http.Response{StatusCode: http.StatusUnauthorized}
	ExpectedResultHttpResponseStatusUnprocessableEntity                         = http.Response{StatusCode: http.StatusUnprocessableEntity}
	ExpectedResultIntegrationListBytes                                          = []byte("[]")
	ExpectedResultIntegrationListWithObjectBytes                                = []byte("[{}]")
	ExpectedResultIntegrationObjectBytes                                        = []byte("{}")
	ExpectedResultUnprocessableEntityBytes, _                                   = json.Marshal([]lib_errors.Item{lib_errors.Item{}})
)

type clientError struct{}

func (c clientError) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return nil, nil, ExpectedErrorClient
}

func (c clientError) AuthenticateWithIam(_ context.Context) error {
	return ExpectedErrorClient
}

type clientSuccessStatusAccepted struct{}

func (c clientSuccessStatusAccepted) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAccepted, nil, nil
}

func (c clientSuccessStatusAccepted) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusAcceptedNoLocation struct{}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusAcceptedNoLocation, nil, nil
}

func (c clientSuccessStatusAcceptedNoLocation) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusBadGateway struct{}

func (c clientSuccessStatusBadGateway) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusBadGateway, nil, nil
}

func (c clientSuccessStatusBadGateway) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusCreated struct{}

func (c clientSuccessStatusCreated) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, nil, nil
}

func (c clientSuccessStatusCreated) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusCreatedNoLocation struct{}

func (c clientSuccessStatusCreatedNoLocation) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, nil, nil
}

func (c clientSuccessStatusCreatedNoLocation) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusCreatedBodyObject struct{}

func (c clientSuccessStatusCreatedBodyObject) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreated, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObject) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusCreatedBodyObjectNoLocation struct{}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusCreatedNoLocation, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusCreatedBodyObjectNoLocation) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusConflict struct{}

func (c clientSuccessStatusConflict) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflict, nil, nil
}

func (c clientSuccessStatusConflict) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusConflictWithLocation struct{}

func (c clientSuccessStatusConflictWithLocation) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusConflictWithLocation, nil, nil
}

func (c clientSuccessStatusConflictWithLocation) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusGatewayTimeout struct{}

func (c clientSuccessStatusGatewayTimeout) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusGatewayTimeout, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusGatewayTimeout) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusInternalServerError struct{}

func (c clientSuccessStatusInternalServerError) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusInternalServerError, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusInternalServerError) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusNoContent struct{}

func (c clientSuccessStatusNoContent) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContent, nil, nil
}

func (c clientSuccessStatusNoContent) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusNotFound struct{}

func (c clientSuccessStatusNotFound) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNotFound, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusNotFound) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusOkBodyList struct{}

func (c clientSuccessStatusOkBodyList) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListBytes, nil
}

func (c clientSuccessStatusOkBodyList) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusOkBodyListWithObject struct{}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationListWithObjectBytes, nil
}

func (c clientSuccessStatusOkBodyListWithObject) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusOkBodyObject struct{}

func (c clientSuccessStatusOkBodyObject) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusOk, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusOkBodyObject) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject struct{}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfOne, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfOneObject) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject struct{}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusNoContentWithHeaderXLcPaginationTotalOfZero, ExpectedResultIntegrationObjectBytes, nil
}

func (c ClientSuccessStatusNoContentWithHeaderXLcPaginationTotalOfZeroObject) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusServiceUnavailable struct{}

func (c clientSuccessStatusServiceUnavailable) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusServiceUnavailable, ExpectedResultIntegrationObjectBytes, nil
}

func (c clientSuccessStatusServiceUnavailable) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusUnauthorized struct{}

func (c clientSuccessStatusUnauthorized) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnauthorized, nil, nil
}

func (c clientSuccessStatusUnauthorized) AuthenticateWithIam(_ context.Context) error {
	return nil
}

type clientSuccessStatusUnprocessableEntity struct{}

func (c clientSuccessStatusUnprocessableEntity) DoRequestUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestNotUsingIamAuthorization(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestNotUsingIamAuthorizationNoRetries(_ context.Context, _ *http.Request, _, _ bool) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(_ context.Context, _ *http.Request, _, _ bool, _ map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	return &ExpectedResultHttpResponseStatusUnprocessableEntity, ExpectedResultUnprocessableEntityBytes, nil
}

func (c clientSuccessStatusUnprocessableEntity) AuthenticateWithIam(_ context.Context) error {
	return nil
}
