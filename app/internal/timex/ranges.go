package timex

import (
	"time"
)

// Range creates a list of times between the start & end time with each entry being increased
// by (increament * interval)
func Range(start time.Time, end time.Time, interval Interval, increment int) (dates []time.Time) {
	var (
		current time.Time
		latest  time.Time
	)
	// if start is after the end date, return empty
	if start.After(end) {
		return
	}
	// setup
	dates = []time.Time{}
	// start of the period
	start = Reset(start, interval)
	// last momment of the last period
	end = End(end, interval)
	current = start
	// keep looping until the current time is after the end date
	for current.Before(end) {
		latest = current                            // track the last date to use
		dates = append(dates, current)              // add to the list
		current = Add(current, interval, increment) // increment the date
	}
	// if the increament jumped the date beyond the end, we do want to capture that
	// so add it directly
	end = Reset(end, interval)
	if end.After(latest) {
		dates = append(dates, end)
	}

	return
}
