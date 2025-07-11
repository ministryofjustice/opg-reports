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

type Granularity string

const (
	GranularityYear  Granularity = "yearly"
	GranularityMonth Granularity = "monthly"
	GranularityDay   Granularity = "daily"
)

var GRANULARITY_TO_FORMAT = map[string]string{
	string(GranularityYear):  "%Y",
	string(GranularityMonth): "%Y-%m",
	string(GranularityDay):   "%Y-%m-%d",
}

// TimeInterval
type TimeInterval string

// Enum values for TimeInterval
const (
	TimeIntervalYear  TimeInterval = "yearly"
	TimeIntervalMonth TimeInterval = "monthly"
	TimeIntervalDay   TimeInterval = "daily"
	TimeIntervalHour  TimeInterval = "hourly"
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

// LastDayOfMonth provides the time for the last minute of the last day of the month
//
// If `t` is `2022-12-01T09:00:00` then should return `2022-12-31T23:59:59`
func LastDayOfMonth(t time.Time) (lastDay time.Time) {
	var nextMonth = t.UTC().AddDate(0, 1, 0)
	var reset = TimeReset(nextMonth, TimeIntervalMonth)

	// fmt.Printf("	month: [%s] reset [%s]\n", nextMonth.Format(DATE_FORMATS.Full), reset.Format(DATE_FORMATS.Full))

	lastDay = reset.Add(-1 * time.Second)
	return
}

// Month returns a YYYY-MM-DD formatted string of the first data of a month. The modifier adjusts
// the month (plus or minus)
func Month(monthModifier int) (m string) {
	m = TimeReset(time.Now().UTC().AddDate(0, monthModifier, 0), TimeIntervalMonth).Format(DATE_FORMATS.YMD)
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

// timeAdd calls `AddDate` on `date` param and increments year / month / day by the quantity
// requested.
//
// Used to generate loop condition for incrementing between two dates
func timeAdd(date time.Time, quantity int, interval TimeInterval) (t time.Time) {
	switch interval {
	case TimeIntervalYear:
		t = date.AddDate(quantity, 0, 0)
	case TimeIntervalMonth:
		t = date.AddDate(0, quantity, 0)
	case TimeIntervalDay:
		t = date.AddDate(0, 0, quantity)
	}
	return
}

// Times generates all the time.Time values between the start and end times passed, incrementing by the
// interval.
func Times(start time.Time, end time.Time, interval TimeInterval, increment int) (times []time.Time) {
	times = []time.Time{}
	start = TimeReset(start, interval)
	end = TimeReset(end, interval)

	for d := start; d.After(end) == false; d = timeAdd(d, increment, interval) {
		times = append(times, d)
	}
	return
}

// Months is an opinated wrapper around Times that fixes the interval to month and returns
// string formatted values instead of time.Time
//
// Generally used for display headers etc
func Months(start string, end string) (months []string) {
	var err error
	var startD, endD time.Time
	var times []time.Time

	startD, err = StringToTime(start)
	if err != nil {
		return
	}
	endD, err = StringToTime(end)
	if err != nil {
		return
	}
	times = Times(startD, endD, TimeIntervalMonth, 1)
	for _, item := range times {
		months = append(months, item.Format(DATE_FORMATS.YM))
	}

	return
}
