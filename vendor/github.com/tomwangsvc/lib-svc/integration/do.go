package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_regexp "github.com/tomwangsvc/lib-svc/regexp"
)

func (c client) DoRequestNotUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody))

	res, resBody, err = c.doRequestNotUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, nil, true)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request not using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c client) DoRequestNotUsingIamAuthorizationNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody))

	res, resBody, err = c.doRequestNotUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, nil, false)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request not using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c client) DoRequestNotUsingIamAuthorizationWithBodyLogRedactions(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody), lib_log.FmtInt("len(bodyLogJsonRedactionsByPath)", len(bodyLogJsonRedactionsByPath)))

	res, resBody, err = c.doRequestNotUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, true)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request not using iam authorizations")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c client) DoRequestNotUsingIamAuthorizationWithBodyLogRedactionsNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody), lib_log.FmtInt("len(bodyLogJsonRedactionsByPath)", len(bodyLogJsonRedactionsByPath)))

	res, resBody, err = c.doRequestNotUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, false)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request not using iam authorizations")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c client) doRequestNotUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction, allowRetries bool) (res *http.Response, resBody []byte, err error) {
	reqBody, err := extractRequestBody(req)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed extracting request body")
		return
	}
	res, resBody, err = c.do(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request")
		return
	}
	return
}

func (c *client) DoRequestUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody))

	res, resBody, err = c.doRequestUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, nil, true)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c *client) DoRequestUsingIamAuthorizationNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody))

	res, resBody, err = c.doRequestUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, nil, false)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c *client) DoRequestUsingIamAuthorizationWithBodyLogRedactions(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody), lib_log.FmtInt("len(bodyLogJsonRedactionsByPath)", len(bodyLogJsonRedactionsByPath)))

	res, resBody, err = c.doRequestUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, true)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c *client) DoRequestUsingIamAuthorizationWithBodyLogRedactionsNoRetries(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) (res *http.Response, resBody []byte, err error) {
	lib_log.Info(ctx, "Doing", lib_log.FmtBool("logRequestBody", logRequestBody), lib_log.FmtBool("logResponseBody", logResponseBody), lib_log.FmtInt("len(bodyLogJsonRedactionsByPath)", len(bodyLogJsonRedactionsByPath)))

	res, resBody, err = c.doRequestUsingIamAuthorization(ctx, req, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, false)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request using iam authorization")
		return
	}

	lib_log.Info(ctx, "Done")
	return
}

func (c *client) doRequestUsingIamAuthorization(ctx context.Context, req *http.Request, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction, allowRetries bool) (res *http.Response, resBody []byte, err error) {
	reqBody, err := extractRequestBody(req)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed extracting request body")
		return
	}
	req = c.setIamHeaders(ctx, req)

	iamToken := c.iamToken
	if lib_context.Test(ctx) {
		iamToken = c.iamTokenTest
	}
	if iamToken == "" {
		if res, resBody, err = c.authenticateWithIamAndDoRequest(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries); err != nil {
			err = lib_errors.Wrap(err, "Failed authenticating with iam and doing request")
			return
		}
		lib_log.Info(ctx, "Done")
		return
	}

	iamTokenUnverifiedClaims, err := c.tokenIamClient.ExtractUnverifiedClaims(ctx, iamToken)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed extracting unverified claims")
		return
	}
	if expiresAt := time.Unix(iamTokenUnverifiedClaims.ExpiresAt, 0); expiresAt.Before(time.Now()) {
		if res, resBody, err = c.authenticateWithIamAndDoRequest(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries); err != nil {
			err = lib_errors.Wrap(err, "Failed authenticating with iam and doing request")
			return
		}
		lib_log.Info(ctx, "Done")
		return
	}

	res, resBody, err = c.do(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request")
		return
	}
	if res.StatusCode == http.StatusUnauthorized {
		if res, resBody, err = c.authenticateWithIamAndDoRequest(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries); err != nil {
			err = lib_errors.Wrap(err, "Failed authenticating with iam and doing request")
			return
		}
	}

	return
}

func extractRequestBody(req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, lib_errors.New("Request cannot be nil")
	}

	var reqBody []byte
	if req.Body != nil {
		var err error
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed reading request body to store for retries")
		}
		req.Body.Close()
	}

	return reqBody, nil
}

func (c *client) authenticateWithIamAndDoRequest(ctx context.Context, req *http.Request, reqBody []byte, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction, allowRetries bool) (res *http.Response, resBody []byte, err error) {
	if err = c.AuthenticateWithIam(ctx); err != nil {
		if lib_errors.IsCustomWithCode(err, http.StatusUnauthorized) {
			lib_log.Error(ctx, "Failed authenticating", lib_log.FmtError(err))
			err = lib_errors.New("Failed authenticating")
			return
		}
		err = lib_errors.Wrap(err, "Failed authenticating")
		return
	}

	req = c.setIamHeaders(ctx, req)
	res, resBody, err = c.do(ctx, req, reqBody, logRequestBody, logResponseBody, bodyLogJsonRedactionsByPath, allowRetries)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request")
		return
	}

	return
}

func (c client) setIamHeaders(ctx context.Context, req *http.Request) *http.Request {
	if lib_context.Test(ctx) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.iamTokenTest))
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.iamToken))
	}
	return req
}

func (c client) do(ctx context.Context, req *http.Request, reqBody []byte, logRequestBody, logResponseBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction, allowRetries bool) (res *http.Response, resBody []byte, err error) {
	req = c.setXLcHeaders(ctx, req)
	if logRequestBody {
		lib_http.LogRequestWithBodyAfterApplyingRedactions(req, reqBody, bodyLogJsonRedactionsByPath)
	} else {
		lib_http.LogRequest(req)
	}
	if reqBody != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	client := lib_http.NewClient(ctx, allowRetries)
	res, err = client.Do(req)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed doing request")
		return
	}
	defer lib_http.CloseBody(ctx, res.Body)

	resBody, err = lib_http.ReadResponseBodyWithLogRedactions(res, client, logResponseBody, bodyLogJsonRedactionsByPath)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed reading response body")
		return
	}

	return
}

func (c client) setXLcHeaders(ctx context.Context, req *http.Request) *http.Request {
	if isToInternalService(req.URL.String()) {
		req.Header.Set(lib_http.HeaderKeyXLcCorrelationId, lib_context.CorrelationId(ctx))
		req.Header.Set(lib_http.HeaderKeyXLcCaller, c.config.Env.SvcId)
		if lib_context.IntegrationTest(ctx) {
			req.Header.Set(lib_http.HeaderKeyXLcSvcIntegrationTest, lib_http.HeaderValueXLcSvcIntegrationTest)
		}
		if lib_context.Test(ctx) {
			req.Header.Set(lib_http.HeaderKeyXLcSvcTest, lib_http.HeaderValueXLcSvcTest)
		}
	}
	return req
}

var internalUrlServiceRegExpCompile *regexp.Regexp

func init() {
	internalUrlServiceRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlServiceRegExp)
}

func isToInternalService(url string) bool {
	if internalUrlServiceRegExpCompile.MatchString(strings.ToLower(url)) {
		return true
	}
	return false
}
