package compliance

import (
	"context"
	"log/slog"
	"opg-reports/shared/data"
	"opg-reports/shared/gh"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
)

const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

type BaseCompliance struct {
	DefaultBranchIsCalledMain bool `json:"default_branch_is_called_main"`
	IssuesAreEnabled          bool `json:"issues_are_enabled"`
	HasADescription           bool `json:"has_a_description"`
	HasALicense               bool `json:"has_a_license"`
	RulesEnforcedForAdmins    bool `json:"rules_enforced_for_admins"`
	RequiresApproval          bool `json:"requires_approval"`
}

func (b *BaseCompliance) Comp() bool {
	return b.DefaultBranchIsCalledMain &&
		b.IssuesAreEnabled &&
		b.HasADescription &&
		b.HasALicense &&
		b.RulesEnforcedForAdmins &&
		b.RequiresApproval
}

type ExtendedCompliance struct {
	RequiresCodeOwnerApproval     bool `json:"requires_code_owner_approval"`
	VulnerabilityAlertsAreEnabled bool `json:"vulnerability_alerts_are_enabled"`
	ReadmeIsPresent               bool `json:"readme_is_present"`
	CodeOfConductIsPresent        bool `json:"code_of_conduct_is_present"`
	ContributingGuideIsPresent    bool `json:"contributing_guide_is_present"`
}

func (b *ExtendedCompliance) Comp() bool {
	return b.RequiresCodeOwnerApproval &&
		b.VulnerabilityAlertsAreEnabled &&
		b.ReadmeIsPresent &&
		b.CodeOfConductIsPresent &&
		b.ContributingGuideIsPresent
}

type Compliance struct {
	BaseCompliance
	ExtendedCompliance
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	FullName          string `json:"full_name"`
	DefaultBranchName string `json:"default_branch_name"`
	LicenseName       string `json:"license"`
	Baseline          bool   `json:"baseline_compliant"`
	Extended          bool   `json:"extended_compliant"`
}

// UID is the unique id (UUID) for this Compliance item
func (i *Compliance) UID() string {
	slog.Debug("[gh/compliance] UID()", slog.String("UID", i.UUID))
	return i.UUID
}

// Valid returns true only if all fields are present with a non-empty value
func (i *Compliance) Valid() (valid bool) {
	mapped, _ := data.ToMap(i)
	for k, val := range mapped {
		if val == "" {
			slog.Debug("[gh/compliance] invalid", slog.String("UID", i.UID()), slog.String(k, val))
			return false
		}
	}
	slog.Debug("[gh/compliance] valid", slog.String("UID", i.UID()))
	return true
}

func (i *Compliance) BaselineComp() bool {
	return i.BaseCompliance.Comp()
}
func (i *Compliance) ExtendedComp() bool {
	return i.ExtendedCompliance.Comp()
}

var _ data.IEntry = &Compliance{}

func New(uid *string) *Compliance {
	c := &Compliance{}
	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	return c
}

func NewFromR(ctx context.Context, uid *string, r *github.Repository, client *github.Client) *Compliance {
	c := New(uid)
	c.Name = r.GetName()
	c.FullName = r.GetFullName()
	c.DefaultBranchName = r.GetDefaultBranch()

	// Set them all to false
	c.DefaultBranchIsCalledMain = false
	c.IssuesAreEnabled = false
	c.HasADescription = false
	c.HasALicense = false
	c.RulesEnforcedForAdmins = false
	c.RequiresApproval = false

	c.RequiresCodeOwnerApproval = false
	c.VulnerabilityAlertsAreEnabled = false
	c.ReadmeIsPresent = false
	c.CodeOfConductIsPresent = false
	c.CodeOfConductIsPresent = false
	c.VulnerabilityAlertsAreEnabled = false
	c.ReadmeIsPresent = false
	c.CodeOfConductIsPresent = false
	c.ContributingGuideIsPresent = false

	// Set them to real values
	c.DefaultBranchIsCalledMain = (c.DefaultBranchName == "main")
	c.IssuesAreEnabled = r.GetHasIssues()
	c.HasADescription = len(r.GetDescription()) > 0

	if l := r.GetLicense(); l != nil {
		c.LicenseName = l.GetName()
		c.HasALicense = true
	}
	owner := r.GetOwner().GetName()
	branch, err := gh.Branch(ctx, client, owner, c.FullName, c.DefaultBranchName)
	if err == nil {
		protection := branch.GetProtection()
		if protection != nil {
			c.RulesEnforcedForAdmins = protection.GetEnforceAdmins().Enabled
			c.RequiresApproval = protection.GetRequiredPullRequestReviews().RequiredApprovingReviewCount > 0
			c.RequiresCodeOwnerApproval = protection.GetRequiredPullRequestReviews().RequireCodeOwnerReviews
		}
	}

	if ve, err := gh.HasVulnerabilityAlerts(ctx, client, owner, c.FullName); err == nil {
		c.VulnerabilityAlertsAreEnabled = ve
	}
	if re, err := gh.HasReadme(ctx, client, owner, c.FullName); err == nil {
		c.ReadmeIsPresent = re
	}
	if cc, err := gh.HasFile(ctx, client, owner, c.FullName, codeOfConductPath); err != nil {
		c.CodeOfConductIsPresent = cc
	}
	if cc, err := gh.HasFile(ctx, client, owner, c.FullName, contributingGuidePath); err != nil {
		c.ContributingGuideIsPresent = cc
	}
	c.Baseline = c.BaselineComp()
	c.Extended = c.ExtendedComp()

	return c
}
