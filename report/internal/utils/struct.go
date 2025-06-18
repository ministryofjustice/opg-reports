package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// StructToFile converts the item passed into json byte array and writes that to
// the file path passed along
func StructToJsonFile[T any](item T, filename string) (err error) {
	var (
		bytes     []byte
		parentDir string = filepath.Dir(filename)
	)

	bytes, err = Marshal(item)
	if err != nil {
		return
	}

	os.MkdirAll(parentDir, os.ModePerm)
	err = os.WriteFile(filename, bytes, os.ModePerm)

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
