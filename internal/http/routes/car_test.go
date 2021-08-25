package routes

import (
	app_mock "car-svc/internal/app/mock"
	parser_mock "car-svc/internal/http/routes/parser/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	lib_mock "github.com/tomwangsvc/lib-svc/mock"
	lib_testing "github.com/tomwangsvc/lib-svc/testing"
)

func Test_client_CreateCar(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	router := chi.NewRouter()

	type expected struct {
		body           string
		code           int
		headerLocation string
	}
	var data = []struct {
		desc string
		client
		expected
	}{
		{
			desc:   "app error",
			client: clientErrorApp,
			expected: expected{
				body:           "",
				code:           app_mock.ExpectedErrorClient.Code,
				headerLocation: "",
			},
		},
		{
			desc:   "parser error",
			client: clientErrorParser,
			expected: expected{
				body:           "",
				code:           parser_mock.ExpectedErrorClient.Code,
				headerLocation: "",
			},
		},
		{
			desc:   "success",
			client: clientSuccess,
			expected: expected{
				body:           "",
				code:           http.StatusCreated,
				headerLocation: lib_mock.ExpectedResultString,
			},
		},
	}

	for i, d := range data {
		router.Post("/", d.client.CreateCar())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if code := rr.Code; code != d.expected.code {
			t.Error(lib_testing.Errorf(lib_testing.Error{
				Unexpected: "code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.code,
				Result:     code,
			}))
		}

		if body := rr.Body.String(); body != d.expected.body {
			t.Error(lib_testing.Errorf(lib_testing.Error{
				Unexpected: "body",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.body,
				Result:     body,
			}))
		}

		if headerLocation, ok := rr.HeaderMap["Location"]; !ok {
			if d.expected.headerLocation != "" {
				t.Error(lib_testing.Errorf(lib_testing.Error{
					Unexpected: "headerLocation exists",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.headerLocation,
					Result:     nil,
				}))
			}
		} else if strings.Join(headerLocation, ",") != d.expected.headerLocation {
			t.Error(lib_testing.Errorf(lib_testing.Error{
				Unexpected: "headerLocation exists",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.headerLocation,
				Result:     nil,
			}))
		}
	}
}
