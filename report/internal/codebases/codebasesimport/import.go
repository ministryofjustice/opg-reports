package codebasesimport

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasesimport/args"
	"opg-reports/report/internal/codebases/codebasesimport/clients"
	"opg-reports/report/internal/codebases/codebasesimport/metrics"
	"opg-reports/report/internal/codebases/codebasesimport/owners"
	"opg-reports/report/internal/codebases/codebasesimport/stats"
	"opg-reports/report/package/cntxt"

	"github.com/google/go-github/v84/github"
)

var ErrFailedGettingRepositoryPage = errors.New("error getting page of repositories")

type Clients struct {
	Teams   clients.TeamClient   // *github.TeamsService
	Repos   clients.RepoClient   // *github.RepositoriesService
	Actions clients.ActionClient // *github.ActionsService
}

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, client *Clients, in *args.Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "Import")
	var list []*github.Repository

	log.Info("starting ...")
	if client.Teams == nil {
		log.Error("teams client is empty")
		return
	}
	// fetch all the repos
	log.Debug("getting repository list ...")
	list, err = getRepositoryList(ctx, client.Teams, in)
	if err != nil {
		return
	}
	// always run the codebase import
	err = handleCodebases(ctx, list, in)
	if err != nil {
		return
	}
	// if enabled, run stats
	if in.IncludeStats && client.Repos != nil {
		if err = stats.HandleCodebaseStats(ctx, client.Repos, list, in); err != nil {
			return
		}
	}
	if in.IncludeMetrics && client.Actions != nil {
		if err = metrics.HandleCodebaseMetrics(ctx, client.Actions, list, in); err != nil {
			return
		}
	}
	// if enabled, run code owners
	if in.IncludeCodeowners && client.Repos != nil {
		if err = owners.HandleCodebaseOwners(ctx, client.Repos, list, in); err != nil {
			return
		}
	}

	log.Info("complete.")
	return
}

// getRepositoryList iterates over paginated data set from github api and merges all data
// into one block
func getRepositoryList(ctx context.Context, client clients.TeamClient, options *args.Args) (repositories []*github.Repository, err error) {

	var (
		page int                 = 1
		opts *github.ListOptions = &github.ListOptions{PerPage: 200}
		log  *slog.Logger        = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "getRepositoryList")
	)
	log.Debug("starting ...")

	for page > 0 {
		var response *github.Response
		var list []*github.Repository
		// set the page to request
		opts.Page = page
		log.With("page", page).Debug("getting page of repositories ...")
		// fetch data from api
		list, response, err = client.ListTeamReposBySlug(ctx, options.OrgSlug, options.ParentSlug, opts)
		if err != nil {
			err = errors.Join(ErrFailedGettingRepositoryPage, err)
			return
		}
		// only add non archived repos
		log.With("page", page, "count", len(list)).Debug("found repositories ...")

		for _, repo := range list {
			log.With("repo", *repo.FullName).Debug("found repository ...")
			if options.FilterByName == "" || (options.FilterByName == *repo.Name) {
				repositories = append(repositories, repo)
				log.With("repo", *repo.FullName).Debug("added repository ...")
			}
		}
		page = response.NextPage
	}

	log.With("count", len(repositories)).Debug("complete.")
	return
}
