package ghs

import (
	"fmt"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/fake"
)

func Fake(id *int, owner *string) (gs *GithubStandard) {
	if id == nil {
		i := fake.Int(1000, 9999)
		id = &i
	}
	if owner == nil {
		o := fake.String(12)
		owner = &o
	}
	name := fake.String(20)
	full := fmt.Sprintf("%s/%s", *owner, name)

	defTeams := []string{"foo", "bar"}
	teams := []string{"my-org", "test", "thisteam"}

	now := time.Now().UTC().Format(dates.Format)

	gs = &GithubStandard{
		ID:             *id,
		Ts:             now,
		DefaultBranch:  fake.Choice[string]([]string{"main", "master"}),
		FullName:       full,
		Name:           name,
		Owner:          *owner,
		License:        fake.Choice[string]([]string{"MIT", "GPL", ""}),
		LastCommitDate: now,
		CreatedAt:      now,
		IsArchived:     fake.Choice[int]([]int{0, 1}),
		Teams:          fmt.Sprintf("#%s#%s#", fake.Choice(teams), fake.Choice(defTeams)),
	}
	// make the odds of baseline compliance higher
	b := fake.Choice[int]([]int{0, 1, 1, 1, 1, 1, 1})
	ComplyB(gs, b)
	e := fake.Choice[int]([]int{0, 1})
	ComplyE(gs, e)

	gs.UpdateCompliance()

	return
}

func ComplyB(g *GithubStandard, truthy int) {
	g.HasDefaultBranchOfMain = truthy
	g.HasLicense = truthy
	g.HasIssues = truthy
	g.HasDescription = truthy
	g.HasRulesEnforcedForAdmins = truthy
	g.HasPullRequestApprovalRequired = truthy
}

func ComplyE(g *GithubStandard, truthy int) {
	g.HasCodeownerApprovalRequired = truthy
	g.HasReadme = truthy
	g.HasCodeOfConduct = truthy
	g.HasContributingGuide = truthy
}
