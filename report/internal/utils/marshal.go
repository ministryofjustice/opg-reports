package utils

import (
	"encoding/json"
	"os"
)

// Marshal converts all items using json marshaling
func Marshal[T any](item T) (bytes []byte, err error) {
	bytes, err = json.MarshalIndent(item, "", "  ")
	return
}

func MustMarshal[T any](item T) (bytes []byte) {
	bytes = []byte{}
	if b, err := json.MarshalIndent(item, "", "  "); err == nil {
		bytes = b
	}
	return
}

func MarshalStr[T any](item T) (str string) {
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
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
		return
	}
	err = Unmarshal[T](content, destination)
	return
}

// Convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = Marshal(source); err == nil {
		err = Unmarshal(bytes, destination)
	}
	return
}
