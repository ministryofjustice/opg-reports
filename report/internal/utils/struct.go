package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StructToFile converts the item passed into json byte array and writes that to
// the file path passed along
func StructToJsonFile[T any](filename string, source T) (err error) {
	var (
		bytes     []byte
		parentDir string = filepath.Dir(filename)
	)

	bytes, err = Marshal(source)
	if err != nil {
		return
	}

	os.MkdirAll(parentDir, os.ModePerm)
	err = os.WriteFile(filename, bytes, os.ModePerm)

	return
}

// StructFromJsonFile wraps UnmarshalFile but add a dedicated check to make
// sure the file exists first
func StructFromJsonFile[T any](filename string, destination T) (err error) {

	if !FileExists(filename) {
		return fmt.Errorf("file does not exist [%s]", filename)
	}
	err = UnmarshalFile(filename, destination)
	return
}

// StructToMap takes a struct, marshals that to json and then unmarshals
// that into a map[string]interface.
// Used to get values out of structs etc
func StructToMap[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	}
	return
}
