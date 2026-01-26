package times

// Interval type is used for increamenting / resetting times
type Interval string

// Enum values for Interval
const (
	YEARLY  Interval = "yearly"
	MONTHLY Interval = "monthly"
	DAILY   Interval = "daily"
	HOURLY  Interval = "hourly"
)
