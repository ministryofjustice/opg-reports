package comp

import (
	"context"
	"log/slog"
	"opg-reports/shared/data"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
)

const readmePath string = "./README.md"
const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

type Compliance struct {
	r *github.Repository `json:"-"`

	UUID string `json:"uuid"`

	Archived      bool   `json:"archived"`
	DefaultBranch string `json:"default_branch"`
	FullName      string `json:"full_name"`
	License       string `json:"license"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`

	CountClones       int `json:"clone_traffic"`
	CountForks        int `json:"forks_count"`
	CountPullRequests int `json:"open_pull_requests"`
	CountWebhooks     int `json:"webhooks"`

	HasDescription bool `json:"has_description"`
	HasDiscussions bool `json:"has_discussions"`
	HasDownloads   bool `json:"has_downloads"`
	HasIssues      bool `json:"has_issues"`
	HasPages       bool `json:"has_pages"`
	HasProjects    bool `json:"has_projects"`
	HasWiki        bool `json:"has_wiki"`
	IsPrivate      bool `json:"is_private"`

	HasCodeOfConduct               bool `json:"has_code_of_conduct"`
	HasCodeOwnerApprovalRequired   bool `json:"has_code_owner_approval_required"`
	HasContributingGuide           bool `json:"has_contributing_guide"`
	HasDefaultBranchProtection     bool `json:"has_default_branch_protection"`
	HasReadme                      bool `json:"has_readme"`
	HasRulesEnforcedForAdmins      bool `json:"has_rules_enforced_for_admins"`
	HasPullRequestApprovalRequired bool `json:"has_pull_request_approval_required"`
	HasVulnerabilityAlerts         bool `json:"has_vulnerability_alerts"`

	Baseline bool `json:"baseline_compliant"`
	Extended bool `json:"extended_compliant"`
}

// UID is the unique id (UUID) for this Compliance item
func (c *Compliance) UID() string {
	slog.Debug("[gh/compliance] UID()", slog.String("UID", c.UUID))
	return c.UUID
}

// Valid returns true only if all fields are present with a non-empty value
func (c *Compliance) Valid() (valid bool) {
	mapped, _ := data.ToMap(c)
	for k, val := range mapped {
		if val == "" {
			slog.Debug("[gh/compliance] invalid",
				slog.String("UID", c.UID()),
				slog.String(k, val.(string)))
			return false
		}
	}
	slog.Debug("[gh/compliance] valid", slog.String("UID", c.UID()))
	return true
}

// CompliesWithBaseline checks if this instance meets all the
// basedline criteria for compliance
func (c *Compliance) CompliesWithBaseline() bool {
	return c.DefaultBranch == "main" &&
		len(c.License) > 0 &&
		c.HasIssues &&
		c.HasDescription &&
		c.HasRulesEnforcedForAdmins &&
		c.HasPullRequestApprovalRequired
}

// CompliesWithExtended checks other features that are nice to have
func (c *Compliance) CompliesWithExtended() bool {
	return c.HasCodeOwnerApprovalRequired &&
		c.HasReadme &&
		c.HasCodeOfConduct &&
		c.HasContributingGuide
}

func (c *Compliance) Complies() (b bool, e bool) {
	c.Baseline = c.CompliesWithBaseline()
	c.Extended = c.CompliesWithExtended()
	b = c.Baseline
	e = c.Extended
	return
}

// setData calls the internal data setting funcs
func (c *Compliance) setData(client *github.Client) {
	c.dataDirectFromR()
	c.dataViaClient(client)
}

// dataViaClient sets data that requires additional calls using the
// github client to fetch more information
func (c *Compliance) dataViaClient(client *github.Client) {
	// Check branch protection
	if branch, _, err := client.Repositories.GetBranch(context.Background(), c.Owner, c.Name, c.DefaultBranch, 1); err == nil {
		c.HasDefaultBranchProtection = branch.GetProtected()
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
func (c *Compliance) dataDirectFromR() {
	r := c.r
	c.Archived = r.GetArchived()
	c.DefaultBranch = r.GetDefaultBranch()
	c.FullName = r.GetFullName()
	c.Name = r.GetName()
	c.Owner = r.GetOwner().GetLogin()
	if l := r.GetLicense(); l != nil {
		c.License = l.GetName()
	}

	c.HasDescription = len(r.GetDescription()) > 0
	c.HasDiscussions = r.GetHasDiscussions()
	c.HasDownloads = r.GetHasDownloads()
	c.HasIssues = r.GetHasIssues()
	c.HasPages = r.GetHasPages()
	c.HasProjects = r.GetHasProjects()
	c.HasWiki = r.GetHasWiki()
	c.IsPrivate = r.GetPrivate()

	c.CountForks = r.GetForksCount()
}

var _ data.IEntry = &Compliance{}

func New(uid *string) *Compliance {
	c := &Compliance{}
	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	c.Complies()
	return c
}

func NewWithR(uid *string, r *github.Repository, client *github.Client) *Compliance {
	c := New(uid)
	c.r = r
	c.setData(client)
	c.Complies()
	return c
}
