package liveworkflowruns

import (
	"context"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"opg-reports/app/internal/ghdata/teamrepositories"
	"opg-reports/app/internal/timex"
	"testing"
	"time"

	"github.com/google/go-github/v87/github"
)

// TestWorkflowRunesGetDataActual uses real api connection and client
// to fetch data
func TestWorkflowRunesGetDataActual(t *testing.T) {
	var (
		// skipped      []any
		err          error
		client       *github.Client
		res          []*github.WorkflowRun
		src          *Source[*github.ActionsService, *github.WorkflowRun]
		now          time.Time            = time.Now().UTC()
		ctx          context.Context      = t.Context()
		token        string               = ghclient.Token()
		repositories []*github.Repository = getRealRepo(t)
		cfg          *Config              = &Config{
			Repositories: repositories,
			Event:        "push",
			Status:       "success",
			DateStart:    timex.Add(timex.Reset(now, timex.DAY), timex.DAY, -7),
			DateEnd:      timex.Reset(now, timex.DAY),
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

	src, err = New[*github.ActionsService, *github.WorkflowRun](ctx, client.Actions, cfg, &ghfilters.FilterWorkfowRunByPartialName{Name: "path to live"})
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	res, _, err = src.GetData()
	if err != nil {
		t.Errorf("unexpected error creating source: %s", err.Error())
		t.FailNow()
	}

	if len(res) <= 0 {
		t.Errorf("failed to find any workflow runs.")
	}

	// for _, r := range res {
	// 	fmt.Println(*r.Repository.FullName, *r.ID, *r.Name)
	// }

	// for _, s := range skipped {
	// 	var r = s.(*github.Repository)
	// 	fmt.Println(*r.FullName)
	// }

	// t.FailNow()
}

// returns just a sub-set of repos for testing
func getRealRepo(t *testing.T) (repos []*github.Repository) {
	var (
		err    error
		client *github.Client
		src    *teamrepositories.Source[*github.TeamsService, *github.Repository]
		token  string                   = ghclient.Token()
		ctx    context.Context          = context.TODO()
		cfg    *teamrepositories.Config = &teamrepositories.Config{
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
	src, err = teamrepositories.New[*github.TeamsService, *github.Repository](
		ctx,
		client.Teams,
		cfg,
		&ghfilters.FilterByRepositoryName{Name: "opg-use-an-lpa"},
	)
	// test fetching the data
	repos, _, err = src.GetData()
	return
}
