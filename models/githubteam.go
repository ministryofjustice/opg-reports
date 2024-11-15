package models

import "github.com/ministryofjustice/opg-reports/internal/dbs"

// GitHubTeam captures just the team names attached to our
// repositories etc that
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type GitHubTeam struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"`
	Name string `json:"name,omitempty" db:"name" faker:"unique, oneof:opg,opg-webops,opg-sirius,opg-use,opg-make,foo,bar"`
}

// TableName returns named table for GitHubTeam - GitHubTeams
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubTeam) TableName() string {
	return "github_teams"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubTeam) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "name": "TEXT NOT NULL UNIQUE"}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubTeam) Indexes() map[string][]string {
	return map[string][]string{
		"gh_team_idx": {"name"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeam) InsertColumns() []string {
	return []string{"name"}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeam) GetID() int {
	return self.ID
}

// SetID allows setting the ID of this row - normally used within the insert calls
// to update the original data passed in with the new id
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeam) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubTeam) New() dbs.Cloneable {
	return &GitHubTeam{}
}
