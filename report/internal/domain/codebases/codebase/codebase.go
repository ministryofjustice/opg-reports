package codebase

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/google/go-github/v81/github"
)

var ErrFailedGettingRepositoryPage = errors.New("error getting page of repositories")

// GithubClient wrapper around *github.TeamsService
type GithubClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// GetCodebasesOptions used to decide / change what repositories to return
// from the full list
type GetCodebasesOptions struct {
	ExcludeArchived bool
}

// fixed values
const (
	githubOrg  string = "ministryofjustice" // github org
	githubTeam string = "opg"               // github root team
)

// GetCodebases finds all github repositories and returns them for the moj/opg team
func GetCodebases[T GithubClient](ctx context.Context, log *slog.Logger, client T, options *GetCodebasesOptions) (repos []*codebasemodels.Codebase, err error) {
	var list []*github.Repository

	log = log.With("package", "codebases", "func", "GetCodebases")
	log.Debug("starting ...")

	// fetch all the repos
	log.Debug("getting repository list ...")
	list, err = getRepositoryList(ctx, log, client, options)
	if err != nil {
		return
	}
	// convert to local structs
	log.Debug("converting to models ...")
	repos, err = toModels(ctx, log, list)
	if err != nil {
		return
	}

	log.With("count", len(repos)).Debug("complete.")
	return
}

// toModels converts the api results into local structs
func toModels(ctx context.Context, log *slog.Logger, list []*github.Repository) (repos []*codebasemodels.Codebase, err error) {

	repos = []*codebasemodels.Codebase{}
	log = log.With("package", "codebases", "func", "toModels")
	log.Debug("starting ...")

	for _, item := range list {
		var repo = &codebasemodels.Codebase{
			Name:     *item.Name,
			FullName: *item.FullName,
			Url:      *item.HTMLURL,
		}
		repos = append(repos, repo)
	}
	log.Debug("complete.")
	return
}

// getRepositoryList iterates over paginated data set from github api and merges all data
// into one block
func getRepositoryList[T GithubClient](ctx context.Context, log *slog.Logger, client T, options *GetCodebasesOptions) (repositories []*github.Repository, err error) {

	var (
		page int                 = 1
		opts *github.ListOptions = &github.ListOptions{PerPage: 200}
	)
	log = log.With("package", "codebases", "func", "getRepositoryList")
	log.Debug("starting ...")

	for page > 0 {
		var response *github.Response
		var list []*github.Repository
		// set the page to request
		opts.Page = page
		log.With("page", page).Debug("getting page of repositories ...")
		// fetch data from api
		list, response, err = client.ListTeamReposBySlug(ctx, githubOrg, githubTeam, opts)
		if err != nil {
			err = errors.Join(ErrFailedGettingRepositoryPage, err)
			return
		}
		// only add non archived repos
		log.With("page", page, "count", len(list)).Debug("found repositories ... ")

		for _, repo := range list {
			var include = includeInResult(repo, options)
			log.With("include", include, "repo", *repo.FullName).Debug("checking repository ...")
			if include {
				repositories = append(repositories, repo)
			}
		}
		page = response.NextPage
	}

	log.Debug("complete.")
	return
}

func includeInResult(repo *github.Repository, criteria *GetCodebasesOptions) (pass bool) {
	pass = true
	if criteria == nil {
		return
	}
	// if we're exlcudeding archived & this is archived, dont include
	if criteria.ExcludeArchived && repo.Archived != nil && *repo.Archived == true {
		pass = false
	}
	return
}
