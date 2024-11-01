// Package convert contains funcs to convert between types
//
// Typically things like a date entered / stored as a string
// back into a time.Time.
//
// Also swapping from struct into bytes for http usage etc
package convert

import (
	"encoding/json"
	"log/slog"
	"os"
)

// Map converts any T item (struct generally) into a map
func Map[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("[convert.Map] failed", slog.String("err", err.Error()))
	}
	return
}

// Marshal
func Marshal[T any](item T) (bytes []byte, err error) {

	return
}

func Cast[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(source); err == nil {
		err = UnmarshalInto(bytes, destination)
	}
	return
}

// Unmarshal wraper json.Unmarshal and handles error messages etc
func Unmarshal[T any](content []byte) (item T, err error) {
	var i T
	err = json.Unmarshal(content, &i)
	if err != nil {
		slog.Error("[convert.Unmarshal] failed", slog.String("err", err.Error()))
	} else {
		item = i
	}
	return
}

func UnmarshalInto[T any](content []byte, item T) (err error) {
	err = json.Unmarshal(content, &item)
	if err != nil {
		slog.Error("[convert.UnmarshalInto] failed", slog.String("err", err.Error()))
	}
	return
}

// UnmarshalFile converts the content of the file into item T
// - reads the content of then file and then uses Unmarshal
func UnmarshalFile[T any](filepath string) (item T, err error) {
	var content []byte

	if content, err = os.ReadFile(filepath); err != nil {
		slog.Error("[convert.UnmarshalFile] failed", slog.String("err", err.Error()))
		return
	}

	item, err = Unmarshal[T](content)
	return
}
