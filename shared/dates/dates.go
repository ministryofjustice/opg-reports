// Package dates provides helpful formatting constants and a convertor
//
// To keep date formats consistent between packages, we share the formats here and include
// a StringToDate func to easily convert a string into a time.Time
package dates

import (
	"fmt"
	"log/slog"
	"time"
)

const Format string = time.RFC3339
const FormatYMD string = "2006-01-02"
const FormatYM string = "2006-01"
const FormatY string = "2006"
const ErrYear string = "0000"

// GetFormat will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func GetFormat(value string) string {
	max := len(Format)
	l := len(value)
	if l > max {
		return Format
	}
	f := Format[:l]
	slog.Debug("[dates] GetFormat", slog.String("format", f))
	return f
}

// StringToDate tries to convert the string value passed into a time.Time,
// will return errors if it doesnt work
func StringToDate(value string) (t time.Time, err error) {
	layout := GetFormat(value)
	t, err = time.Parse(layout, value)
	if err == nil {
		t = t.UTC()
	}
	slog.Debug("[dates] StringToDate", slog.String("t", t.String()), slog.String("err", fmt.Sprintf("%v", err)))
	return t, err
}

// Reformat takes a date string, converts to a time.Time and the returns
// a string in outFormat
// Allows quick conversions from YYYY-MM-DD to YYYY-MM and so on
func Reformat(value string, outFormat string) string {
	d, err := StringToDate(value)
	if err != nil {
		return ""
	}
	out := d.Format(outFormat)
	slog.Debug("[dates] Reformat", slog.String("out", out))
	return out
}

// StringToDateDefault checks if the date string passed matches a known value, if so it uses
// the defaultValue passed as the source and converts that to a time.Time
// Used to handle empty string versions of dates that typicall want to be "now"
func StringToDateDefault(value string, comp string, defaultV string) (t time.Time, err error) {
	if value == comp {
		value = defaultV
	}
	return StringToDate(value)
}

// Strings converts a series of time.Times into string versions using dateFormat
func Strings(dates []time.Time, dateFormat string) []string {
	strs := []string{}
	for _, d := range dates {
		strs = append(strs, d.Format(dateFormat))
	}
	return strs
}
