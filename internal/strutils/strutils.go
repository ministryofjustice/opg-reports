package strutils

import (
	"fmt"
	"strconv"
)

// Float converts a string to a float, swallowing any errors
func Float(s string) (f float64) {
	f, _ = strconv.ParseFloat(s, 10)
	return
}

// FloatF converts to a string and then formats to 4 decimals
func FloatF(s string) (f string) {
	val, _ := strconv.ParseFloat(s, 10)
	f = fmt.Sprintf("%.4f", val)
	return
}

// Adds converts strings to floats, add them together and returns
// the value
// Uses Add, and swallows the error
func Adds(a string, b string) (sum string) {
	sum, _ = Add(a, b)
	return
}

// Add converts the strings into floats and adds those together as
// numbers before returning a string version
//
// If there is an error converting from a string to
// a float it will return an error
//
// Uses `%g` for format to return float without trailing 0s
func Add(a string, args ...any) (result string, err error) {

	result = a
	floated, err := strconv.ParseFloat(a, 10)
	if err != nil {
		return
	}
	for _, arg := range args {
		if val, err := strconv.ParseFloat(arg.(string), 10); err == nil {
			floated += val
		}
	}

	result = fmt.Sprintf("%g", floated)

	return
}
