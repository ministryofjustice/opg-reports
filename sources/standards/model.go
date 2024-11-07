package standards

import (
	"fmt"
	"strings"

	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/record"
)

type Standard struct {
	ID                             int    `json:"id,omitempty" db:"id" faker:"-" doc:"Database primary key."`                   // ID is a generated primary key
	Ts                             string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	CompliantBaseline              uint8  `json:"compliant_baseline,omitempty" db:"compliant_baseline" faker:"oneof: 0, 1"`
	CompliantExtended              uint8  `json:"compliant_extended,omitempty" db:"compliant_extended" faker:"oneof: 0, 1"`
	CountOfClones                  int    `json:"count_of_clones,omitempty" db:"count_of_clones" faker:"oneof: 0, 1"`
	CountOfForks                   int    `json:"count_of_forks,omitempty" db:"count_of_forks" faker:"oneof: 0, 1"`
	CountOfPullRequests            int    `json:"count_of_pull_requests,omitempty" db:"count_of_pull_requests" faker:"oneof: 0, 1"`
	CountOfWebHooks                int    `json:"count_of_web_hooks,omitempty" db:"count_of_web_hooks" faker:"oneof: 0, 1"`
	CreatedAt                      string `json:"created_at,omitempty" db:"created_at" faker:"date_string"`
	DefaultBranch                  string `json:"default_branch,omitempty" db:"default_branch" faker:"oneof: main, master"`
	FullName                       string `json:"full_name,omitempty" db:"full_name"`
	HasCodeOfConduct               uint8  `json:"has_code_of_conduct,omitempty" db:"has_code_of_conduct" faker:"oneof: 0, 1"`
	HasCodeownerApprovalRequired   uint8  `json:"has_codeowner_approval_required,omitempty" db:"has_codeowner_approval_required" faker:"oneof: 0, 1"`
	HasContributingGuide           uint8  `json:"has_contributing_guide,omitempty" db:"has_contributing_guide" faker:"oneof: 0, 1"`
	HasDefaultBranchOfMain         uint8  `json:"has_default_branch_of_main,omitempty" db:"has_default_branch_of_main" faker:"oneof: 0, 1"`
	HasDefaultBranchProtection     uint8  `json:"has_default_branch_protection,omitempty" db:"has_default_branch_protection" faker:"oneof: 0, 1"`
	HasDeleteBranchOnMerge         uint8  `json:"has_delete_branch_on_merge,omitempty" db:"has_delete_branch_on_merge" faker:"oneof: 0, 1"`
	HasDescription                 uint8  `json:"has_description,omitempty" db:"has_description" faker:"oneof: 0, 1"`
	HasDiscussions                 uint8  `json:"has_discussions,omitempty" db:"has_discussions" faker:"oneof: 0, 1"`
	HasDownloads                   uint8  `json:"has_downloads,omitempty" db:"has_downloads" faker:"oneof: 0, 1"`
	HasIssues                      uint8  `json:"has_issues,omitempty" db:"has_issues" faker:"oneof: 0, 1"`
	HasLicense                     uint8  `json:"has_license,omitempty" db:"has_license" faker:"oneof: 0, 1"`
	HasPages                       uint8  `json:"has_pages,omitempty" db:"has_pages" faker:"oneof: 0, 1"`
	HasPullRequestApprovalRequired uint8  `json:"has_pull_request_approval_required,omitempty" db:"has_pull_request_approval_required" faker:"oneof: 0, 1"`
	HasReadme                      uint8  `json:"has_readme,omitempty" db:"has_readme" faker:"oneof: 0, 1"`
	HasRulesEnforcedForAdmins      uint8  `json:"has_rules_enforced_for_admins,omitempty" db:"has_rules_enforced_for_admins" faker:"oneof: 0, 1"`
	HasVulnerabilityAlerts         uint8  `json:"has_vulnerability_alerts,omitempty" db:"has_vulnerability_alerts" faker:"oneof: 0, 1"`
	HasWiki                        uint8  `json:"has_wiki,omitempty" db:"has_wiki" faker:"oneof: 0, 1"`
	IsArchived                     uint8  `json:"is_archived,omitempty" db:"is_archived" faker:"oneof: 0, 1"`
	IsPrivate                      uint8  `json:"is_private,omitempty" db:"is_private" faker:"oneof: 0, 1"`
	License                        string `json:"license,omitempty" db:"license" faker:"oneof: MIT, GPL"`
	LastCommitDate                 string `json:"last_commit_date,omitempty" db:"last_commit_date" faker:"date_string"`
	Name                           string `json:"name" db:"name"`
	Owner                          string `json:"owner" db:"owner" faker:"oneof: ministryofjusice"`
	Teams                          string `json:"teams" db:"teams" faker:"oneof: #unitA#, #unitB#, #unitC#"`
}

