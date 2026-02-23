// Package response provides functions for handling a response
//
// Typically data conversion to structs etc
package response

import (
	"encoding/json"
	"io"
	"net/http"
)

// As is a helper to read and unmarshal resposne body returned from a http
// call into a struct / map
//
// Used when fetching data from the api
func As[T any](response *http.Response, destination T) (err error) {
	var content []byte
	defer response.Body.Close()

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, destination)
	return
}

// AsString returns response body as a string, used more for fetching raw
// html / plain text content
func AsString(response *http.Response) (destination string, err error) {
	var content []byte
	defer response.Body.Close()

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	destination = string(content)

	return
}
