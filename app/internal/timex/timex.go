// Package timex contains extended functions relating to times & dates that are used regulary
package timex

import "time"

// Format is representation of date formats
type Format string

// Time formats
const (
	FULL Format = time.RFC3339
	YMD  Format = "2006-01-02"
	YM   Format = "2006-01"
	Y    Format = "2006"
)

// Intervals
type Interval string

// Time intervals
const (
	YEAR   Interval = "year"
	MONTH  Interval = "month"
	DAY    Interval = "day"
	HOUR   Interval = "hour"
	MINUTE Interval = "minute"
	SECOND Interval = "second"
)
