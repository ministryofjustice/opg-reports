package convert

import (
	"encoding/json"
	"log/slog"
	"os"
)

func Map[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("[convert.Map] failed", slog.String("err", err.Error()))
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
