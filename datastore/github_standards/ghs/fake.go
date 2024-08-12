package ghs

import (
	"fmt"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/fake"
)

func Fake() (gs *GithubStandard) {
	owner := fake.String(12)
	name := fake.String(20)
	full := fmt.Sprintf("%s/%s", owner, name)

	defTeams := []string{"foo", "bar"}
	teams := []string{"my-org", "test", "thisteam"}

	now := time.Now().UTC().Format(dates.Format)

	gs = &GithubStandard{
		Ts:             now,
		DefaultBranch:  fake.Choice[string]([]string{"main", "master"}),
		FullName:       full,
		Name:           name,
		Owner:          owner,
		License:        fake.Choice[string]([]string{"MIT", "GPL", ""}),
		LastCommitDate: now,
		CreatedAt:      now,
		IsArchived:     fake.Choice[int]([]int{0, 1}),
		Teams:          fmt.Sprintf("%s#%s#", fake.Choice(teams), fake.Choice(defTeams)),
		// -- compliance
		HasDefaultBranchOfMain:         fake.Choice[int]([]int{0, 1}),
		HasLicense:                     fake.Choice[int]([]int{0, 1}),
		HasIssues:                      fake.Choice[int]([]int{0, 1}),
		HasDescription:                 fake.Choice[int]([]int{0, 1}),
		HasRulesEnforcedForAdmins:      fake.Choice[int]([]int{0, 1}),
		HasPullRequestApprovalRequired: fake.Choice[int]([]int{0, 1}),
		HasCodeownerApprovalRequired:   fake.Choice[int]([]int{0, 1}),
		HasReadme:                      fake.Choice[int]([]int{0, 1}),
		HasCodeOfConduct:               fake.Choice[int]([]int{0, 1}),
		HasContributingGuide:           fake.Choice[int]([]int{0, 1}),
	}

	gs.UpdateCompliance()

	return
}
