package githubr

import (
	"context"
	"opg-reports/report/config"
	"opg-reports/report/internal/utils"
	"testing"

	"github.com/google/go-github/v75/github"
)

func TestGithubrGetRepositoryOwners(t *testing.T) {
	var (
		err    error
		res    []string
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		rp     = Default(ctx, log, conf)
		client = DefaultClient(conf).Repositories
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-modernising-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-modernising-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
		// client = &mockClientTeamListRepositories{}
	)

	res, err = rp.GetRepositoryOwners(client, repo, &GetTeamsForRepositoryOptions{FilterByParent: "opg"})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	utils.Dump(res)

	t.FailNow()
}

func TestGithubrGetTeamsForRepository(t *testing.T) {

	var (
		err    error
		res    []*github.Team
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		rp     = Default(ctx, log, conf)
		client = DefaultClient(conf).Repositories
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-use-an-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-use-an-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
		// client = &mockClientTeamListRepositories{}
	)

	res, err = rp.GetTeamsForRepository(client, repo, &GetTeamsForRepositoryOptions{FilterByParent: "opg"})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	utils.Dump(len(res))

	// t.FailNow()
}

func TestGithubrGetCodeOwnersForRepository(t *testing.T) {
	var (
		err    error
		res    []string
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		rp     = Default(ctx, log, conf)
		client = DefaultClient(conf).Repositories
		repo   = &github.Repository{
			Name:     utils.Ptr("opg-modernising-lpa"),
			FullName: utils.Ptr("ministryofjustice/opg-modernising-lpa"),
			Owner: &github.User{
				Login: utils.Ptr("ministryofjustice"),
			},
		}
		// client = &mockClientTeamListRepositories{}
	)

	res, err = rp.GetCodeOwnersForRepository(client, repo)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	utils.Dump(res)

	// t.FailNow()
}

type testCodeOwnerLines struct {
	Lines    []string
	Expected int
}

func TestGithubrCodeOwnersFromLines(t *testing.T) {

	var tests = []*testCodeOwnerLines{
		{
			Expected: 2,
			Lines: []string{
				"* @ministryofjustice/opg @ministryofjustice/opg-webops",
				"/.github/  @ministryofjustice/opg-webops",
			},
		},
		{
			Expected: 1,
			Lines: []string{
				"* @ministryofjustice/opg-webops",
			},
		},
	}

	for _, test := range tests {
		res := codeOwnersFromLines(test.Lines)
		if test.Expected != len(res) {
			t.Errorf("unexpected number of codeowners returned expected [%v] actual [%v]", test.Expected, len(res))
		}
	}

}
