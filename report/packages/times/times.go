package times

import (
	"time"
)

// Interval type is used for increamenting / resetting times
type Interval string

// Enum values for Interval
const (
	YEAR   Interval = "yearly"
	MONTH  Interval = "monthly"
	DAY    Interval = "daily"
	HOUR   Interval = "hourly"
	MINUTE Interval = "minute"
	SECOND Interval = "second"
)

// Format is representation of date formats
type Format string

// Enum for time related format strings
const (
	FULL string = time.RFC3339
	YMD  string = "2006-01-02"
	YM   string = "2006-01"
	Y    string = "2006"
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
		layout string = FULL
		max    int    = len(layout)
		l      int    = len(value)
	)
	// return FULL version f the length is longer
	if l > max {
		return layout
	}
	return layout[:l]
}

// FromString converts string into a time where possible
func FromString(str string) (t time.Time) {
	t, _ = time.Parse(GetFormat(str), str)
	return
}

// Time converts string / time to time
func Time(src any) (res time.Time) {
	switch any(src).(type) {
	case time.Time:
		res = src.(time.Time)
	case string:
		res = FromString(src.(string))
	}
	return
}
