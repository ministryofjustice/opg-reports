package codebasesimport

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/repos"

	"github.com/google/go-github/v84/github"
)

// Raw codebase entry
const InsertCodebaseStatement string = `
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
)
ON CONFLICT (full_name) DO UPDATE SET
	name=excluded.name,
	url=excluded.url,
	archived=excluded.archived
RETURNING id
;
`

// TeamClient wrapper around *github.TeamsService
type teamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params

	OrgSlug      string `json:"org_slug"`       // github org name
	ParentSlug   string `json:"parent_slug"`    // parent slug
	FilterByName string `json:"filter_by_name"` // used to limit the repos to those that exactly match this name
}

// Codebase represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type Codebase struct {
	Name     string `json:"name,omitempty"`       // short name of codebase (without owner)
	FullName string `json:"full_name,omitempty" ` // full name including the owner
	Url      string `json:"url,omitempty" `       // url to access the codebase
	Archived int    `json:"archived"`             // int version of the archived flag on the repo
}

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, client teamClient, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "Import")
	var list []*github.Repository

	log.Info("starting ...")
	// fetch all the repos
	log.Debug("getting repository list ...")
	list, err = repos.GetList(ctx, client, &repos.Args{
		OrgSlug:      in.OrgSlug,
		ParentSlug:   in.ParentSlug,
		FilterByName: in.FilterByName,
	})
	if err != nil {
		return
	}
	// always run the codebase import
	err = handleCodebases(ctx, list, in)
	if err != nil {
		return
	}

	log.Info("complete.")
	return
}

func handleCodebases(ctx context.Context, repositories []*github.Repository, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "handleCodebases")
	var data []*Codebase = []*Codebase{}
	log.Info("starting codebase import ...")
	// convert to local structs
	log.Debug("converting to codebase models ...")
	data, err = generateCodebases(ctx, repositories)
	if err != nil {
		return
	}
	// now write to db
	err = dbx.Insert(ctx, InsertCodebaseStatement, data, &dbx.InsertArgs{
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

// toCodebases converts the api results into local structs
func generateCodebases(ctx context.Context, list []*github.Repository) (repos []*Codebase, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toCodebases")

	repos = []*Codebase{}
	log.Debug("starting ...")

	for _, item := range list {
		var archived = 0
		if *item.Archived {
			archived = 1
		}
		var repo = &Codebase{
			Name:     *item.Name,
			FullName: *item.FullName,
			Url:      *item.HTMLURL,
			Archived: archived,
		}
		repos = append(repos, repo)
		log.Debug("adding codebase", "full_name", repo.FullName, "archived", repo.Archived)
	}
	log.Debug("complete.")
	return
}
