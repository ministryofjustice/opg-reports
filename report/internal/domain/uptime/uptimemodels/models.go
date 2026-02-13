package uptimemodels

// Uptime struct for importing / seeding the db
type Uptime struct {
	ID          int    `json:"id,omitempty" db:"id"`
	Date        string `json:"date,omitempty" db:"date" `              // The data the uptime value was for
	Average     string `json:"average,omitempty" db:"average" `        // uptime average as a percentage
	Granularity string `json:"granularity,omitempty" db:"granularity"` // the time period in seconds used for this metric
	AccountID   string `json:"account_id" db:"account_id"`             // the account id reference
}

type UptimeData struct {
	Date     string  `json:"date" db:"date" `       // The month for this cost
	Average  float64 `json:"average" db:"average" ` // uptime average as a percentage
	TeamName string  `json:"team" db:"team"`        // the team this costs is attached to
}
