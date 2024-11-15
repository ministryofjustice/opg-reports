package models

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// GitHubRepositoryGitHubTeam represents the join between a github repo and
// a git hub team providing structs for both sides of the many to many join
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
type GitHubRepositoryGitHubTeam struct {
	ID                 int `json:"id,omitempty" db:"id" faker:"-"`
	GitHubRepositoryID int `json:"github_repository_id,omitempty" db:"github_repository_id" faker:"-"`
	GitHubTeamID       int `json:"github_team_id,omitempty" db:"github_team_id" faker:"-"`
}

// TableName returns named table for GitHubRepositoryGitHubTeam - GitHubRepositoryGitHubTeams
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubRepositoryGitHubTeam) TableName() string {
	return "github_repositories_github_teams"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepositoryGitHubTeam) Columns() map[string]string {
	return map[string]string{
		"id":                   "INTEGER PRIMARY KEY",
		"github_repository_id": "INTEGER NOT NULL",
		"github_team_id":       "INTEGER NOT NULL",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepositoryGitHubTeam) Indexes() map[string][]string {
	return map[string][]string{
		"ghgt_t_idx": {"github_team_id"},
		"ghgt_r_idx": {"github_repository_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRepositoryGitHubTeam) InsertColumns() []string {
	return []string{"github_repository_id", "github_team_id"}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRepositoryGitHubTeam) GetID() int {
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
func (self *GitHubRepositoryGitHubTeam) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubRepositoryGitHubTeam) New() dbs.Cloneable {
	return &GitHubRepositoryGitHubTeam{}
}

// GitHubTeams is to be used on the struct that needs to pull in
// the github teams via a many to many join select statement and provides
// the Scan method so sqlx will handle the result correctly
//
// Interfaces:
//   - sql.Scanner
type GitHubRepositories []*GitHubRepository

// Scan converts the json aggregate result from a select statement into
// a series of GitHubTeams attached to the main struct and will be called
// directly by sqlx
//
// Interfaces:
//   - sql.Scanner
func (self *GitHubRepositories) Scan(src interface{}) (err error) {
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
