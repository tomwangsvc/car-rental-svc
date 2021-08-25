package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_svc "github.com/tomwangsvc/lib-svc/svc"
	lib_token "github.com/tomwangsvc/lib-svc/token"
)

type AuthenticateWithIam func(ctx context.Context) error

func (c *client) AuthenticateWithIam(ctx context.Context) error {
	body, err := json.Marshal(map[string]interface{}{
		"client_id":   c.config.Env.RuntimeId,
		"client_name": c.config.Env.SvcId,
	})
	if err != nil {
		return lib_errors.Wrap(err, "Failed generating body for authentication")
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://iam-svc.%s.tomwang.cc/iam-svc/v1/authenticate", c.config.Env.Id), bytes.NewBuffer(body))
	if err != nil {
		return lib_errors.Wrap(err, "Failed creating authentication request")
	}
	req = req.WithContext(ctx)

	lib_log.Info(ctx, "Creating token")
	token, err := c.tokenSvcClient.NewToken(ctx, lib_svc.IamId)
	if err != nil {
		return lib_errors.Wrap(err, fmt.Sprintf("Failed creating new svc token %q", lib_svc.IamId))
	}
	lib_log.Info(ctx, "Created token", lib_log.FmtString("token", token))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token)) // IAM uses this svc token to generate the authenticated token
	req.Header.Set("Content-Type", "application/json")
	req = c.setXLcHeaders(ctx, req)

	lib_http.LogRequestWithBodyAfterApplyingRedactions(req, body, lib_token.JsonRedactionsByPath)

	client := lib_http.NewClient(ctx, true)
	res, err := client.Do(req) // This doesn't use this pkg's client func because this will be called by said func and therefore would cause recursion
	if err != nil {
		return lib_errors.Wrap(err, "Failed issuing authentication request")
	}
	defer lib_http.CloseBody(ctx, res.Body)

	body, err = lib_http.ReadResponseBodyWithLogRedactions(res, client, true, lib_token.JsonRedactionsByPath)
	if err != nil {
		return lib_errors.Wrap(err, "Failed reading response body")
	}

	if res.StatusCode != http.StatusOK {
		return lib_errors.NewCustomf(http.StatusBadGateway, "Failed authentication request: Request resulted in HTTP status code %d", res.StatusCode)
	}

	type authenticationResponse struct {
		Token string `json:"token"`
	}
	var a authenticationResponse
	if err := json.Unmarshal(body, &a); err != nil {
		return lib_errors.Wrap(err, "Failed unmarshalling response into authenticationResponse")
	}
	if lib_context.Test(ctx) {
		c.iamTokenTest = a.Token
	} else {
		c.iamToken = a.Token
	}

	return nil
}
