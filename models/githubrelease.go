package models

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

type GitHubReleaseType string

const (
	GitHubWorkflowRelease GitHubReleaseType = "workflow_run"
	GitHubPRMergeRelease  GitHubReleaseType = "pull_request"
)

// GitHubRelease tracks workflow runs or merge requests that act as
// a production release
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type GitHubRelease struct {
	ID                 int                         `json:"id,omitempty" db:"id" faker:"-"`
	Ts                 string                      `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Name               string                      `json:"name,omitempty" db:"name" faker:"word"`
	Count              int                         `json:"count,omitempty" db:"count" faker:"oneof: 1"`
	RelaseType         GitHubReleaseType           `json:"release_type,omitempty" db:"release_type" faker:"oneof:workflow_run, pull_request" enum:"oneof:workflow_run, pull_request"`
	SourceURL          string                      `json:"source_url" db:"source_url" faker:"uri"`
	Date               string                      `json:"date,omitempty" db:"date" faker:"date_string"`
	GitHubRepositoryID int                         `json:"github_repository_id,omitempty" db:"github_repository_id" faker:"-"`
	GitHubRepository   *GitHubRepositoryForeignKey `json:"github_repository,omitempty" db:"github_repository" faker:"-"`
	GitHubTeams        GitHubTeams                 `json:"github_teams,omitempty" db:"github_teams" faker:"-"`
}

// TableName returns named table for GitHubRelease - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubRelease) TableName() string {
	return "github_releases"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRelease) Columns() map[string]string {
	return map[string]string{
		"id":                   "INTEGER PRIMARY KEY",
		"ts":                   "TEXT NOT NULL",
		"name":                 "TEXT NOT NULL",
		"count":                "INTEGER NOT NULL",
		"release_type":         "TEXT NOT NULL",
		"source_url":           "TEXT NOT NULL",
		"date":                 "TEXT NOT NULL",
		"github_repository_id": "INTEGER NOT NULL",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRelease) Indexes() map[string][]string {
	return map[string][]string{
		"ghr_date_idx":      {"date"},
		"ghr_repo_idx":      {"github_repository_id"},
		"ghr_date_repo_idx": {"date", "github_repository_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRelease) InsertColumns() []string {
	return []string{
		"ts",
		"name",
		"count",
		"release_type",
		"source_url",
		"date",
		"github_repository_id",
	}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRelease) GetID() int {
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
func (self *GitHubRelease) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubRelease) New() dbs.Cloneable {
	return &GitHubRelease{}
}