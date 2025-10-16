package utils

import "fmt"

// Dump is a helper function that runs printf against a json
// string version of the item passed.
// Used for testing only.
func Dump[T any](item T) {
	fmt.Printf("%+v\n", MarshalStr(item))
}
