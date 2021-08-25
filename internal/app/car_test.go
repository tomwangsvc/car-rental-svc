package app

import (
	"car-svc/internal/lib/dto"
	spanner_mock "car-svc/internal/lib/spanner/mock"
	"context"
	"reflect"
	"testing"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
	lib_testing "github.com/tomwangsvc/lib-svc/testing"
)

func Test_client_CreateCar(t *testing.T) {
	type expected struct {
		result string
		err    error
	}
	var data = []struct {
		desc string
		client
		expected
	}{
		{
			desc:   "spanner error",
			client: clientErrorSpanner,
			expected: expected{
				err: lib_errors.Wrap(spanner_mock.ExpectedErrorClient, "Failed creating car"),
			},
		},
		{
			desc:   "success",
			client: clientSuccess,
			expected: expected{
				result: lib_mock.ExpectedResultString,
				err:    nil,
			},
		},
	}

	for i, d := range data {
		result, err := d.client.CreateCar(context.Background(), dto.CarCreate{})

		if d.expected.err != nil {
			if !reflect.DeepEqual(err, d.expected.err) {
				var r interface{} = err
				if err != nil {
					r = err.Error()
				}
				t.Error(lib_testing.Errorf(lib_testing.Error{
					Unexpected: "err not equal",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.err.Error(),
					Result:     r,
				}))
			}
		} else if err != nil {
			t.Error(lib_testing.Errorf(lib_testing.Error{
				Unexpected: "err exists",
				Desc:       d.desc,
				At:         i,
				Expected:   nil,
				Result:     err.Error(),
			}))

		} else {
			if !reflect.DeepEqual(result, d.expected.result) {
				t.Error(lib_testing.Errorf(lib_testing.Error{
					Unexpected: "result",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.result,
					Result:     result,
				}))
			}
		}
	}
}
