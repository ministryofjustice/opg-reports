// Package dates igroups a series of funcs converting between time.Time and string
//
// Provides series of helpers for string to become time.Time and to parse / convert
// existing times.Time into common formats and structs
package dates

import (
	"log/slog"
	"time"
)

// 2006-01-02T15:04:05Z07:00
const Format string = time.RFC3339
const FormatYMDHMS string = "2006-01-02T15:04:05"
const FormatYMD string = "2006-01-02"
const FormatYM string = "2006-01"
const FormatY string = "2006"
const ErrorYear string = "0000"

var ErrorTime time.Time = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

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

func IsDate(str string) (valid bool) {
	valid = true
	df := GetFormat(str)
	if _, err := time.Parse(df, str); err != nil {
		valid = false
	}
	return
}
