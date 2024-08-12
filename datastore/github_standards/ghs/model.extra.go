package ghs

import (
	"fmt"
)

func (g *GithubStandard) UID() string {
	return g.FullName
}

func (g *GithubStandard) UpdateCompliance() (baseline int, extended int) {
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

	return
}

func (g *GithubStandard) ToCSV() (line string) {

	line = fmt.Sprintf(`%d,"%s",%d,%d,%d,%d,%d,%d,"%s","%s","%s",%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,"%s","%s","%s","%s","%s"`,
		g.ID,
		g.Ts,
		g.CompliantBaseline,
		g.CompliantExtended,
		g.CountOfClones,
		g.CountOfForks,
		g.CountOfPullRequests,
		g.CountOfWebHooks,
		g.CreatedAt,
		g.DefaultBranch,
		g.FullName,
		g.HasCodeOfConduct,
		g.HasCodeownerApprovalRequired,
		g.HasContributingGuide,
		g.HasDefaultBranchOfMain,
		g.HasDefaultBranchProtection,
		g.HasDeleteBranchOnMerge,
		g.HasDescription,
		g.HasDiscussions,
		g.HasDownloads,
		g.HasIssues,
		g.HasLicense,
		g.HasPages,
		g.HasPullRequestApprovalRequired,
		g.HasReadme,
		g.HasRulesEnforcedForAdmins,
		g.HasVulnerabilityAlerts,
		g.HasWiki,
		g.IsArchived,
		g.IsPrivate,
		g.License,
		g.LastCommitDate,
		g.Name,
		g.Owner,
		g.Teams,
	) + "\n"
	return
}
