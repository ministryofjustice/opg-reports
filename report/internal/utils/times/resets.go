package times

import "time"

// ResetYear switches the time to the first minute of the times passed year:
//
//	`2024-06-05T23:59` => `2024-01-01T00:00`
func ResetYear(t time.Time) (reset time.Time) {
	return Reset(t, YEAR)
}

// ResetMonth switches the time to the first minute of this times month:
//
//	`2024-06-05T23:59` => `2024-06-01T00:00`
func ResetMonth(t time.Time) (reset time.Time) {
	return Reset(t, MONTH)
}

// ResetDay switches the time to the first minute of this times day:
//
//	`2024-06-05T23:59` => `2024-06-05T00:00`
func ResetDay(t time.Time) (reset time.Time) {
	return Reset(t, DAY)
}

// ResetHour switches the time to the first minute of this times hour:
//
//	`2024-06-05T23:59` => `2024-06-05T23:00`
func ResetHour(t time.Time) (reset time.Time) {
	return Reset(t, HOUR)
}

// Reset changes the time to be the start of the interval type asked for to avoid
// errors with date addition (rounding when adding days to month etc).
//
//	HOURLY:  `2024-06-05T23:59` => `2024-06-05T23:00`
//	DAILY: 	 `2024-06-05T23:59` => `2024-06-05T00:00`
//	MONTHLY: `2024-06-05T23:59` => `2024-06-01T00:00`
//	YEARLY:	 `2024-06-05T23:59` => `2024-01-01T00:00`
func Reset(t time.Time, interval Interval) (reset time.Time) {
	switch interval {
	case YEAR:
		reset = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case MONTH:
		reset = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	case DAY:
		reset = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	case HOUR:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	}
	return
}
