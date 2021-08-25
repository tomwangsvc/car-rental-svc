package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_misc "github.com/tomwangsvc/lib-svc/misc"
	lib_reflect "github.com/tomwangsvc/lib-svc/reflect"
	lib_token "github.com/tomwangsvc/lib-svc/token"
)

const (
	RecovererMaxPanics   = 10
	RecovererPauseInMSec = 5000
)

// LogRequest logs a http request at info level
// -> We require all HTTP requests to be logged at INFO so we can correlate business transactions across microservices
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func LogRequest(r *http.Request) {
	fields := []lib_log.Field{
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Method),
		lib_log.FmtAny("HTTPRequestHeader", filterHeader(r.Header)),
		lib_log.FmtString("At", lib_reflect.At(2)),
	}
	lib_log.Info(r.Context(), "HTTP Request", fields...)
}

// LogRequestWithBody logs a http request with a body at info level
// -> We require all HTTP requests to be logged at INFO so we can correlate business transactions across microservices
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func LogRequestWithBody(r *http.Request, body []byte) {
	fields := []lib_log.Field{
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Method),
		lib_log.FmtAny("HTTPRequestHeader", filterHeader(r.Header)),
		lib_log.FmtBytes("HTTPRequestBody", body),
		lib_log.FmtString("At", lib_reflect.At(2)),
	}
	lib_log.Info(r.Context(), "HTTP Request with body", fields...)
}

func LogRequestWithBodyAfterApplyingRedactions(r *http.Request, body []byte, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) {
	fields := []lib_log.Field{
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Method),
		lib_log.FmtAny("HTTPRequestHeader", filterHeader(r.Header)),
		lib_log.FmtString("At", lib_reflect.At(2)),
	}
	body, err := lib_json.ApplyRedactions(body, bodyLogJsonRedactionsByPath)
	if err != nil {
		lib_log.Error(r.Context(), "Failed filtering body, will not log body", lib_log.FmtError(err))
	} else {
		fields = append(fields, lib_log.FmtBytes("HTTPRequestBody", body))
	}
	lib_log.Info(r.Context(), "HTTP Request with body", fields...)
}

// LogResponse logs a http response at an appropriate level for the status code
// -> We require all HTTP responses to be logged at INFO so we can correlate business transactions across microservices
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func LogResponse(ctx context.Context, statusCode int, header http.Header, body interface{}, causedBy error) {
	fields := []lib_log.Field{
		lib_log.FmtInt("HTTPResponseStatusCode", statusCode),
		lib_log.FmtString("HTTPResponseStatusText", http.StatusText(statusCode)),
		lib_log.FmtAny("HTTPResponseHeaders", filterHeader(header)),
		lib_log.FmtAny("HTTPResponseBody", body),
		lib_log.FmtString("At", lib_reflect.At(2)),
	}

	unexpectedFailureInAutomatedProcess := (statusCode == http.StatusNotFound || statusCode == http.StatusConflict) &&
		(lib_context.CloudSchedulerPush(ctx) ||
			(lib_context.CloudTasksPush(ctx) && lib_context.CloudTaskCreatedDate(ctx).After(time.Now().Add(RequestTimeout))) ||
			(lib_context.PubsubPush(ctx) && lib_context.PubsubMessagePublishTime(ctx).After(time.Now().Add(RequestTimeout))))
	if unexpectedFailureInAutomatedProcess {
		fields = append(
			fields,
			lib_log.FmtBool("lib_context.CloudSchedulerPush(ctx)", lib_context.CloudSchedulerPush(ctx)),
			lib_log.FmtBool("lib_context.CloudTasksPush(ctx)", lib_context.CloudTasksPush(ctx)),
			lib_log.FmtTime("lib_context.CloudTaskCreatedDate(ctx)", lib_context.CloudTaskCreatedDate(ctx)),
			lib_log.FmtBool("lib_context.PubsubPush(ctx)", lib_context.PubsubPush(ctx)),
			lib_log.FmtTime("lib_context.PubsubMessagePublishTime(ctx)", lib_context.PubsubMessagePublishTime(ctx)),
		)
	}

	if causedBy != nil {
		fields = append(
			fields,
			lib_log.FmtError(causedBy),
		)
	}

	message := fmt.Sprintf("HTTP Response: %d %s", statusCode, http.StatusText(statusCode))

	if statusCode >= 500 && statusCode <= 599 && statusCode != http.StatusGatewayTimeout {
		lib_log.Error(ctx, message, fields...)

	} else if statusCode == http.StatusBadRequest || statusCode == http.StatusGatewayTimeout || unexpectedFailureInAutomatedProcess {
		lib_log.Warn(ctx, message, fields...)

	} else {
		lib_log.Info(ctx, message, fields...)
	}
}

