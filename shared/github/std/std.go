package std

import (
	"context"
	"log/slog"
	"opg-reports/services/front/cnf"
	"opg-reports/shared/data"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
)

const readmePath string = "./README.md"
const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

var defaultBaselineCompliance = []string{
	"has_default_branch_of_main",
	"has_license",
	"has_issues",
	"has_description",
	"has_rules_enforced_for_admins",
	"has_pull_request_approval_required",
}
var defaultExtendedCompliance = []string{
	"has_code_owner_approval_required",
	"has_readme",
	"has_code_of_conduct",
	"has_contributing_guide",
}
var defaultInformation = []string{
	"archived",
	"branch_name",
	"has_delete_branch_on_merge",
	"has_pages",
	"has_downloads",
	"has_discussions",
	"has_wiki",
	"forks",
	"webhooks",
	"open_pull_requests",
	"clone_traffic",
	"is_private",
	"last_commit_date",
}

type Repository struct {
	r    *github.Repository
	stds *cnf.RepoStandards

	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"time"`

	Archived       bool      `json:"archived"`
	DefaultBranch  string    `json:"default_branch"`
	FullName       string    `json:"full_name"`
	License        string    `json:"license"`
	Name           string    `json:"name"`
	Owner          string    `json:"owner"`
	LastCommitDate time.Time `json:"last_commit_date"`

	CountClones       int `json:"clone_traffic"`
	CountForks        int `json:"forks"`
	CountPullRequests int `json:"open_pull_requests"`
	CountWebhooks     int `json:"webhooks"`

	HasCodeOfConduct               bool `json:"has_code_of_conduct"`
	HasCodeOwnerApprovalRequired   bool `json:"has_code_owner_approval_required"`
	HasContributingGuide           bool `json:"has_contributing_guide"`
	HasDefaultBranchOfMain         bool `json:"has_default_branch_of_main"`
	HasDefaultBranchProtection     bool `json:"has_default_branch_protection"`
	HasDeleteBranchOnMerge         bool `json:"has_delete_branch_on_merge"`
	HasDescription                 bool `json:"has_description"`
	HasDiscussions                 bool `json:"has_discussions"`
	HasDownloads                   bool `json:"has_downloads"`
	HasIssues                      bool `json:"has_issues"`
	HasLicense                     bool `json:"has_license"`
	HasPages                       bool `json:"has_pages"`
	HasProjects                    bool `json:"has_projects"`
	HasPullRequestApprovalRequired bool `json:"has_pull_request_approval_required"`
	HasReadme                      bool `json:"has_readme"`
	HasRulesEnforcedForAdmins      bool `json:"has_rules_enforced_for_admins"`
	HasVulnerabilityAlerts         bool `json:"has_vulnerability_alerts"`
	HasWiki                        bool `json:"has_wiki"`
	IsPrivate                      bool `json:"is_private"`
}

// UID is the unique id (UUID) for this Compliance item
func (c *Repository) UID() string {
	slog.Debug("[github/std/repo] UID()", slog.String("UID", c.UUID))
	return c.UUID
}

func (c *Repository) TS() time.Time {
	return c.Timestamp
}

// Valid returns true only if all fields are present with a non-empty value
func (c *Repository) Valid() (valid bool) {
	mapped, _ := data.ToMap(c)
	for k, val := range mapped {
		if val == "" {
			slog.Debug("[github/std/repo] invalid",
				slog.String("UID", c.UID()),
				slog.String(k, val.(string)))
			return false
		}
	}
	slog.Debug("[github/std/repo] valid", slog.String("UID", c.UID()))
	return true
}

func (c *Repository) SetStandards(stds *cnf.RepoStandards) {
	c.stds = stds
}
func (c *Repository) GetStandards() *cnf.RepoStandards {
	return c.stds
}

// Compliant checks that the fieldnames passed have true values, if they dont, then this does not comply
// Allows dynamic setting of compliance that be draw from a config
func (c *Repository) Compliant(booleanFields []string) (complies bool, values map[string]bool, err error) {
	complies = false
	values = map[string]bool{}
	mapped, err := data.ToMap(c)
	if err != nil {
		return
	}
	complies = true
	for _, key := range booleanFields {
		values[key] = true
		if v, ok := mapped[key]; !ok || !v.(bool) {
			complies = false
			values[key] = false
		}
	}
	return
}

// setData calls the internal data setting funcs
func (c *Repository) setData(client *github.Client) {
	c.dataDirectFromR()
	c.dataViaClient(client)
}

// dataViaClient sets data that requires additional calls using the
// github client to fetch more information
func (c *Repository) dataViaClient(client *github.Client) {

	// Check branch protection
	if branch, _, err := client.Repositories.GetBranch(context.Background(), c.Owner, c.Name, c.DefaultBranch, 1); err == nil {
		c.HasDefaultBranchProtection = branch.GetProtected()
		c.LastCommitDate = branch.GetCommit().GetCommitter().GetCreatedAt().Time
	}

	// check branch protection rules
	if protection, _, err := client.Repositories.GetBranchProtection(context.Background(), c.Owner, c.Name, c.DefaultBranch); err == nil {
		c.HasRulesEnforcedForAdmins = protection.EnforceAdmins.Enabled
		c.HasPullRequestApprovalRequired = protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0
		c.HasCodeOwnerApprovalRequired = protection.RequiredPullRequestReviews.RequireCodeOwnerReviews

	}
	// vuln alerts enabled
	if alerts, _, err := client.Repositories.GetVulnerabilityAlerts(context.Background(), c.Owner, c.Name); err == nil {
		c.HasVulnerabilityAlerts = alerts
	}
	// readme present
	if _, _, _, err := client.Repositories.GetContents(context.Background(), c.Owner, c.Name, readmePath, nil); err == nil {
		c.HasReadme = true
	}
	// code of conduct
	if _, _, _, err := client.Repositories.GetContents(context.Background(), c.Owner, c.Name, codeOfConductPath, nil); err == nil {
		c.HasCodeOfConduct = true
	}
	// contributing guide
	if _, _, _, err := client.Repositories.GetContents(context.Background(), c.Owner, c.Name, contributingGuidePath, nil); err == nil {
		c.HasContributingGuide = true
	}
	// open pull requests
	if prs, _, err := client.PullRequests.List(context.Background(), c.Owner, c.Name, &github.PullRequestListOptions{State: "open"}); err == nil {
		c.CountPullRequests = len(prs)
	}
	// clone traffic
	if clones, _, err := client.Repositories.ListTrafficClones(context.Background(), c.Owner, c.Name, &github.TrafficBreakdownOptions{Per: "day"}); err == nil {
		c.CountClones = clones.GetCount()
	}
	// webhooks
	if hooks, _, err := client.Repositories.ListHooks(context.Background(), c.Owner, c.Name, &github.ListOptions{PerPage: 100}); err == nil {
		c.CountWebhooks = len(hooks)
	}
}

// dataDirectFromR sets data that can be found directly on the
// repository without making extra calls
func (c *Repository) dataDirectFromR() {
	r := c.r
	c.Archived = r.GetArchived()
	c.DefaultBranch = r.GetDefaultBranch()
	c.FullName = r.GetFullName()
	c.Name = r.GetName()
	c.Owner = r.GetOwner().GetLogin()

	if l := r.GetLicense(); l != nil {
		c.License = l.GetName()
	}

	c.HasLicense = (len(c.License) > 0)
	c.HasDefaultBranchOfMain = (c.DefaultBranch == "main")
	c.HasDescription = len(r.GetDescription()) > 0
	c.HasDiscussions = r.GetHasDiscussions()
	c.HasDeleteBranchOnMerge = r.GetDeleteBranchOnMerge()
	c.HasDownloads = r.GetHasDownloads()
	c.HasIssues = r.GetHasIssues()
	c.HasPages = r.GetHasPages()
	c.HasProjects = r.GetHasProjects()
	c.HasWiki = r.GetHasWiki()
	c.IsPrivate = r.GetPrivate()

	c.CountForks = r.GetForksCount()

}

var _ data.IEntry = &Repository{}

func New(uid *string) *Repository {
	c := &Repository{Timestamp: time.Now().UTC()}
	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	return c
}

func NewWithR(uid *string, r *github.Repository, client *github.Client) *Repository {
	c := New(uid)
	c.r = r
	c.setData(client)
	return c
}
