package times

import "time"

// Months returns a slice of times from start, upto and including end
// adding a month on at a time.
//
// `start` is reset to the first day of a month
// `end` is reset to the last day of the month
func Months(start time.Time, end time.Time) (months []time.Time) {
	months = []time.Time{}
	start = FirstDayOfMonth(start)
	end = LastDayOfMonth(end)

	for d := start; d.After(end) == false; d = Add(d, 1, MONTH) {
		months = append(months, d)
	}

	return
}

// Days returns all days between the start & end times
func Days(start time.Time, end time.Time) (days []time.Time) {
	days = []time.Time{}
	for d := start; d.After(end) == false; d = Add(d, 1, DAY) {
		days = append(days, d)
	}
	return
}

// provides all days, but with n days intervals
func DaysN(start time.Time, end time.Time, n int) (days []time.Time) {
	var current time.Time
	var last time.Time
	days = []time.Time{}

	start = ResetDay(start)
	end = EndOfDay(end)

	current = start
	for current.Before(end) {
		last = current
		days = append(days, current)
		current = Add(current, n, DAY)
	}
	if end.After(last) {
		days = append(days, ResetDay(end))
	}

	return
}
