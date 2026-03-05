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

func EndOfDay(t time.Time) (reset time.Time) {
	reset = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.UTC)
	return
}

// FirstDayOfMonth provides the time for the first minute of the first day of the month
//
// If `t` is `2022-12-10T09:00:00` then should return `2022-12-01T00:00:00`
func FirstDayOfMonth(t time.Time) (firstDay time.Time) {
	firstDay = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return
}

// LastDayOfMonth provides the time for the last minute of the last day of the month
//
// If `t` is `2022-12-01T09:00:00` then should return `2022-12-31T23:59:59`
func LastDayOfMonth(t time.Time) (lastDay time.Time) {
	var nextMonth = Add(ResetMonth(t), 1, MONTH)

	lastDay = nextMonth.Add(-1 * time.Second)
	return
}

// Returns the last fully billed month that can be used for data.
//
// Example: Cost data for June will not be available until the 15th
// July, so on the 14th this will give you May, but 15th will
// be June
func BillingMonth(t time.Time, billingDay int) (billing time.Time) {

	if t.Day() < billingDay {
		billing = time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, time.UTC).Add(-1 * time.Second)
	} else {
		billing = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-1 * time.Second)
	}
	return
}

func LastBillingMonth(billingDay int) (billing time.Time) {
	return BillingMonth(time.Now().UTC(), billingDay)
}
