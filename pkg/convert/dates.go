package convert

import (
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/pkg/consts"
)

// DateAddYear adds `mod` number of years onto the time
func DateAddYear(d time.Time, mod int) time.Time {
	return d.AddDate(mod, 0, 0)
}

// DateAddMonth adds `mod` number of months onto the time
func DateAddMonth(d time.Time, mod int) time.Time {
	return d.AddDate(0, mod, 0)
}

// DateAddDay adds `mod` number of days onto the time
func DateAddDay(d time.Time, mod int) time.Time {
	return d.AddDate(0, 0, mod)
}

// ResetMonth resets the day of the month to the 1st and zeros the time
func DateResetYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
}

// ResetMonth resets the day of the month to the 1st and zeros the time
func DateResetMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ResetDay resets time of the day to zeros
func DateResetDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// GetDateFormat will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func GetDateFormat(value string) string {
	max := len(consts.DateFormat)
	l := len(value)
	if l > max {
		return consts.DateFormat
	}
	f := consts.DateFormat[:l]
	slog.Debug("[dates] GetFormat", slog.String("format", f))
	return f
}

// ToTime will try to convert the string passed into a time.Time
func ToTime(s string) (t time.Time, err error) {
	layout := GetDateFormat(s)
	t, err = time.Parse(layout, s)
	return
}
