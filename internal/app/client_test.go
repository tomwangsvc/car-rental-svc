package app

import (
	spanner_mock "car-svc/internal/lib/spanner/mock"
)

var (
	clientErrorSpanner = client{
		spannerClient: spanner_mock.ClientError,
	}
	clientErrorSpannerTransform = client{
		spannerClient: spanner_mock.ClientErrorTransform,
	}
	clientSuccess = client{
		spannerClient: spanner_mock.ClientSuccess,
	}
)
