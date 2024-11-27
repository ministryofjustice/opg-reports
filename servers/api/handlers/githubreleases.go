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
	"github.com/ministryofjustice/opg-reports/internal/dbs/statements"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/api/inputs"
)

var (
	GitHubReleasesSegment string   = "github/releases"
	GitHubReleasesTags    []string = []string{"github", "releases"}
)

type GitHubReleasesListBody struct {
	Operation string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.OptionalGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRelease               `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}
type GitHubReleasesListResponse struct {
	Body *GitHubReleasesListBody
}

const GitHubReleasesListOperationID string = "get-github-releases-list"
const GitHubReleasesListDescription string = `Returns all github releases within the database.

Apply a start and end date or a unit name filter to restrict the data set.
`
const gitHubReleasesListSQL string = `
SELECT
	github_releases.*,
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
	) as github_repository,
	json_group_array(
		DISTINCT json_object(
			'id', units.id,
			'name', units.name
		)
	) filter ( where units.id is not null) as units
FROM github_releases
LEFT JOIN github_repositories ON github_repositories.id = github_releases.github_repository_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_releases.github_repository_id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN units ON units.id = github_teams_units.unit_id
GROUP BY github_releases.id
ORDER BY github_releases.date ASC;
`

// ApiGitHubTeamsListHandler queries the database for all github teams and returns a list including
// joins with github teams and aws accounts. There is no option to filter of limit the results.
//
// Endpoints:
//
//	/version/github/releases/list?unit=<unit>&start_date=<date>&end_date=<date>&interval=<interval>
func ApiGitHubReleasesListHandler(ctx context.Context, input *inputs.OptionalGroupedDateRangeInput) (response *GitHubReleasesListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease = []*models.GitHubRelease{}
		dbPath  string                  = ctx.Value(dbPathKey).(string)
		// replace string                  = "{WHERE}"
		sqlStmt string                  = gitHubReleasesListSQL
		param   statements.Named        = input
		body    *GitHubReleasesListBody = &GitHubReleasesListBody{
			Request:   input,
			Operation: GitHubReleasesListOperationID,
		}
	)
	// setup response
	response = &GitHubReleasesListResponse{}
	// // if there is a unit, setup the where clause
	// // otherwise remove it
	// if input.Unit != "" {
	// 	sqlStmt = strings.ReplaceAll(sqlStmt, replace, "WHERE units.name = :unit ")
	// } else {
	// 	sqlStmt = strings.ReplaceAll(sqlStmt, replace, "")
	// }

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] github releases list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubRelease](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] github releases list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github releases list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

func RegisterGitHubRelases(api huma.API) {
	var uri string = "/{version}/" + GitHubReleasesSegment + "/list"

	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubReleasesListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List GitHub releases",
		Description:   GitHubReleasesListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubReleasesTags,
	}, ApiGitHubReleasesListHandler)

}
