package testing

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	letters []rune
)

func init() {
	rand.Seed(time.Now().UnixNano())

	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

type Error struct {
	Unexpected    string      `json:"unexpected"`
	Desc          string      `json:"desc"`
	At            int         `json:"at"`
	Input         interface{} `json:"input"`
	Expected      interface{} `json:"expected"`
	Result        interface{} `json:"result"`
	MarshalErrors string      `json:"marshal_errors,omitempty"`
}

func Errorf(e Error) string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		var marshalErrors []string
		input, ierr := json.MarshalIndent(e.Input, "", "\t")
		if ierr != nil {
			input = []byte("potentially not marshallable")
			marshalErrors = append(marshalErrors, fmt.Sprintf("input: %s", ierr))
		}
		expected, eerr := json.MarshalIndent(e.Expected, "", "\t")
		if eerr != nil {
			expected = []byte("potentially not marshallable")
			marshalErrors = append(marshalErrors, fmt.Sprintf("expected: %s", eerr))
		}
		result, rerr := json.MarshalIndent(e.Result, "", "\t")
		if rerr != nil {
			result = []byte("potentially not marshallable")
			marshalErrors = append(marshalErrors, fmt.Sprintf("result: %s", rerr))
		}
		return fmt.Sprintf(fmtString, e.Unexpected, e.Desc, e.At, string(input), string(expected), string(result), strings.Join(marshalErrors, " | "))
	}
	return string(b)
}

const fmtString = `{
	"unexpected": %q,
	"desc": %q,
	"at": %d,
	"input": %q,
	"expected": %q,
	"result": %q,
	"marshal_errors": %q
}`

func RandomByteArray(n int) []byte {
	token := make([]byte, n)
	rand.Read(token)
	return token
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
