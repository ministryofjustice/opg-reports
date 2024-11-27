package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/api/inputs"
)

var (
	GitHubTeamsSegment string   = "github/teams"
	GitHubTeamsTags    []string = []string{"github", "teams"}
)

// --- github teams list

type GitHubTeamsListBody struct {
	Operation string               `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.VersionInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubTeam `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error              `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}
type GitHubTeamsResponse struct {
	Body *GitHubTeamsListBody
}

const GitHubTeamsDescription string = `Returns all github teams within the database.`
const GitHubTeamsOperationID string = "get-github-teams-list"
const gitHubTeamsSQL string = `
SELECT
	github_teams.*,
	json_group_array(
		json_object(
			'id', units.id,
			'name', units.name
		)
	) filter ( where units.id is not null) as units,
	 json_group_array(
		json_object(
			'id', github_repositories.id,
			'ts', github_repositories.ts,
			'owner', github_repositories.owner,
			'name', github_repositories.name,
			'full_name', github_repositories.full_name,
			'created_at', github_repositories.created_at,
			'default_branch', github_repositories.default_branch,
			'archived', github_repositories.archived,
			'private', github_repositories.private,
			'license', github_repositories.license,
			'last_commit_date', github_repositories.last_commit_date
		)
	) filter ( where github_repositories.id is not null) as github_repositories
FROM github_teams
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_teams.id
LEFT JOIN units ON units.id = github_teams_units.unit_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_team_id = github_teams.id
LEFT JOIN github_repositories ON github_repositories.id = github_repositories_github_teams.github_repository_id
GROUP BY github_teams.id
ORDER BY github_teams.slug ASC;
`

// ApiGitHubTeamsListHandler queries the database for all github teams and returns a list including
// joins with github teams and aws accounts. There is no option to filter of limit the results.
//
// Endpoints:
//
//	/version/github/teams/list
func ApiGitHubTeamsListHandler(ctx context.Context, input *inputs.VersionInput) (response *GitHubTeamsResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubTeam = []*models.GitHubTeam{}
		dbPath  string               = ctx.Value(dbPathKey).(string)
		body    *GitHubTeamsListBody = &GitHubTeamsListBody{
			Request:   input,
			Operation: GitHubTeamsOperationID,
		}
	)
	// setup response
	response = &GitHubTeamsResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] github teams list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubTeam](ctx, adaptor, gitHubTeamsSQL, nil)
	if err != nil {
		slog.Error("[api] github teams list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github teams list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

// Register attaches the handler to the main api
func RegisterGitHubTeams(api huma.API) {
	var uri string = "/{version}/" + GitHubTeamsSegment + "/list"

	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubTeamsOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List all units",
		Description:   GitHubTeamsDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubTeamsTags,
	}, ApiGitHubTeamsListHandler)

}
