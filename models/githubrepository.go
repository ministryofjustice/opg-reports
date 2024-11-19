package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/bools"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// GitHubRepository stores all the infor we want locally about our
// current github repos
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type GitHubRepository struct {
	ID             int         `json:"id,omitempty" db:"id" faker:"-"`
	Ts             string      `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Owner          string      `json:"owner,omitempty" db:"owner" faker:"oneof: ministryofjusice"`
	Name           string      `json:"name,omitempty" db:"name" faker:"word"`
	FullName       string      `json:"full_name,omitempty" db:"full_name" faker:"unique"`
	CreatedAt      string      `json:"created_at,omitempty" db:"created_at" faker:"date_string"`
	DefaultBranch  string      `json:"default_branch,omitempty" db:"default_branch" faker:"oneof: main, master"`
	Archived       uint8       `json:"archived,omitempty" db:"archived" faker:"oneof: 0, 1"`
	Private        uint8       `json:"private,omitempty" db:"private" faker:"oneof: 0, 1"`
	License        string      `json:"license,omitempty" db:"license" faker:"oneof: MIT, GPL"`
	LastCommitDate string      `json:"last_commit_date,omitempty" db:"last_commit_date" faker:"date_string"`
	GitHubTeams    GitHubTeams `json:"github_teams,omitempty" db:"github_teams" faker:"-"`
}

// TableName returns named table for GitHubRepository - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubRepository) TableName() string {
	return "github_repositories"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepository) Columns() map[string]string {
	return map[string]string{
		"id":               "INTEGER PRIMARY KEY",
		"ts":               "TEXT NOT NULL",
		"owner":            "TEXT NOT NULL",
		"name":             "TEXT NOT NULL",
		"full_name":        "TEXT NOT NULL UNIQUE",
		"created_at":       "TEXT NOT NULL",
		"default_branch":   "TEXT NOT NULL",
		"archived":         "INTEGER NOT NULL DEFAULT 0",
		"private":          "INTEGER NOT NULL DEFAULT 0",
		"license":          "TEXT NOT NULL DEFAULT ''",
		"last_commit_date": "TEXT NOT NULL",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepository) Indexes() map[string][]string {
	return map[string][]string{
		"repo_name_idx":    {"full_name"},
		"repo_archive_idx": {"archived"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRepository) InsertColumns() []string {
	return []string{
		"ts",
		"owner",
		"name",
		"full_name",
		"created_at",
		"default_branch",
		"archived",
		"private",
		"license",
		"last_commit_date",
	}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRepository) GetID() int {
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
func (self *GitHubRepository) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubRepository) New() dbs.Cloneable {
	return &GitHubRepository{}
}

// GitHubRepositories is to be used on the struct that needs to pull in
// the repos via a join select statement and provides
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
//   - sql.Scanner1
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

// GitHubRepositoryForeignKey is to be used on the struct that needs to pull in
// the repo via one to many join (being used on the `one` side).
//
// To swap a GitHubRepository to a GitHubRepositoryID:
//
//	var join = models.GitHubRepositoryForeignKey(&GitHubRepository{})
//
// Interfaces:
//   - sql.Scanner
type GitHubRepositoryForeignKey GitHubRepository

// Scan converts the json aggregate result from a select statement into
// a series of GitHubTeams attached to the main struct and will be called
// directly by sqlx
//
// Interfaces:
//   - sql.Scanner
func (self *GitHubRepositoryForeignKey) Scan(src interface{}) (err error) {
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

// NewRepositoryFromRemote converts a github.Repository over to local version and fetches some additional
// innformation like license name and list of teams
func NewRepository(ctx context.Context, client *github.Client, r *github.Repository) (repo *GitHubRepository) {
	var ts = time.Now().UTC().Format(dateformats.Full)
	repo = &GitHubRepository{
		Ts:            ts,
		Owner:         r.GetOwner().GetLogin(),
		Name:          r.GetName(),
		FullName:      r.GetFullName(),
		CreatedAt:     r.GetCreatedAt().Format(dateformats.Full),
		DefaultBranch: r.GetDefaultBranch(),
		Archived:      bools.Int(r.GetArchived()),
		Private:       bools.Int(r.GetPrivate()),
	}
	// get the license and attach the name
	if l := r.GetLicense(); l != nil {
		repo.License = l.GetName()
	}
	// get the default branch and grab the dates
	if branch, _, err := client.Repositories.GetBranch(ctx, repo.Owner, repo.Name, repo.DefaultBranch, 1); err == nil {
		repo.LastCommitDate = branch.Commit.Commit.Author.Date.Time.String()
	}

	repo.GitHubTeams, _ = NewGitHubTeams(ctx, client, repo.Owner, repo.Name)

	return
}
