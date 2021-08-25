package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_testing "github.com/tomwangsvc/lib-svc/testing"
)

func Test_client_Health(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/?health=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(client{config: Config{Env: lib_env.Env{Id: "dev"}}}.Health()).ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Error(lib_testing.Errorf(lib_testing.Error{
			Unexpected: "status",
			Expected:   expectedStatus,
			Result:     status,
		}))
	}

	expectedBody := fmt.Sprintf("{\"build_date\":\"%s\",\"build_number\":\"%s\",\"commit_id\":\"%s\",\"env\":\"dev\",\"status\":\"OK\",\"svc\":\"\"}\n", "@foo_BUILD_DATE@", "@foo_BUILD_NUMBER@", "@foo_COMMIT_ID@")
	if body := rr.Body.String(); body != expectedBody {
		t.Error(lib_testing.Errorf(lib_testing.Error{
			Unexpected: "body",
			Expected:   expectedBody,
			Result:     body,
		}))
	}
}
