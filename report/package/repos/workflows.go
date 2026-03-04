package repos

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/times"
	"slices"
	"strings"

	"github.com/google/go-github/v84/github"
)

// GetWorkflowRuns fetches workflows for the reporsitory between the dates stipulated via
// `DateStart` & `DateEnd`.
//
// `pathToLiveOnly` detemines if we are fetching only workflows that ran the main branch called
// path to live. If its false then pull requests are fetched as well - which is a heavier
// operation.
func GetWorkflowRuns(ctx context.Context, client actionClient, repo *github.Repository, in *Args, pathToLiveOnly bool) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		log            *slog.Logger                    = cntxt.GetLogger(ctx).With("package", "reports", "func", "GetWorkflows")
		pathToLiveComp string                          = in.FilterByName
		dateIntervals  []string                        = []string{}
		opts           *github.ListWorkflowRunsOptions = &github.ListWorkflowRunsOptions{
			ExcludePullRequests: pathToLiveOnly,
			Event:               "pull_request",
		}
	)
	fmt.Println("getting workflow runs ...")
	// filter event to push for path to live
	if pathToLiveOnly {
		opts.Branch = *repo.DefaultBranch
		opts.Event = "push"
		opts.Status = "success"
	}
	// use date intervals to call the api in smaller, 3 day chunks
	// as the api result will contain a max of 1k, regardless of
	// pagination
	dateIntervals = intervalChunks(ctx, in)
	workflowRuns = []*github.WorkflowRun{}

	log.Info("fetching workflow runs ...", "repository", *repo.Name, "date_start", in.DateStart, "date_end", in.DateEnd)
	// generate the created string for the data range, doing one call per
	for _, date := range dateIntervals {
		var wr = []*github.WorkflowRun{}
		log.Debug("date range ... ", "range", date)
		opts.Created = date
		// fmt.Printf("[%s] workflow runs for [%s]\n", *repo.Name, date)
		// get just the releases for this time period
		wr, err = paginatedWorkflowRuns(ctx, client, repo, in, opts)
		if err != nil {
			log.Error("error getting workflow runs", "err", err.Error())
			return
		}
		log.Debug("found workflows ...", "count", len(wr))
		// filter out path to live if asked to
		for _, r := range wr {
			name := strings.ToLower(*r.Name)
			if !pathToLiveOnly || (pathToLiveOnly && strings.Contains(name, pathToLiveComp)) {
				workflowRuns = append(workflowRuns, r)
			}
		}
	}
	log.Debug("complete.", "count", len(workflowRuns))
	return
}

// paginatedWorkflowRuns iterrates over github results and merges results
// avoiding duplicates by checking the the workflow run id
func paginatedWorkflowRuns(ctx context.Context, client actionClient, repo *github.Repository, in *Args, opts *github.ListWorkflowRunsOptions) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		page    int                           = 1
		allRuns map[int64]*github.WorkflowRun = map[int64]*github.WorkflowRun{}
		log     *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "repos", "func", "paginatedWorkflowRuns")
	)
	workflowRuns = []*github.WorkflowRun{}
	// force the max per page
	opts.PerPage = 100
	log.Debug("starting ....")
	for page > 0 {
		var runs *github.WorkflowRuns
		var response *github.Response
		log.Debug("getting page of results ...", "page", page)

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
	log.With("count", len(workflowRuns)).Debug("complete.")
	return
}

// intervalChunks generates smaller chunk of date ranges as the api call is limited
// to 1k results in one go, so need to make the window small enough to find them.
//
// use 3 days as chunk window, should cover the most frequently used repos
func intervalChunks(ctx context.Context, in *Args) (dates []string) {

	var chunks = times.DaysN(in.DateStart, in.DateEnd, 3)
	var l = len(chunks)
	var list = []string{}

	for i, date := range chunks {
		var end = date
		if i+1 < l {
			end = chunks[i+1]
		}
		list = append(list, fmt.Sprintf("%s..%s", times.AsYMDString(date), times.AsYMDString(end)))
	}
	slices.Sort(list)
	slices.Compact(list)
	// remove any empty values
	for _, item := range list {
		if item != "" {
			dates = append(dates, item)
		}
	}
	return
}
