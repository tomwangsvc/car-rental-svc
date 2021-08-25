package gcp

import (
	"context"
	"encoding/json"
	"net/http"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_integration "github.com/tomwangsvc/lib-svc/integration"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_token "github.com/tomwangsvc/lib-svc/token"
)

type Client interface {
	VerifyForCloudTasksPush(ctx context.Context, jwtToken string) error
	VerifyForPubsubPush(ctx context.Context, jwtToken string) error
	VerifyForCloudSchedulerPush(ctx context.Context, jwtToken string) error
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, integrationClient lib_integration.Client) Client {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))
	lib_log.Info(ctx, "Initialized")
	return client{
		config:            config,
		integrationClient: integrationClient,
	}
}

type client struct {
	config            Config
	integrationClient lib_integration.Client
}

func (c client) VerifyForCloudTasksPush(ctx context.Context, idToken string) error {
	lib_log.Info(ctx, "Verifying", lib_log.FmtString("lib_token.Redact(idToken)", lib_token.Redact(idToken)))

	tokenInfo, err := c.verifyWithGoogleApi(ctx, idToken)
	if err != nil {
		return lib_errors.Wrap(err, "Failed verifying with google api")
	}

	if err := c.checkTokenInfoForCloudTasksPush(*tokenInfo); err != nil {
		return lib_errors.Wrap(err, "Failed checking audience")
	}

	lib_log.Info(ctx, "Verified")
	return nil
}

func (c client) VerifyForPubsubPush(ctx context.Context, idToken string) error {
	lib_log.Info(ctx, "Verifying", lib_log.FmtString("lib_token.Redact(idToken)", lib_token.Redact(idToken)))

	tokenInfo, err := c.verifyWithGoogleApi(ctx, idToken)
	if err != nil {
		return lib_errors.Wrap(err, "Failed verifying with google api")
	}

	if err := c.checkTokenInfoForPubsubPush(*tokenInfo); err != nil {
		return lib_errors.Wrap(err, "Failed checking audience")
	}

	lib_log.Info(ctx, "Verified")
	return nil
}

func (c client) VerifyForCloudSchedulerPush(ctx context.Context, idToken string) error {
	lib_log.Info(ctx, "Verifying", lib_log.FmtString("lib_token.Redact(idToken)", lib_token.Redact(idToken)))

	tokenInfo, err := c.verifyWithGoogleApi(ctx, idToken)
	if err != nil {
		return lib_errors.Wrap(err, "Failed verifying with google api")
	}

	if err := c.checkTokenInfoForCloudSchedulerPush(*tokenInfo); err != nil {
		return lib_errors.Wrap(err, "Failed checking audience")
	}

	lib_log.Info(ctx, "Verified")
	return nil
}

type tokenInfo struct {
	Audience      string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

func (c client) verifyWithGoogleApi(ctx context.Context, idToken string) (*tokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://oauth2.googleapis.com/tokeninfo", nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating payee read request")
	}
	req = req.WithContext(ctx)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Set("id_token", idToken)
	req.URL.RawQuery = q.Encode()

	res, resBody, err := c.integrationClient.DoRequestNotUsingIamAuthorization(ctx, req, true, true)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed doing request using iam authorization")
	}

	switch res.StatusCode {
	default:
		return nil, lib_errors.NewCustomf(http.StatusUnauthorized, "Request resulted in http status code %d", res.StatusCode)

	case http.StatusBadGateway:
		return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Request resulted in http status code %d", res.StatusCode)

	case http.StatusGatewayTimeout:
		return nil, lib_errors.NewCustomf(http.StatusGatewayTimeout, "Request resulted in http status code %d", res.StatusCode)

	case http.StatusOK:

	case http.StatusServiceUnavailable:
		return nil, lib_errors.NewCustomf(http.StatusServiceUnavailable, "Request resulted in http status code %d", res.StatusCode)
	}

	var tokenInfo tokenInfo
	if err := json.Unmarshal(resBody, &tokenInfo); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling resBody into tokenInfo")
	}

	return &tokenInfo, nil
}

func (c client) checkTokenInfoForCloudTasksPush(tokenInfo tokenInfo) error {
	if tokenInfo.Audience != c.config.Env.SvcId {
		return lib_errors.Errorf("Claim 'aud' %q does not match expected aud %q", tokenInfo.Audience, c.config.Env.SvcId)
	}

	if tokenInfo.Email != c.config.Env.CloudTasksPushServiceAccount {
		return lib_errors.Errorf("Claim 'email' %q does not match expected cloud tasks push service account %q", tokenInfo.Email, c.config.Env.CloudTasksPushServiceAccount)
	}
	if tokenInfo.EmailVerified != "true" {
		return lib_errors.Errorf("Claim 'email_verified' %q must be true", tokenInfo.EmailVerified)
	}
	return nil
}

func (c client) checkTokenInfoForPubsubPush(tokenInfo tokenInfo) error {
	if tokenInfo.Audience != c.config.Env.SvcId {
		return lib_errors.Errorf("Claim 'aud' %q does not match expected aud %q", tokenInfo.Audience, c.config.Env.SvcId)
	}

	if tokenInfo.Email != c.config.Env.PubsubPushServiceAccount {
		return lib_errors.Errorf("Claim 'email' %q does not match expected pubsub push service account %q", tokenInfo.Email, c.config.Env.PubsubPushServiceAccount)
	}
	if tokenInfo.EmailVerified != "true" {
		return lib_errors.Errorf("Claim 'email_verified' %q must be true", tokenInfo.EmailVerified)
	}
	return nil
}

func (c client) checkTokenInfoForCloudSchedulerPush(tokenInfo tokenInfo) error {
	if tokenInfo.Audience != c.config.Env.SvcId {
		return lib_errors.Errorf("Claim 'aud' %q does not match expected aud %q", tokenInfo.Audience, c.config.Env.SvcId)
	}

	if tokenInfo.Email != c.config.Env.CloudSchedulerPushServiceAccount {
		return lib_errors.Errorf("Claim 'email' %q does not match expected cloud scheduler service push account %q", tokenInfo.Email, c.config.Env.CloudSchedulerPushServiceAccount)
	}
	if tokenInfo.EmailVerified != "true" {
		return lib_errors.Errorf("Claim 'email_verified' %q must be true", tokenInfo.EmailVerified)
	}
	return nil
}
