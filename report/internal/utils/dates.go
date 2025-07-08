package utils

import (
	"time"
)

type dateFormats struct {
	Full string
	Y    string
	YM   string
	YMD  string
}

var DATE_FORMATS = &dateFormats{
	Full: time.RFC3339,
	Y:    "2006",
	YM:   "2006-01",
	YMD:  "2006-01-02",
}

var GRANULARITY_TO_FORMAT = map[string]string{
	"year":  "%Y",
	"month": "%Y-%m",
	"day":   "%Y-%m-%d",
}

// TimeInterval
type TimeInterval string

// Enum values for TimeInterval
const (
	TimeIntervalYear  TimeInterval = "year"
	TimeIntervalMonth TimeInterval = "month"
	TimeIntervalDay   TimeInterval = "day"
	TimeIntervalHour  TimeInterval = "hour"
)

func (TimeInterval) Values() []TimeInterval {
	return []TimeInterval{
		TimeIntervalYear,
		TimeIntervalMonth,
		TimeIntervalDay,
		TimeIntervalHour,
	}
}

// GetTimeFormat will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func GetDateTimeFormat(value string) string {
	var layout = DATE_FORMATS.Full
	max := len(layout)
	l := len(value)
	if l > max {
		return layout
	}
	f := layout[:l]
	return f
}

// StringToTime converts string to time.Time via time.Parse and
// GetDateTimeFormat.
func StringToTime(str string) (t time.Time, err error) {
	t, err = time.Parse(GetDateTimeFormat(str), str)
	return
}

// StringToTimeResetAsString converts the string `s` to a time, resets that time and
// then returns a YYYY-MM-DD formatted string
func StringToTimeResetAsString(s string, interval TimeInterval) (t string) {

	if dt, err := StringToTime(s); err == nil {
		reset := TimeReset(dt, interval)
		t = reset.Format(DATE_FORMATS.YMD)
	}
	return
}

func StringToTimeReset(s string, interval TimeInterval) (t time.Time) {
	if dt, err := StringToTime(s); err == nil {
		t = TimeReset(dt, interval)
	}
	return
}

// TimeReset changes the time to be the start of the interval type asked for to avoid
// errors with date addition (rounding when adding days to month etc).
//
//	`2024-06-05T23:00` with interval of `day` => `2024-06-05T00:00`
//	`2024-06-05T23:00` with interval of `month` => `2024-06-01T00:00`
//	`2024-06-05T23:00` with interval of `year` => `2024-01-01T00:00`
func TimeReset(t time.Time, interval TimeInterval) (reset time.Time) {
	switch interval {
	case TimeIntervalYear:
		reset = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case TimeIntervalMonth:
		reset = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	case TimeIntervalDay:
		reset = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	case TimeIntervalHour:
		reset = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	}
	return
}

// Month returns a YYYY-MM-DD formatted string of the first data of a month. The modifier adjusts
// the month (plus or minus)
func Month(modifier int) (m string) {
	m = TimeReset(time.Now().UTC().AddDate(0, modifier, 0), TimeIntervalMonth).Format(DATE_FORMATS.YMD)
	return
}
