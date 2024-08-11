package github_standards_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
)

func TestServersApiGithubStandardsCompliance(t *testing.T) {

	g := ghs.Fake()

	// -- baseline
	g.HasDefaultBranchOfMain = 1
	g.HasLicense = 1
	g.HasIssues = 1
	g.HasDescription = 1
	g.HasRulesEnforcedForAdmins = 1
	g.HasPullRequestApprovalRequired = 1

	if b, _ := github_standards.Compliant(g); !b {
		t.Errorf("baseline failed")
	}

	// -- extended
	g.HasCodeownerApprovalRequired = 1
	g.HasReadme = 1
	g.HasCodeOfConduct = 1
	g.HasContributingGuide = 1
	if _, e := github_standards.Compliant(g); !e {
		t.Errorf("ext failed")
	}
}
