package convert

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const indentWith string = "  "

// Marshal is a local wrapper around json.MarshalIndent for
// consistency
func Marshal[T any](item T) (content []byte, err error) {
	return json.MarshalIndent(item, "", indentWith)
}

// Marshals is also a wraper around MarshalIndent, used to keep naming
// convention of unmarshaling versions
func Marshals[T any](items []T) (content []byte, err error) {
	return json.MarshalIndent(items, "", indentWith)
}

// Unmarshal wraper json.Unmarshal and handles error messages etc
func Unmarshal[T any](content []byte) (item T, err error) {
	var i T
	err = json.Unmarshal(content, &i)
	if err != nil {
		slog.Error("unmarshal failed", slog.String("err", err.Error()))
	} else {
		item = i
	}
	return
}

// Unmarshals wrapper for mutliple types for shorthand and deals with error
// logging
func Unmarshals[T any](content []byte) (items []T, err error) {
	var i []T
	err = json.Unmarshal(content, &i)
	if err != nil {
		slog.Error("unmarshals failed", slog.String("err", err.Error()))
	} else {
		items = i
	}
	return
}

// Map converts a structure into the map of interface by using
// json marshaling to convert between types
func Map[T any](item T) (m map[string]interface{}, err error) {
	byt, err := json.Marshal(item)
	if err == nil {
		err = json.Unmarshal(byt, &m)
	} else {
		slog.Error("map failed", slog.String("err", err.Error()))
	}
	return
}

// Maps converts a slice of structures into the map of interface by using
// json marshaling to convert between types
func Maps[T any](item []T) (m []map[string]interface{}, err error) {
	bytes, err := Marshals(item)
	if err == nil {
		m, err = Unmarshals[map[string]interface{}](bytes)
	} else {
		slog.Error("map failed", slog.String("err", err.Error()))
	}
	return
}

// Unmap uses json marshaling to convert from a map back to a struct.
func Unmap[T any](m map[string]interface{}) (item T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(jBytes, &item)
	} else {
		slog.Error("unmap failed", slog.String("err", err.Error()))
	}

	return
}

// Unmaps uses json marshaling to convert from a map back to a struct.
func Unmaps[T any](m []map[string]interface{}) (items []T, err error) {
	jBytes, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(jBytes, &items)
	} else {
		slog.Error("unmaps failed", slog.String("err", err.Error()))
	}

	return
}

// Stringify returns the body content of a http.Response as both a string and []byte.
// Very helpful for debugging, testing and converting back and forth from the api.
func Stringify(r *http.Response) (s string, b []byte) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("stringify failed", slog.String("err", err.Error()))
	}
	s = string(b)

	return
}

// String uses the json marshal to quickly convert any
// struct into a string for display
func String[T any](item T) (s string) {
	bytes, _ := Marshal(item)
	s = string(bytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, indentWith, "")
	return
}

// PrettyString is similar to String, but retains the indentations
func PrettyString[T any](item T) (s string) {
	bytes, _ := Marshal(item)
	s = string(bytes)
	return
}

// IntToBool helper used with sql conversion as sqlite has no
// boolean type, they are stored as 1 (true) or 0, this maps them back to
// a bool
func IntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

// IntToBool helper used with sql conversion as sqlite has no
// boolean type, they are stored as 1 (true) or 0
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolStringToInt helper to deal with get param bools
// that convert over to 1 | 0 for the db
func BoolStringToInt(s string) int {
	b, err := strconv.ParseBool(s)
	if err == nil && b {
		return 1
	}
	return 0
}

// Title is used by template helpers to make a "pretty" versions
// of a column name
//
//	team_a => Team A
func Title(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

// Dict generates a map from an even number of arguments (key, value)
// Used within templates to merge data
func Dict(values ...any) (dict map[string]any) {
	dict = map[string]any{}
	if len(values)%2 != 0 {
		return
	}
	// if the key isnt a string, this will crash!
	for i := 0; i < len(values); i += 2 {
		var key string = values[i].(string)
		var v any = values[i+1]
		dict[key] = v
	}
	return
}

// Curr generates a currency formated string from
// either a string or float64 value and prepends
// the symbol to the start
// Will round to nearest .2
// Returns <symbol>0.0 by default
func Curr(s interface{}, symbol string) string {
	p := message.NewPrinter(language.English)
	switch s.(type) {
	case string:
		f, _ := strconv.ParseFloat(s.(string), 10)
		return symbol + p.Sprintf("%.2f", f)
	case float64:
		return symbol + p.Sprintf("%.2f", s.(float64))
	}
	return symbol + "0.0"
}

// StripIntPrefix removes chunk of string before the first '.'
// - this is used in templates to remove integer indexes from strings
// - will remove all other . from the source string
func StripIntPrefix(s string) string {
	sp := strings.Split(s, ".")
	if len(sp) > 1 {
		return strings.Join(sp[1:], "")
	}
	return s
}

// Percent works out % value of items
// Used in templates
// Returns 2 decimals
func Percent(got int, total int) string {
	x := float64(got)
	y := float64(total)
	p := x / (y / 100)
	return fmt.Sprintf("%.2f", p)
}
