// Package fmtx is an extended fmt package containing common helpers
package fmtx

import (
	"encoding/json"
	"fmt"
)

// Sprintj tries to convert the item passed into a json string
// via json.Marshal and then fmt.Sprintf.
//
// Used to create a output friently version of a complext struct
// which can be viewed in logs / cli
func Sprintj[T any](item T) (s string) {
	var str = ""
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err == nil {
		str = string(bytes)
	}
	s = fmt.Sprintf("%+v\n", str)

	return
}

// Printj tries to convert the item passed to a string via
// json.Marshal, fmt.Sprintf & and then directly output
// the result with fmt.Println.
//
// Intended to be used for debugging / outputting complex
// structs to see nested values etc
func Printj(item any) {
	fmt.Println(Sprintj(item))
}
