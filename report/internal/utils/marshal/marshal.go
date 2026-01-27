package marshal

import "encoding/json"

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
