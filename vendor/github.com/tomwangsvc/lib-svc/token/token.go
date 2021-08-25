package token

import (
	"strings"

	lib_constants "github.com/tomwangsvc/lib-svc/constants"
	lib_json "github.com/tomwangsvc/lib-svc/json"
)

func Redact(token string) string {
	if s := strings.Split(token, "."); len(s) == 3 {
		return strings.Join(append(s[:2], lib_constants.Redacted), ".")
	}
	return lib_constants.Redacted
}

var (
	JsonRedactionsByPath = map[string]lib_json.Redaction{
		"token": func(result string) string {
			return Redact(result)
		},
	}
)
