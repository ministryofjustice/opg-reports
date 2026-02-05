package marshal

import (
	"encoding/json"
)

// Marshal converts all items using json marshaling
func Marshal[T any](item T) (bytes []byte, err error) {
	bytes, err = json.MarshalIndent(item, "", "  ")
	return
}

// ToString converts the json.Marshal result into a string to easier
// usage.
//
// If an error occurs the returned string is empty.
func ToString[T any](item T) (str string) {
	str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	return
}

// Convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func Convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	if bytes, err = Marshal(source); err == nil {
		err = json.Unmarshal(bytes, destination)
	}
	return
}
