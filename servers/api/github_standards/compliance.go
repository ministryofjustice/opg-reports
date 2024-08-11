package github_standards

import (
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
)

func Compliant(g ghs.GithubStandard) (baseline bool, extended bool) {
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

	baseline = true
	extended = true

	for _, is := range baselineChecks {
		if !is {
			baseline = false
		}
	}
	for _, is := range extendedChecks {
		if !is {
			extended = false
		}
	}

	return
}
