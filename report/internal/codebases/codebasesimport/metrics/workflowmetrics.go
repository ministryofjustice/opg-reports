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
	"strings"

	"github.com/google/go-github/v84/github"
)

type wrCounts struct {
	Total    int
	Security int
}
type wrAvg struct {
	Total   int64
	Count   int
	Average float64
}

// workflowMetrics gets all of our metrics based on workflow runs
//
// If the repository doesnt have any workflow runs within the time period then an
// error is returned. This error triggers a fallback to looking at pull request
// data instead, but there doesnt get average duration data.
func workflowMetrics(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args) (data []*CodebaseMetric, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "metrics", "func", "toCodebaseMetrics")
	var byMonth = map[string]*CodebaseMetric{}

	data = []*CodebaseMetric{}
	log.Debug("starting ... ")

	log.Info("getting metric data for repository ...", "repo", *repo.Name)
	// get try to use workflow data fore release counters
	//
	// - gets path to live workflows
	// - get count of releases & security-ish releases
	// - gets the average runtime of path to live worklows
	// byMonth, err = releaseWorkflowMetrics(ctx, client, repo, in, byMonth)
	// if err != nil {
	// 	log.Error("error getting release workflow metrics", "err", err.Error())
	// }

	// get pull request run times over the same time periods
	byMonth, err = prWorkflowMetrics(ctx, client, repo, in, byMonth)
	if err != nil {
		log.Error("error getting pr workflow metrics", "err", err.Error())
	}

	dump.Now(byMonth)

	log.Debug("complete.")
	return
}

// prWorkflowMetrics returns only the average run times for workflows
func prWorkflowMetrics(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, metrics map[string]*CodebaseMetric) (data map[string]*CodebaseMetric, err error) {
	var (
		log          *slog.Logger = cntxt.GetLogger(ctx).With("package", "metrics", "func", "releaseWorkflowMetrics", "repo", *repo.Name)
		workflowRuns []*github.WorkflowRun
		averages     map[string]*wrAvg
	)
	// set to be the same as a base line
	data = metrics

	log.Info("getting pr workflow runs ...", "repo", *repo.Name)
	workflowRuns, err = getWorkflowsForRepo(ctx, client, repo, in, false)
	if err != nil {
		log.Error("error getting path to live workflows", "err", err.Error())
		return
	}
	if len(workflowRuns) == 0 {
		err = ErrNoWorkflows
		log.Error("no workflow runs for this repo, erroring and falling back to pull requests ...")
		return
	}

	// now we just work out the average run times
	log.Info("getting release runtime metrics ...", "repo", *repo.Name)
	averages, err = averageWorkflowRunTimeByMonth(ctx, client, repo, workflowRuns)
	if err != nil {
		log.Error("error getting release avg run time", "err", err.Error())
		return
	}
	dump.Now(averages)
	// add average data
	for month, v := range averages {
		if _, ok := data[month]; !ok {
			data[month] = emptyMetric(*repo.FullName, month)
		}
		data[month].AverageTimePR = fmt.Sprintf("%g", v.Average)
	}

	return
}

// releaseWorkflowMetrics handles getting the workflow related data. If no workflows are found
// an error is triggered so the code then fallsback to using pull requests in the level above
//
// - gets path to live workflow runs
// - get count of releases & security-ish releases
// - get average run workflow run time in ms
func releaseWorkflowMetrics(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, metrics map[string]*CodebaseMetric) (data map[string]*CodebaseMetric, err error) {
	var (
		log          *slog.Logger = cntxt.GetLogger(ctx).With("package", "metrics", "func", "releaseWorkflowMetrics", "repo", *repo.Name)
		workflowRuns []*github.WorkflowRun
		counters     map[string]*wrCounts
		averages     map[string]*wrAvg
	)
	data = metrics

	log.Info("getting path to live workflow runs ...", "repo", *repo.Name)
	workflowRuns, err = getWorkflowsForRepo(ctx, client, repo, in, true)
	if err != nil {
		log.Error("error getting path to live workflows", "err", err.Error())
		return
	}
	if len(workflowRuns) == 0 {
		err = ErrNoWorkflows
		log.Error("no workflow runs for this repo, erroring and falling back to pull requests ...")
		return
	}
	log.Info("getting release metrics ...", "repo", *repo.Name)
	counters, err = workflowReleasesByMonth(ctx, client, repo, workflowRuns)
	if err != nil {
		log.Error("error getting release counters", "err", err.Error())
		return
	}

	log.Info("getting release runtime metrics ...", "repo", *repo.Name)
	averages, err = averageWorkflowRunTimeByMonth(ctx, client, repo, workflowRuns)
	if err != nil {
		log.Error("error getting release avg run time", "err", err.Error())
		return
	}

	// add release count to the map to return
	for month, v := range counters {
		if _, ok := data[month]; !ok {
			data[month] = emptyMetric(*repo.FullName, month)
		}
		data[month].Releases = v.Total
		data[month].ReleasesSecurityish = v.Security
	}

	// add average data
	for month, v := range averages {
		if _, ok := data[month]; !ok {
			data[month] = emptyMetric(*repo.FullName, month)
		}
		data[month].AverageTimeLive = fmt.Sprintf("%g", v.Average)
	}
	return
}

