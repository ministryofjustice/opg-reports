package dates

import (
	"slices"
	"time"
)

type Interval string

const (
	DAY   Interval = "DAY"
	MONTH Interval = "MONTH"
)

func addDay(d time.Time) time.Time {
	return d.AddDate(0, 0, 1)
}

func addMonth(d time.Time) time.Time {
	return d.AddDate(0, 1, 0)
}

// ResetMonth resets the day of the month to the 1st and zeros the time
func ResetMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ResetDay resets time of the day to zeros
func ResetDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// Range provides a slice of times between the 2 dates with a variable interval
// allowing it to iterate by either day or month
func Range(start time.Time, end time.Time, interval Interval) []time.Time {
	var intF = addDay

	if interval == DAY {
		start = ResetDay(start)
		end = ResetDay(end)
	} else if interval == MONTH {
		start = ResetMonth(start)
		end = ResetMonth(end)
		intF = addMonth
	}
	times := []time.Time{}
	for d := start; d.After(end) == false; d = intF(d) {
		times = append(times, d)
	}

	return times
}

func Strings(dates []time.Time, format string) (strs []string) {
	strs = []string{}
	for _, d := range dates {
		strs = append(strs, d.Format(format))
	}
	return
}

// MaxTime finds the max time in slice of times
func MaxTime(times []time.Time) time.Time {
	max := slices.MaxFunc(times, func(a time.Time, b time.Time) int {
		if a.After(b) {
			return 1
		}
		if a.Before(b) {
			return -1
		}
		return 0
	})
	return max
}

func YearToBillingDate(when time.Time, billingDay int) (s time.Time, e time.Time) {
	e = BillingEndDate(when, billingDay)
	s = time.Date(e.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	return
}

func BillingEndDate(when time.Time, billingDay int) (m time.Time) {
	if when.Day() < billingDay {
		m = ResetMonth(when).AddDate(0, -1, 0)
	} else {
		m = ResetMonth(when).AddDate(0, 0, 0)
	}
	return
}

// BillingDates returns the first day of start month and
// the last day of the billing month
// - the end date date is midnight the day after to be used in less than (<) queries
func BillingDates(when time.Time, billingDay int, months int) (s time.Time, e time.Time) {

	e = BillingEndDate(when, billingDay)

	diff := (0 - (months))
	s = e.AddDate(0, diff, 1)
	// always reset the day of the month
	s = ResetMonth(s)

	return
}
