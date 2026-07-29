// Package convert contains helpful functions to convert between data types
package convert

import (
	"encoding/json"
	"errors"
)

var (
	ErrMarshaling   = errors.New("error marshaling source data.")
	ErrUnMarshaling = errors.New("error unmarshaling into destination.")
)

// Ptr converts an item to the pointer version of itself.
func Ptr[T any](item T) *T {
	var ptr = &item
	return ptr
}

// Convert takes original T and by using json marshaling and then unmarshaling applies its
// content to destination R
//
// destination should always be passed like `Convert(srcStruct, &destinationMap)` so added
// the *D to enforce
func Convert[S any, D any](source S, destination *D) (err error) {
	var (
		e     error  // localised error used when marshaling from source
		bytes []byte // bytes captures the ouput of the marshal before converting to destination
	)

	if bytes, e = json.MarshalIndent(source, "", "  "); e == nil {
		err = json.Unmarshal(bytes, destination)
	}
	// error handling to try and give more details on which step failed
	if e != nil {
		err = errors.Join(ErrMarshaling, e)
	}
	if err != nil {
		err = errors.Join(ErrUnMarshaling, err)
	}

	return
}
