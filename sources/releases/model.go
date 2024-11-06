package releases

import "fmt"

type Release struct {
	ID         int    `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."` // ID is a generated primary key
	Ts         string `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`                             // TS is timestamp when the record was created
	Repository string `json:"repository,omitempty" db:"repository"`
	Name       string `json:"name,omitempty" db:"name"`
	Source     string `json:"source,omitempty" db:"source" doc:"url of source event."`
	Date       string `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this release happened."`
	Type       string `json:"type,omitempty" db:"type" faker:"oneof: workflow, merge" enum:"workflow,merge" doc:"tracks if this was a workflow run or a merge"`
	Teams      string `json:"teams" db:"teams" faker:"oneof: #unitA#, #unitB#, #unitC#"`
}

// UID
// Record interface
func (self *Release) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}

// TDate
// Transformable interface
func (self *Release) TDate() string {
	return self.Date
}

// TValue
// Transformable interface
// Always 1 so it increments per period
func (self *Release) TValue() string {
	return "1"
}
