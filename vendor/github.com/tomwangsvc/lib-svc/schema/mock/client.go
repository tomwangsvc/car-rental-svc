package mock

import (
	"context"
	"encoding/binary"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_schema "github.com/tomwangsvc/lib-svc/schema"
)

var (
	ClientError         lib_schema.Client = clientError{}
	ClientSuccess       lib_schema.Client = clientSuccess{}
	ExpectedErrorClient                   = lib_errors.NewCustom(int(binary.BigEndian.Uint64([]byte("LIB_SCHEMA_CLIENT"))), "Mock miscellaneous error")
)

type clientError struct{}

func (c clientError) CheckContentAgainstSchema(_ context.Context, _ string, _ interface{}) error {
	return ExpectedErrorClient
}

func (c clientError) CheckBodyAgainstSchema(_ context.Context, _ string, _ []byte) error {
	return ExpectedErrorClient
}

type clientSuccess struct{}

func (c clientSuccess) CheckContentAgainstSchema(_ context.Context, _ string, _ interface{}) error {
	return nil
}

func (c clientSuccess) CheckBodyAgainstSchema(_ context.Context, _ string, _ []byte) error {
	return nil
}
