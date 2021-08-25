package http

import (
	"net/http"
	"strings"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

// ProxyHeaders copies headers striping any that should not accompany a request that will be proxied which are typically those added by CDNs
// Note that it is difficult to backlist all headers that may cause problems and so it is normally safer to include only the required headers in a proxied request
func ProxyHeaders(incomingHeaders http.Header) http.Header {
	proxiedHeaders := make(http.Header)

	for k, v := range incomingHeaders {
		if strings.HasPrefix(strings.ToLower(k), "cf-") {
			continue
		}
		if strings.ToLower(k) == "x-cloud-trace-context" {
			continue
		}
		if strings.ToLower(k) == "__cfduid" {
			continue
		}
		if strings.ToLower(k) == "cdn-loop" {
			continue
		}

		for _, vv := range v {
			if proxiedHeaders.Get(k) == "" {
				proxiedHeaders.Set(k, vv)
			} else {
				proxiedHeaders.Add(k, vv)
			}
		}
	}

	return proxiedHeaders
}

func TokenFromAuthorizationHeader(header http.Header) (string, error) {
	a := header.Get("Authorization")
	if a == "" {
		return "", lib_errors.New("Missing Authorization header")
	}
	w := strings.Split(a, " ")
	if len(w) != 2 || w[0] != "Bearer" {
		return "", lib_errors.Errorf("Invalid Authorization header: expected format 'Bearer {token}', found %q", a)
	}
	return w[1], nil
}
