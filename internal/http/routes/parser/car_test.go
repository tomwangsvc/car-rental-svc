package parser

import (
	"bytes"
	"car-svc/internal/lib/dto"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_schema_mock "github.com/tomwangsvc/lib-svc/schema/mock"
	lib_testing "github.com/tomwangsvc/lib-svc/testing"
)

func Test_ParseCreateCar(t *testing.T) {
	carCreate := dto.CarCreate{
		Test: true,
		UserInput: dto.CarCreateUserInput{
			BrandName: "brand_name",
			ModelName: "model_name",
		},
	}

	ctx := context.Background()
	ctx = lib_context.WithTest(ctx, carCreate.Test)
	body, err := json.Marshal(carCreate.UserInput)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("", "", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)

	type expected struct {
		err      error
		hasError bool
		result   *dto.CarCreate
	}
	var data = []struct {
		desc string
		client
		input *http.Request
		expected
	}{
		{
			desc:   "success",
			client: clientSuccess,
			input:  req,
			expected: expected{
				result: &carCreate,
			},
		},
		{
			desc:   "schema error",
			client: clientErrorLibSchema,
			input:  req,
			expected: expected{
				err:      lib_errors.Wrap(lib_schema_mock.ExpectedErrorClient, "Failed checking body against schema"),
				hasError: true,
				result:   nil,
			},
		},
	}

	for i, d := range data {
		result, err := d.client.ParseCreateCar(d.input)

		if d.expected.hasError {
			if err == nil {
				t.Error(lib_testing.Errorf(lib_testing.Error{
					Unexpected: "err not exist",
					Desc:       d.desc,
					At:         i,
					Expected:   nil,
					Result:     d.expected,
				}))
			}

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
			if !reflect.DeepEqual(*result, *d.expected.result) {
				t.Error(lib_testing.Errorf(lib_testing.Error{
					Unexpected: "result",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected,
					Result:     result,
				}))
			}
		}
	}
}
