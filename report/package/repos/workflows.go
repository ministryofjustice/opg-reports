package repos

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/times"
	"slices"
	"strings"
	"time"

	"github.com/google/go-github/v84/github"
)

// GetWorkflowRuns fetches workflows for the reporsitory between the dates stipulated via
// `DateStart` & `DateEnd`.
//
// `pathToLiveOnly` detemines if we are fetching only workflows that ran the main branch called
// path to live. If its false then pull requests are fetched as well - which is a heavier
// operation.
//
// use a map based on workflow id to avoid duplicates in date ranges
//
// note: the workflow run api treats the end date as and up to and including, so we
// change that to second before (previous day)
func GetWorkflowRuns(ctx context.Context, client actionClient, repo *github.Repository, in *Args, pathToLiveOnly bool) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		log            *slog.Logger                    = cntxt.GetLogger(ctx).With("package", "repos", "func", "GetWorkflowRuns", "repo", *repo.Name)
		pathToLiveComp string                          = in.FilterByName
		dateIntervals  []string                        = []string{}
		byID           map[int64]*github.WorkflowRun   = map[int64]*github.WorkflowRun{}
		opts           *github.ListWorkflowRunsOptions = &github.ListWorkflowRunsOptions{
			ExcludePullRequests: pathToLiveOnly,
			Event:               "pull_request",
		}
	)
	// the workflow run api treats the end date as and up to and including, so we need to change that to the day
	// before
	in.DateEnd = times.Add(in.DateEnd, -1, times.SECOND)

	log.Debug("getting workflow runs ...", "date_start", times.AsYMDString(in.DateStart), "date_end", times.AsYMDString(in.DateEnd))
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

	// generate the created string for the data range, doing one call per
	for _, date := range dateIntervals {
		var wr = []*github.WorkflowRun{}
		log.Debug("date range ... ", "range", date)
		opts.Created = date
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
				byID[*r.ID] = r
			}
		}
	}
	//
	for _, wr := range byID {
		workflowRuns = append(workflowRuns, wr)
	}

	log.With("count", len(workflowRuns)).Debug("complete.")
	return
}

// paginatedWorkflowRuns iterrates over github results and merges results
// avoiding duplicates by checking the the workflow run id
func paginatedWorkflowRuns(ctx context.Context, client actionClient, repo *github.Repository, in *Args, opts *github.ListWorkflowRunsOptions) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		maxRetry int                           = 3
		page     int                           = 1
		allRuns  map[int64]*github.WorkflowRun = map[int64]*github.WorkflowRun{}
		log      *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "repos", "func", "paginatedWorkflowRuns", "repo", *repo.Name)
	)
	workflowRuns = []*github.WorkflowRun{}
	// force the max per page
	opts.PerPage = 100
	log.Debug("starting ....")
	for page > 0 {
		var runs *github.WorkflowRuns
		var response *github.Response
		var retry = 0
		log.Debug("getting page of results ...", "page", page)

		opts.Page = page
		runs, response, err = client.ListRepositoryWorkflowRuns(ctx, *repo.Owner.Login, *repo.Name, opts)

		// simple re-try loop as we get sporadic failures
		for err != nil && retry < maxRetry {
			retry += 1
			log.Warn("error getting pull request data, retrying in 1 second ...", "err", err.Error())
			time.Sleep(time.Second * 1)
			runs, response, err = client.ListRepositoryWorkflowRuns(ctx, *repo.Owner.Login, *repo.Name, opts)
		}

		if err != nil {
			log.Warn("error getting workflow runs, retrying in 1 second.", "err", err.Error())
			time.Sleep(time.Second * 1)
			runs, response, err = client.ListRepositoryWorkflowRuns(ctx, *repo.Owner.Login, *repo.Name, opts)
		}

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
// use 4 days as chunk window, should cover the most frequently used repos
func intervalChunks(ctx context.Context, in *Args) (dates []string) {

	var chunks = times.DaysN(in.DateStart, in.DateEnd, 4)
	var l = len(chunks)
	var list = []string{}

	for i, date := range chunks {
		var end = date
		if i+1 < l {
			// reduce the end date by 1 as the range is upto & including on the api
			end = times.Add(chunks[i+1], -1, times.SECOND)
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
