package dates

import "time"

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

func ResetMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func ResetDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

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
