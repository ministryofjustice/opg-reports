package structs

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

// ToFile converts the item passed into json byte array and writes that to
// the file path passed along
func ToFile[T any](item T, filename string) (err error) {
	var (
		bytes     []byte
		parentDir string = filepath.Dir(filename)
	)
	os.MkdirAll(parentDir, os.ModePerm)

	bytes, err = json.MarshalIndent(item, "", "  ")
	if err != nil {
		return
	}
	err = os.WriteFile(filename, bytes, os.ModePerm)

	return
}

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

// Convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(source); err == nil {
		err = Unmarshal(bytes, destination)
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

func Jsonify[T any](item T) (str string) {
	bytes, _ := json.MarshalIndent(item, "", "  ")
	str = string(bytes)
	return
}
