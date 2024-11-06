package releases

import (
	"fmt"
	"strconv"

	"github.com/ministryofjustice/opg-reports/pkg/record"
)

type Release struct {
	ID         int     `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."` // ID is a generated primary key
	Ts         string  `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`                             // TS is timestamp when the record was created
	Repository string  `json:"repository,omitempty" db:"repository" faker:"word"`
	Name       string  `json:"name,omitempty" db:"name" faker:"word"`
	Source     string  `json:"source,omitempty" db:"source" doc:"url of source event." faker:"uri"`
	Type       string  `json:"type,omitempty" db:"type" faker:"oneof: workflow_run, pull_request" enum:"workflow,merge" doc:"tracks if this was a workflow run or a merge"`
	Date       string  `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this release happened."`
	Count      int     `json:"count,omitempty" db:"count" faker:"oneof: 1" enum:"1" doc:"Number of releases"`
	TeamList   []*Team `json:"team_list,omitempty" db:":teams" doc:"pulled from a many to many join table"`
}

// New
// Record interface
func (self *Release) New() record.Record {
	return &Release{}
}

// UID
// Record interface
func (self *Release) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}

// SetID
// Record interface
func (self *Release) SetID(id int) {
	self.ID = id
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
	return strconv.Itoa(self.Count)
}

type Team struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`
	Name string `json:"name,omitempty" db:"name" faker:"oneof: A, B, C, D"`
}

// New
// Record interface
func (self *Team) New() record.Record {
	return &Team{}
}

// UID
// Record interface
func (self *Team) UID() string {
	return fmt.Sprintf("%s-%d", "team", self.ID)
}

// SetID
// Record interface
func (self *Team) SetID(id int) {
	self.ID = id
}

// Join is many to many table between both team and releases
type Join struct {
	ID        int `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`
	TeamID    int `json:"team_id,omitempty" db:"team_id"`
	ReleaseID int `json:"release_id,omitempty" db:"release_id"`
}

// New
// Record interface
func (self *Join) New() record.Record {
	return &Join{}
}

// UID
// Record interface
func (self *Join) UID() string {
	return fmt.Sprintf("%s-%d", "join", self.ID)
}

// SetID
// Record interface
func (self *Join) SetID(id int) {
	self.ID = id
}
