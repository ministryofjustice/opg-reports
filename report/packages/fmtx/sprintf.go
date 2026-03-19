// Package fmtx contains extended methods for fmt, such as named parameters
package fmtx

import (
	"fmt"
	"opg-reports/report/packages/convert"
	"reflect"
	"regexp"
	"strings"
)

const pattern string = `(?m):[[:alnum:]_-]+`

var re *regexp.Regexp = regexp.MustCompile(pattern)

// SprintfNamed reviews the source string for any place holder values in the format
// of `:${name}` and replaces those with real values from the data source.
//
// The `${name}` must match the pattern `:[[:alnum]_-]+`, so it can only be made
// from alphanumeric chars, hyphen or underscores.
//
// The `${name}` is used as the key to the data map:
//
//	`:${name}` => data[${name}]
//
// Calls itself reursively to replace every instance of `:${name}` it can find.
//
// `:${name}` patterns that dont have a key in the map are removed & ignored.
func SprintfNamed(str string, data map[string]interface{}, sql bool) (updated string, ordered []interface{}) {
	var (
		s, e     = pos(str)
		expanded = expandedValues(data)
		key      = ""
		replace  = ""
		params   = []interface{}{}
	)
	// if a pattern is found and end is more than 0 then we should replace
	// it..
	if s >= 0 && e > 0 {
		// find the key for the map
		key = strings.TrimPrefix(str[s:e], `:`)
		// append the ordered values
		params = append(params, expanded[key]...)
		// get the replacement
		replace = replacement(key, expanded, sql)
		// insert
		str = str[0:s] + replace + str[e:]
	}
	// set the values
	updated = str
	ordered = params
	// look for next replacements... if found, recurse and update
	// values
	s, e = pos(str)
	if e > 0 {
		u, o := SprintfNamed(str, data, sql)
		updated = u
		ordered = append(ordered, o...)
	}

	return
}

// replacement creates the string replacement to be used in string re-formatting
//
// If sel is true, then values become ? placeholders for binding later
func replacement(key string, ex map[string][]interface{}, sql bool) (val string) {
	val = ""
	for _, v := range ex[key] {
		if sql {
			val += "?,"
		} else {
			val += fmt.Sprintf("%s,", convert.String(v))
		}
	}
	val = strings.TrimSuffix(val, `,`)
	return
}

// pos returns the the positions of the first `:${name}` pattern
// within the passed string.
func pos(str string) (s int, e int) {
	var idx = re.FindStringIndex(str)
	s = -1
	e = -1

	if len(idx) >= 2 {
		s = idx[0]
		e = idx[1]
	}
	return
}

// expandedValues converts each value in the map into a slice of interfaces.
//
// This allows uniform joining for a slice of somethign as well a a single value
// via slices.Join when doing the data sub
func expandedValues(data map[string]interface{}) (expanded map[string][]interface{}) {
	expanded = map[string][]interface{}{}

	for k, v := range data {
		expanded[k] = asSlice(v)
	}

	return
}

// asSlice uses reflection to expand val into multiples if its a p[slice etc]
func asSlice[T any](val T) (values []interface{}) {
	var v reflect.Value
	var t reflect.Type
	values = []interface{}{}

	v = reflect.ValueOf(val)
	t = v.Type()
	if t.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface().(T))
		}
	} else {
		values = append(values, val)
	}
	return
}
