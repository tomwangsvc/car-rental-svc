package http

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	lib_appnative "github.com/tomwangsvc/lib-svc/appnative"
	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_regexp "github.com/tomwangsvc/lib-svc/regexp"
	lib_svc "github.com/tomwangsvc/lib-svc/svc"
)

const (
	RequestTimeout = time.Second * 30
)

var internalUrlAppDevRegExpCompile *regexp.Regexp
var internalUrlAppPrdRegExpCompile *regexp.Regexp
var internalUrlAppStgRegExpCompile *regexp.Regexp
var internalUrlAppUatRegExpCompile *regexp.Regexp

// TODO: Delete these direct firebase domains once goboxer.com domain is live
var internalUrlBoxerAppDevRegExpCompile *regexp.Regexp
var internalUrlBoxerAppPrdRegExpCompile *regexp.Regexp
var internalUrlBoxerAppStgRegExpCompile *regexp.Regexp
var internalUrlBoxerAppUatRegExpCompile *regexp.Regexp

func init() {
	internalUrlAppDevRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlAppDevRegExp)
	internalUrlAppPrdRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlAppPrdRegExp)
	internalUrlAppStgRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlAppStgRegExp)
	internalUrlAppUatRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlAppUatRegExp)
	internalUrlBoxerAppDevRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlBoxerAppDevRegExp)
	internalUrlBoxerAppPrdRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlBoxerAppPrdRegExp)
	internalUrlBoxerAppStgRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlBoxerAppStgRegExp)
	internalUrlBoxerAppUatRegExpCompile = regexp.MustCompile(lib_regexp.InternalUrlBoxerAppUatRegExp)
}

func GeneralMiddleware(r *chi.Mux, env string, maintenanceMode bool, logExcludePaths []string) {
	GeneralMiddlewareWithTimeout(r, env, maintenanceMode, logExcludePaths, RequestTimeout)
}

// Prefer 'GeneralMiddleware' over this func. This func is used for long running requests (minutes) e.g. backup and ML workloads
func GeneralMiddlewareWithTimeout(r *chi.Mux, env string, maintenanceMode bool, logExcludePaths []string, timeout time.Duration) {
	// Middleware order is important
	r.Use(middleware.Timeout(timeout))
	GeneralMiddlewareWithNoTimeout(r, env, maintenanceMode, logExcludePaths)
}

// Prefer 'GeneralMiddleware' over this func. This func is used for routers that need different length requests (minutes) per route e.g. schedule-svc
func GeneralMiddlewareWithNoTimeout(r *chi.Mux, env string, maintenanceMode bool, logExcludePaths []string) {
	// Middleware order is important
	r.Use(checkMaintenanceMode(maintenanceMode))
	r.Use(RequireHttpsInGcp)
	r.Use(middleware.RequestID)
	r.Use(PopulateContext(env))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(NewCors().Handler)
	r.Use(middleware.Compress(5, "application/json"))
	r.Use(NoCacheResponse)
	r.Use(NewRequestLogger(logExcludePaths).InfoHandler)
}

