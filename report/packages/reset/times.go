package reset

import (
	"opg-reports/report/packages/times"
	"time"
)

// Time changes the time to be the start of the interval type asked for to avoid
// errors with date addition (rounding when adding days to month etc).
//
//	YEAR:	`2024-06-05T23:59:55` => `2024-01-01T00:00:00`
//	MONTH:	`2024-06-05T23:59:55` => `2024-06-01T00:00:00`
//	DAY:	`2024-06-05T23:59:55` => `2024-06-05T00:00:00`
//	HOUR:	`2024-06-05T23:59:55` => `2024-06-05T23:00:00`
//	MINUTE:	`2024-06-05T23:59:55` => `2024-06-05T23:59:00`
//	SECOND:	`2024-06-05T23:59:55` => `2024-06-05T23:59:55`
//
// If the interval is not supported then no change is made to the time.
func Time(t time.Time, interval times.Interval) (reset time.Time) {
	switch interval {
	case times.YEAR:
		reset = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case times.MONTH:
		reset = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	case times.DAY:
		reset = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	case times.HOUR:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	case times.MINUTE:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	case times.SECOND:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	default:
		reset = t
	}
	return
}

// Month switches the time to the first minute of this times month. Shorthand
// wrapper for `Time(t, MONTH)`.
//
// If t is nil, then current time is used
func Month(t *time.Time) (reset time.Time) {
	if t == nil {
		ti := time.Now().UTC()
		t = &ti
	}
	return Time(*t, times.MONTH)
}

// Day switches the time to the first minute of this times day. Shorthand
// wrapper for `Time(t, DAY)`.
//
// If t is nil, then current time is used
func Day(t *time.Time) (reset time.Time) {
	if t == nil {
		ti := time.Now().UTC()
		t = &ti
	}
	return Time(*t, times.DAY)
}
