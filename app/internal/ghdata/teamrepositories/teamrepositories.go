// teamrepositories provides a struct and method for fetching all repositories owned
// by a team within an organisation.
//
// Pagination of the api call is handled within the call, but as the github api does
// not expose any filtering, none is applied
package teamrepositories

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/app/internal/ghdata/ghconfig"
	"opg-reports/app/internal/logx"
	"time"

	"github.com/google/go-github/v87/github"
)

// ErrGettingRepositoryList indicates the ListTeamReposBySlug function returned an error for
// the requested page
var ErrGettingRepositoryList = errors.New("error getting page of repositories from api")

// errDefaultLoop is a dummy error used to handle the fail & retry loop within the paginated
// api call
var errDefaultLoop = errors.New("dummy error for rety loop logic")

// Result is an alias for *github.Repository
type Result interface {
	*github.Repository
}

// Client is an interface for *github.TeamsService
type Client interface {
	ListTeamReposBySlug(ctx context.Context, orgSlug string, teamSlug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// Source is the data source to fetch repositories from the
// github api that are owned by the configured team
type Source[C Client, R Result] struct {
	client C                // the *github.TeamsService compatible interface
	ctx    context.Context  // ctx is the context to use
	cfg    *ghconfig.Config // configuration values to use
	log    *slog.Logger     // logger
}

// GetData calls the github api, iterates over all of the paginated results, fetching
// 200 repositories per page and merging those together into a slice.
//
// It adds a very simple retry loop per page incase of simple errors or rate limiting
// triggering.
//
// As the api endpoint doesnt not provide filtering of any kind, neither does this
// function.
func (self *Source[C, R]) GetData() (results []R, err error) {
	results, err = self.getPaginatedData()
	return
}

// getPaginatedData handles the paginated data calls to the github api and merges the results.
//
// There is no filtering at this point as the API does not provide a means to
func (self *Source[C, R]) getPaginatedData() (results []R, err error) {
	var (
		page     int                 = 1                                       // first page to fetch data from
		maxRetry int                 = 3                                       // due to api errors we have a catch / retry loop as well
		org      string              = self.cfg.OrganisationSlug               // organisation slug
		team     string              = self.cfg.TeamSlug                       // team slug
		options  *github.ListOptions = &github.ListOptions{PerPage: 200}       // set the default options
		lg       *slog.Logger        = self.log.With("team", team, "org", org) // localised logger with config values added
	)
	lg.Debug("getting repositories ...")
	results = []R{}

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
			//
			fetched, response, e = self.client.ListTeamReposBySlug(self.ctx, org, team, options)
			retry += 1
			time.Sleep(1) // pause for a second - as error might be rate limiting
		}
		// if the error persits, then return
		if e != nil {
			lg.Error("failed to get list of repositories", "err", e.Error())
			err = errors.Join(e, ErrGettingRepositoryList)
			return
		}

		// now merge the data into the main result set
		for _, repo := range fetched {
			results = append(results, repo)
		}
		// increment the page
		page = response.NextPage
	}
	lg.With("count", len(results)).Debug("getting repositories completed.")
	return
}

// New
func New[C Client, R Result](ctx context.Context, client C, config *ghconfig.Config) (source *Source[C, R], err error) {
	// get logger
	ctx, lg := logx.Get(ctx)
	source = &Source[C, R]{
		ctx:    ctx,
		client: client,
		cfg:    config,
		log:    lg,
	}

	return
}
