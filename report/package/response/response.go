package response

import (
	"encoding/json"
	"io"
	"net/http"
)

// AsT is a helper to read and unmarshal content returned from a http
// call, such as fetching data from the api
func AsT[T any](response *http.Response, destination T) (err error) {
	var content []byte
	defer response.Body.Close()

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, destination)
	return
}
