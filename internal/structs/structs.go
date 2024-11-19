package structs

import (
	"encoding/json"
	"log/slog"
	"os"
)

// ToMap takes a struct, marshals that to json and then unmarshals
// that into a map[string]interface.
// Used to get values out of structs etc
func ToMap[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("[structs] failed to covnert to a map", slog.String("err", err.Error()))
	}
	return
}

// Unmarshal takes []byte - likely from file content or similar - and uses
// json.Unmarshal to convert that to the struct passed in destination
func Unmarshal[T any](content []byte, destination T) (err error) {
	err = json.Unmarshal(content, destination)
	return
}

// UnmarshalFile takes a file path, reads its content, and the uses the []byte
// with Unmarshal
func UnmarshalFile[T any](filepath string, destination T) (err error) {
	var content []byte

	if content, err = os.ReadFile(filepath); err != nil {
		slog.Error("[structs] UmarshalFile failed to read a file", slog.String("err", err.Error()))
		return
	}
	err = Unmarshal[T](content, destination)
	return
}
