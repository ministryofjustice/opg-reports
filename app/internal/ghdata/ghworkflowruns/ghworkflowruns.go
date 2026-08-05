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
	"strings"
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

// Result interface
type Result interface {
	*ResultData
}

// ResultData is the result type which contains a combination of the
// repository and workflow run data for other uses (filteringer /
// sorting etc)
type ResultData struct {
	Repository  *github.Repository
	WorkflowRun *github.WorkflowRun
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
	filters []Filter        // set of filter functions to run against each result
	lg      *slog.Logger    // logger
}

// GetData returns all workflow runs between the start & end date for each repository
// that is within the configuration list.
//
// Filtering is done per repository, executing after all the workflow runs are fetched.
//
// `skipped` becomes a list of repository structs that have no results and likely need
// additional processing (checking for merges)
func (self *Source[C, R]) GetData() (results []R, skipped []any, err error) {
	self.lg.Info("getting data.")

	results, skipped, err = self.allWorkflowRuns()

	self.lg.Info("got data.",
		"count", len(results),
		"skipped", len(skipped),
		"err", err)

	return
}

// allWorkflowRuns
func (self *Source[C, R]) allWorkflowRuns() (results []R, skipped []any, err error) {
	var (
		total int      = len(self.cfg.Repositories)
		dates []string = self.cfg.Dates()
	)
	self.lg.Debug("getting all workflow runs.")
	results = []R{}
	skipped = []any{}

	// loop over each repo
	for i, repo := range self.cfg.Repositories {
		self.lg.Debug("getting workflow runs for repository.",
			"i", i, "total", total,
			"repository", *repo.FullName)

		// fetch the runs within the date range we've worked out
		res, e := self.workflowRunsWithinDateRanges(repo, dates)
		if e != nil {
			err = e
			self.lg.Error("error getting workflow runs.", "err", err.Error())
			return
		}
		// run the filters against the found workflows
		self.lg.Debug("filtering workflow runs.")
		filtered := self.filter(repo, res)
		// merge filtered set into main results
		for _, wfr := range filtered {
			results = append(results, &ResultData{Repository: repo, WorkflowRun: wfr})
		}
		// if there are no workflows found, add the repo to the missing list
		if len(filtered) <= 0 {
			skipped = append(skipped, repo)
		}
	}
	self.lg.Debug("got workflow runs.",
		"count", len(results),
		"skipped", len(skipped),
		"err", err)
	return
}

// filter handles running the filters from the config against this repo &
// workflow run list
func (self *Source[C, R]) filter(repo *github.Repository, workflowruns []*github.WorkflowRun) (filtered []*github.WorkflowRun) {

	filtered = []*github.WorkflowRun{}
	// now check each workflow against the configured filters
	for _, workflowrun := range workflowruns {
		var include = true
		// check each filter, break on the first fail
		for _, f := range self.filters {

			self.lg.Debug("checking filter.",
				"repository", *repo.FullName,
				"workflowRunName", *workflowrun.Name,
				"T", fmt.Sprintf("%T", f))

			// if the filter is ever not true, break the loop
			if include = f.Filter(self.ctx, workflowrun); !include {
				break
			}
		}

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
		allRuns map[int64]*github.WorkflowRun // using id based map as sometimes the date range campture overlaps
	)
	self.lg.Debug("getting workflow runs witin date ranges.",
		"repository", *repo.FullName,
		"dateRanges", strings.Join(dateRanges, ","))

	allRuns = map[int64]*github.WorkflowRun{}
	results = []*github.WorkflowRun{}

	// for each repo, loop over date range parameters and grab the api content for that date range
	for _, dateRange := range dateRanges {
		runs, e := self.paginatedWorkflowRunsForDate(repo, dateRange)
		if e != nil {
			self.lg.Error("error getting workflow runs.", "e", e.Error())
			err = e
			return
		}
		// add local results to the main list
		for i, r := range runs {
			allRuns[i] = r
		}
	}
	// map to slice
	for _, run := range allRuns {
		results = append(results, run)
	}

	self.lg.Debug("got workflow runs witin date ranges.",
		"repository", *repo.FullName,
		"count", len(results))

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
		page     int                             = 1 // first page to fetch data from
		maxRetry int                             = 3 // max retry counter
		options  *github.ListWorkflowRunsOptions = &github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{PerPage: 100},
			Branch:      *repo.DefaultBranch,
			Status:      self.cfg.Status,
			Event:       self.cfg.Event,
		}
	)
	self.lg.Debug("getting workflow runs witin date range.",
		"repository", *repo.FullName,
		"dateRange", dateRange)
	// using id based map as sometimes the date range campture overlaps
	runs = map[int64]*github.WorkflowRun{}
	// if config overwrites the branch name, change the options here
	if self.cfg.OverwriteBranch != "" {
		options.Branch = self.cfg.OverwriteBranch
	}
	// now call the paginated api end points to fetch all within the date range
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
			self.lg.Debug("getting page of workflow runs,",
				"page", page, "try", retry, "repository", *repo.FullName)
			// make the api call
			fetched, response, e = self.client.ListRepositoryWorkflowRuns(self.ctx, *repo.Owner.Login, *repo.Name, options)
			retry += 1
			// if theres an error pause for a second - as error might be rate limiting
			if e != nil {
				self.lg.Warn("error fetching page - sleeping & retrying.", "e", e.Error())
				time.Sleep(1)
			}
		}
		// if the error persits, then return
		if e != nil {
			self.lg.Error("error after retries.", "e", e.Error())
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

	self.lg.Debug("got workflow runs witin date range.",
		"count", len(runs),
		"err", err,
		"repository", *repo.FullName)

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
		defaultEvent  string       = "push"
		defaultStatus string       = "success"
		lg            *slog.Logger = logx.Default()
	)
	// get logger
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
		filters: filters,
		lg:      lg,
	}

	return
}
