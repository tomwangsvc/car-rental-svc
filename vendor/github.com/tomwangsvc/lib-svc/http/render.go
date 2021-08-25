package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_misc "github.com/tomwangsvc/lib-svc/misc"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
)

// RenderError adds an custom error to the response
func RenderError(ctx context.Context, w http.ResponseWriter, err error) {
	if cerr, ok := err.(lib_errors.Custom); ok {
		RenderCustomError(ctx, w, cerr)
	} else {
		header := SetHeaderWithXLc(ctx, w, nil, nil, nil)
		w.WriteHeader(http.StatusInternalServerError)
		LogResponse(ctx, http.StatusInternalServerError, header, nil, err)
	}
}

// RenderCustomError adds a custom error to the response
func RenderCustomError(ctx context.Context, w http.ResponseWriter, cerr lib_errors.Custom) {
	if body, warnings := cerr.Render(); body != nil {
		header := SetHeaderWithXLc(ctx, w, nil, http.Header{
			"Content-Type": []string{
				"application/json",
			},
		}, nil)
		if len(warnings) > 0 {
			lib_log.Error(ctx, "Warning returned when rendering custom error, most likely due to improper usage of custom errors", lib_log.FmtAny("warnings", warnings))
		}
		writeErrorHeader(ctx, w, cerr)
		if _, err := w.Write(body); err != nil {
			lib_log.Error(ctx, "Error writing response, HTTP response status code already set", lib_log.FmtError(err))
		}
		LogResponse(ctx, cerr.Code, header, string(body), cerr)
		return
	}
	// ELSE -> FOR ALL OTHER ERROR CODES REASON IS NOT EXPOSED
	header := SetHeaderWithXLc(ctx, w, nil, nil, nil)
	writeErrorHeader(ctx, w, cerr)
	LogResponse(ctx, cerr.Code, header, nil, cerr)
}

// RenderCustomErrorWithBody adds a custom error's code and a body to the response
func RenderCustomErrorWithBody(ctx context.Context, w http.ResponseWriter, cerr lib_errors.Custom, body interface{}) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	writeErrorHeader(ctx, w, cerr)
	LogResponse(ctx, cerr.Code, header, body, cerr)
	if err := lib_misc.NewJsonEncoder(w).Encode(body); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

func writeErrorHeader(ctx context.Context, w http.ResponseWriter, cerr lib_errors.Custom) {
	if (cerr.Code < 400 || cerr.Code >= 500) && cerr.Code != http.StatusBadGateway && cerr.Code != http.StatusGatewayTimeout {
		lib_log.Warn(ctx, "THIS SHOULD BE CHANGED, we should only use 400's, 502, and 504 for custom errors", lib_log.FmtInt("cerr.Code", cerr.Code))
	}
	if cerr.Code >= 200 && cerr.Code < 600 {
		w.WriteHeader(cerr.Code)
	} else {
		lib_log.Warn(ctx, "HTTP status code not supported, will use 500 in order to avoid a panic in the middleware", lib_log.FmtInt("cerr.Code", cerr.Code))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RenderAuthError adds an authorization error to the response
func RenderAuthError(ctx context.Context, w http.ResponseWriter, err error) {
	cerr, ok := err.(lib_errors.Custom)
	if !ok || !(cerr.Code == http.StatusForbidden || cerr.Code == http.StatusUnprocessableEntity || cerr.Code == http.StatusBadGateway || cerr.Code == http.StatusBadRequest) {
		lib_log.Info(ctx, "Auth error, should not render directly, will render Unauthorized", lib_log.FmtError(err))
		cerr = lib_errors.NewCustom(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	} else {
		lib_log.Info(ctx, "Auth error, will render directly", lib_log.FmtError(err))
	}
	RenderCustomError(ctx, w, cerr)
}

// RenderIamAuthError adds an iam authorization error to the response
func RenderIamAuthError(ctx context.Context, w http.ResponseWriter, err error) {
	cerr, ok := err.(lib_errors.Custom)
	if !ok || !(cerr.Code == http.StatusForbidden || cerr.Code == http.StatusUnprocessableEntity || cerr.Code == http.StatusBadGateway || cerr.Code == http.StatusBadRequest) {
		lib_log.Info(ctx, "IAM auth error, should not render directly, will render Unauthorized", lib_log.FmtError(err))
		cerr = lib_errors.NewCustom(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	} else {
		lib_log.Info(ctx, "IAM auth error, will render directly", lib_log.FmtError(err))
	}
	RenderCustomError(ctx, w, cerr)
}

// RenderOkWithNoBody adds a Ok status to the response
func RenderOkWithNoBody(ctx context.Context, w http.ResponseWriter) {
	header := SetHeaderWithXLc(ctx, w, nil, nil, nil)
	w.WriteHeader(http.StatusOK)
	LogResponse(ctx, http.StatusOK, header, nil, nil)
}

// RenderNoContent adds a NoContent status to the response
func RenderNoContent(ctx context.Context, w http.ResponseWriter) {
	header := SetHeaderWithXLc(ctx, w, nil, nil, nil)
	w.WriteHeader(http.StatusNoContent)
	LogResponse(ctx, http.StatusNoContent, header, nil, nil)
}

// RenderNoContentWithPagination adds a NoContent status to the response with pagination
func RenderNoContentWithPagination(ctx context.Context, w http.ResponseWriter, pagination lib_pagination.Pagination) {
	header := SetHeaderWithXLc(ctx, w, &pagination, nil, nil)
	w.WriteHeader(http.StatusNoContent)
	LogResponse(ctx, http.StatusNoContent, header, nil, nil)
}

// RenderNoContentWithLocation adds a NoContent status to the response with a location header
func RenderNoContentWithLocation(ctx context.Context, w http.ResponseWriter, location string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusNoContent)
	LogResponse(ctx, http.StatusNoContent, header, nil, nil)
}

// RenderNoContentWithTestMetadata adds a NoContent status to the response with a test metadata header containing the content passed
func RenderNoContentWithTestMetadata(ctx context.Context, w http.ResponseWriter, content string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		HeaderKeyXLcSvcTestMetadata: []string{
			content,
		},
	}, nil)
	w.WriteHeader(http.StatusNoContent)
	LogResponse(ctx, http.StatusNoContent, header, nil, nil)
}

// RenderCreated adds a Created status to the response with a location header
func RenderCreated(ctx context.Context, w http.ResponseWriter, location string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusCreated)
	LogResponse(ctx, http.StatusCreated, header, nil, nil)
}

