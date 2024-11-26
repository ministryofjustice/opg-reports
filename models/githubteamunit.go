package models

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// GitHubTeamUnit represents the join between a unit and ad github team
// providing structs for both sides of the many to many join
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
type GitHubTeamUnit struct {
	ID           int `json:"id,omitempty" db:"id" faker:"-"`
	GitHubTeamID int `json:"github_team_id,omitempty" db:"github_team_id" faker:"-"`
	UnitID       int `json:"unit_id,omitempty" db:"unit_id" faker:"-"`
}

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *GitHubTeamUnit) UniqueValue() string {
	return fmt.Sprintf("%d,%d", self.GitHubTeamID, self.UnitID)
}

// Interfaces:
//   - dbs.Insertable
func (self *GitHubTeamUnit) UniqueField() string {
	return "github_team_id, unit_id"
}

func (self *GitHubTeamUnit) UpsertUpdate() string {
	return "github_team_id=excluded.github_team_id, unit_id=excluded.unit_id"
}

// TableName returns named table for GitHubTeamUnit - GitHubTeamUnits
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubTeamUnit) TableName() string {
	return "github_teams_units"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubTeamUnit) Columns() map[string]string {
	return map[string]string{
		"id":             "INTEGER PRIMARY KEY",
		"github_team_id": "INTEGER NOT NULL",
		"unit_id":        "INTEGER NOT NULL",
		"UNIQUE":         "(github_team_id, unit_id)",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubTeamUnit) Indexes() map[string][]string {
	return map[string][]string{
		"ghu_gh_idx": {"github_team_id"},
		"ghu_u_idx":  {"unit_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeamUnit) InsertColumns() []string {
	return []string{"github_team_id", "unit_id"}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeamUnit) GetID() int {
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
func (self *GitHubTeamUnit) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubTeamUnit) New() dbs.Cloneable {
	return &GitHubTeamUnit{}
}