// New
// Record interface
func (self *Standard) New() record.Record {
	return &Standard{}
}

// UID
// Record interface
func (self *Standard) UID() string {
	return fmt.Sprintf("%s-%d", "standard", self.ID)
}

// SetID
// Record interface
func (self *Standard) SetID(id int) {
	self.ID = id
}

// IsCompliantBaseline checks itself and returns
func (self *Standard) IsCompliantBaseline() bool {
	return convert.IntToBool(self.CompliantExtended)
}

// IsCompliantExtended checks itself and returns bool
func (self *Standard) IsCompliantExtended() bool {
	return convert.IntToBool(self.CompliantExtended)

}

func (self *Standard) Archived() bool {
	return convert.IntToBool(self.IsArchived)
}

func (self *Standard) Private() bool {
	return convert.IntToBool(self.IsPrivate)
}

// TeamList retusn a slice of stirngs of the team names without the `#`
// db seperators
func (self *Standard) TeamList() (teams []string) {
	teams = []string{}
	for _, t := range strings.Split(self.Teams, "#") {
		if t != "" {
			teams = append(teams, t)
		}
	}

	return
}

// Info returns the infomational standards
func (g *Standard) Info() map[string]string {

	return map[string]string{
		"archived":                   fmt.Sprintf("%t", convert.IntToBool(g.IsArchived)),
		"created_at":                 fmt.Sprintf("%s", g.CreatedAt),
		"branch_name":                g.DefaultBranch,
		"has_delete_branch_on_merge": fmt.Sprintf("%t", convert.IntToBool(g.HasDeleteBranchOnMerge)),
		"has_pages":                  fmt.Sprintf("%t", convert.IntToBool(g.HasPages)),
		"has_downloads":              fmt.Sprintf("%t", convert.IntToBool(g.HasDownloads)),
		"has_discussions":            fmt.Sprintf("%t", convert.IntToBool(g.HasDiscussions)),
		"has_wiki":                   fmt.Sprintf("%t", convert.IntToBool(g.HasWiki)),
		"forks":                      fmt.Sprintf("%d", g.CountOfForks),
		"webhooks":                   fmt.Sprintf("%d", g.CountOfWebHooks),
		"open_pull_requests":         fmt.Sprintf("%d", g.CountOfPullRequests),
		"clone_traffic":              fmt.Sprintf("%d", g.CountOfClones),
		"is_private":                 fmt.Sprintf("%t", convert.IntToBool(g.IsPrivate)),
		"last_commit_date":           fmt.Sprintf("%s", g.LastCommitDate),
	}
}

func (g *Standard) Baseline() map[string]bool {
	return map[string]bool{
		"has_default_branch_of_main":         convert.IntToBool(g.HasDefaultBranchOfMain),
		"has_license":                        convert.IntToBool(g.HasLicense),
		"has_issues":                         convert.IntToBool(g.HasIssues),
		"has_description":                    convert.IntToBool(g.HasDescription),
		"has_rules_enforced_for_admins":      convert.IntToBool(g.HasRulesEnforcedForAdmins),
		"has_pull_request_approval_required": convert.IntToBool(g.HasPullRequestApprovalRequired),
	}
}

func (g *Standard) Extended() map[string]bool {
	return map[string]bool{
		"has_code_owner_approval_required": convert.IntToBool(g.HasCodeownerApprovalRequired),
		"has_readme":                       convert.IntToBool(g.HasReadme),
		"has_code_of_conduct":              convert.IntToBool(g.HasCodeOfConduct),
		"has_contributing_guide":           convert.IntToBool(g.HasContributingGuide),
	}
}

func (g *Standard) UpdateCompliance() (baseline uint8, extended uint8) {
	baselineChecks := map[string]bool{
		"has_default_branch_of_main":         (g.HasDefaultBranchOfMain == 1),
		"has_license":                        (g.HasLicense == 1),
		"has_issues":                         (g.HasIssues == 1),
		"has_description":                    (g.HasDescription == 1),
		"has_rules_enforced_for_admins":      (g.HasRulesEnforcedForAdmins == 1),
		"has_pull_request_approval_required": (g.HasPullRequestApprovalRequired == 1),
	}
	extendedChecks := map[string]bool{
		"has_code_owner_approval_required": (g.HasCodeownerApprovalRequired == 1),
		"has_readme":                       (g.HasReadme == 1),
		"has_code_of_conduct":              (g.HasCodeOfConduct == 1),
		"has_contributing_guide":           (g.HasContributingGuide == 1),
	}

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
	g.CompliantBaseline = baseline
	g.CompliantExtended = extended

	return
}
