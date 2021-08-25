package integration

import (
	"context"
	"net/http"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_token_iam "github.com/tomwangsvc/lib-svc/token/iam"
	lib_token_svc "github.com/tomwangsvc/lib-svc/token/svc"
)

type Client interface {
	DoRequestUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error)
	DoRequestUsingIamAuthorizationNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error)

	DoRequestUsingIamAuthorizationWithBodyLogRedactions(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error)
	DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error)

	DoRequestNotUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error)
	DoRequestNotUsingIamAuthorizationNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error)

	DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error)
	DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error)

	AuthenticateWithIam(ctx context.Context) error
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, tokenIamClient lib_token_iam.Client, tokenSvcClient lib_token_svc.Client) Client {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))
	lib_log.Info(ctx, "Initialized")
	return &client{
		config:         config,
		tokenIamClient: tokenIamClient,
		tokenSvcClient: tokenSvcClient,
	}
}

type client struct {
	config         Config
	iamToken       string
	iamTokenTest   string
	tokenIamClient lib_token_iam.Client
	tokenSvcClient lib_token_svc.Client
}
