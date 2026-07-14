package ghmergedprs

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/app/internal/convert"
	"opg-reports/app/internal/fmtx"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"opg-reports/app/internal/ghdata/ghteamrepositories"
	"testing"
	"time"

	"github.com/google/go-github/v87/github"
)

var (
	ownerA = &github.User{
		Login: convert.Ptr("ministryofjustice-test-a"),
		Name:  convert.Ptr("ministryofjustice-test-a"),
	}
	repoA = &github.Repository{
		FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *ownerA.Login, "opg-example-a")),
		Name:          convert.Ptr("opg-example-a"),
		DefaultBranch: convert.Ptr("main"),
		Archived:      convert.Ptr(false),
		Owner:         ownerA,
	}
)

type mockPullRequestsService struct{}

func (self *mockPullRequestsService) List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) (prs []*github.PullRequest, resp *github.Response, err error) {

	resp = &github.Response{
		Response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		},
		NextPage: 0,
	}

	prs = []*github.PullRequest{
		// august
		{
			ID:    convert.Ptr(int64(2026080101)),
			State: convert.Ptr("closed"),
			CreatedAt: &github.Timestamp{
				Time: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			},
			MergeCommitSHA: convert.Ptr("AUGUST1"),
			MergedAt: &github.Timestamp{
				Time: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		// july 15th
		{
			ID:    convert.Ptr(int64(2026070501)),
			State: convert.Ptr("closed"),
			CreatedAt: &github.Timestamp{
				Time: time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC),
			},
			MergeCommitSHA: convert.Ptr("JULY15"),
			MergedAt: &github.Timestamp{
				Time: time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
			},
		},
		// july 3rd
		{
			ID:    convert.Ptr(int64(2026070301)),
			State: convert.Ptr("closed"),
			CreatedAt: &github.Timestamp{
				Time: time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC),
			},
			MergeCommitSHA: convert.Ptr("JULY3"),
			MergedAt: &github.Timestamp{
				Time: time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
			},
		},
		// june 29th
		{
			ID:    convert.Ptr(int64(2026070301)),
			State: convert.Ptr("closed"),
			CreatedAt: &github.Timestamp{
				Time: time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC),
			},
			MergeCommitSHA: convert.Ptr("JUNE"),
			MergedAt: &github.Timestamp{
				Time: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	return
}

func TestMergedPrsGetDataMocked(t *testing.T) {
	var (
		err          error
		client       *mockPullRequestsService = &mockPullRequestsService{}
		res          []*ResultData
		src          *Source[*mockPullRequestsService, *ResultData]
		ctx          context.Context      = t.Context()
		repositories []*github.Repository = []*github.Repository{repoA}
		cfg          *Config              = &Config{
			// 2026-07-01 - 2026-08-01 = should be 2 pr's in that time
			DateStart:    time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			DateEnd:      time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			Repositories: repositories,
		}
	)

	src, err = New[*mockPullRequestsService, *ResultData](ctx, client, cfg)
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	res, _, err = src.GetData()
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	if len(res) != 2 {
		t.Errorf("expected more pull requests in time period")
		fmtx.Printj(res)
	}
}

func TestMergedPrsGetDataActual(t *testing.T) {
	var (
		// skipped      []any
		err          error
		client       *github.Client
		res          []*ResultData
		src          *Source[*github.PullRequestsService, *ResultData]
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

	src, err = New[*github.PullRequestsService, *ResultData](ctx, client.PullRequests, cfg)
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