type HttpClientLogger interface {
	LogString() string
}

// LogResponseReceived logs a http response at info level without the body
// -> We require all HTTP responses to be logged at INFO so we can correlate business transactions across microservices
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func LogResponseReceived(r *http.Response, httpClientLogger HttpClientLogger) {
	fields := []lib_log.Field{
		lib_log.FmtInt("HTTPResponseStatusCode", r.StatusCode),
		lib_log.FmtString("HTTPResponseStatusText", http.StatusText(r.StatusCode)),
		lib_log.FmtAny("HTTPResponseHeaders", filterHeader(r.Header)),
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.Request.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Request.Method),
	}
	if requestErrors := httpClientLogger.LogString(); requestErrors != "" {
		fields = append(fields, lib_log.FmtString("HTTPRequestErrors", requestErrors))
	}
	fields = append(fields, lib_log.FmtString("At", lib_reflect.At(2)))
	lib_log.Info(r.Request.Context(), "HTTP Response", fields...)
}

// LogResponseReceivedWithBody logs a http response at info level with the body
// -> We require all HTTP responses to be logged at INFO so we can correlate business transactions across microservices
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func LogResponseReceivedWithBody(r *http.Response, httpClientLogger HttpClientLogger, body []byte) {
	fields := []lib_log.Field{
		lib_log.FmtInt("HTTPResponseStatusCode", r.StatusCode),
		lib_log.FmtString("HTTPResponseStatusText", http.StatusText(r.StatusCode)),
		lib_log.FmtAny("HTTPResponseHeaders", filterHeader(r.Header)),
		lib_log.FmtBytes("HTTPResponseBody", body),
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.Request.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Request.Method),
	}
	if requestErrors := httpClientLogger.LogString(); requestErrors != "" {
		fields = append(fields, lib_log.FmtString("HTTPRequestErrors", requestErrors))
	}
	fields = append(fields, lib_log.FmtString("At", lib_reflect.At(2)))
	lib_log.Info(r.Request.Context(), "HTTP Response", fields...)
}

func LogResponseReceivedWithBodyAfterApplyingRedactions(r *http.Response, httpClientLogger HttpClientLogger, body []byte, bodyLogJsonRedactionsByPath map[string]lib_json.Redaction) {
	fields := []lib_log.Field{
		lib_log.FmtInt("HTTPResponseStatusCode", r.StatusCode),
		lib_log.FmtString("HTTPResponseStatusText", http.StatusText(r.StatusCode)),
		lib_log.FmtAny("HTTPResponseHeaders", filterHeader(r.Header)),
		lib_log.FmtString("HTTPRequestURL", lib_misc.RedactTokenFromUrl(r.Request.URL.String())),
		lib_log.FmtString("HTTPRequestMethod", r.Request.Method),
	}
	body, err := lib_json.ApplyRedactions(body, bodyLogJsonRedactionsByPath)
	if err != nil {
		lib_log.Error(r.Request.Context(), "Failed filtering body, will not log body", lib_log.FmtError(err))
	} else {
		fields = append(fields, lib_log.FmtBytes("HTTPResponseBody", body))
	}
	if requestErrors := httpClientLogger.LogString(); requestErrors != "" {
		fields = append(fields, lib_log.FmtString("HTTPRequestErrors", requestErrors))
	}
	fields = append(fields, lib_log.FmtString("At", lib_reflect.At(2)))
	lib_log.Info(r.Request.Context(), "HTTP Response", fields...)
}

func filterHeader(header http.Header) http.Header {
	h := make(http.Header)
	for k, v := range header {
		switch k {
		default:
			h[k] = v

		case "Authorization":
			h["Authorization"] = filterAuthorizationHeader(v)
		}
	}
	return h
}

func filterAuthorizationHeader(values []string) []string {
	var s []string
	for _, v := range values {
		s = append(s, lib_token.Redact(v))
	}
	return s
}
