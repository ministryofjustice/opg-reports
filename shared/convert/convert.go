package convert

import (
	"encoding/json"
	"io"
	"net/http"
)

func ToJson[T any](item T) (content []byte, err error) {
	return json.MarshalIndent(item, "", "  ")
}

func ListToJson[T any](items []T) (content []byte, err error) {
	return json.MarshalIndent(items, "", "  ")
}

func ToMap[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	}
	return
}

// FromMap uses json marshaling to convert from a map back to a struct.
func FromMap[T any](m map[string]interface{}) (item T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(jBytes, &item)
	}

	return
}

// Stringify returns the body content of a http.Response as both a string and []byte.
// Very helpful for debugging, testing and converting back and forth from the api.
func Stringify(r *http.Response) (s string, b []byte) {
	b, _ = io.ReadAll(r.Body)
	s = string(b)

	return
}
