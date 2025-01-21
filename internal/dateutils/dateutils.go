package dateutils

import (
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
)

// Format will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func Format(value string) string {
	var layout = dateformats.Full
	max := len(layout)
	l := len(value)
	if l > max {
		return layout
	}
	f := layout[:l]
	slog.Debug("[dates] Format", slog.String("format", f))
	return f
}

// Time converts string to a time.Time using time.Parse
// and Format
func Time(s string) (t time.Time, err error) {
	layout := Format(s)
	t, err = time.Parse(layout, s)
	return
}

// Reformat will convert a date string from its current format into the new layout
// passed along.
//
// It will try and parse the date using `time.RFC3339` as the base format, so if
// its current form does not match this it will fail and return empty string.
//
// In that case, use `Convert` instead passing along both old and new formats.
func Reformat(s string, layout string) (date string) {
	if t, err := time.Parse(Format(s), s); err == nil {
		date = t.Format(layout)
	}
	return
}

// Convert takes a string versio of a date and using `oldlayout` will parse that
// into a time.Time and then proced to conver that time into a string in the format
// specified by `newlayout`
//
// Use this when your original date string is not based on `time.RFC3339`.
func Convert(s string, oldlayout string, newlayout string) (date string) {
	if t, err := time.Parse(oldlayout, s); err == nil {
		date = t.Format(newlayout)
	}
	return
}

// Reset changes the time to be the start of the interval type asked for to avoid
// errors with date addition (rounding when adding days to month etc).
//
//	`2024-06-05T23:00` with interval of `day` => `2024-06-05T00:00`
//	`2024-06-05T23:00` with interval of `month` => `2024-06-01T00:00`
//	`2024-06-05T23:00` with interval of `year` => `2024-01-01T00:00`
func Reset(date time.Time, interval dateintervals.Interval) (t time.Time) {
	switch interval {
	case dateintervals.Year:
		t = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case dateintervals.Month:
		t = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	case dateintervals.Day:
		t = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	}
	return
}

// Add calls `AddDate` on `date` param and increments year / month / day by the quantity
// requested.
//
// Used to generate loop condition for incrementing between two dates
func Add(date time.Time, quantity int, interval dateintervals.Interval) (t time.Time) {
	switch interval {
	case dateintervals.Year:
		t = date.AddDate(quantity, 0, 0)
	case dateintervals.Month:
		t = date.AddDate(0, quantity, 0)
	case dateintervals.Day:
		t = date.AddDate(0, 0, quantity)
	}
	return
}

// Dates returns a slice of strings of all intervals between the start and end date passed.
// These are formatted based on the interval type used -see `dateintervals.Format`
//
// Calls `Times` and then uses that dateset to generate the strings
func Dates(start time.Time, end time.Time, interval dateintervals.Interval) (dates []string) {
	var times = Times(start, end, interval)
	var layout = dateintervals.Format(interval)

	dates = []string{}
	for _, t := range times {
		dates = append(dates, t.Format(layout))
	}

	return
}

// Range returns a list of times using NOW and the start / end ints to adjust range
func Range(start int, end int, interval dateintervals.Interval) (times []time.Time) {
	var now = time.Now().UTC()
	var s, e time.Time

	s = Add(now, start, interval)
	e = Add(now, end, interval)
	return Times(s, e, interval)

}

// Times generates all the time.Time values between the start and end times passed, incrementing by the
// interval.
func Times(start time.Time, end time.Time, interval dateintervals.Interval) (times []time.Time) {
	times = []time.Time{}
	start = Reset(start, interval)
	end = Reset(end, interval)

	for d := start; d.Equal(end) == false && d.After(end) == false; d = Add(d, 1, interval) {
		times = append(times, d)
	}
	return
}

// TimesI is like Times, but allows custom interval count - allowing counting every X days / months
func TimesI(start time.Time, end time.Time, interval dateintervals.Interval, inc int) (times []time.Time) {
	times = []time.Time{}
	start = Reset(start, interval)
	end = Reset(end, interval)

	for d := start; d.Equal(end) == false && d.After(end) == false; d = Add(d, inc, interval) {
		times = append(times, d)
	}
	return
}

// CountInRange returns how many of interval exists between start and end
func CountInRange(start time.Time, end time.Time, interval dateintervals.Interval) (count int) {
	count = 0

	for d := start; d.Equal(end) == false && d.After(end) == false; d = Add(d, 1, interval) {
		count += 1
	}

	return
}
