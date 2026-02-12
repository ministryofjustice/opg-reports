package debugger

import (
	"encoding/json"
	"fmt"
)

// Dump is a helper function that runs printf against a json
// string version of the item passed.
// Used for testing only.
func Dump[T any](item T) {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	fmt.Printf("%+v\n", str)
}

func DumpStr[T any](item T) string {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	return fmt.Sprintf("%+v\n", str)
}
