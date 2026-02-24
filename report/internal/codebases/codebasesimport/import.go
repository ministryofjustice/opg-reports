package codebasesimport

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"

	"github.com/google/go-github/v81/github"
)

const InsertStatement string = `
INSERT INTO codebases (
	name,
	full_name,
	url
) VALUES (
	:name,
	:full_name,
	:url
)
ON CONFLICT (full_name) DO UPDATE SET name=excluded.name, url=excluded.url
RETURNING id
;
`

var ErrFailedGettingRepositoryPage = errors.New("error getting page of repositories")

// Client wrapper around *github.TeamsService
type Client interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

type Args struct {
	DB            string `json:"db"`             // database path
	Driver        string `json:"driver"`         // database driver
	Params        string `json:"params"`         // database connection params
	MigrationFile string `json:"migration_file"` // database migrations

	OrgSlug    string `json:"org_slug"`
	ParentSlug string `json:"parent_slug"`
}

// Model represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type Model struct {
	Name     string `json:"name,omitempty"`       // short name of codebase (without owner)
	FullName string `json:"full_name,omitempty" ` // full name including the owner
	Url      string `json:"url,omitempty" `       // url to access the codebase
}

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, client Client, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "Import")
	var list []*github.Repository
	var data []*Model = []*Model{}

	log.Debug("starting ...")
	// fetch all the repos
	log.Debug("getting repository list ...")
	list, err = getRepositoryList(ctx, client, in)
	if err != nil {
		return
	}
	// convert to local structs
	log.Debug("converting to models ...")
	data, err = toModels(ctx, list)
	if err != nil {
		return
	}

	// now write to db
	err = dbx.Insert(ctx, InsertStatement, data, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.With("count", len(data)).Debug("complete.")
	return
}

// toModels converts the api results into local structs
func toModels(ctx context.Context, list []*github.Repository) (repos []*Model, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toModels")

	repos = []*Model{}
	log.Debug("starting ...")

	for _, item := range list {
		var repo = &Model{
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
func getRepositoryList(ctx context.Context, client Client, options *Args) (repositories []*github.Repository, err error) {

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

// included if its not archived
func includeInResult(repo *github.Repository, args *Args) (pass bool) {
	pass = !*repo.Archived
	return
}
