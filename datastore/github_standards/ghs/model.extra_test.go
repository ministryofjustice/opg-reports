package ghs_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
)

func TestServersApiGithubStandardsCompliance(t *testing.T) {

	g := ghs.Fake(nil, nil)

	// -- baseline
	g.HasDefaultBranchOfMain = 1
	g.HasLicense = 1
	g.HasIssues = 1
	g.HasDescription = 1
	g.HasRulesEnforcedForAdmins = 1
	g.HasPullRequestApprovalRequired = 1

	if b, _ := g.UpdateCompliance(); b != 1 {
		t.Errorf("baseline failed")
	}

	// -- extended
	g.HasCodeownerApprovalRequired = 1
	g.HasReadme = 1
	g.HasCodeOfConduct = 1
	g.HasContributingGuide = 1
	if _, e := g.UpdateCompliance(); e != 1 {
		t.Errorf("ext failed")
	}
}
