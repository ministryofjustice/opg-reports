package convert

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func Marshal[T any](item T) (content []byte, err error) {
	return json.MarshalIndent(item, "", "  ")
}

func Unmarshal[T any](content []byte, i T) (item T, err error) {
	err = json.Unmarshal(content, &i)
	if err != nil {
		slog.Error("unmarshal failed", slog.String("err", err.Error()))
	} else {
		item = i
	}
	return
}

func Unmarshals[T any](content []byte, i []T) (items []T, err error) {
	err = json.Unmarshal(content, &i)
	if err != nil {
		slog.Error("unmarshals failed", slog.String("err", err.Error()))
	} else {
		items = i
	}
	return
}

func Map[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("map failed", slog.String("err", err.Error()))
	}
	return
}

// Unmap uses json marshaling to convert from a map back to a struct.
func Unmap[T any](m map[string]interface{}) (item T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(jBytes, &item)
	} else {
		slog.Error("unmap failed", slog.String("err", err.Error()))
	}

	return
}

// Stringify returns the body content of a http.Response as both a string and []byte.
// Very helpful for debugging, testing and converting back and forth from the api.
func Stringify(r *http.Response) (s string, b []byte) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("stringify failed", slog.String("err", err.Error()))
	}
	s = string(b)

	return
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolStringToInt(s string) int {
	b, err := strconv.ParseBool(s)
	if err == nil && b {
		return 1
	}
	return 0
}
