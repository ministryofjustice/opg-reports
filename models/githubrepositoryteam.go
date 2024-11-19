package models

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
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
	ID int `json:"id,omitempty" db:"id" faker:"-"`
	// Join to the repo
	GitHubRepositoryID int                         `json:"github_repository_id,omitempty" db:"github_repository_id" faker:"-"`
	GitHubRepository   *GitHubRepositoryForeignKey `json:"github_repository,omitempty" db:"github_repository" faker:"-"`
	// Join to the team
	GitHubTeamID int                   `json:"github_team_id,omitempty" db:"github_team_id" faker:"-"`
	GitHubTeam   *GitHubTeamForeignKey `json:"github_team,omitempty" db:"github_team" faker:"-"`
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

// NewGitHubTeams returns a GitHubTeams ([]*GitHubTeam) of all the teams attached to the repo whose org and name is passed.
func NewGitHubTeams(ctx context.Context, client *github.Client, organisation string, repoName string) (teams GitHubTeams, err error) {
	teams = GitHubTeams{}
	opts := &github.ListOptions{PerPage: 100}

	if teamList, _, err := client.Repositories.ListTeams(ctx, organisation, repoName, opts); err == nil {
		for _, team := range teamList {
			var ts = time.Now().UTC().Format(dateformats.Full)
			teams = append(teams, &GitHubTeam{
				Ts:   ts,
				Slug: strings.ToLower(*team.Slug),
			})
		}
	}
	return
}
