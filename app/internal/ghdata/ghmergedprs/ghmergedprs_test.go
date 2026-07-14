package ghmergedprs

import (
	"context"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"opg-reports/app/internal/ghdata/ghteamrepositories"
	"testing"
	"time"

	"github.com/google/go-github/v87/github"
)

type mockPullRequestsService struct{}

func (self *mockPullRequestsService) List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) (prs []*github.PullRequest, resp *github.Response, err error) {

	return
}

func TestMergedPrsGetDataActual(t *testing.T) {
	var (
		// skipped      []any
		err          error
		client       *github.Client
		res          []*MergedPullRequest
		src          *Source[*github.PullRequestsService, *MergedPullRequest]
		ctx          context.Context      = t.Context()
		token        string               = ghclient.Token()
		repositories []*github.Repository = getRealRepos(t)
		cfg          *Config              = &Config{
			// 2026-06-01 - 2026-07-01 = should be 2 pr's in that time
			DateStart:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			DateEnd:      time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			Repositories: repositories,
		}
	)

	// if theres no github token, skip this test
	if token == "" {
		t.SkipNow()
	}

	// create the client
	client, err = ghclient.New(ctx, token)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	src, err = New[*github.PullRequestsService, *MergedPullRequest](ctx, client.PullRequests, cfg)
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	res, _, err = src.GetData()
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	if len(res) < 2 {
		t.Errorf("expected more pull requests in time period")
	}

	// for _, r := range res {
	// 	fmt.Printf("[%s] - %s\n", *r.Repository.FullName, r.PullRequest.CreatedAt.Time)
	// }
	// t.FailNow()

}

// returns just one repo for testing
func getRealRepos(t *testing.T) (repos []*github.Repository) {
	var (
		err    error
		client *github.Client
		src    *ghteamrepositories.Source[*github.TeamsService, *github.Repository]
		token  string                     = ghclient.Token()
		ctx    context.Context            = context.TODO()
		cfg    *ghteamrepositories.Config = &ghteamrepositories.Config{
			OrganisationSlug: "ministryofjustice",
			TeamSlug:         "opg",
		}
	)
	// if theres no github token, skip this test
	if token == "" {
		t.SkipNow()
	}
	// create the client
	client, err = ghclient.New(ctx, token)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	// now test this with a filter removes archived values
	src, err = ghteamrepositories.New[*github.TeamsService, *github.Repository](
		ctx,
		client.Teams,
		cfg,
		&ghfilters.FilterByRepositoryName{Name: "opg-reports"},
	)
	// test fetching the data
	repos, _, err = src.GetData()
	return
}
