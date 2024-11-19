package models

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

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
	ID    int    `json:"id,omitempty" db:"id" faker:"-"`
	Ts    string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Slug  string `json:"slug,omitempty" db:"slug" faker:"unique, oneof:opg,opg-webops,opg-sirius,opg-use,opg-make,foo,bar"`
	Units Units  `json:"units,omitempty" db:"units" faker:"-"`
	// Many to many join with the repositories
	GitHubRepositories GitHubRepositories `json:"github_repositories,omitempty" db:"github_repositories" faker:"-"`
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
	return map[string]string{
		"id":   "INTEGER PRIMARY KEY",
		"ts":   "TEXT NOT NULL",
		"slug": "TEXT NOT NULL UNIQUE",
	}
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
		"gh_team_idx": {"slug"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubTeam) InsertColumns() []string {
	return []string{
		"ts",
		"slug",
	}
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

// GitHubTeams is to be used on the struct that needs to pull in
// the github teams via a many to many join select statement and provides
// the Scan method so sqlx will handle the result correctly
//
// Interfaces:
//   - sql.Scanner
type GitHubTeams []*GitHubTeam

// Scan converts the json aggregate result from a select statement into
// a series of GitHubTeams attached to the main struct and will be called
// directly by sqlx
//
// Interfaces:
//   - sql.Scanner
func (self *GitHubTeams) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = structs.Unmarshal(src.([]byte), self)
	case string:
		err = structs.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// GitHubTeamForeignKey is to be used on the struct that needs to pull in
// the team via one to many join (being used on the `one` side).
//
// To swap a GitHubTeam to a GitHubTeamForeignKey:
//
//	var join = GitHubRepositoryForeignKey(&GitHubTeam{})
//
// or
//
//	var join = (*GitHubRepositoryForeignKey)(&GitHubTeam)
//
// Interfaces:
//   - sql.Scanner
type GitHubTeamForeignKey GitHubTeam

// Scan converts the json aggregate result from a select statement into
// attached to the main struct and will be called directly by sqlx
//
// Interfaces:
//   - sql.Scanner
func (self *GitHubTeamForeignKey) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = structs.Unmarshal(src.([]byte), self)
	case string:
		err = structs.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}
