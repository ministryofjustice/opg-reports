package ghmergedprs

import (
	"context"
	"fmt"
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
		res          []*MergedPR
		src          *Source[*github.PullRequestsService, *MergedPR]
		ctx          context.Context      = t.Context()
		token        string               = ghclient.Token()
		repositories []*github.Repository = getRealRepos(t)
		// now          time.Time            = time.Now().UTC()
		cfg *Config = &Config{
			Repositories: repositories,
			DateStart:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			DateEnd:      time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			// DateStart:    timex.Add(timex.Reset(now, timex.DAY), timex.DAY, -7),
			// DateEnd:      timex.Reset(now, timex.DAY),
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

	src, err = New[*github.PullRequestsService, *MergedPR](ctx, client.PullRequests, cfg)
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	res, _, err = src.GetData()
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	for _, r := range res {
		fmt.Printf("[%s] - %s\n", *r.Repository.FullName, r.PullRequest.CreatedAt.Time)
	}

	t.FailNow()

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
