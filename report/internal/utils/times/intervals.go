package times

import "time"

// Interval type is used for increamenting / resetting times
type Interval string

// Enum values for Interval
const (
	YEARLY  Interval = "yearly"
	MONTHLY Interval = "monthly"
	DAILY   Interval = "daily"
	HOURLY  Interval = "hourly"
)

// Add increases (or decreases if modifier is negative) the time passed by the
// modifier * interval
func Add(t time.Time, modifier int, interval Interval) (added time.Time) {

	added = t
	switch interval {
	case YEARLY:
		added = t.AddDate(modifier, 0, 0)
	case MONTHLY:
		added = t.AddDate(0, modifier, 0)
	case DAILY:
		added = t.AddDate(0, 0, modifier)
	case HOURLY:
		dur := time.Hour * time.Duration(modifier)
		added = t.Add(dur)
	}
	return
}
