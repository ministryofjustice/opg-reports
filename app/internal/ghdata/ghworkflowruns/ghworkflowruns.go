// Package ghworkflowruns provides a struct and method for fetching all workflow runs that
// ran against the default branch between two dates based on the configured event &
// status type (set via *Config).
//
// Intention is this data is used to determine the number releases during a given time period.
//
// Pagination of the api call is handled within this package.
//
// Filtering is done at the repository level, after all workflow runs are pulled.
package ghworkflowruns

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/app/internal/logx"
	"opg-reports/app/internal/timex"
	"slices"
	"time"

	"github.com/google/go-github/v87/github"
)

// ErrGettingList indicates the ListTeamReposBySlug function returned an error for
// the requested page
var ErrGettingList = errors.New("error getting page of workflow runes from api")

// ErrNoRepositoriesConfigured is triggered when the Config passed to New does not have any
// repositories attached
var ErrNoRepositoriesConfigured = errors.New("no repositories have been set on the configuration struct.")

// errDefaultLoop is a dummy error used to handle the fail & retry loop within the paginated
// api call
var errDefaultLoop = errors.New("dummy error for retry loop logic")

// Result is an alias for *github.WorkflowRuns
type Result interface {
	*github.WorkflowRun
}

// Client is an interface for *github.ActionsService
type Client interface {
	ListRepositoryWorkflowRuns(ctx context.Context, owner string, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
}

// Filter iterface is used for running functions against a set of results
//
// Generally should be:
//   - ghfilters.FilterWorkfowRunByPartialName
type Filter interface {
	Filter(ctx context.Context, result *github.WorkflowRun) (include bool)
}

// Config is used to capture information needed for the api/sdk call
type Config struct {
	Repositories    []*github.Repository // list of all repositories we want to fetch workflow runs for
	Event           string               // the event type to filter against - generally "push".
	Status          string               // the workflow result status to filter on - generally "success".
	OverwriteBranch string               // if present, will use this instead of the default branch of the repo
	DateStart       time.Time            // start of date range to fetch workflow runs for each repo.
	DateEnd         time.Time            // end of date range to fetch workflow runs for each repo.
}

// Dates converts the start & end date in to a list of strings (each
// formatted as `YYYY-MM-DD..YYYY-MM-DD`).
//
// Required as the api endpoint being called has a hard limit on the number
// of results (max 1000) that can be returned for any given call, which is not
// enough for the busier repos or a longer time frame.
//
// To work around this limit we chop the date range into 4 day chunks, which
// will then be iterated over
func (self *Config) Dates() (dates []string) {
	var (
		dateRange []time.Time = timex.Range(self.DateStart, self.DateEnd, timex.DAY, 4) // create 4 day chunks
		count     int         = len(dateRange)                                          // number of dates ranges
		list      []string    = []string{}                                              // temp list of ranges
	)
	dates = []string{}

	// now loop over and create the string ranges
	for i, date := range dateRange {
		var endDate = date
		if i+1 < count {
			endDate = timex.Add(dateRange[i+1], timex.SECOND, -1)
		}
		list = append(list,
			fmt.Sprintf("%s..%s", timex.ToString(date, timex.YMD), timex.ToString(endDate, timex.YMD)))
	}
	// sort and remove duplicates
	slices.Sort(list)
	slices.Compact(list)
	// remove empties
	for _, dr := range list {
		if dr != "" {
			dates = append(dates, dr)
		}
	}
	return
}

// Source is the data source to fetch workflow run data from the api
// for repositories in the configuration
type Source[C Client, R Result] struct {
	client  C               // the *github.TeamsService compatible interface
	ctx     context.Context // ctx is the context to use
	cfg     *Config         // configuration values to use
	log     *slog.Logger    // logger
	filters []Filter        // set of filter functions to run against each result
}

// GetData returns all workflow runs between the start & end date for each repository
// that is within the configuration list.
//
// Filtering is done per repository, executing after all the workflow runs are fetched.
//
// `skipped` becomes a list of repository structs that have no results and likely need
// additional processing (checking for merges)
func (self *Source[C, R]) GetData() (results []R, skipped []any, err error) {
	results, skipped, err = self.allWorkflowRuns()
	return
}

// allWorkflowRuns
func (self *Source[C, R]) allWorkflowRuns() (results []R, skipped []any, err error) {
	var (
		total int          = len(self.cfg.Repositories)
		dates []string     = self.cfg.Dates()
		lg    *slog.Logger = self.log.With("date_start", self.cfg.DateStart, "date_end", self.cfg.DateEnd) // localised logger with config values added
	)
	lg.Debug("getting all workflow runs for all repositories ...")
	results = []R{}
	skipped = []any{}

	// loop over each repo
	for i, repo := range self.cfg.Repositories {
		lg.Debug(fmt.Sprintf("[%d/%d] (%s)", i+1, total, *repo.FullName))
		// fetch the runs within the date range we've worked out
		res, e := self.workflowRunsWithinDateRanges(repo, dates)
		if e != nil {
			err = e
			return
		}
		// run the filters against the found workflows
		filtered := self.filter(repo, res)
		// merge filtered set into main results
		for _, wfr := range filtered {
			results = append(results, wfr)
		}
		// if there are no workflows found, add the repo to the missing list
		if len(filtered) <= 0 {
			skipped = append(skipped, repo)
		}
	}

	lg.With("count", len(results)).Debug("getting workflow runs completed.")
	return
}

// filter handles running the filters from the config against this repo &
// workflow run list
func (self *Source[C, R]) filter(repo *github.Repository, workflowruns []*github.WorkflowRun) (filtered []*github.WorkflowRun) {
	var lg *slog.Logger = self.log // localised logger with config values added

	filtered = []*github.WorkflowRun{}
	// now check each workflow against the configured filters
	for _, workflowrun := range workflowruns {
		var include = true
		// check each filter, break on the first fail
		for _, f := range self.filters {
			// var inc = true
			lg.Debug(fmt.Sprintf("[%s][%s]:(%d)(%T) checking filter...", *repo.FullName, *workflowrun.Name, *workflowrun.ID, f))
			// if the filter is ever true, break the loop as we cant include & add to skipped list
			if include = f.Filter(self.ctx, workflowrun); !include {
				break
			}
		}
		// log if the run should be included or not
		lg.Debug(fmt.Sprintf("[%s][%s]:(%d) include workflow run? [%v]", *repo.FullName, *workflowrun.Name, *workflowrun.ID, include))
		if include {
			filtered = append(filtered, workflowrun)
		}
	}
	return
}

// workflowRunsWithinDateRanges iterates over each date range for this repository and fetches the workflow
// runs for that period (handling api pagination).
//
// WorkflowRun that are found are tracked via ID as the date ranges can hit an exact time / overlap
// and possibly return the same run more than once.
func (self *Source[C, R]) workflowRunsWithinDateRanges(repo *github.Repository, dateRanges []string) (results []*github.WorkflowRun, err error) {
	var (
		allRuns map[int64]*github.WorkflowRun                                         // using id based map as sometimes the date range campture overlaps
		lg      *slog.Logger                  = self.log.With("repo", *repo.FullName) // localised logger with config values added
	)

	lg.Debug("getting workflow runs for repo ...")
	allRuns = map[int64]*github.WorkflowRun{}
	results = []*github.WorkflowRun{}

	// for each repo, loop over date range parameters and grab the api content for that date range
	for _, dateRange := range dateRanges {
		runs, e := self.paginatedWorkflowRunsForDate(repo, dateRange)
		if e != nil {
			err = e
			return
		}
		// add local results to the main list
		for i, r := range runs {
			allRuns[i] = r
		}
	}
	// flattern the runs to a slice
	for _, run := range allRuns {
		results = append(results, run)
	}

	return

}

// paginatedWorkflowRunsForDate finds all successful workflow runs that ran against the repositories
// default branch.
//
// Handles paginated api responses, iterating over pages of 100 at a time.
//
// It runs the localised filters, which should be used to check workflow type / name (path to live).
func (self *Source[C, R]) paginatedWorkflowRunsForDate(repo *github.Repository, dateRange string) (runs map[int64]*github.WorkflowRun, err error) {

	var (
		page     int                             = 1                                                              // first page to fetch data from
		maxRetry int                             = 3                                                              // max retry counter
		lg       *slog.Logger                    = self.log.With("date_range", dateRange, "repo", *repo.FullName) // localised logger
		options  *github.ListWorkflowRunsOptions = &github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{PerPage: 100},
			Branch:      *repo.DefaultBranch,
			Status:      self.cfg.Status,
			Event:       self.cfg.Event,
		}
	)
	// if config overwrites the branch name, change the options here
	if self.cfg.OverwriteBranch != "" {
		options.Branch = self.cfg.OverwriteBranch
	}

	lg.Debug("getting workflow runs for repo in date range ...")
	runs = map[int64]*github.WorkflowRun{} // using id based map as sometimes the date range campture overlaps

	lg.Debug(fmt.Sprintf("[%s] (%s)", *repo.FullName, dateRange))
	// now call the paginated help to fetch all workflow runs for this range
	page = 1
	options.Created = dateRange
	for page > 0 {
		// pagination vars
		var (
			response *github.Response
			fetched  *github.WorkflowRuns
			e        error = errDefaultLoop // give error a default, non nil value so the for loop runs
			retry    int   = 0
		)
		// set the page
		options.Page = page
		// max of 3 attempts to call the same data set before failing.
		// 	- e has a default value so will always run at least once
		for e != nil && retry < maxRetry {
			// log
			lg.With("page", page, "try", retry).Debug("getting list of repository workflow runs in range ...")
			// make the api call
			fetched, response, e = self.client.ListRepositoryWorkflowRuns(self.ctx, *repo.Owner.Login, *repo.Name, options)
			retry += 1
			// if theres an error pause for a second - as error might be rate limiting
			if e != nil {
				time.Sleep(1)
			}
		}
		// if the error persits, then return
		if e != nil {
			lg.Error("failed to get workflow runs", "err", e.Error())
			err = errors.Join(e, ErrGettingList)
			return
		}
		// add workflow to list
		for _, workflowrun := range fetched.WorkflowRuns {
			runs[*workflowrun.ID] = workflowrun
		}
		// increment page
		page = response.NextPage
	}

	return
}

// New creates a source thats capable of fetching workflow runs for each repository.
//
// If config.Repositories is empty an error (ErrNoRepositoriesConfigured) will be returned.
//
// Notes:
//   - slog instance is pulled from the context.
//   - client is a *github.ActionService or mock version
//   - config contains parameters for the sdk / api call & repos
//   - config.Event will be set to "push" if blank
//   - config.Status will be set to "success" if blank
//   - filters is optional way of reducing the dataset afterwards
func New[C Client, R Result](ctx context.Context, client C, config *Config, filters ...Filter) (source *Source[C, R], err error) {
	var (
		defaultEvent  string = "push"
		defaultStatus string = "success"
	)
	// get logger
	ctx, lg := logx.Get(ctx)
	// if no repositories, return an error
	if len(config.Repositories) <= 0 {
		err = ErrNoRepositoriesConfigured
		return
	}
	// set default values on config event & status
	if config.Event == "" {
		config.Event = defaultEvent
	}
	if config.Status == "" {
		config.Status = defaultStatus
	}

	source = &Source[C, R]{
		ctx:     ctx,
		client:  client,
		cfg:     config,
		log:     lg,
		filters: filters,
	}

	return
}
