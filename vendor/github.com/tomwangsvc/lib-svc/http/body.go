package http

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

func ReadRequestBody(r *http.Request, shouldLogBody bool) ([]byte, error) {
	body, err := ReadRequestBodyWithLogRedactions(r, shouldLogBody, nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading request body with redactions")
	}
	return body, nil
}

func ReadRequestBodyWithLogRedactions(r *http.Request, shouldLogBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) ([]byte, error) {
	if r.Body == nil {
		LogRequest(r)
		return nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		LogRequest(r)
		lib_log.Error(r.Context(), "Failed reading body", lib_log.FmtError(err))
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Bad request: malformed body")
	}

	if shouldLogBody {
		if len(bodyLogJsonRedactionsByPath) > 0 {
			LogRequestWithBodyAfterApplyingRedactions(r, body, bodyLogJsonRedactionsByPath)
		} else {
			LogRequestWithBody(r, body)
		}
	} else {
		LogRequest(r)
	}

	return body, nil
}

// ReadResponseBody returns the body as bytes
// Use this function after an http request is made.
// -> For proxies, render the bytes in the response.
// -> To use the response, unmarshal the bytes into a struct using the json.Unmarshal function.
func ReadResponseBody(r *http.Response, httpClientLogger HttpClientLogger, shouldLogBody bool) ([]byte, error) {
	body, err := ReadResponseBodyWithLogRedactions(r, httpClientLogger, shouldLogBody, nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading response body with redactions")
	}
	return body, nil
}

func ReadResponseBodyWithLogRedactions(r *http.Response, httpClientLogger HttpClientLogger, shouldLogBody bool, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) ([]byte, error) {
	if r.Body == nil {
		LogResponseReceived(r, httpClientLogger)
		return nil, nil
	}

	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			LogResponseReceived(r, httpClientLogger)
			return nil, lib_errors.Wrap(err, "Failed setting up gzip reader for response body")
		}
		defer reader.Close()
	default:
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		LogResponseReceived(r, httpClientLogger)
		return nil, lib_errors.Wrap(err, "Failed reading response body")
	}

	if shouldLogBody {
		if len(bodyLogJsonRedactionsByPath) > 0 {
			LogResponseReceivedWithBodyAfterApplyingRedactions(r, httpClientLogger, body, bodyLogJsonRedactionsByPath)
		} else {
			LogResponseReceivedWithBody(r, httpClientLogger, body)
		}
	} else {
		LogResponseReceived(r, httpClientLogger)
	}

	return body, nil
}

// CloseBody is used to close an http.Request or http.Response body
func CloseBody(ctx context.Context, body io.Closer) {
	if err := body.Close(); err != nil {
		lib_log.Error(ctx, "Error closing response body", lib_log.FmtError(err))
	}
}
