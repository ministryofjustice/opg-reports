// Package ghteamrepositories provides a struct and method for fetching all repositories owned
// by a team within an organisation.
//
// Pagination of the api call is handled within this package.
//
// Data filtering is done after all pages are fetched.
package ghteamrepositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/app/internal/logx"
	"time"

	"github.com/google/go-github/v87/github"
)

// ErrGettingList indicates the ListTeamReposBySlug function returned an error for
// the requested page
var ErrGettingList = errors.New("error getting page of repositories from api")

// ErrNoOrgSlug is triggered in New when config struct does not have an OrgansiationSlug present
var ErrNoOrgSlug = errors.New("no organistion slug has been configured")

// ErrNoTeamSlug is triggered in New when config struct does not have an TeamSlug present
var ErrNoTeamSlug = errors.New("no team slug has been configured")

// errDefaultLoop is a dummy error used to handle the fail & retry loop within the paginated
// api call
var errDefaultLoop = errors.New("dummy error for rety loop logic")

// Result is an alias for *github.Repository
type Result interface {
	*github.Repository
}

// Client is an interface for *github.TeamsService
type Client interface {
	// https://docs.github.com/en/rest/teams/teams#list-team-repositories
	ListTeamReposBySlug(ctx context.Context, orgSlug string, teamSlug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// Filter interface to allow filtering of the resulting data set
//
// Generally should be one of:
//   - ghfilters.ExcludeArchivedRepository
//   - ghfilters.FilterByRepositoryName
type Filter interface {
	Filter(ctx context.Context, result *github.Repository) (include bool)
}

// Config is used to capture information needed for the api/sdk call
type Config struct {
	OrganisationSlug string // the slug of the organisation
	TeamSlug         string // slug of the team whose repositories we're returning
}

// Source is the data source to fetch repositories from the
// github api that are owned by the configured team
type Source[C Client, R Result] struct {
	client  C               // the *github.TeamsService compatible interface
	ctx     context.Context // ctx is the context to use
	cfg     *Config         // configuration values to use
	log     *slog.Logger    // logger
	filters []Filter        // set of filter functions to run against each result
}

// GetData calls the github api, iterates over all of the paginated results, fetching
// 200 repositories per page and merging those together into a slice.
//
// It adds a very simple retry loop per page incase of simple errors or rate limiting
// triggering.
//
// After all original data is fetched, the results are then filtered against the
// configured filter functions
//
// `skipped` is unused but kept for interface
func (self *Source[C, R]) GetData() (results []R, skipped []any, err error) {
	results, skipped, err = self.repositories()
	return
}

func (self *Source[C, R]) filter(repositories []*github.Repository) (filtered []R, skipped []any) {
	var lg *slog.Logger = self.log // localised logger with config values added

	filtered = []R{}
	skipped = []any{}
	// now merge the data into the main result set based on filters
	for _, repo := range repositories {
		var include bool = true
		// check each filter, break on the first fail
		for _, f := range self.filters {
			lg.Debug(fmt.Sprintf("[%s] [%T] checking filter...", *repo.FullName, f))
			// if the filter is ever true, break the loop as we cant include & add to skipped list
			if include = f.Filter(self.ctx, repo); !include {
				skipped = append(skipped, *repo.FullName)
				break
			}
		}
		// log if the repo should be included or not
		lg.Debug(fmt.Sprintf("[%s] include repository? [%v]", *repo.FullName, include))
		if include {
			filtered = append(filtered, repo)
		}
	}

	return
}

// repositories handles the paginated data calls to the github api and merges the results.
//
// The API doesnt handle status or name filtering etc, so we run set of filter functions within
// the processing here to allow limiting of the data set by archive / name
func (self *Source[C, R]) repositories() (results []R, skipped []any, err error) {
	var (
		page     int                  = 1                         // first page to fetch data from
		maxRetry int                  = 3                         // due to api errors we have a catch / retry loop as well
		org      string               = self.cfg.OrganisationSlug // organisation slug
		team     string               = self.cfg.TeamSlug         // team slug
		allRepos []*github.Repository = []*github.Repository{}
		options  *github.ListOptions  = &github.ListOptions{PerPage: 200}       // set the default options
		lg       *slog.Logger         = self.log.With("team", team, "org", org) // localised logger with config values added
	)
	lg.Debug("getting repositories ...")

	for page > 0 {
		var (
			response *github.Response
			fetched  []*github.Repository
			e        error = errDefaultLoop // give error a default, non nil value for the loop
			retry    int   = 0
		)
		// set the page number for the api call
		options.Page = page
		// handle fetching the data with a rety loop.
		// max of 3 attempts to call the same data set before failing.
		for e != nil && retry < maxRetry {
			// log
			lg.With("page", page, "try", retry).Debug("getting list of repositories for team ...")
			// make the api call
			fetched, response, e = self.client.ListTeamReposBySlug(self.ctx, org, team, options)
			retry += 1
			// if theres an error pause for a second - as error might be rate limiting
			if e != nil {
				time.Sleep(1)
			}
		}
		// if the error persits, then return
		if e != nil {
			lg.Error("failed to get list of repositories", "err", e.Error())
			err = errors.Join(e, ErrGettingList)
			return
		}
		// attach to main list
		allRepos = append(allRepos, fetched...)
		// increment the page
		page = response.NextPage
	}

	results, skipped = self.filter(allRepos)

	lg.With("count", len(results)).Debug("getting repositories completed.")
	return
}

// New creates a source thats capable of fetching all repositories
//
// Will return an error if either the OrgansiationSlug or TeamSlug are empty.
//
// Notes:
// - slog instance is pulled from the context.
// - client is a *github.TeamsService or mock version
// - config contains parameters for the sdk / api call
// - filters is optional way of reducing the dataset
func New[C Client, R Result](ctx context.Context, client C, config *Config, filters ...Filter) (source *Source[C, R], err error) {
	// get logger
	ctx, lg := logx.Get(ctx)

	// if no org slug, throw an error
	if config.OrganisationSlug == "" {
		err = ErrNoOrgSlug
		return
	}
	// if no team slug, throw an error
	if config.TeamSlug == "" {
		err = ErrNoTeamSlug
		return
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
