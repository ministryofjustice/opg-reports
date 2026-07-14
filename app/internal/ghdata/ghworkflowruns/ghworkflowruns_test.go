package ghworkflowruns

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/app/internal/convert"
	"opg-reports/app/internal/ghdata/ghclient"
	"opg-reports/app/internal/ghdata/ghfilters"
	"opg-reports/app/internal/ghdata/teamrepositories"
	"opg-reports/app/internal/timex"
	"testing"
	"time"

	"github.com/google/go-github/v87/github"
)

var (
	ownerA = &github.User{
		Login: convert.Ptr("ministryofjustice-test-a"),
		Name:  convert.Ptr("ministryofjustice-test-a"),
	}
	ownerB = &github.User{
		Login: convert.Ptr("ministryofjustice-test-b"),
		Name:  convert.Ptr("ministryofjustice-test-b"),
	}
	repoA = &github.Repository{
		FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *ownerA.Login, "opg-example-a")),
		Name:          convert.Ptr("opg-example-a"),
		DefaultBranch: convert.Ptr("main"),
		Archived:      convert.Ptr(false),
		Owner:         ownerA,
	}
	repoB = &github.Repository{
		FullName:      convert.Ptr(fmt.Sprintf("%s/%s", *ownerB.Login, "opg-example-b")),
		Name:          convert.Ptr("opg-example-b"),
		DefaultBranch: convert.Ptr("main"),
		Archived:      convert.Ptr(false),
		Owner:         ownerB,
	}
)

type mockActionService struct{}

func (self *mockActionService) ListRepositoryWorkflowRuns(ctx context.Context, owner string, repo string, opts *github.ListWorkflowRunsOptions) (runs *github.WorkflowRuns, response *github.Response, err error) {
	var allRuns *github.WorkflowRuns = &github.WorkflowRuns{}

	runs = &github.WorkflowRuns{}

	response = &github.Response{
		Response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		},
		NextPage: 0,
	}

	runsA := &github.WorkflowRuns{
		TotalCount: convert.Ptr(3),
		WorkflowRuns: []*github.WorkflowRun{
			{
				ID:             convert.Ptr(int64(1234908)),
				Name:           convert.Ptr("[Workflow] Path to Live a1"),
				Event:          convert.Ptr("push"),
				HeadBranch:     convert.Ptr("main"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("success"),
				Repository:     repoA,
				HeadRepository: repoA,
			},
			{
				ID:             convert.Ptr(int64(1235711)),
				Name:           convert.Ptr("[composite] Synk runner a2"),
				Event:          convert.Ptr("pull_request"),
				HeadBranch:     convert.Ptr("test_branch_name"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("success"),
				Repository:     repoA,
				HeadRepository: repoA,
			},
			{
				ID:             convert.Ptr(int64(2285711)),
				Name:           convert.Ptr("[workflow] path to live a3"),
				Event:          convert.Ptr("push"),
				HeadBranch:     convert.Ptr("main"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("success"),
				Repository:     repoA,
				HeadRepository: repoA,
			},
		},
	}

	runsB := &github.WorkflowRuns{
		TotalCount: convert.Ptr(3),
		WorkflowRuns: []*github.WorkflowRun{
			{
				ID:             convert.Ptr(int64(6295908)),
				Name:           convert.Ptr("[Workflow] Path to Live b1"),
				Event:          convert.Ptr("push"),
				HeadBranch:     convert.Ptr("main"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("success"),
				Repository:     repoB,
				HeadRepository: repoB,
			},
			{
				ID:             convert.Ptr(int64(1235711)),
				Name:           convert.Ptr("[composite] Synk runner b2"),
				Event:          convert.Ptr("pull_request"),
				HeadBranch:     convert.Ptr("test_branch_name"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("success"),
				Repository:     repoB,
				HeadRepository: repoB,
			},
			{
				ID:             convert.Ptr(int64(5128571)),
				Name:           convert.Ptr("[workflow] path to live b3"),
				Event:          convert.Ptr("push"),
				HeadBranch:     convert.Ptr("main"),
				Status:         convert.Ptr("completed"),
				Conclusion:     convert.Ptr("failed"),
				Repository:     repoB,
				HeadRepository: repoB,
			},
		},
	}

	if owner == *ownerA.Login && repo == *repoA.Name {
		allRuns = runsA
	} else if owner == *ownerB.Login && repo == *repoB.Name {
		allRuns = runsB
	}

	for _, wfr := range allRuns.WorkflowRuns {
		var add = true
		// ignore if event doesnt match
		if opts.Event != "" && opts.Event != *wfr.Event {
			add = false
		}
		// ignore if bramch mismatches
		if opts.Branch != "" && opts.Branch != *wfr.HeadBranch {
			add = false
		}
		if opts.Status != "" && opts.Status != *wfr.Conclusion {
			add = false
		}
		if add {
			runs.WorkflowRuns = append(runs.WorkflowRuns, wfr)
		}
	}
	runs.TotalCount = convert.Ptr(len(runs.WorkflowRuns))

	return
}

// TestWorkflowRunesGetDataMocked
func TestWorkflowRunsGetDataMocked(t *testing.T) {
	var (
		err          error
		client       *mockActionService = &mockActionService{}
		res          []*github.WorkflowRun
		src          *Source[*mockActionService, *github.WorkflowRun]
		ctx          context.Context      = t.Context()
		repositories []*github.Repository = []*github.Repository{repoA, repoB}
		cfg          *Config              = &Config{
			Repositories: repositories,
			Event:        "push",
			Status:       "success",
		}
	)

	src, err = New[*mockActionService, *github.WorkflowRun](ctx, client, cfg, &ghfilters.FilterWorkfowRunByPartialName{Name: "path to live"})
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
		t.FailNow()
	}
	// known results
	expected := []int64{
		int64(1234908),
		int64(2285711),
		int64(6295908),
	}

	for _, item := range res {
		var found = false
		for _, ex := range expected {
			if ex == *item.ID {
				found = true
			}
		}
		if !found {
			t.Error("failed to find expected workflow run")
		}
	}
	// fmtx.Printj(res)
	// t.FailNow()
}

// TestWorkflowRunesGetDataActual uses real api connection and client
// to fetch data
func TestWorkflowRunsGetDataActual(t *testing.T) {
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
