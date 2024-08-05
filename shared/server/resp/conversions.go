package resp

import (
	"encoding/json"
	"io"
	"net/http"
)

// Stringify returns the body content of a http.Response as both a string and []byte.
// Very helpful for debugging, testing and converting back and forth from the api.
func Stringify(r *http.Response) (s string, b []byte) {
	b, _ = io.ReadAll(r.Body)
	s = string(b)
	return
}

// ToJson converts a response into json friendly []bye that is indented for readability.
// This is used for passing the data back from the api
func ToJson(r *Response) (content []byte, err error) {
	return json.MarshalIndent(r, "", "  ")
}

// FromJson converts a []byte back into an IResponse by using json unmarshaling
func FromJson(content []byte, r *Response) (err error) {
	err = json.Unmarshal(content, r)
	return
}

// FromHttp is similar to FromJson, but first fetches the content from the http.Repsonse body
// and then converts using that
func FromHttp(content *http.Response, r *Response) (err error) {
	_, bytes := Stringify(content)
	return FromJson(bytes, r)
}
