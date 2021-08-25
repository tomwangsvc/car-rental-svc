package http

import (
	"context"
	"net/http"
	"time"

	"github.com/sethgrid/pester"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	HttpClientLogger
}

// NewClient returns a new http client.
// A new client should only be used for one request.
func NewClient(ctx context.Context, allowRetries bool) HttpClient {
	lib_log.Info(ctx, "Initializing", lib_log.FmtBool("allowRetries", allowRetries))

	var client HttpClient
	if allowRetries {
		client = newHttpClientWithRetries(ctx)
	} else {
		client = newHttpClientNoRetries()
	}

	lib_log.Info(ctx, "Initialized")
	return client
}

const (
	outgoingRequestTimeout = 90 * time.Second
)

// NewClient returns a new http client with retries configured.
func newHttpClientWithRetries(ctx context.Context) HttpClient {

	client := pester.New()

	// Pester will run a log hook when an error occurs with an attempt
	client.LogHook = func(pesterErrEntry pester.ErrEntry) {
		lib_log.Info(ctx, "Pester http client encountered an error, this usually results in a retry when possible", lib_log.FmtAny("pesterErrEntry", pesterErrEntry), lib_log.FmtError(pesterErrEntry.Err))
	}

	// We do not want to implement a timeout on outgoing requests because the code sending the request will have a context with a timeout however because it is dangerous to never timeout we configure a fallback
	client.Timeout = outgoingRequestTimeout

	// Pester's timeout will cause retries even if the initial request is successful, which we don't want.
	// We only want to retry if we get 429 or > 500 status
	client.RetryOnHTTP429 = true
	return client
}

// NewClient returns a new http client with retries disabled.
func newHttpClientNoRetries() HttpClient {

	client := http.Client{}

	// We do not want to implement a timeout on outgoing requests because the code sending the request will have a context with a timeout however because it is dangerous to never timeout we configure a fallback
	client.Timeout = outgoingRequestTimeout

	return httpClientNoRetries{
		client: client,
	}
}

type httpClientNoRetries struct {
	client http.Client
}

func (c httpClientNoRetries) Do(req *http.Request) (resp *http.Response, err error) {
	return c.client.Do(req)
}

func (c httpClientNoRetries) LogString() string {
	return ""
}
