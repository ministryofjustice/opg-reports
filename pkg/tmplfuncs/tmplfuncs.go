// Package tmplfuncs contains series of funcs to use within
// front end templates
//
// Exposed as a map (`All`)
package tmplfuncs

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type nums interface {
	float32 | float64 | int
}
type adders interface {
	nums | string
}

func add[T adders](a T, args ...any) (result T) {
	result = a
	for _, arg := range args {
		// check the T casting works before doing the +
		if val, ok := arg.(T); ok {
			result += val
		}
	}
	return
}

// addString converts the strings into floats and adds
// those together as numbers before returning a string
// version
// If there is an error converting  from a string to
// a float it will return an error
func addString(a string, args ...any) (result string, err error) {
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

	result = fmt.Sprintf("%.4f", floated)

	return
}

// Add handles "adding" floats, ints and strings being added.
//
// For strings, it will try to treat them as floats first (via
// `addString`) but if that fails due to parsing errors it will
// instead concatenate them (via `add`).
//
// Examples:
//
//	Add(1, 2, 3 ) 	// 6
//	Add("1", "2")	// 3
//	Add(1.0, 2.0)	// 3.0
//	Add("A", "b")	// "Ab"
func Add(a interface{}, args ...interface{}) (result interface{}) {
	switch a.(type) {
	case float64:
		result = add(a.(float64), args...)
	case int:
		result = add(a.(int), args...)
	case string:
		v, err := addString(a.(string), args...)
		if err != nil {
			result = add(a.(string), args...)
		} else {
			result = v
		}
	default:
		result = ""
	}
	return
}

func Increment(i interface{}) (result interface{}) {
	result = i
	switch i.(type) {
	case float64:
		result = add(i.(float64), 1)
	case int:
		result = add(i.(int), 1)
	}

	return
}

func Title(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

// func Col(i string, mapped map[string]string) string {
// 	return mapped[i]
// }

var All map[string]interface{} = map[string]interface{}{
	"add":       Add,
	"increment": Increment,
	"title":     Title,
}
