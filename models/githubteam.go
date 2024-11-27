package models

import (
	"fmt"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
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
	ID   int    `json:"id,omitempty" db:"id" faker:"-"`
	Ts   string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Slug string `json:"slug,omitempty" db:"slug" faker:"unique, oneof:digideps,opg-lpa-team,opg-modernising-lpa-team,opg-use-a-lpa-team,opg-refund,serve-opg,sirius,opg-webops,opg-admin,opg_intergrations,opg-metrics-team"`
	// Join to units - team has many units, unit has many teams
	Units Units `json:"units,omitempty" db:"units" faker:"-"`
	// Many to many join with the repositories
	GitHubRepositories GitHubRepositories `json:"github_repositories,omitempty" db:"github_repositories" faker:"-"`
}

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *GitHubTeam) UniqueValue() string {
	return self.Slug
}

// Interfaces:
//   - dbs.Insertable
func (self *GitHubTeam) UniqueField() string {
	return "slug"
}
func (self *GitHubTeam) UpsertUpdate() string {
	return "slug=excluded.slug"
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

func (self *GitHubTeam) StandardUnits() (units []*Unit) {
	now := time.Now().UTC().Format(dateformats.Full)
	units = []*Unit{}
	switch self.Slug {
	case "digideps":
		units = append(units, &Unit{
			Ts:   now,
			Name: "digideps",
		})
	case "opg-lpa-team":
		units = append(units, &Unit{
			Ts:   now,
			Name: "make",
		})
	case "opg-modernising-lpa-team":
		units = append(units, &Unit{
			Ts:   now,
			Name: "modernise",
		})
	case "opg-use-a-lpa-team":
		units = append(units, &Unit{
			Ts:   now,
			Name: "use",
		})
	case "opg-refund":
		units = append(units, &Unit{
			Ts:   now,
			Name: "refunds",
		})
	case "serve-opg":
		units = append(units, &Unit{
			Ts:   now,
			Name: "serve",
		})
	case "sirius":
		units = append(units, &Unit{
			Ts:   now,
			Name: "sirius",
		})
	default:
		units = append(units, &Unit{
			Ts:   now,
			Name: "org",
		})
	}
	return
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
