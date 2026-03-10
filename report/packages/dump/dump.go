// Package dump provides some helper methods for debugging
//
// Generally used in testing and logging only.
package dump

import (
	"encoding/json"
	"fmt"
)

func Any(item any) string {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	return fmt.Sprintf("%+v\n", str)
}

func Now(item any) {
	fmt.Println(Any(item))
}
