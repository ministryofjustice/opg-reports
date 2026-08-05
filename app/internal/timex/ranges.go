package timex

import (
	"log/slog"
	"opg-reports/app/internal/logx"
	"time"
)

// Range creates a list of times between the start & end time with each entry being increased
// by (increament * interval)
func Range(start time.Time, end time.Time, interval Interval, increment int) (dates []time.Time) {
	var (
		current time.Time
		latest  time.Time
		lg      *slog.Logger = logx.Default().With(
			"start", start.String(),
			"end", end.String(),
			"interval", string(interval),
			"increment", increment,
		)
	)
	// if start is after the end date, return empty
	if start.After(end) {
		lg.Error("start time is after end time.")
		return
	}
	// setup
	dates = []time.Time{}

	// start of the period
	lg.Debug("resetting start time.")
	start = Reset(start, interval)
	// last momment of the last period
	lg.Debug("getting end of end time.")
	end = End(end, interval)
	current = start
	// keep looping until the current time is after the end date
	for current.Before(end) {
		latest = current                            // track the last date to use
		dates = append(dates, current)              // add to the list
		current = Add(current, interval, increment) // increment the date
		lg.Debug("added to range", "added", latest.String())
	}
	// if the increament jumped the date beyond the end, we do want to capture that
	// so add it directly
	lg.Debug("checking end time against last item as interval might skip.")
	end = Reset(end, interval)
	if end.After(latest) {
		dates = append(dates, end)
	}

	return
}
