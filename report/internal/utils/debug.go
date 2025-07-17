package utils

import "fmt"

// Debug is a helper function that runs printf against a json
// string version of the item passed.
// Used for testing only.
func Debug[T any](item T) {
	fmt.Printf("%+v\n", MarshalStr(item))
}
