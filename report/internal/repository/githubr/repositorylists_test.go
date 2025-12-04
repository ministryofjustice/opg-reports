package githubr

import (
	"context"
	"opg-reports/report/config"
	"opg-reports/report/internal/utils"
	"testing"

	"github.com/google/go-github/v77/github"
)

var mockRepos = []*github.Repository{
	{
		FullName: utils.Ptr("test/A"),
		Archived: utils.Ptr(true),
	},
	{
		FullName: utils.Ptr("test/B"),
		Archived: utils.Ptr(false),
	},
	{
		FullName: utils.Ptr("test/C"),
		Archived: utils.Ptr(false),
	},
}

type mockClientTeamListRepositories struct{}

func (m *mockClientTeamListRepositories) ListTeamReposBySlug(ctx context.Context, org string, team string, opts *github.ListOptions) (repos []*github.Repository, resp *github.Response, err error) {
	resp = &github.Response{NextPage: 0}
	repos = mockRepos
	err = nil
	return
}

// TestGithubrGetRepositoriesForTeam tests that can find service
func TestGithubrGetRepositoriesForTeam(t *testing.T) {

	var (
		err   error
		res   []*github.Repository
		ctx   = context.TODO()
		log   = utils.Logger("ERROR", "TEXT")
		conf  = config.NewConfig()
		owner = "ministryofjustice"
		team  = "opg"
		rp    = Default(ctx, log, conf)
		// client = DefaultClient(conf).Teams
		client = &mockClientTeamListRepositories{}
	)

	// check without fitlering
	res, err = rp.GetRepositoriesForTeam(client, owner, team, &GetRepositoriesForTeamOptions{ExcludeArchived: false})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 3 {
		t.Errorf("expected all 3 repos to return")
	}

	// check with fitlering
	res, err = rp.GetRepositoriesForTeam(client, owner, team, &GetRepositoriesForTeamOptions{ExcludeArchived: true})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("expected only 2 repos to return")
	}

}
