package timex

import "time"

// Add is used to add or remove time from an existing date.
//
// To remove time use a negative quantity
func Add(t time.Time, interval Interval, quantity int) (updated time.Time) {
	updated = t
	switch interval {
	case YEAR:
		updated = t.AddDate(quantity, 0, 0)
	case MONTH:
		updated = t.AddDate(0, quantity, 0)
	case DAY:
		updated = t.AddDate(0, 0, quantity)
	case HOUR:
		dur := time.Hour * time.Duration(quantity)
		updated = t.Add(dur)
	case MINUTE:
		dur := time.Minute * time.Duration(quantity)
		updated = t.Add(dur)
	case SECOND:
		dur := time.Second * time.Duration(quantity)
		updated = t.Add(dur)
	}

	return
}

// End returns the last nsec of the time
func End(t time.Time, interval Interval) (updated time.Time) {
	updated = t
	switch interval {
	case YEAR:
		updated = time.Date(t.Year(), 12, 31, 23, 59, 59, 59, time.UTC)
	case MONTH:
		nextMonth := Add(Reset(t, MONTH), MONTH, 1)
		lastDay := Add(nextMonth, DAY, -1)
		updated = time.Date(t.Year(), t.Month(), lastDay.Day(), 23, 59, 59, 59, time.UTC)
	case DAY:
		updated = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.UTC)
	case HOUR:
		updated = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 59, time.UTC)
	case MINUTE:
		updated = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 59, time.UTC)
	case SECOND:
		updated = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 59, time.UTC)
	}
	return
}

// Reset changes the time to be the start of the interval type asked for to avoid
// errors with date addition (rounding when adding days to month etc).
//
//	YEAR:	`2024-06-05T23:59:59` => `2024-01-01T00:00:00`
//	MONTH:	`2024-06-05T23:59:59` => `2024-06-01T00:00:00`
//	DAY: 	`2024-06-05T23:59:59` => `2024-06-05T00:00:00`
//	HOUR:	`2024-06-05T23:59:59` => `2024-06-05T23:00:00`
//	MINUTE:	`2024-06-05T23:59:59` => `2024-06-05T23:59:00`
//	SECOND:	`2024-06-05T23:59:59` => `2024-06-05T23:59:59`
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
	case MINUTE:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	case SECOND:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	}
	return
}