// averageWorkflowRunTimeByMonth works out the total duration of all workflows in milliseconds.
//
// This make an extra api call to github to find the workflow run duration (via GetWorkflowRunUsageByID),
// so is therefore a bit slower
//
// groups all workflow runs by the month and returns the average duration
func averageWorkflowRunTimeByMonth(ctx context.Context, client clients.ActionClient, repo *github.Repository, workflowRuns []*github.WorkflowRun) (byMonth map[string]*wrAvg, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "metrics", "func", "averageWorkflowRunTimeByMonth", "repo", *repo.Name)
	byMonth = map[string]*wrAvg{}

	for _, wr := range workflowRuns {
		var usage *github.WorkflowRunUsage
		var when = times.AsYMString(wr.CreatedAt.Time)
		fmt.Printf("[%s] %s %v \n", *repo.Name, when, wr.ID)

		usage, _, err = client.GetWorkflowRunUsageByID(ctx, *repo.Owner.Login, *repo.Name, *wr.ID)
		if err != nil {
			log.Error("error getting runtime of a workflow", "err", err.Error())
			return
		}

		if _, ok := byMonth[when]; !ok {
			byMonth[when] = &wrAvg{Total: 0, Count: 0, Average: 0.0}
		}
		byMonth[when].Total += *usage.RunDurationMS
		byMonth[when].Count += 1
	}

	for k, v := range byMonth {
		byMonth[k].Average = float64(v.Total) / float64(v.Count)
	}
	return
}

// workflowReleasesByMonth determines the number of releases (and security ish ones) per
// month for the workflow list passed along
//
// Securityish is determined on the commit message having keywords (security, vuln) or the
// commit author being a bot (so an automated patch)
//
// As only releases are passed to this, all are counted as a release
func workflowReleasesByMonth(ctx context.Context, client clients.ActionClient, repo *github.Repository, workflowRuns []*github.WorkflowRun) (byMonth map[string]*wrCounts, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "metrics", "func", "workflowReleasesByMonth", "repo", *repo.Name)

	byMonth = map[string]*wrCounts{}

	for _, wr := range workflowRuns {
		var security = false
		var when = times.AsYMString(wr.CreatedAt.Time)

		if _, ok := byMonth[when]; !ok {
			byMonth[when] = &wrCounts{Total: 0, Security: 0}
		}
		// should only be getting release workflows...
		byMonth[when].Total += 1
		// now look to see if its security related ...
		// - the the commit has vuln / security keywords
		msg := strings.ToLower(*wr.HeadCommit.Message)
		if strings.Contains(msg, "vuln") || strings.Contains(msg, "security") {
			security = true
		}
		// - if the commit was by a bot
		author := strings.ToLower(*wr.HeadCommit.Author.Name)
		if strings.Contains(author, "renovate") || strings.Contains(author, "dependabot") {
			security = true
		}
		if security {
			byMonth[when].Security += 1
		}
		log.Debug("processed workflow run ...", "id", *wr.ID, "when", when, "securityish", security)
	}

	return
}

// getWorkflowsForRepo returns workflow run data between the dats asked for this repo
func getWorkflowsForRepo(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, pathToLiveOnly bool) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		log            *slog.Logger                    = cntxt.GetLogger(ctx).With("package", "metrics", "func", "getWorkflowsForRepo")
		pathToLiveComp string                          = "path to live"
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
		fmt.Printf("[%s] workflow runs for [%s]\n", *repo.Name, date)
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

// paginatedWorkflowRuns iterrates over github results and merges results avoiding duplicates
func paginatedWorkflowRuns(ctx context.Context, client clients.ActionClient, repo *github.Repository, in *args.Args, opts *github.ListWorkflowRunsOptions) (workflowRuns []*github.WorkflowRun, err error) {
	var (
		page    int                           = 1
		allRuns map[int64]*github.WorkflowRun = map[int64]*github.WorkflowRun{}
		log     *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "metrics", "func", "paginatedWorkflowRuns")
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
func intervalChunks(ctx context.Context, in *args.Args) (dates []string) {
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