// NewCors returns customized chi cors configuration for chi cors middleware
func NewCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Authorization",
			"Content-Disposition",
			"Content-Type",
			"If-Match",
			"If-None-Match",
			"If-Unmodified-Since",
			"X-CsrF-Token",
			HeaderKeyXLcSvcTest,
			HeaderKeyXLcCorrelationId,
		},
		ExposedHeaders: []string{
			"Content-Disposition",
			"Content-Type",
			"ETag",
			"Link",
			"Location",
			HeaderKeyXLcCorrelationId,
			HeaderKeyXLcETagsForObjects,
			HeaderKeyXLcPaginationCursor,
			HeaderKeyXLcPaginationLimit,
			HeaderKeyXLcPaginationOffset,
			HeaderKeyXLcPaginationReadTimestamp,
			HeaderKeyXLcPaginationTotal,
			HeaderKeyXLcLocationCity,
			HeaderKeyXLcLocationCityLatLng,
			HeaderKeyXLcLocationCountry,
			HeaderKeyXLcLocationIp,
			HeaderKeyXLcLocationRegion,
			HeaderKeyXLcLocationRegionSubdivision,
			HeaderKeyXLcSvcTest,
			HeaderKeyXLcSvcTestMetadata,
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

// RequireHttpsInGcp require the use of HTTPS
// -> E.g. Cloud Run is similar but for GAE see https://cloud.google.com/appengine/docs/flexible/go/how-requests-are-handled
func RequireHttpsInGcp(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if v, ok := r.Header["X-Forwarded-Proto"]; ok {
			for i, vv := range v {
				if vv == "http" {
					// For performance reasons we only log if we intend to interrupt the request
					lib_log.Info(r.Context(), "RequireHttpsInGcp", lib_log.FmtString(fmt.Sprintf("X-Forwarded-Proto[%d]", i), vv), lib_log.FmtString("X-Forwarded-Proto", strings.Join(v, ",")))

					// We don't redirect because it's tricky in practice and we don't need to, this is an API
					//http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func PopulateContext(env string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := lib_context.WithCorrelationId(r.Context(), correlationId(r.Header))

			ctx = lib_context.WithHttpRequestHeaderIfNoneMatch(ctx, r.Header.Get("If-None-Match"))

			integrationTest := integrationTest(ctx, r.Header)
			ctx = lib_context.WithIntegrationTest(ctx, integrationTest)

			ctx = lib_context.WithIntegrationTestPubsubAutoAckDisable(ctx, integrationTestPubsubAutoAckDisable(ctx, r.Header, integrationTest))

			ctx = lib_context.WithTest(ctx, test(ctx, env, r.Header, integrationTest))

			ctx = lib_context.WithTestMetadata(ctx, testMetadata(r.Header))

			ctx = lib_context.WithLcCaller(ctx, lcCaller(r.Header, env))

			ctx = lib_context.WithXLcLocationHeaders(ctx, xLcLocationHeaders(r.Header))

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func correlationId(header http.Header) string {
	correlationId := uuid.New().String()
	if v, ok := header[HeaderKeyXLcCorrelationId]; ok {
		correlationId = fmt.Sprintf("%s,%s", correlationId, strings.Join(v, ","))
	}
	return correlationId
}

func integrationTest(ctx context.Context, header http.Header) bool {
	_, postmanTokenHeaderExists := header["Postman-Token"]
	return postmanTokenHeaderExists ||
		strings.Contains(header.Get("User-Agent"), "PostmanRuntime") ||
		strings.Contains(header.Get(HeaderKeyXLcSvcIntegrationTest), HeaderValueXLcSvcIntegrationTest) ||
		lib_context.IntegrationTest(ctx)
}

func integrationTestPubsubAutoAckDisable(ctx context.Context, header http.Header, integrationTest bool) bool {
	integrationTestPubsubAutoAckDisable, err := strconv.ParseBool(header.Get(HeaderKeyXLcSvcIntegrationTestPubsubAutoAckDisable))
	if err != nil {
		return false
	}
	return (integrationTest && integrationTestPubsubAutoAckDisable) ||
		lib_context.IntegrationTestPubsubAutoAckDisable(ctx)
}

// Adds test mode to the context if these conditions are met:
// 1. If running an integration test
// 2. If passed the test HTTP header
// 3. If using the test app in production i.e. the test app URL is passed in the 'Origin' HTTP header for a production environment
// 4. If the context is already in test mode
func test(ctx context.Context, env string, header http.Header, integrationTest bool) bool {
	return integrationTest ||
		strings.Contains(header.Get(HeaderKeyXLcSvcTest), HeaderValueXLcSvcTest) ||
		isTestAppInProduction(header, env) ||
		lib_context.Test(ctx)
}

func isTestAppInProduction(header http.Header, env string) bool {
	switch env {
	case lib_env.Prd:
		return strings.Contains(strings.ToLower(header.Get("Origin")), "test.tomwang.com")
	default:
		return false
	}
}

func testMetadata(header http.Header) string {
	return header.Get(HeaderKeyXLcSvcTestMetadata)
}

func lcCaller(header http.Header, env string) bool {
	var lcCaller bool
	if v, ok := header["Origin"]; ok {
		for _, vv := range v {
			if isSupportedLcCallerHeader(vv, env) {
				lcCaller = true
				break
			}
		}
	}
	if !lcCaller {
		if v, ok := header[HeaderKeyXLcCaller]; ok {
			for _, vv := range v {
				if isSupportedLcCallerHeader(vv, env) {
					lcCaller = true
					break
				}
			}
		}
	}
	return lcCaller
}

func isSupportedLcCallerHeader(headerValue, env string) bool {
	if lib_svc.IsRecognizedId(headerValue) {
		return true
	}

	if lib_appnative.IsRecognizedId(headerValue, env) {
		return true
	}

	switch env {
	case lib_env.Dev:
		return internalUrlAppDevRegExpCompile.MatchString(strings.TrimSpace(headerValue)) || internalUrlBoxerAppDevRegExpCompile.MatchString(strings.TrimSpace(headerValue))
	case lib_env.Prd:
		// Go RegExp does not support lookaround  `(?!` so we go the long way round
		return !strings.Contains(headerValue, "."+lib_env.Dev+".") && !strings.Contains(headerValue, "."+lib_env.Uat+".") && internalUrlAppPrdRegExpCompile.MatchString(strings.TrimSpace(headerValue)) || internalUrlBoxerAppPrdRegExpCompile.MatchString(strings.TrimSpace(headerValue))
	case lib_env.Stg:
		return internalUrlAppStgRegExpCompile.MatchString(strings.TrimSpace(headerValue)) || internalUrlBoxerAppStgRegExpCompile.MatchString(strings.TrimSpace(headerValue))
	case lib_env.Uat:
		return internalUrlAppUatRegExpCompile.MatchString(strings.TrimSpace(headerValue)) || internalUrlBoxerAppUatRegExpCompile.MatchString(strings.TrimSpace(headerValue))
	default:
		return false
	}
}

func xLcLocationHeaders(header http.Header) http.Header {
	xLcLocationHeaders := make(http.Header)
	if v, ok := header["X-Lc-Location-City"]; ok {
		xLcLocationHeaders[HeaderKeyXLcLocationCity] = v
	}
	if v, ok := header["X-Lc-Location-City-Lat-Lng"]; ok {
		xLcLocationHeaders[HeaderKeyXLcLocationCityLatLng] = v
	}
	if v, ok := header["X-Lc-Location-Region"]; ok {
		xLcLocationHeaders[HeaderKeyXLcLocationCountry] = v
		xLcLocationHeaders[HeaderKeyXLcLocationRegion] = v
	}
	if v, ok := header["X-Lc-Location-Region-Subdivision"]; ok {
		xLcLocationHeaders[HeaderKeyXLcLocationRegionSubdivision] = v
	}
	if v, ok := header["X-Forwarded-For"]; ok {
		xLcLocationHeaders[HeaderKeyXLcLocationIp] = v
	}
	return xLcLocationHeaders
}

// NewRequestLogger creates a new RequestLogger
func NewRequestLogger(excludePaths []string) *RequestLogger {
	return &RequestLogger{
		excludePaths,
	}
}

// RequestLogger logs a http request at info level
type RequestLogger struct {
	excludePaths []string
}

// InfoHandler is middleware to log a http request at info level
func (l *RequestLogger) InfoHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		shouldLog := true
		for _, v := range l.excludePaths {
			if r.URL.String() == v {
				shouldLog = false
				break
			}
		}

		if shouldLog {
			LogRequest(r)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NoCacheResponse sets a number of HTTP headers to prevent a router (or subrouter) caching requests.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
// -> Although other directives may be set, this alone is the only directive you need in preventing cached responses on MODERN BROWSERS.
func NoCacheResponse(next http.Handler) http.Handler {
	var noCacheHeaders = map[string]string{
		"Expires":         time.Unix(0, 0).Format(time.RFC1123),
		"Cache-Control":   "no-cache, no-store, must-revalidate, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type ContextKey int

func AddUrlParamToContext(urlParamKey string, key ContextKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, chi.URLParam(r, urlParamKey))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthorizeLcCaller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !lib_context.LcCaller(ctx) {
			RenderError(ctx, w, lib_errors.NewCustom(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		lib_log.Info(ctx, "Authorized")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkMaintenanceMode(maintenanceMode bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if maintenanceMode {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
