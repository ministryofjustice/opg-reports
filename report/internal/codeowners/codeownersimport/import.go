package codeownersimport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/files"
	"opg-reports/report/package/repos"
	"slices"
	"strings"

	"github.com/google/go-github/v84/github"
)

// code owner entry
const InsertOwnersStatement string = `
INSERT INTO codebase_owners (
	owner,
	codebase,
	team_name
) VALUES (
	:owner,
	:codebase,
	:team_name
)
ON CONFLICT (owner,codebase,team_name) DO UPDATE SET
	owner=excluded.owner,
	codebase=excluded.codebase,
	team_name=excluded.team_name
RETURNING id
;
`
const truncateStmt string = `DELETE FROM codebase_owners;`

// teamClient wrapper around *github.TeamsService
type teamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// repoClient wrapper around *github.RepositoriesService
type repoClient interface {
	// fetch attached teams (*github.RepositoriesService)
	ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	// fetch file content (*github.RepositoriesService)
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
	// get contenst method to fetch directory content
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

type Clients struct {
	Teams teamClient // *github.TeamsService
	Repos repoClient // *github.RepositoriesService
}

type Args struct {
	DB           string `json:"db"`             // database path
	Driver       string `json:"driver"`         // database driver
	Params       string `json:"params"`         // database connection params
	OrgSlug      string `json:"org_slug"`       // github org name
	ParentSlug   string `json:"parent_slug"`    // parent slug
	FilterByName string `json:"filter_by_name"` // used to limit the repos to those that exactly match this name
}

type CodebaseOwner struct {
	Owner    string `json:"owner,omitempty"`
	Codebase string `json:"codebase,omitempty"` // full name of codebase
	TeamName string `json:"team_name"`
}

// mapping of codeowner / github teams to service teams (teams)
var codeOwnerToTeamName map[string]string = map[string]string{
	"ministryofjustice/digideps":                 "digideps",
	"ministryofjustice/opg-lpa-team":             "make",
	"ministryofjustice/opg-modernising-lpa-team": "modernise",
	"ministryofjustice/opg-sirius-poas":          "sirius",
	"ministryofjustice/opg-sirius-supervision":   "sirius",
	"ministryofjustice/opg-use-a-lpa-team":       "use",
	"ministryofjustice/serve-opg":                "serve",
	"ministryofjustice/sirius":                   "sirius",
}
var ErrFailedGettingRepositoryTeams = errors.New("failed to get team details for repository.")

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, client *Clients, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "Import")
	var list []*github.Repository

	log.Info("starting ...")
	// fetch all the repos
	log.Debug("getting repository list ...")
	list, err = repos.GetList(ctx, client.Teams, &repos.Args{
		OrgSlug:      in.OrgSlug,
		ParentSlug:   in.ParentSlug,
		FilterByName: in.FilterByName,
	})
	if err != nil {
		return
	}

	if err = handleCodebaseOwners(ctx, client.Repos, list, in); err != nil {
		return
	}

	log.Info("complete.")
	return
}

// handleCodebaseOwners
func handleCodebaseOwners(ctx context.Context, client repoClient, repositories []*github.Repository, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "handleCodebaseOwners")
	var data []*CodebaseOwner = []*CodebaseOwner{}
	log.Info("starting codebase owner import ...")
	// convert to local structs
	log.Debug("converting to codeowner models ...")
	data, err = generateCodebaseOwners(ctx, client, repositories, in)
	if err != nil {
		return
	}
	// trucate table first
	err = dbx.Exec(ctx, truncateStmt, &dbx.ExecArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error truncating table", "err", err.Error())
		return
	}
	// now write to db
	err = dbx.Insert(ctx, InsertOwnersStatement, data, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}
	log.With("count", len(data)).Info("complete.")
	return
}

// generateCodebaseOwners converts the repo data and then fetches extra data via api;
// in this case we pull teams and content of CODEOWNER files to determine all of our
// codeowner information
func generateCodebaseOwners(ctx context.Context, client repoClient, list []*github.Repository, in *Args) (data []*CodebaseOwner, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toCodebaseOwners")

	data = []*CodebaseOwner{}
	log.Debug("starting ...")

	for _, item := range list {
		log.Info("getting codeowners ...", "codebase", *item.FullName, "archived", *item.Archived)
		var teams []*github.Team = []*github.Team{}
		var owners []string = []string{}
		var merged []string = []string{}
		// only do this for active code bases, so if its archived, skip
		if *item.Archived {
			log.Info("archived, skipping", "codebase", *item.FullName)
			continue
		}
		// fetch teams for this code base
		teams, err = getTeams(ctx, client, item, in)
		if err != nil {
			return
		}
		// fetch content from code owner files
		owners, err = getCodeownersFromFiles(ctx, client, item, in)
		if err != nil {
			return
		}
		merged = filter(merge(teams, owners))
		// now make entry for each codeowner found
		for _, row := range merged {
			data = append(data, &CodebaseOwner{
				Codebase: *item.FullName,
				Owner:    row,
				TeamName: strings.ToLower(ownerToServiceTeam(row)),
			})
		}
	}
	log.Debug("complete.")
	return
}

