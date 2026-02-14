// Package codeownergetter is used to fetch & combine github team and CODEOWNER info as general code ownership.
package codeownergetter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/files"
	"slices"
	"strings"

	"github.com/google/go-github/v81/github"
)

var (
	ErrFailedGettingRepositoryTeams = errors.New("failed to get team details for repository.")
)

// GitHubClient
//
// Wrapper around *github.RepositoriesService
type GitHubClient interface {
	// fetch attached teams (*github.RepositoriesService)
	ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	// fetch file content (*github.RepositoriesService)
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
}

// Input struct contains the options and required data for the function
type Input struct {
	OrgSlug    string                     // github org slug
	ParentTeam string                     // github required parent team slug
	Codebases  []*codebasemodels.Codebase // list of codebases to fetch code ownership details about
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

// GetCodeowners uses a list of repositories (`in.Codebases`) to find all code owners attached to those and
// will also try to map those to a specific team
func GetCodeowners[T GitHubClient](ctx context.Context, log *slog.Logger, client T, in *Input) (result []*codeownermodels.Codeowner, err error) {
	var lg *slog.Logger = log.With("func", "codeowner.GetCodeowners")

	lg.With("codebases", len(in.Codebases)).Debug("starting ...")

	for _, code := range in.Codebases {
		var (
			teams  []*github.Team = []*github.Team{}
			owners []string       = []string{}
			merged []string       = []string{}
		)

		// fetch team info from the repo
		lg.With("codebase", code.FullName).Debug("getting teams ...")
		teams, err = getTeams(ctx, log, client, in, code)
		if err != nil {
			return
		}
		// fetch content from code owner files
		lg.With("codebase", code.FullName).Debug("getting codeowners ...")
		owners, err = getCodeownersFromFiles(ctx, log, client, in, code)
		if err != nil {
			return
		}
		// merge teams and owners together in consistent way for slug values
		merged = merge(teams, owners)
		// now create entries to return
		for _, row := range merged {
			result = append(result, &codeownermodels.Codeowner{
				Name:             row,
				CodebaseFullName: code.FullName,
				TeamName:         strings.ToLower(ownerToServiceTeam(row)),
			})
		}
	}
	lg.With("count", len(result)).Debug("complete.")
	return
}

// getCodeownersFromFiles tries to fetch CODEOWNER file content from set locations and
// will process the content into just the team names, removing duplicates.
func getCodeownersFromFiles[T GitHubClient](ctx context.Context, log *slog.Logger, client T, in *Input, code *codebasemodels.Codebase) (owners []string, err error) {
	var (
		lg            *slog.Logger = log.With("func", "codeowner.getCodeownersFromFiles")
		fileLocations []string     = []string{
			"./CODEOWNERS",
			"./.github/CODEOWNERS",
		}
	)
	owners = []string{}

	lg.With("codebase", code.FullName).Debug("starting ...")
	for _, filename := range fileLocations {
		var (
			e     error
			lines []string
			buff  io.ReadCloser
		)
		lg.With("codeowner", filename).Debug("getting codeowner file ...")
		// fetch
		buff, _, e = client.DownloadContents(ctx, in.OrgSlug, code.Name, filename, nil)
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
	lg.With("count", len(owners)).Debug("complete.")
	return
}

// getTeams returns all attached teams for this code repository and deals with pagination
// of github results
//
// Filters based on team having a parent of `opg`
func getTeams[T GitHubClient](ctx context.Context, log *slog.Logger, client T, in *Input, code *codebasemodels.Codebase) (teams []*github.Team, err error) {
	var (
		lg   *slog.Logger        = log.With("func", "codeowner.getTeams")
		page int                 = 1
		opts *github.ListOptions = &github.ListOptions{PerPage: 200}
	)
	teams = []*github.Team{}

	lg.With("codebase", code.FullName).Debug("starting ...")
	for page > 0 {
		var response *github.Response
		var list []*github.Team
		opts.Page = page

		lg = log.With("page", page)
		lg.Debug("getting team list ... ")
		// fetch team data
		list, response, err = client.ListTeams(ctx, in.OrgSlug, code.Name, opts)
		if err != nil {
			log.Error("error getting team list")
			err = errors.Join(ErrFailedGettingRepositoryTeams, err)
			return
		}
		lg.With("count", len(list)).Debug("found teams ...")
		// attach teams to the list
		for _, team := range list {
			if team.Parent != nil && *team.Parent.Slug == in.ParentTeam {
				teams = append(teams, team)
			}
		}
		// next loop
		page = response.NextPage
	}

	lg.With("count", len(teams)).Debug("complete.")
	return
}

// ownerToServiceTeam fetches service team where possible, or returns empty
func ownerToServiceTeam(owner string) (serviceTeam string) {
	serviceTeam = ""
	if team, ok := codeOwnerToTeamName[owner]; ok {
		serviceTeam = team
	}
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
