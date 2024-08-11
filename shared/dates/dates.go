package dates

import (
	"log/slog"
	"time"
)

const Format string = time.RFC3339
const FormatYMD string = "2006-01-02"
const FormatYM string = "2006-01"
const FormatY string = "2006"

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
