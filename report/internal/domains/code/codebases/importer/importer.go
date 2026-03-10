package importer

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domains/code/types"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/logger"

	"github.com/google/go-github/v84/github"
)

const InsertStatement string = `
INSERT INTO codebases (
	name,
	full_name,
	url,
	archived
) VALUES (
	:name,
	:full_name,
	:url,
	:archived
) ON CONFLICT (full_name) DO UPDATE SET
	name=excluded.name,
	url=excluded.url,
	archived=excluded.archived
RETURNING id
;
`

// Get returns the repositories from the github api
func Get(ctx context.Context, client *github.Client, opts *args.Import, previous ...*github.Repository) (found []*github.Repository, err error) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	log.Info("getting list of repositories from github ...")
	found, err = paginated(ctx, client.Teams, opts)
	if err != nil {
		log.Error("error getting paginated list of data.", "err", err.Error())
		return
	}

	log.Info("repository list completed.", "count", len(found))

	return
}

// Filter
func Filter(ctx context.Context, items []*github.Repository, filters *args.Filters) (included []*github.Repository) {
	included = []*github.Repository{}

	for _, item := range items {
		if filters.Filter == "" || (*item.Name == filters.Filter) {
			included = append(included, item)
		}
	}
	return
}

// Transform converts the original data into record for local database insertion
func Transform(ctx context.Context, data []*github.Repository, opts *args.Import) (results []*types.Codebase, err error) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)
	results = []*types.Codebase{}

	log.Info("transforming repositories to local types.Codebases ...", "count", len(data))
	for _, repo := range data {
		var archived = 0
		if *repo.Archived {
			archived = 1
		}
		results = append(results, &types.Codebase{
			Name:     *repo.Name,
			FullName: *repo.FullName,
			Url:      *repo.HTMLURL,
			Archived: archived,
		})
	}

	log.Info("repository transformation completed.", "count", len(results))
	return
}

// internal type used to allow mocked testing.
//
// *github.TeamsService
type teamsService interface {
	// https://docs.github.com/en/rest/teams/teams#list-team-repositories
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

func paginated(ctx context.Context, client teamsService, opts *args.Import) (repositories []*github.Repository, err error) {
	var (
		log     *slog.Logger
		options *github.ListOptions
		info    *args.GitHub = opts.Github
		page    int          = 1
	)
	ctx, log = logger.Get(ctx)
	repositories = []*github.Repository{}
	log.Debug("getting list of repositories ... starting ...")

	for page > 0 {
		var response *github.Response
		var list []*github.Repository

		options = &github.ListOptions{
			PerPage: 200,
			Page:    page,
		}

		log.Debug("getting page of repositories ... ", "page", page)
		list, response, err = client.ListTeamReposBySlug(ctx, info.Organisation, info.Parent, options)
		if err != nil {
			log.Error("error getting paginated data", "page", page, "err", err.Error())
			return
		}
		// add data
		repositories = append(repositories, list...)
		// go to next page
		page = response.NextPage
	}

	log.Debug("complete.")
	return
}
