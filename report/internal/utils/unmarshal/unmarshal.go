package unmarshal

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// Unmarshal takes []byte - likely from file content or similar - and uses
// json.Unmarshal to convert that to the struct passed in destination
func Unmarshal[T any](content []byte, destination T) (err error) {
	err = json.Unmarshal(content, destination)
	return
}

// FromFile takes a file path, reads its content, and the uses the []byte
// with Unmarshal
func FromFile[T any](filepath string, destination T) (err error) {
	var content []byte

	if content, err = os.ReadFile(filepath); err != nil {
		return
	}
	err = Unmarshal[T](content, destination)
	return
}

// FromResponse is a helper to read and unmarshal content returned from a http
// call, such as fetching data from the api
func FromResponse[T any](response *http.Response, destination T) (err error) {
	var content []byte
	defer response.Body.Close()

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = Unmarshal[T](content, destination)
	return
}
