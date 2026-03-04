package repos

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/package/cntxt"

	"github.com/google/go-github/v84/github"
)

// GetList iterates over paginated data set from github api and merges all data
// into one block
func GetList(ctx context.Context, client teamClient, options *Args) (repositories []*github.Repository, err error) {
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
