package times

import "time"

// Format is representation of date formats
type Format string

// Time related format strings
const (
	FULL Format = time.RFC3339
	YMD  Format = "2006-01-02"
	YM   Format = "2006-01"
	Y    Format = "2006"
)

// GetFormat will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
//
// Used in coverting between string and time.Time for the Parse call
func GetFormat(value string) string {
	var (
		layout string = string(FULL)
		max    int    = len(layout)
		l      int    = len(value)
	)
	// return FULL version f the length is longer
	if l > max {
		return layout
	}
	return layout[:l]
}