// getTeams returns all attached teams for this code repository and deals with pagination
// of github results
func getTeams(ctx context.Context, client repoClient, code *github.Repository, in *Args) (teams []*github.Team, err error) {
	var (
		log  *slog.Logger        = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "getTeams")
		page int                 = 1
		opts *github.ListOptions = &github.ListOptions{PerPage: 200}
	)
	teams = []*github.Team{}

	log.With("codebase", code.FullName).Debug("starting ...")
	for page > 0 {
		var response *github.Response
		var list []*github.Team
		opts.Page = page

		log = log.With("page", page)
		log.Debug("getting team list ... ")
		// fetch team data
		list, response, err = client.ListTeams(ctx, in.OrgSlug, *code.Name, opts)
		if err != nil {
			log.Error("error getting team list")
			err = errors.Join(ErrFailedGettingRepositoryTeams, err)
			return
		}
		log.With("count", len(list)).Debug("found teams ...")
		// attach teams to the list
		for _, team := range list {
			if team.Parent != nil && *team.Parent.Slug == in.ParentSlug {
				teams = append(teams, team)
			}
		}
		// next loop
		page = response.NextPage
	}

	log.With("count", len(teams)).Debug("complete.")
	return
}

// getCodeownersFromFiles tries to fetch CODEOWNER file content from set locations and
// will process the content into just the team names, removing duplicates.
func getCodeownersFromFiles(ctx context.Context, client repoClient, code *github.Repository, in *Args) (owners []string, err error) {
	var (
		log           *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "getCodeownersFromFiles")
		fileLocations []string     = []string{
			"./CODEOWNERS",
			"./.github/CODEOWNERS",
		}
	)
	owners = []string{}

	log.With("codebase", code.FullName).Debug("starting ...")
	for _, filename := range fileLocations {
		var (
			e     error
			lines []string
			buff  io.ReadCloser
		)
		log.With("codeowner", filename).Debug("getting codeowner file ...")
		// fetch
		buff, _, e = client.DownloadContents(ctx, in.OrgSlug, *code.Name, filename, nil)
		lines = files.Lines(buff)
		// if there is an error, file might not be present, so ignore rather than return
		if e == nil && len(lines) > 0 {
			owners = append(owners, ownersFromLines(lines)...)
			break
		}

	}
	// remove duplicates
	slices.Sort(owners)
	owners = slices.Compact(owners)
	log.With("count", len(owners)).Debug("complete.")
	return
}

// ownersFromLines find all the code owners from lines in the codeowners file
// Note: strips lead @ from the team slug name and removes duplicates
func ownersFromLines(lines []string) (owners []string) {
	owners = []string{}
	for _, line := range lines {
		exploded := strings.Split(line, " ")
		if len(exploded) > 1 {
			for _, segment := range exploded[1:] {
				if len(segment) > 0 && segment != " " {
					owners = append(owners, segment)
				}
			}
		}
	}
	// remove duplicates
	slices.Sort(owners)
	owners = slices.Compact(owners)
	// remove the @prefix
	for i, o := range owners {
		if o[0] == '@' {
			owners[i] = o[1:]
		}
	}
	return
}

// filter removes org level codeowner / teams and replaces it with
// NONE string for easier querying
func filter(owners []string) (list []string) {
	var exclude = []string{
		"ministryofjustice/opg",
		"ministryofjustice/opg-webops",
	}
	list = []string{}
	for _, owner := range owners {
		if !slices.Contains(exclude, owner) {
			list = append(list, owner)
		}
	}
	// if there is owner found, then append none as a holder
	if len(list) == 0 {
		list = append(list, "none")
	}

	slices.Sort(list)
	list = slices.Compact(list)
	return
}

// merge teams and owners together in consistent way for slug values
func merge(teams []*github.Team, owners []string) (merged []string) {
	merged = []string{}

	for _, t := range teams {
		merged = append(merged, teamSlug(t))
	}

	for _, o := range owners {
		merged = append(merged, o)
	}
	// remove dups
	slices.Sort(merged)
	merged = slices.Compact(merged)
	return
}

// teamSlug helps create a standard team slug which will include the
// organsiation of the team to align with the content of CODEOWNERS
func teamSlug(team *github.Team) string {
	var teamSlug = *team.Slug
	// check the url structure for an org as sometimes the team.GetOrganization comes back nil...
	if team.HTMLURL != nil && strings.Contains(*team.HTMLURL, "/orgs/") {
		var url = *team.HTMLURL
		stripped := strings.ReplaceAll(url, "https://github.com/orgs/", "")
		org := strings.Split(stripped, "/")[0]
		teamSlug = fmt.Sprintf("%s/%s", org, *team.Slug)
	} else if o := team.GetOrganization(); o != nil && o.Login != nil {
		teamSlug = fmt.Sprintf("%s/%s", *o.Login, *team.Slug)
	}
	return teamSlug
}

// ownerToServiceTeam fetches service team where possible, or returns empty
func ownerToServiceTeam(owner string) (serviceTeam string) {
	serviceTeam = "none"
	if team, ok := codeOwnerToTeamName[owner]; ok {
		serviceTeam = team
	}
	return
}
