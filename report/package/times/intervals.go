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
