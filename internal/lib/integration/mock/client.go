package mock

import (
	"car-svc/internal/lib/integration"
	"encoding/binary"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

var (
	ClientError   integration.Client = clientError{}
	ClientSuccess integration.Client = clientSuccess{}

	ExpectedErrorClient = lib_errors.NewCustom(int(binary.BigEndian.Uint64([]byte("INTEGRATION_CLIENT"))), "")
)

type clientError struct{}

type clientSuccess struct{}
