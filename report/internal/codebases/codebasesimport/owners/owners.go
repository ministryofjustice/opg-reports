package owners

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasesimport/args"
	"opg-reports/report/internal/codebases/codebasesimport/clients"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/files"
	"slices"
	"strings"

	"github.com/google/go-github/v81/github"
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

type CodebaseOwner struct {
	Owner    string `json:"owner,omitempty"`
	Codebase string `json:"codebase,omitempty"` // full name of codebase
	TeamName string `json:"team_name"`
}

var ErrFailedGettingRepositoryTeams = errors.New("failed to get team details for repository.")

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

func HandleCodebaseOwners(ctx context.Context, client clients.RepoClient, repositories []*github.Repository, in *args.Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "handleCodebaseOwners")
	var data []*CodebaseOwner = []*CodebaseOwner{}
	log.Info("starting codebase owner import ...")
	// convert to local structs
	log.Debug("converting to codeowner models ...")
	data, err = toCodebaseOwners(ctx, client, repositories, in)
	if err != nil {
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

// toCodebaseOwners converts the api results into local structs
func toCodebaseOwners(ctx context.Context, client clients.RepoClient, list []*github.Repository, in *args.Args) (data []*CodebaseOwner, err error) {
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
func getTeams(ctx context.Context, client clients.RepoClient, code *github.Repository, in *args.Args) (teams []*github.Team, err error) {
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
func getCodeownersFromFiles(ctx context.Context, client clients.RepoClient, code *github.Repository, in *args.Args) (owners []string, err error) {
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
