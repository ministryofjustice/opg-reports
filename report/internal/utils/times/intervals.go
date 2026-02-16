package times

import "time"

// Interval type is used for increamenting / resetting times
type Interval string

// Enum values for Interval
const (
	YEAR  Interval = "yearly"
	MONTH Interval = "monthly"
	DAY   Interval = "daily"
	HOUR  Interval = "hourly"
)

// Add increases (or decreases if modifier is negative) the time passed by the
// modifier * interval
func Add(t time.Time, modifier int, interval Interval) (added time.Time) {

	added = t
	switch interval {
	case YEAR:
		added = t.AddDate(modifier, 0, 0)
	case MONTH:
		added = t.AddDate(0, modifier, 0)
	case DAY:
		added = t.AddDate(0, 0, modifier)
	case HOUR:
		dur := time.Hour * time.Duration(modifier)
		added = t.Add(dur)
	}
	return
}

func Ago(t time.Time, ago int, interval Interval) time.Time {
	return Add(t, 0-ago, interval)
}

// Yesterday returns the start of yesterday as a time
func Yesterday() (t time.Time) {
	return Add(ResetDay(time.Now().UTC()), -1, DAY)
}

// Today returns todays time - reset to midnight
func Today() (t time.Time) {
	return ResetDay(time.Now().UTC())
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
