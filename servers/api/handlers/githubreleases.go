package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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

// -- Release listing

// GitHubReleasesListBody contains the resposne body to send back
// for a request to the /list endpoint
type GitHubReleasesListBody struct {
	Operation string                         `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.OptionalDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRelease        `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                        `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// GitHubReleasesListResponse is the main response struct for the
// /list endpoint
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
{WHERE}
LEFT JOIN github_repositories ON github_repositories.id = github_releases.github_repository_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_releases.github_repository_id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN units ON units.id = github_teams_units.unit_id
GROUP BY github_releases.id
ORDER BY github_releases.date ASC;
`

// ApiGitHubReleasesListHandler accepts and processes requests to the below endpoints.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/github/releases/list?unit=<unit>&start_date=<date>&end_date=<date>
func ApiGitHubReleasesListHandler(ctx context.Context, input *inputs.OptionalDateRangeInput) (response *GitHubReleasesListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease = []*models.GitHubRelease{}
		dbPath  string                  = ctx.Value(dbPathKey).(string)
		where   string                  = ""
		replace string                  = "{WHERE}"
		sqlStmt string                  = gitHubReleasesListSQL
		param   statements.Named        = input
		body    *GitHubReleasesListBody = &GitHubReleasesListBody{
			Request:   input,
			Operation: GitHubReleasesListOperationID,
		}
	)
	// setup response
	response = &GitHubReleasesListResponse{}

	// check for start, end and unit being passed
	if input.StartDate != "" && input.EndDate != "" && input.Unit != "" {
		where = "WHERE github_releases.date >= :start_date AND github_releases.date < :end_date AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else if input.Unit != "" {
		where = "WHERE units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

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

// -- Release count

// GitHubReleasesCountBody contains the resposne details for a request to the /count
// endpoint
type GitHubReleasesCountBody struct {
	Operation string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.RequiredGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRelease               `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// GitHubReleasesCountResponse is used by the /count endpoint
type GitHubReleasesCountResponse struct {
	Body *GitHubReleasesCountBody
}

const GitHubReleasesCountOperationID string = "get-github-releases-count"
const GitHubReleasesCountDescription string = `Returns count of github releases within the database between start_date and end_date.

Can also be filtered by unit name.`
const gitHubReleasesCountSQL string = `
SELECT
	COUNT(DISTINCT github_releases.id) as count,
	strftime(:date_format, github_releases.date) as date
FROM github_releases
LEFT JOIN github_repositories ON github_repositories.id = github_releases.github_repository_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_releases.github_repository_id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN units ON units.id = github_teams_units.unit_id
WHERE
	github_releases.date >= :start_date
	AND github_releases.date < :end_date
	{WHERE}
GROUP BY strftime(:date_format, github_releases.date)
ORDER BY strftime(:date_format, github_releases.date) ASC;
`

// ApiGitHubReleasesCountHandler accepts and processes requests to the below endpoints.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query.
//
// Endpoints:
//
//	/version/github/releases/count/{interval}/{start_date}/{end_date}?unit=<unit>
func ApiGitHubReleasesCountHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeInput) (response *GitHubReleasesCountResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease  = []*models.GitHubRelease{}
		dbPath  string                   = ctx.Value(dbPathKey).(string)
		where   string                   = ""
		replace string                   = "{WHERE}"
		sqlStmt string                   = gitHubReleasesCountSQL
		param   statements.Named         = input
		body    *GitHubReleasesCountBody = &GitHubReleasesCountBody{
			Request:   input,
			Operation: GitHubReleasesCountOperationID,
		}
	)
	// setup response
	response = &GitHubReleasesCountResponse{}

	// check for unit
	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] github releases count adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubRelease](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] github releases count select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github releases count selection failed."))
	} else {
		body.Result = results
	}
	// blank the date format
	body.Request.DateFormat = ""
	response.Body = body
	return
}

// const gitHubReleasesCountSQL string = `
// SELECT
// 	COUNT(DISTINCT github_releases.id) as count,
// 	strftime(:date_format, github_releases.date) as date,
// 	json_object(
// 		'id', github_repositories.id,
// 		'ts', github_repositories.ts,
// 		'owner', github_repositories.owner,
// 		'name', github_repositories.name,
// 		'full_name', github_repositories.full_name,
// 		'created_at', github_repositories.created_at,
// 		'default_branch', github_repositories.default_branch,
// 		'archived', github_repositories.archived,
// 		'private', github_repositories.private,
// 		'license', github_repositories.license,
// 		'last_commit_date', github_repositories.last_commit_date
// 	) as github_repository,
// 	json_group_array(
// 		DISTINCT json_object(
// 			'id', units.id,
// 			'name', units.name
// 		)
// 	) filter ( where units.id is not null) as units
// FROM github_releases
// LEFT JOIN github_repositories ON github_repositories.id = github_releases.github_repository_id
// LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_releases.github_repository_id
// LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
// LEFT JOIN units ON units.id = github_teams_units.unit_id
// WHERE
// 	github_releases.date >= :start_date
// 	AND github_releases.date < :end_date
// 	{WHERE}
// GROUP BY github_releases.id, strftime(:date_format, github_releases.date)
// ORDER BY strftime(:date_format, github_releases.date) ASC;
// `

func RegisterGitHubRelases(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + GitHubReleasesSegment + "/list"
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

	uri = "/{version}/" + GitHubReleasesSegment + "/count/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubReleasesCountOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Count GitHub releases",
		Description:   GitHubReleasesCountDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubReleasesTags,
	}, ApiGitHubReleasesCountHandler)

}
