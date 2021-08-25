package json

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

var (
	null = []byte("null")
)

type NullBool struct {
	exists bool
	valid  bool
	value  bool
}

func NewNullBool(value bool) NullBool {
	return NullBool{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullBoolExistsNotValid() NullBool {
	return NullBool{
		exists: true,
	}
}

func (n NullBool) Exists() bool {
	return n.exists
}

func (n NullBool) Valid() bool {
	return n.valid
}

func (n NullBool) Value() (value bool, ok bool) {
	return n.value, n.valid
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf("%v", n.value)), nil
	}
	return null, nil
}

func (n *NullBool) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = false
		return nil
	}

	b, err := strconv.ParseBool(string(payload))
	if err != nil {
		return lib_errors.Errorf("Failed converting payload to bool: got %q", string(payload))
	}
	n.exists = true
	n.valid = true
	n.value = b

	return nil
}

type NullDate struct {
	exists bool
	valid  bool
	value  civil.Date
}

func NewNullDate(value civil.Date) NullDate {
	return NullDate{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullDateExistsNotValid() NullDate {
	return NullDate{
		exists: true,
	}
}

func (n NullDate) Exists() bool {
	return n.exists
}

func (n NullDate) Valid() bool {
	return n.valid
}

func (n NullDate) Value() (value civil.Date, ok bool) {
	return n.value, n.valid
}

func (n NullDate) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf("%q", n.value)), nil
	}
	return null, nil
}

func (n *NullDate) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = civil.Date{}
		return nil
	}

	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}

	d, err := civil.ParseDate(string(payload))
	if err != nil {
		return lib_errors.Errorf("Failed converting payload to civil.Date: got %q", string(payload))
	}
	n.exists = true
	n.valid = true
	n.value = d

	return nil
}

type NullFloat64 struct {
	exists bool
	valid  bool
	value  float64
}

func NewNullFloat64(value float64) NullFloat64 {
	return NullFloat64{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullFloat64ExistsNotValid() NullFloat64 {
	return NullFloat64{
		exists: true,
	}
}

func (n NullFloat64) Exists() bool {
	return n.exists
}

func (n NullFloat64) Valid() bool {
	return n.valid
}

func (n NullFloat64) Value() (value float64, ok bool) {
	return n.value, n.valid
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf("%v", n.value)), nil
	}
	return null, nil
}

func (n *NullFloat64) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = 0
		return nil
	}

	f, err := strconv.ParseFloat(string(payload), 64)
	if err != nil {
		return lib_errors.Errorf("Failed converting payload to float64: got %q", string(payload))
	}
	n.exists = true
	n.valid = true
	n.value = f

	return nil
}

type NullInt64 struct {
	exists bool
	valid  bool
	value  int64
}

func NewNullInt64(value int64) NullInt64 {
	return NullInt64{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullInt64ExistsNotValid() NullInt64 {
	return NullInt64{
		exists: true,
	}
}

func (n NullInt64) Exists() bool {
	return n.exists
}

func (n NullInt64) Valid() bool {
	return n.valid
}

func (n NullInt64) Value() (value int64, ok bool) {
	return n.value, n.valid
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf("%v", n.value)), nil
	}
	return null, nil
}

func (n *NullInt64) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = 0
		return nil
	}

	i, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return lib_errors.Errorf("Failed converting payload to int64: got %q", string(payload))
	}
	n.exists = true
	n.valid = true
	n.value = i

	return nil
}

type NullString struct {
	exists bool
	valid  bool
	value  string
}

func NewNullString(value string) NullString {
	return NullString{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullStringExistsNotValid() NullString {
	return NullString{
		exists: true,
	}
}

func (n NullString) Exists() bool {
	return n.exists
}

func (n NullString) Valid() bool {
	return n.valid
}

func (n NullString) Value() (value string, ok bool) {
	return n.value, n.valid
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf(`"%s"`, n.value)), nil
	}
	return null, nil
}

func (n *NullString) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = ""
		return nil
	}

	s, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	n.exists = true
	n.valid = true
	n.value = string(s)

	return nil
}

type NullTime struct {
	exists bool
	valid  bool
	value  time.Time
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{
		exists: true,
		valid:  true,
		value:  value,
	}
}

func NewNullTimeExistsNotValid() NullTime {
	return NullTime{
		exists: true,
	}
}

func (n NullTime) Exists() bool {
	return n.exists
}

func (n NullTime) Valid() bool {
	return n.valid
}

func (n NullTime) Value() (value time.Time, ok bool) {
	return n.value, n.valid
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.valid {
		return []byte(fmt.Sprintf("%q", n.value.Format(time.RFC3339Nano))), nil
	}
	return null, nil
}

func (n *NullTime) UnmarshalJSON(payload []byte) error {
	if len(payload) == 0 {
		return lib_errors.New("Expected payload, found none")
	}

	if bytes.Equal(payload, null) {
		n.exists = true
		n.valid = false
		n.value = time.Time{}
		return nil
	}

	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339Nano, string(payload))
	if err != nil {
		return lib_errors.Errorf("Failed converting payload to time: got %q", string(payload))
	}
	n.exists = true
	n.valid = true
	n.value = t

	return nil
}

func trimDoubleQuotes(payload []byte) ([]byte, error) {
	if len(payload) <= 1 || payload[0] != '"' || payload[len(payload)-1] != '"' {
		return nil, lib_errors.Errorf("Expected payload to be wrapped with double quotes: got %q", string(payload))
	}
	return payload[1 : len(payload)-1], nil
}
