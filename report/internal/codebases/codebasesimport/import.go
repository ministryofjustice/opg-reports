package codebasesimport

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"opg-reports/report/package/cntxt"

	"github.com/google/go-github/v81/github"
)

var ErrFailedGettingRepositoryPage = errors.New("error getting page of repositories")

// TeamClient wrapper around *github.TeamsService
type TeamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// RepoClient wrapper around *github.RepositoriesService
type RepoClient interface {
	// fetch attached teams (*github.RepositoriesService)
	ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	// fetch file content (*github.RepositoriesService)
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
}

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params

	OrgSlug    string `json:"org_slug"`    // github org name
	ParentSlug string `json:"parent_slug"` // parent slug

	IncludeStats      bool `json:"include_stats"`      // run the code base stats handler - stats are non-time boxed details
	IncludeCodeowners bool `json:"include_codeowners"` // option to fetch all codebases and then fetch codeowner data as well
}

type Clients struct {
	Teams TeamClient // *github.TeamsService
	Repos RepoClient // *github.RepositoriesService
}

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, client *Clients, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "Import")
	var list []*github.Repository

	log.Info("starting ...")
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
	if in.IncludeStats {
		if err = handleCodebaseStats(ctx, list, in); err != nil {
			return
		}
	}
	// if enabled, run code owners
	if in.IncludeCodeowners {
		if err = handleCodebaseOwners(ctx, client.Repos, list, in); err != nil {
			return
		}
	}

	log.Info("complete.")
	return
}

// getRepositoryList iterates over paginated data set from github api and merges all data
// into one block
func getRepositoryList(ctx context.Context, client TeamClient, options *Args) (repositories []*github.Repository, err error) {

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
			log.With("repo", *repo.FullName).Debug("checking repository ...")
			repositories = append(repositories, repo)

		}
		page = response.NextPage
	}

	log.Debug("complete.")
	return
}
