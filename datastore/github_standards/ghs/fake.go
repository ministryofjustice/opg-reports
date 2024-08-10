package ghs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ministryofjustice/opg-reports/shared/fake"
)

func Fake() (gs *GithubStandard) {
	owner := fake.String(12)
	name := fake.String(20)
	full := fmt.Sprintf("%s/%s", owner, name)

	defTeams := []string{"foo", "bar"}
	teams := []string{"my-org", "test", "thisteam"}

	gs = &GithubStandard{
		Uuid:           uuid.NewString(),
		Ts:             time.Now().UTC().String(),
		DefaultBranch:  fake.Choice[string]([]string{"main", "master"}),
		FullName:       full,
		Name:           name,
		Owner:          owner,
		License:        fake.Choice[string]([]string{"MIT", "GPL", ""}),
		LastCommitDate: time.Now().String(),
		CreatedAt:      time.Now().String(),
		IsArchived:     fake.Choice[int]([]int{0, 1}),
		Teams:          fmt.Sprintf("%s#%s#", fake.Choice(teams), fake.Choice(defTeams)),
	}

	return
}
