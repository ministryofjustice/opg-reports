package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/bools"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/ints"
)

// GitHubRepositoryStandard stores all the infor we want locally about our
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
type GitHubRepositoryStandard struct {
	ID                             int                         `json:"id,omitempty" db:"id" faker:"-"`
	Ts                             string                      `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	CompliantBaseline              uint8                       `json:"compliant_baseline,omitempty" db:"compliant_baseline" faker:"oneof: 0, 1"`
	CompliantExtended              uint8                       `json:"compliant_extended,omitempty" db:"compliant_extended" faker:"oneof: 0, 1"`
	CountOfClones                  int                         `json:"count_of_clones,omitempty" db:"count_of_clones" faker:"oneof: 0, 1"`
	CountOfForks                   int                         `json:"count_of_forks,omitempty" db:"count_of_forks" faker:"oneof: 0, 1"`
	CountOfPullRequests            int                         `json:"count_of_pull_requests,omitempty" db:"count_of_pull_requests" faker:"oneof: 0, 1"`
	CountOfWebHooks                int                         `json:"count_of_web_hooks,omitempty" db:"count_of_web_hooks" faker:"oneof: 0, 1"`
	DefaultBranch                  string                      `json:"default_branch,omitempty" db:"default_branch" faker:"oneof: main, master"`
	HasCodeOfConduct               uint8                       `json:"has_code_of_conduct,omitempty" db:"has_code_of_conduct" faker:"oneof: 0, 1"`
	HasCodeownerApprovalRequired   uint8                       `json:"has_codeowner_approval_required,omitempty" db:"has_codeowner_approval_required" faker:"oneof: 0, 1"`
	HasContributingGuide           uint8                       `json:"has_contributing_guide,omitempty" db:"has_contributing_guide" faker:"oneof: 0, 1"`
	HasDefaultBranchOfMain         uint8                       `json:"has_default_branch_of_main,omitempty" db:"has_default_branch_of_main" faker:"oneof: 0, 1"`
	HasDefaultBranchProtection     uint8                       `json:"has_default_branch_protection,omitempty" db:"has_default_branch_protection" faker:"oneof: 0, 1"`
	HasDeleteBranchOnMerge         uint8                       `json:"has_delete_branch_on_merge,omitempty" db:"has_delete_branch_on_merge" faker:"oneof: 0, 1"`
	HasDescription                 uint8                       `json:"has_description,omitempty" db:"has_description" faker:"oneof: 0, 1"`
	HasDiscussions                 uint8                       `json:"has_discussions,omitempty" db:"has_discussions" faker:"oneof: 0, 1"`
	HasDownloads                   uint8                       `json:"has_downloads,omitempty" db:"has_downloads" faker:"oneof: 0, 1"`
	HasIssues                      uint8                       `json:"has_issues,omitempty" db:"has_issues" faker:"oneof: 0, 1"`
	HasLicense                     uint8                       `json:"has_license,omitempty" db:"has_license" faker:"oneof: 0, 1"`
	HasPages                       uint8                       `json:"has_pages,omitempty" db:"has_pages" faker:"oneof: 0, 1"`
	HasPullRequestApprovalRequired uint8                       `json:"has_pull_request_approval_required,omitempty" db:"has_pull_request_approval_required" faker:"oneof: 0, 1"`
	HasReadme                      uint8                       `json:"has_readme,omitempty" db:"has_readme" faker:"oneof: 0, 1"`
	HasRulesEnforcedForAdmins      uint8                       `json:"has_rules_enforced_for_admins,omitempty" db:"has_rules_enforced_for_admins" faker:"oneof: 0, 1"`
	HasVulnerabilityAlerts         uint8                       `json:"has_vulnerability_alerts,omitempty" db:"has_vulnerability_alerts" faker:"oneof: 0, 1"`
	HasWiki                        uint8                       `json:"has_wiki,omitempty" db:"has_wiki" faker:"oneof: 0, 1"`
	IsArchived                     uint8                       `json:"is_archived,omitempty" db:"is_archived" faker:"oneof: 0, 1"`
	IsPrivate                      uint8                       `json:"is_private,omitempty" db:"is_private" faker:"oneof: 0, 1"`
	License                        string                      `json:"license,omitempty" db:"license" faker:"oneof: MIT, GPL"`
	LastCommitDate                 string                      `json:"last_commit_date,omitempty" db:"last_commit_date" faker:"date_string"`
	GitHubRepositoryFullName       string                      `json:"github_repository_full_name,omitempty" db:"github_repository_full_name" faker:"-"`
	GitHubRepositoryID             int                         `json:"github_repository_id,omitempty" db:"github_repository_id" faker:"-"`
	GitHubRepository               *GitHubRepositoryForeignKey `json:"github_repository,omitempty" db:"github_repository" faker:"-"`
}

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *GitHubRepositoryStandard) UniqueValue() string {
	return self.GitHubRepositoryFullName
}

// Interfaces:
//   - dbs.Insertable
func (self *GitHubRepositoryStandard) UniqueField() string {
	return "github_repository_full_name"
}

func (self *GitHubRepositoryStandard) UpsertUpdate() string {
	return `compliant_baseline=excluded.compliant_baseline, compliant_extended=excluded.compliant_extended, count_of_clones=excluded.count_of_clones, count_of_forks=excluded.count_of_forks, count_of_pull_requests=excluded.count_of_pull_requests, count_of_web_hooks=excluded.count_of_web_hooks, default_branch=excluded.default_branch, has_code_of_conduct=excluded.has_code_of_conduct, has_codeowner_approval_required=excluded.has_codeowner_approval_required, has_contributing_guide=excluded.has_contributing_guide, has_default_branch_of_main=excluded.has_default_branch_of_main, has_default_branch_protection=excluded.has_default_branch_protection, has_delete_branch_on_merge=excluded.has_delete_branch_on_merge, has_description=excluded.has_description, has_discussions=excluded.has_discussions, has_downloads=excluded.has_downloads, has_issues=excluded.has_issues, has_license=excluded.has_license, has_pages=excluded.has_pages, has_pull_request_approval_required=excluded.has_pull_request_approval_required, has_readme=excluded.has_readme, has_rules_enforced_for_admins=excluded.has_rules_enforced_for_admins, has_vulnerability_alerts=excluded.has_vulnerability_alerts, has_wiki=excluded.has_wiki, is_archived=excluded.is_archived, is_private=excluded.is_private, license=excluded.license, last_commit_date=excluded.last_commit_date, github_repository_id=excluded.github_repository_id`
}

// TableName returns named table for GitHubRepositoryStandard - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *GitHubRepositoryStandard) TableName() string {
	return "github_repository_standards"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepositoryStandard) Columns() map[string]string {
	return map[string]string{
		"id":                                 "INTEGER PRIMARY KEY",
		"ts":                                 "TEXT NOT NULL",
		"compliant_baseline":                 "INTEGER NOT NULL DEFAULT 0",
		"compliant_extended":                 "INTEGER NOT NULL DEFAULT 0",
		"count_of_clones":                    "INTEGER NOT NULL DEFAULT 0",
		"count_of_forks":                     "INTEGER NOT NULL DEFAULT 0",
		"count_of_pull_requests":             "INTEGER NOT NULL DEFAULT 0",
		"count_of_web_hooks":                 "INTEGER NOT NULL DEFAULT 0",
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
		"github_repository_full_name":        "TEXT NOT NULL UNIQUE",
		"github_repository_id":               "INTEGER NOT NULL",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *GitHubRepositoryStandard) Indexes() map[string][]string {
	return map[string][]string{
		"ghs_fullname_idx": {"github_repository_full_name"},
		"ghs_baseline_idx": {"compliant_baseline"},
		"ghs_extended_idx": {"compliant_extended"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *GitHubRepositoryStandard) InsertColumns() []string {
	return []string{
		"ts",
		"compliant_baseline",
		"compliant_extended",
		"count_of_clones",
		"count_of_forks",
		"count_of_pull_requests",
		"count_of_web_hooks",
		"default_branch",
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
		"github_repository_full_name",
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
func (self *GitHubRepositoryStandard) GetID() int {
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
func (self *GitHubRepositoryStandard) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *GitHubRepositoryStandard) New() dbs.Cloneable {
	return &GitHubRepositoryStandard{}
}

// Info returns the infomational standards
func (self *GitHubRepositoryStandard) Info() map[string]string {

	return map[string]string{
		"archived":                   fmt.Sprintf("%t", ints.Bool(self.IsArchived)),
		"created_at":                 fmt.Sprintf("%s", self.GitHubRepository.CreatedAt),
		"branch_name":                self.DefaultBranch,
		"has_delete_branch_on_merge": fmt.Sprintf("%t", ints.Bool(self.HasDeleteBranchOnMerge)),
		"has_pages":                  fmt.Sprintf("%t", ints.Bool(self.HasPages)),
		"has_downloads":              fmt.Sprintf("%t", ints.Bool(self.HasDownloads)),
		"has_discussions":            fmt.Sprintf("%t", ints.Bool(self.HasDiscussions)),
		"has_wiki":                   fmt.Sprintf("%t", ints.Bool(self.HasWiki)),
		"forks":                      fmt.Sprintf("%d", self.CountOfForks),
		"webhooks":                   fmt.Sprintf("%d", self.CountOfWebHooks),
		"open_pull_requests":         fmt.Sprintf("%d", self.CountOfPullRequests),
		"clone_traffic":              fmt.Sprintf("%d", self.CountOfClones),
		"is_private":                 fmt.Sprintf("%t", ints.Bool(self.IsPrivate)),
		"last_commit_date":           fmt.Sprintf("%s", self.LastCommitDate),
	}
}

func (self *GitHubRepositoryStandard) Baseline() map[string]bool {
	return map[string]bool{
		"has_default_branch_of_main":         ints.Bool(self.HasDefaultBranchOfMain),
		"has_license":                        ints.Bool(self.HasLicense),
		"has_issues":                         ints.Bool(self.HasIssues),
		"has_description":                    ints.Bool(self.HasDescription),
		"has_rules_enforced_for_admins":      ints.Bool(self.HasRulesEnforcedForAdmins),
		"has_pull_request_approval_required": ints.Bool(self.HasPullRequestApprovalRequired),
	}
}

func (self *GitHubRepositoryStandard) Extended() map[string]bool {
	return map[string]bool{
		"has_code_owner_approval_required": ints.Bool(self.HasCodeownerApprovalRequired),
		"has_readme":                       ints.Bool(self.HasReadme),
		"has_code_of_conduct":              ints.Bool(self.HasCodeOfConduct),
		"has_contributing_guide":           ints.Bool(self.HasContributingGuide),
	}
}

func (self *GitHubRepositoryStandard) UpdateCompliance() (baseline uint8, extended uint8) {
	baselineChecks := self.Baseline()
	extendedChecks := self.Extended()

	baseline = 1
	extended = 1

	for _, is := range baselineChecks {
		if !is {
			baseline = 0
		}
	}
	for _, is := range extendedChecks {
		if !is {
			extended = 0
		}
	}
	self.CompliantBaseline = baseline
	self.CompliantExtended = extended

	return
}

// NewRepositoryStandard uses the repository details passed to create a new Local GitHubRepositoryStandard with
// an attached GitHubRepository (which in turn will have GitHubTeams) and populates all the compliance related
// fields
func NewRepositoryStandard(ctx context.Context, client *github.Client, r *github.Repository) (g *GitHubRepositoryStandard) {
	const (
		readmePath            string = "./README.md"
		codeOfConductPath     string = "./CODE_OF_CONDUCT.md"
		contributingGuidePath string = "./CONTRIBUTING.md"
	)
	var repo = NewRepository(ctx, client, r)

	g = &GitHubRepositoryStandard{
		Ts:                       time.Now().UTC().Format(dateformats.Full),
		DefaultBranch:            repo.DefaultBranch,
		LastCommitDate:           repo.LastCommitDate,
		GitHubRepository:         (*GitHubRepositoryForeignKey)(repo),
		IsArchived:               repo.Archived,
		IsPrivate:                repo.Private,
		GitHubRepositoryFullName: repo.FullName,
	}

	//
	g.HasDefaultBranchProtection = 0
	if branch, _, err := client.Repositories.GetBranch(ctx, repo.Owner, repo.Name, g.DefaultBranch, 1); err == nil {
		if branch.GetProtected() {
			g.HasDefaultBranchProtection = 1
		}
	}
	// -- counters
	if clones, _, err := client.Repositories.ListTrafficClones(
		context.Background(), repo.Owner, repo.Name,
		&github.TrafficBreakdownOptions{Per: "day"}); err == nil {

		g.CountOfClones = clones.GetCount()
	}
	g.CountOfForks = r.GetForksCount()
	if prs, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name,
		&github.PullRequestListOptions{State: "open"}); err == nil {
		g.CountOfPullRequests = len(prs)
	}
	if hooks, _, err := client.Repositories.ListHooks(ctx, repo.Owner, repo.Name,
		&github.ListOptions{PerPage: 100}); err == nil {
		g.CountOfWebHooks = len(hooks)
	}
	// -- has
	g.HasCodeOfConduct = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name,
		codeOfConductPath, nil); err == nil {
		g.HasCodeOfConduct = 1
	}
	g.HasContributingGuide = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name,
		contributingGuidePath, nil); err == nil {
		g.HasContributingGuide = 1
	}
	g.HasReadme = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name,
		readmePath, nil); err == nil {
		g.HasReadme = 1
	}

	g.HasVulnerabilityAlerts = 0
	if alerts, _, err := client.Repositories.GetVulnerabilityAlerts(ctx, repo.Owner, repo.Name); err == nil {
		g.HasVulnerabilityAlerts = bools.Int(alerts)
	}

	g.HasDefaultBranchOfMain = bools.Int((g.DefaultBranch == "main"))
	g.HasDescription = bools.Int((len(r.GetDescription()) > 0))
	g.HasDiscussions = bools.Int(r.GetHasDiscussions())
	g.HasDownloads = bools.Int(r.GetHasDownloads())
	g.HasIssues = bools.Int(r.GetHasIssues())
	g.HasLicense = bools.Int((len(repo.License) > 0))
	g.HasPages = bools.Int(r.GetHasPages())
	g.HasWiki = bools.Int(r.GetHasWiki())

	if protection, _, err := client.Repositories.GetBranchProtection(ctx, repo.Owner, repo.Name,
		g.DefaultBranch); err == nil {
		g.HasRulesEnforcedForAdmins = bools.Int(protection.EnforceAdmins.Enabled)
		g.HasPullRequestApprovalRequired = bools.Int(protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0)
		g.HasCodeownerApprovalRequired = bools.Int(protection.RequiredPullRequestReviews.RequireCodeOwnerReviews)
	}

	// the GetDeleteBranchOnMerge seems to be empty and have to re-fetch the api to get result
	re, _, _ := client.Repositories.Get(ctx, repo.Owner, repo.Name)
	g.HasDeleteBranchOnMerge = bools.Int(re.GetDeleteBranchOnMerge())

	g.UpdateCompliance()
	return
}