// RenderCreatedWithBody adds a Created status to the response with a location header and body
func RenderCreatedWithBody(ctx context.Context, w http.ResponseWriter, location string, res interface{}) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusCreated)
	LogResponse(ctx, http.StatusCreated, header, res, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(res); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderCreatedWithJsonBytes adds a Created status to the response with a location header and bytes body
func RenderCreatedWithJsonBytes(ctx context.Context, w http.ResponseWriter, location string, bytes []byte) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusCreated)
	LogResponse(ctx, http.StatusCreated, header, bytes, nil)
	lenBytes, err := w.Write(bytes)
	lib_log.Info(ctx, fmt.Sprintf("Sending HTTP Response as bytes with length %d", lenBytes), lib_log.FmtInt("HTTPResponseStatusCode", http.StatusOK), lib_log.FmtString("HTTPResponseStatusText", http.StatusText(http.StatusOK)), lib_log.FmtAny("HTTPResponseHeaders", header))
	if err != nil {
		lib_log.Error(ctx, "Error writing response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderCreatedWithBodyNoLocation adds a Created status to the response with a body but no location header
func RenderCreatedWithBodyNoLocation(ctx context.Context, w http.ResponseWriter, res interface{}) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	w.WriteHeader(http.StatusCreated)
	LogResponse(ctx, http.StatusCreated, header, res, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(res); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderResponse adds a Json encoded body to the response
// TODO: deprecate
func RenderResponse(ctx context.Context, w http.ResponseWriter, res interface{}) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	LogResponse(ctx, http.StatusOK, header, res, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(res); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderResponseWithPagination adds a Json encoded body to the response
// TODO: deprecate
func RenderResponseWithPagination(ctx context.Context, w http.ResponseWriter, res interface{}, pagination lib_pagination.Pagination) {
	header := SetHeaderWithXLc(ctx, w, &pagination, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	LogResponse(ctx, http.StatusOK, header, res, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(res); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderResponseNoBodyLog adds a Json encoded body to the response
// TODO: deprecate
func RenderResponseNoBodyLog(ctx context.Context, w http.ResponseWriter, res interface{}) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	LogResponse(ctx, http.StatusOK, header, "DO NOT LOG BODY", nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(res); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderJsonBytes writes the bytes to the response writer
func RenderJsonBytes(ctx context.Context, w http.ResponseWriter, bytes []byte) {
	RenderBytes(ctx, w, bytes, "application/json")
}

// RenderBytesWithFileNameAsAttachment writes the bytes to the response writer
func RenderBytesWithFileNameAsAttachment(ctx context.Context, w http.ResponseWriter, bytes []byte, contentType, fileName string) {
	renderBytesWithHeader(ctx, w, bytes, nil, nil, http.Header{
		"Content-Type": []string{
			contentType,
		},
		"Content-Disposition": []string{
			fmt.Sprintf("attachment; filename=%s", fileName),
		},
	})
}

// RenderBytes writes the bytes to the response writer
func RenderBytes(ctx context.Context, w http.ResponseWriter, bytes []byte, contentType string) {
	renderBytesWithHeader(ctx, w, bytes, nil, nil, http.Header{
		"Content-Type": []string{
			contentType,
		},
	})
}

// RenderJsonBytesWithPagination writes the bytes to the response writer with pagination
func RenderJsonBytesWithPagination(ctx context.Context, w http.ResponseWriter, bytes []byte, pagination lib_pagination.Pagination) {
	RenderBytesWithPagination(ctx, w, bytes, pagination, "application/json")
}

// RenderJsonBytesWithPaginationWithEtagsForObjects writes the bytes to the response writer with pagination with etags for individual objecgts
func RenderJsonBytesWithPaginationWithEtagsForObjects(ctx context.Context, w http.ResponseWriter, bytes []byte, pagination lib_pagination.Pagination, etagsForObjects []string) {
	RenderBytesWithPaginationWithEtagsForObjects(ctx, w, bytes, pagination, "application/json", etagsForObjects)
}

// RenderBytesWithPagination writes the bytes to the response writer with pagination
func RenderBytesWithPagination(ctx context.Context, w http.ResponseWriter, bytes []byte, pagination lib_pagination.Pagination, contentType string) {
	renderBytesWithHeader(ctx, w, bytes, nil, &pagination, http.Header{
		"Content-Type": []string{
			contentType,
		},
	})
}

// RenderBytesWithPaginationWithEtagsForObjects writes the bytes to the response writer with pagination with etags for individual objects
func RenderBytesWithPaginationWithEtagsForObjects(ctx context.Context, w http.ResponseWriter, bytes []byte, pagination lib_pagination.Pagination, contentType string, etagsForObjects []string) {
	renderBytesWithHeader(ctx, w, bytes, etagsForObjects, &pagination, http.Header{
		"Content-Type": []string{
			contentType,
		},
	})
}

func renderBytesWithHeader(ctx context.Context, w http.ResponseWriter, bytes []byte, xLcEtagsForObjects []string, pagination *lib_pagination.Pagination, header http.Header) {
	etag := lib_misc.GenerateEtag(bytes)
	header["ETag"] = []string{etag}
	header = SetHeaderWithXLc(ctx, w, pagination, header, xLcEtagsForObjects)
	if lib_context.HttpRequestHeaderIfNoneMatch(ctx) == etag {
		w.WriteHeader(http.StatusNotModified)
		LogResponse(ctx, http.StatusNotModified, header, nil, nil)
	} else {
		lenBytes, err := w.Write(bytes)
		// We don't call 'LogResponse' in this method because it is not sensible to log bytes, there will be alot of them
		lib_log.Info(ctx, fmt.Sprintf("Sending HTTP Response as bytes with length %d", lenBytes), lib_log.FmtInt("HTTPResponseStatusCode", http.StatusOK), lib_log.FmtString("HTTPResponseStatusText", http.StatusText(http.StatusOK)), lib_log.FmtAny("HTTPResponseHeaders", header))
		if err != nil {
			lib_log.Error(ctx, "Error writing response, HTTP response status code already set", lib_log.FmtError(err))
		}
	}
}

type accepted struct {
	MaxRetries int `json:"max_retries"`
	WaitMSec   int `json:"wait_msec"`
}

var (
	a = accepted{
		MaxRetries: 5,
		WaitMSec:   2500,
	}
)

// RenderAccepted for rendering an http.StatusAccepted with a location headerdy
func RenderAccepted(ctx context.Context, w http.ResponseWriter, location string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusAccepted)
	LogResponse(ctx, http.StatusAccepted, header, nil, nil)
}

// RenderAcceptedWithRetry for rendering an http.StatusAccepted with a location header and a retry body
func RenderAcceptedWithRetry(ctx context.Context, w http.ResponseWriter, location string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
		"Location": []string{
			location,
		},
	}, nil)
	w.WriteHeader(http.StatusAccepted)
	LogResponse(ctx, http.StatusAccepted, header, a, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(a); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderAcceptedNoLocation for rendering an http.StatusAccepted with no location header
func RenderAcceptedNoLocation(ctx context.Context, w http.ResponseWriter) {
	header := SetHeaderWithXLc(ctx, w, nil, nil, nil)
	w.WriteHeader(http.StatusAccepted)
	LogResponse(ctx, http.StatusAccepted, header, nil, nil)
}

// RenderAcceptedNoLocationWithRetry for rendering an http.StatusAccepted with no location header and a retry body
func RenderAcceptedNoLocationWithRetry(ctx context.Context, w http.ResponseWriter) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		"Content-Type": []string{
			"application/json",
		},
	}, nil)
	w.WriteHeader(http.StatusAccepted)
	LogResponse(ctx, http.StatusAccepted, header, a, nil)
	if err := lib_misc.NewJsonEncoder(w).Encode(a); err != nil {
		lib_log.Error(ctx, "Error encoding response, HTTP response status code already set", lib_log.FmtError(err))
	}
}

// RenderAcceptedWithTestMetadata adds a Accepted status to the response with a test metadata header containing the content passed
func RenderAcceptedWithTestMetadata(ctx context.Context, w http.ResponseWriter, content string) {
	header := SetHeaderWithXLc(ctx, w, nil, http.Header{
		HeaderKeyXLcSvcTestMetadata: []string{
			content,
		},
	}, nil)
	w.WriteHeader(http.StatusAccepted)
	LogResponse(ctx, http.StatusAccepted, header, nil, nil)
}

//revive:disable:cyclomatic
func SetHeaderWithXLc(ctx context.Context, w http.ResponseWriter, pagination *lib_pagination.Pagination, header http.Header, xLcEtagsForObjects []string) http.Header {
	if header == nil {
		header = make(http.Header)
	}

	if len(xLcEtagsForObjects) > 0 {
		header[HeaderKeyXLcETagsForObjects] = xLcEtagsForObjects
	}

	header[HeaderKeyXLcCorrelationId] = []string{lib_context.CorrelationId(ctx)}
	if lib_context.IntegrationTest(ctx) {
		header[HeaderKeyXLcSvcIntegrationTest] = []string{HeaderValueXLcSvcIntegrationTest}
		header[HeaderKeyXLcSvcIntegrationTestPubsubAutoAckDisable] = []string{strconv.FormatBool(lib_context.IntegrationTestPubsubAutoAckDisable(ctx))}
	}
	if lib_context.Test(ctx) {
		header[HeaderKeyXLcSvcTest] = []string{HeaderValueXLcSvcTest}
	}

	if pagination != nil {
		// Headers must only contain ASCII
		// -> string(pagination.Offset) sometimes causes nginx error "upstream sent invalid header while reading response header from upstream"
		header[HeaderKeyXLcPaginationLimit] = []string{strconv.Itoa(pagination.Limit)}
		header[HeaderKeyXLcPaginationOffset] = []string{strconv.Itoa(pagination.Offset)}
		header[HeaderKeyXLcPaginationCursor] = []string{pagination.Cursor}
		header[HeaderKeyXLcPaginationOrder] = []string{pagination.Order}
		if pagination.ReadTimestamp != nil {
			header[HeaderKeyXLcPaginationReadTimestamp] = []string{(*pagination.ReadTimestamp).Format(time.RFC3339)}
		}
		if pagination.Total != nil {
			header[HeaderKeyXLcPaginationTotal] = []string{strconv.Itoa(int(*pagination.Total))}
		}
	}

	for k, v := range header {
		if k == HeaderKeyXLcSvcTestMetadata && !lib_context.Test(ctx) {
			continue
		}

		for i, vv := range v {
			if i == 0 {
				w.Header().Set(k, vv)
			} else {
				w.Header().Add(k, vv)
			}
		}
	}

	for k, v := range lib_context.XLcLocationHeaders(ctx) {
		for i, vv := range v {
			if i == 0 {
				w.Header().Set(k, vv)
				header.Set(k, vv)
			} else {
				w.Header().Add(k, vv)
				header.Add(k, vv)
			}
		}
	}

	return header
	//revive:enable:cyclomatic
}
