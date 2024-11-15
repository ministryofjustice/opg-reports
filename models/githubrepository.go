package models

import (
	"fmt"

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
	ID                             int         `json:"id,omitempty" db:"id" faker:"-"`
	Ts                             string      `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	CompliantBaseline              uint8       `json:"compliant_baseline,omitempty" db:"compliant_baseline" faker:"oneof: 0, 1"`
	CompliantExtended              uint8       `json:"compliant_extended,omitempty" db:"compliant_extended" faker:"oneof: 0, 1"`
	CountOfClones                  int         `json:"count_of_clones,omitempty" db:"count_of_clones" faker:"oneof: 0, 1"`
	CountOfForks                   int         `json:"count_of_forks,omitempty" db:"count_of_forks" faker:"oneof: 0, 1"`
	CountOfPullRequests            int         `json:"count_of_pull_requests,omitempty" db:"count_of_pull_requests" faker:"oneof: 0, 1"`
	CountOfWebHooks                int         `json:"count_of_web_hooks,omitempty" db:"count_of_web_hooks" faker:"oneof: 0, 1"`
	CreatedAt                      string      `json:"created_at,omitempty" db:"created_at" faker:"date_string"`
	DefaultBranch                  string      `json:"default_branch,omitempty" db:"default_branch" faker:"oneof: main, master"`
	FullName                       string      `json:"full_name,omitempty" db:"full_name"`
	HasCodeOfConduct               uint8       `json:"has_code_of_conduct,omitempty" db:"has_code_of_conduct" faker:"oneof: 0, 1"`
	HasCodeownerApprovalRequired   uint8       `json:"has_codeowner_approval_required,omitempty" db:"has_codeowner_approval_required" faker:"oneof: 0, 1"`
	HasContributingGuide           uint8       `json:"has_contributing_guide,omitempty" db:"has_contributing_guide" faker:"oneof: 0, 1"`
	HasDefaultBranchOfMain         uint8       `json:"has_default_branch_of_main,omitempty" db:"has_default_branch_of_main" faker:"oneof: 0, 1"`
	HasDefaultBranchProtection     uint8       `json:"has_default_branch_protection,omitempty" db:"has_default_branch_protection" faker:"oneof: 0, 1"`
	HasDeleteBranchOnMerge         uint8       `json:"has_delete_branch_on_merge,omitempty" db:"has_delete_branch_on_merge" faker:"oneof: 0, 1"`
	HasDescription                 uint8       `json:"has_description,omitempty" db:"has_description" faker:"oneof: 0, 1"`
	HasDiscussions                 uint8       `json:"has_discussions,omitempty" db:"has_discussions" faker:"oneof: 0, 1"`
	HasDownloads                   uint8       `json:"has_downloads,omitempty" db:"has_downloads" faker:"oneof: 0, 1"`
	HasIssues                      uint8       `json:"has_issues,omitempty" db:"has_issues" faker:"oneof: 0, 1"`
	HasLicense                     uint8       `json:"has_license,omitempty" db:"has_license" faker:"oneof: 0, 1"`
	HasPages                       uint8       `json:"has_pages,omitempty" db:"has_pages" faker:"oneof: 0, 1"`
	HasPullRequestApprovalRequired uint8       `json:"has_pull_request_approval_required,omitempty" db:"has_pull_request_approval_required" faker:"oneof: 0, 1"`
	HasReadme                      uint8       `json:"has_readme,omitempty" db:"has_readme" faker:"oneof: 0, 1"`
	HasRulesEnforcedForAdmins      uint8       `json:"has_rules_enforced_for_admins,omitempty" db:"has_rules_enforced_for_admins" faker:"oneof: 0, 1"`
	HasVulnerabilityAlerts         uint8       `json:"has_vulnerability_alerts,omitempty" db:"has_vulnerability_alerts" faker:"oneof: 0, 1"`
	HasWiki                        uint8       `json:"has_wiki,omitempty" db:"has_wiki" faker:"oneof: 0, 1"`
	IsArchived                     uint8       `json:"is_archived,omitempty" db:"is_archived" faker:"oneof: 0, 1"`
	IsPrivate                      uint8       `json:"is_private,omitempty" db:"is_private" faker:"oneof: 0, 1"`
	License                        string      `json:"license,omitempty" db:"license" faker:"oneof: MIT, GPL"`
	LastCommitDate                 string      `json:"last_commit_date,omitempty" db:"last_commit_date" faker:"date_string"`
	Name                           string      `json:"name" db:"name"`
	Owner                          string      `json:"owner" db:"owner" faker:"oneof: ministryofjusice"`
	GitHubTeams                    GitHubTeams `json:"github_teams,omitempty" db:"github_teams" faker:"-"`
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
		"id":                                 "INTEGER PRIMARY KEY",
		"ts":                                 "TEXT NOT NULL",
		"full_name":                          "TEXT NOT NULL UNIQUE",
		"compliant_baseline":                 "INTEGER NOT NULL DEFAULT 0",
		"compliant_extended":                 "INTEGER NOT NULL DEFAULT 0",
		"count_of_clones":                    "INTEGER NOT NULL DEFAULT 0",
		"count_of_forks":                     "INTEGER NOT NULL DEFAULT 0",
		"count_of_pull_requests":             "INTEGER NOT NULL DEFAULT 0",
		"count_of_web_hooks":                 "INTEGER NOT NULL DEFAULT 0",
		"created_at":                         "TEXT NOT NULL",
		"default_branch":                     "TEXT NOT NULL",
		"has_code_of_conduct":                "INTEGER NOT NULL DEFAULT 0",
		"has_codeowner_approval_required":    "INTEGER NOT NULL DEFAULT 0",
		"has_contributing_guide":             "INTEGER NOT NULL DEFAULT 0",
		"has_default_branch_of_main":         "INTEGER NOT NULL DEFAULT 0",
		"has_default_branch_protection":      "INTEGER NOT NULL DEFAULT 0",
		"has_delete_branch_on_merge":         "INTEGER NOT NULL DEFAULT 0",
		"has_description":                    "INTEGER NOT NULL DEFAULT 0",
		"has_discussions":                    "INTEGER NOT NULL DEFAULT 0",
		"has_downloads":                      "INTEGER NOT NULL DEFAULT 0",
		"has_issues":                         "INTEGER NOT NULL DEFAULT 0",
		"has_license":                        "INTEGER NOT NULL DEFAULT 0",
		"has_pages":                          "INTEGER NOT NULL DEFAULT 0",
		"has_pull_request_approval_required": "INTEGER NOT NULL DEFAULT 0",
		"has_readme":                         "INTEGER NOT NULL DEFAULT 0",
		"has_rules_enforced_for_admins":      "INTEGER NOT NULL DEFAULT 0",
		"has_vulnerability_alerts":           "INTEGER NOT NULL DEFAULT 0",
		"has_wiki":                           "INTEGER NOT NULL DEFAULT 0",
		"is_archived":                        "INTEGER NOT NULL DEFAULT 0",
		"is_private":                         "INTEGER NOT NULL DEFAULT 0",
		"license":                            "TEXT NOT NULL DEFAULT ''",
		"last_commit_date":                   "TEXT NOT NULL",
		"name":                               "TEXT NOT NULL",
		"owner":                              "TEXT NOT NULL",
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
		"repo_name_idx":     {"full_name"},
		"repo_baseline_idx": {"compliant_baseline"},
		"repo_extended_idx": {"compliant_extended"},
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
		"compliant_baseline",
		"compliant_extended",
		"count_of_clones",
		"count_of_forks",
		"count_of_pull_requests",
		"count_of_web_hooks",
		"created_at",
		"default_branch",
		"full_name",
		"has_code_of_conduct",
		"has_codeowner_approval_required",
		"has_contributing_guide",
		"has_default_branch_of_main",
		"has_default_branch_protection",
		"has_delete_branch_on_merge",
		"has_description",
		"has_discussions",
		"has_downloads",
		"has_issues",
		"has_license",
		"has_pages",
		"has_pull_request_approval_required",
		"has_readme",
		"has_rules_enforced_for_admins",
		"has_vulnerability_alerts",
		"has_wiki",
		"is_archived",
		"is_private",
		"license",
		"last_commit_date",
		"name",
		"owner",
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
