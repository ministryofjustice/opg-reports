package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasesimport/args"
	"opg-reports/report/internal/codebases/codebasesimport/clients"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dump"
	"opg-reports/report/package/times"
	"slices"

	"github.com/google/go-github/v84/github"
)

// Raw stats entry
const InsertMetricsStatement string = `
INSERT INTO codebase_metrics (
	codebase,
	month,
	releases,
	releases_securityish,
	pr_count,
	pr_count_securityish,
	pr_stale_count,
	average_time_live,
	average_time_pr
) VALUES (
	:codebase,
	:month,
	:releases,
	:releases_securityish,
	:pr_count,
	:pr_count_securityish,
	:pr_stale_count,
	:average_time_live,
	:average_time_pr
)
ON CONFLICT (codebase,month) DO UPDATE SET
	releases=excluded.releases,
	releases_securityish=excluded.releases_securityish,
	pr_count=excluded.pr_count,
	pr_count_securityish=excluded.pr_count_securityish,
	pr_stale_count=excluded.pr_stale_count,
	average_time_live=excluded.average_time_live,
	average_time_live=excluded.average_time_live
RETURNING id
;
`

type CodebaseMetric struct {
	Codebase            string `json:"codebase,omitempty"`    // full name of codebase
	Month               string `json:"month,omitempty"`       // month as YYYY-MM string
	Releases            int    `json:"releases,omitempty"`    // count of releases for this month
	ReleasesSecurityish int    `json:"securityish,omitempty"` // count of releases for this month that seem to be security related
	PRCount             int    `json:"pr_count,omitempty"`    // count of all pull requests for the month
	PRCountSecurityish  int    `json:"pr_count_securityish"`  // count of all pr's that roughly relate to security (bots / keywords)
	PRStaleCount        int    `json:"pr_count_stale"`        // count of stale pull requests - open for longer than x days
	AverageTimeLive     string `json:"average_time_live"`     // average time path to live workflow took
	AverageTimePR       string `json:"average_time_pr"`       // average time a pull request workflow took
}

func HandleCodebaseMetrics(ctx context.Context, client clients.ActionClient, repositories []*github.Repository, in *args.Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "HandleCodebaseMetrics")
	var data []*CodebaseMetric = []*CodebaseMetric{}
	log.With("count", len(repositories)).Info("starting codebase metrics import ...")

	toCodebaseMetrics(ctx, client, repositories, in)

	log.With("count", len(data)).Info("complete.")
	return
}

func toCodebaseMetrics(ctx context.Context, client clients.ActionClient, list []*github.Repository, in *args.Args) (data []*CodebaseMetric, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toCodebaseMetrics")

	data = []*CodebaseMetric{}
	log.Debug("starting ... ")

	for _, repo := range list {
		var pathToLiveWorkflowRuns []*github.WorkflowRun

		// path to live workflows are used to work out....
		// - releases
		// - releases_securityish
		// - average_time_live
		pathToLiveWorkflowRuns, err = getWorkflowsForRepo(ctx, client, repo, in, true)
		if err != nil {
			log.Error("error getting path to live workflows", pathToLiveWorkflowRuns)
			return
		}
		releasesViaWorkflowRuns(ctx, client, repo, pathToLiveWorkflowRuns)
		// averageTimeToLive

		// pr workflows are used to work out...
		// - average_time_pr
		// var prWorkflowRuns []*github.WorkflowRun
		// prWorkflowRuns, err = getWorkflowsForRepo(ctx, client, repo, in, false)
	}

	log.Debug("complete.")
	return
}

// releasesViaWorkflowRuns tries to get release counts from workflow runs
func releasesViaWorkflowRuns(ctx context.Context, client clients.ActionClient, repo *github.Repository, workflowRuns []*github.WorkflowRun) (err error) {
	// var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "viaWorkflowRuns")
	// var workflowRuns = []*github.WorkflowRun{}

	return
}

// getWorkflowsForRepo returns workflow run data between the dats asked for this repo
func getWorkflowsForRepo(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, pathToLiveOnly bool) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		log           *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "getWorkflowsForRepo")
		dateIntervals []string     = []string{}
		exlcudePR     bool         = pathToLiveOnly
		event         string       = "push"
	)
	if !pathToLiveOnly {
		event = "pull_request"
	}
	// use date intervals to call the api in smaller, 3 day chunks
	// as the api result will contain a max of 1k, regardless of
	// pagination
	dateIntervals = weekIntervals(ctx, in)
	workflowRuns = []*github.WorkflowRun{}

	log.Info("fetching workflow run data ...", "repository", *repo.Name)
	// generate the created string for the data range, doing one call per
	for _, date := range dateIntervals {
		var wr = []*github.WorkflowRun{}
		log.Info("date range ... ", "range", date)
		// get just the releases for this time period
		wr, err = paginatedWorkflowRuns(ctx, client, repo, in, &github.ListWorkflowRunsOptions{
			Branch:              *repo.DefaultBranch,
			ExcludePullRequests: exlcudePR,
			Event:               event,
			Status:              "success",
			Created:             date,
		})
		if err != nil {
			log.Error("error getting workflow runs", "err", err.Error())
			return
		}
		//
		// merge in runs
		workflowRuns = append(workflowRuns, wr...)
	}
	return
}

func paginatedWorkflowRuns(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, opts *github.ListWorkflowRunsOptions) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		allRuns map[int64]*github.WorkflowRun = map[int64]*github.WorkflowRun{}
		page    int                           = 1
		log     *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "paginatedWorkflowRuns")
	)
	workflowRuns = []*github.WorkflowRun{}
	// force the max per page
	opts.PerPage = 100

	for page > 0 {
		var runs *github.WorkflowRuns
		var response *github.Response
		log.Info("getting page of results", "page", page)

		opts.Page = page
		runs, response, err = client.ListRepositoryWorkflowRuns(ctx, *repo.Owner.Login, *repo.Name, opts)
		if err != nil {
			log.Error("error getting workflow runs", "err", err.Error())
			return
		}
		// all runs should have unique id
		for _, wr := range runs.WorkflowRuns {
			allRuns[*wr.ID] = wr
		}
		page = response.NextPage
	}
	// push from map to slice
	for _, wr := range allRuns {
		workflowRuns = append(workflowRuns, wr)
	}
	dump.Now(workflowRuns)
	return
}

func weekIntervals(ctx context.Context, in *args.Args) (dates []string) {
	var chunks = times.DaysN(in.DateStart, in.DateEnd, 3)
	var l = len(chunks)

	dates = []string{}
	for i, date := range chunks {
		var end = date
		if i+1 < l {
			end = chunks[i+1]
		}
		dates = append(dates, fmt.Sprintf("%s..%s", times.AsYMDString(date), times.AsYMDString(end)))
	}
	slices.Sort(dates)
	slices.Compact(dates)
	return
}
