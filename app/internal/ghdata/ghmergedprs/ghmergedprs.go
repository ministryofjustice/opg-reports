// Package ghMergedPullRequests provides a struct and method for fetching all pull requests merged to
// the default branch between two dates.
//
// Intention is this data is used to determine the number releases during a given time period
// when workflow runs are not present (so external ci/cd usage).
//
// Pagination of the api call is handled within this package.
//
// The returned type contains both repository and pull request for later
// filtering / sorting etc.
package ghMergedPullRequests

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/app/internal/logx"
	"time"

	"github.com/google/go-github/v87/github"
)

// ErrNoRepositoriesConfigured is triggered when the Config passed to New does not have any
// repositories attached
var ErrNoRepositoriesConfigured = errors.New("no repositories have been set on the configuration struct.")

// ErrGettingList indicates the List call returned an error
var ErrGettingList = errors.New("error getting page of pull requests from api")

// errDefaultLoop is a dummy error used to handle the fail & retry loop within the paginated
// api call
var errDefaultLoop = errors.New("dummy error for retry loop logic")

// Result is an alias for MergedPullRequest
type Result interface {
	*MergedPullRequest
}

// Client is an interface for *github.PullRequestsService
type Client interface {
	// api docs - https://docs.github.com/rest/pulls/pulls#list-pull-requests
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

// Filter iterface is used for running functions against a set of results
//
// Generally should be...
type Filter interface {
	Filter(ctx context.Context, result *github.WorkflowRun) (include bool)
}

// Config is used to capture information needed for the api/sdk call
type Config struct {
	Repositories []*github.Repository // list of all repositories we want to fetch workflow runs for
	DateStart    time.Time            // start of date range to fetch workflow runs for each repo.
	DateEnd      time.Time            // end of date range to fetch workflow runs for each repo.
	State        string               // defaults to "closed"
}

// MergedPullRequest is the result type which contains a combination of the
// repository and pull request data for other uses
type MergedPullRequest struct {
	Repository  *github.Repository
	PullRequest *github.PullRequest
}

// Source is the data source to fetch repositories from the
// github api that are owned by the configured team
type Source[C Client, R Result] struct {
	client  C               // the *github.PullRequestsService compatible interface
	ctx     context.Context // ctx is the context to use
	cfg     *Config         // configuration values to use
	log     *slog.Logger    // logger
	filters []Filter        // set of filter functions to run against results
}

// GetData returns all pull requests between the start & end date for each repository
// that is within the configuration list.
//
// Filtering is done per repository, executing after all the workflow runs are fetched.
//
// `skipped` becomes a list of repository structs that have no results for the time period.
func (self *Source[C, R]) GetData() (results []R, skipped []any, err error) {
	results, skipped, err = self.allMergedPullRequests()
	return
}

// allMergedPullRequests fetches all pull requests for each repository within the
// configured date period
func (self *Source[C, R]) allMergedPullRequests() (results []R, skipped []any, err error) {
	var (
		all   map[int64]*MergedPullRequest = map[int64]*MergedPullRequest{}
		total int                          = len(self.cfg.Repositories)
		lg    *slog.Logger                 = self.log.With("date_start", self.cfg.DateStart, "date_end", self.cfg.DateEnd) // localised logger with config values added
	)
	lg.Debug("getting all pull requests for all repositories ...")
	results = []R{}
	skipped = []any{}

	// loop over each repo
	for i, repo := range self.cfg.Repositories {
		lg.Debug(fmt.Sprintf("[%d/%d] (%s)", i+1, total, *repo.FullName))

		prs, e := self.paginatedPullRequests(repo)
		if e != nil {
			err = e
			return
		}

		// add results into the main map, ignoring any duplicates
		for i, pr := range prs {
			all[i] = &MergedPullRequest{
				Repository:  repo,
				PullRequest: pr,
			}
		}

	}
	// flattern for the result
	for _, pr := range all {
		results = append(results, pr)
	}

	return
}

// paginatedPullRequests handles the pagination and date resitrctions of the api calls.
//
// As the api call does not have a date filter we sort all pull requests based on the created
// date descending and then stop fetching data when we move past the
func (self *Source[C, R]) paginatedPullRequests(repo *github.Repository) (prs map[int64]*github.PullRequest, err error) {

	var (
		page     int                            = 1                                                                                                     // first page to fetch data from
		maxRetry int                            = 3                                                                                                     // max retry counter
		lg       *slog.Logger                   = self.log.With("repo", *repo.FullName, "date_start", self.cfg.DateStart, "date_end", self.cfg.DateEnd) // localised logger
		options  *github.PullRequestListOptions = &github.PullRequestListOptions{
			ListOptions: github.ListOptions{PerPage: 100},
			Base:        *repo.DefaultBranch,
			State:       self.cfg.State,
			Sort:        "created", // fixed to created
			Direction:   "desc",    // fixed ordering for logic to work
		}
	)
	lg.Debug("getting pull requests within date range ...")

	page = 1
	prs = map[int64]*github.PullRequest{}
	// the pagination loop will fetch from the api and has a retry loop for errors
	// to try and recover from timesouts etc
	for page > 0 {
		//  setup pagination vars
		var (
			response *github.Response
			fetched  []*github.PullRequest
			e        error = errDefaultLoop
			retry    int   = 0 // retry counter
		)
		// set the page number
		options.Page = page
		// retry loop
		for e != nil && retry < maxRetry {
			lg.With("page", page, "try", retry).Debug("getting list of pull requests ...")
			// api call
			fetched, response, e = self.client.List(self.ctx, *repo.Owner.Login, *repo.Name, options)
			retry += 1
			// if theres an error pause for a second - as error might be rate limiting
			if e != nil {
				time.Sleep(1)
			}
		}
		// if there is an error, retunr
		if e != nil {
			lg.Error("failed to get workflow runs", "err", e.Error())
			err = errors.Join(e, ErrGettingList)
			return
		}

		// incremeant page here, so if we go outside of the range we can overwrite the
		// page and not fetch any more data
		page = response.NextPage
		// otherwise, merge in results
		for _, pr := range fetched {
			var valid, before, _ bool = self.includePR(lg, pr)
			//
			if valid {
				prs[*pr.ID] = pr
			}
			// if we're outside of the date range, set next page as 0 to
			// break the outer loop and also break this loop so we
			// dont include any more pull requests
			if before {
				lg.Debug("pr created before date range, breaking loops ...")
				page = 0
				break
			}
		}

	}

	return
}

// includePR is used to decide if the pr should be included in the result set and
// tracks if it was within the date range (via before & after)
//
// pr is only valid when ...
//   - createdAt is present
//   - createdAt is between the start & end date configured
//   - merge commit SHA & time are present
//
// valid, before, after all default to false
func (self *Source[C, R]) includePR(lg *slog.Logger, pr *github.PullRequest) (valid bool, before bool, after bool) {
	var (
		created time.Time
		start   time.Time = self.cfg.DateStart
		end     time.Time = self.cfg.DateEnd
	)

	valid = false
	before = false
	after = false

	// check the date of pr
	if pr.CreatedAt == nil {
		lg.Warn("pr createdAt is nil ...")
		return
	}

	created = pr.CreatedAt.Time
	before = created.Before(start)
	after = created.After(end)

	// check merge commit details
	if pr.MergeCommitSHA == nil || len(*pr.MergeCommitSHA) <= 0 || pr.MergedAt == nil {
		lg.Debug("pr merge commit data is missing ...")
		return
	}

	// if pr is within the date range, mark as valid
	if !before && !after {
		valid = true
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
//   - filters is optional way of reducing the dataset afterwards
func New[C Client, R Result](ctx context.Context, client C, config *Config, filters ...Filter) (source *Source[C, R], err error) {
	var (
		defaultState string = "closed"
	)
	// get logger
	ctx, lg := logx.Get(ctx)
	// if no repositories, return an error
	if len(config.Repositories) <= 0 {
		err = ErrNoRepositoriesConfigured
		return
	}

	// set default values on config event & status
	if config.State == "" {
		config.State = defaultState
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
