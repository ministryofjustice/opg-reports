package codebasesimport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/rest"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v81/github"
)

// compliance_level TEXT,
//
//	compliance_report_url TEXT
//	compliance_badge TEXT,
const InsertStatement string = `
INSERT INTO codebases (
	name,
	full_name,
	url,
	compliance_level,
	compliance_report_url,
	compliance_badge
) VALUES (
	:name,
	:full_name,
	:url,
	:compliance_level,
	:compliance_report_url,
	:compliance_badge
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
	Name                string `json:"name,omitempty"`                  // short name of codebase (without owner)
	FullName            string `json:"full_name,omitempty" `            // full name including the owner
	Url                 string `json:"url,omitempty" `                  // url to access the codebase
	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url

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
	var base string = "https://github-community.service.justice.gov.uk/repository-standards"
	// var timeout = (2 * time.Second)

	repos = []*Model{}
	log.Debug("starting ...")

	for _, item := range list {
		var repo = &Model{
			Name:                *item.Name,
			FullName:            *item.FullName,
			Url:                 *item.HTMLURL,
			ComplianceLevel:     "unknown",
			ComplianceReportUrl: fmt.Sprintf("%s/%s", base, *item.Name),
			ComplianceBadge:     fmt.Sprintf("%s/api/%s/badge", base, *item.Name),
		}
		// set the compliance level
		if lvl, e := complianceLevelFromBadge(ctx, repo.ComplianceBadge); e == nil {
			repo.ComplianceLevel = lvl
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

// the badge layout puts the value in the title
var complianceRe = regexp.MustCompile(`(?m)<title>MOJ COMPLIANT:(.*)</title>`)

func complianceLevelFromBadge(ctx context.Context, badge string) (level string, err error) {
	var timeout = (2 * time.Second)
	level = "unknown"
	res, _, err := rest.GetStr(ctx, nil, &rest.Request{Host: badge, Timeout: timeout})
	if err != nil {
		return
	}
	// find a match
	for _, match := range complianceRe.FindAllString(res, 1) {
		level = match
	}
	// trim the extras
	level = strings.ReplaceAll(level, "<title>MOJ COMPLIANT:", "")
	level = strings.ReplaceAll(level, "</title>", "")
	// swap out not foudn for not_found to make parsing easier
	level = strings.ReplaceAll(level, "NOT FOUND", "unknown")
	level = strings.Trim(level, " ")
	// split on space
	levels := strings.Split(level, " ")
	// use. the last part only
	level = levels[len(levels)-1]
	level = strings.ToLower(level)
	return
}

// included if its not archived
func includeInResult(repo *github.Repository, args *Args) (pass bool) {
	pass = !*repo.Archived
	return
}
