package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/cols"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/dbs/statements"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

var (
	GitHubReleasesSegment string   = "github/release"
	GitHubReleasesTags    []string = []string{"GitHub releases"}
)

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
WHERE
	github_releases.date >= :start_date
	AND github_releases.date < :end_date
	{WHERE}
GROUP BY github_releases.id
ORDER BY github_releases.date ASC;
`

// ApiGitHubReleasesListHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/github/release/list/{start_date}/{end_date}?unit=<unit>
func ApiGitHubReleasesListHandler(ctx context.Context, input *inout.DateRangeUnitInput) (response *inout.GitHubReleasesListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease       = []*models.GitHubRelease{}
		dbPath  string                        = ctx.Value(dbPathKey).(string)
		where   string                        = ""
		replace string                        = "{WHERE}"
		sqlStmt string                        = gitHubReleasesListSQL
		param   statements.Named              = input
		body    *inout.GitHubReleasesListBody = &inout.GitHubReleasesListBody{
			Request:   input,
			Operation: GitHubReleasesListOperationID,
		}
	)
	// setup response
	response = &inout.GitHubReleasesListResponse{}

	// check for start, end and unit being passed
	if input.Unit != "" {
		where = "AND units.Name = :unit "
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

const GitHubReleasesCountOperationID string = "get-github-releases-count"
const GitHubReleasesCountDescription string = `Returns count of github releases within the database between start_date and end_date.

Can also be filtered by unit name.`
const gitHubReleasesCountSQL string = `
SELECT
	'Count' as unit_name,
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

// ApiGitHubReleasesCountHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query.
//
// Endpoints:
//
//	/version/github/release/count/{interval}/{start_date}/{end_date}?unit=<unit>
func ApiGitHubReleasesCountHandler(ctx context.Context, input *inout.RequiredGroupedDateRangeUnitInput) (response *inout.GitHubReleasesCountResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease        = []*models.GitHubRelease{}
		dbPath  string                         = ctx.Value(dbPathKey).(string)
		where   string                         = ""
		replace string                         = "{WHERE}"
		sqlStmt string                         = gitHubReleasesCountSQL
		param   statements.Named               = input
		body    *inout.GitHubReleasesCountBody = &inout.GitHubReleasesCountBody{
			Request:     input,
			Operation:   GitHubReleasesCountOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
			// hard code the unit column to only have the word count
			ColumnValues: map[string][]interface{}{
				"unit_name": {"Count"},
			},
		}
	)
	// setup response
	response = &inout.GitHubReleasesCountResponse{}

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

const GitHubReleasesCountPerUnitOperationID string = "get-github-releases-count-per-unit"
const GitHubReleasesCountPerUnitDescription string = `Returns count of github releases within the database between start_date and end_date grouped by the unit name.

Can also be filtered by unit name.`

// gitHubReleasesCountPerUnitSQL starts from units and follows joins to the releases data:
//
//	unit -> github_teams_units -> github_repositories_github_teams -> github_repositories -> github_releases
const gitHubReleasesCountPerUnitSQL string = `
SELECT
	units.name as unit_name,
	COUNT(DISTINCT github_releases.id) as count,
	strftime(:date_format, github_releases.date) as date
FROM units
LEFT JOIN github_teams_units ON github_teams_units.unit_id = units.id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_team_id = github_teams_units.github_team_id
LEFT JOIN github_repositories ON github_repositories.id = github_repositories_github_teams.github_repository_id
LEFT JOIN github_releases ON github_releases.github_repository_id = github_repositories.id
WHERE
	github_releases.date >= :start_date
	AND github_releases.date < :end_date
GROUP BY units.id, strftime(:date_format, github_releases.date)
ORDER BY strftime(:date_format, github_releases.date), units.name ASC;
`

// ApiGitHubReleasesCountPerUnitHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query.
//
// Endpoints:
//
//	/version/github/release/count-per-unit/{interval}/{start_date}/{end_date}
func ApiGitHubReleasesCountPerUnitHandler(ctx context.Context, input *inout.RequiredGroupedDateRangeInput) (response *inout.GitHubReleasesCountPerUnitResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRelease               = []*models.GitHubRelease{}
		dbPath  string                                = ctx.Value(dbPathKey).(string)
		sqlStmt string                                = gitHubReleasesCountPerUnitSQL
		param   statements.Named                      = input
		body    *inout.GitHubReleasesCountPerUnitBody = &inout.GitHubReleasesCountPerUnitBody{
			Request:     input,
			Operation:   GitHubReleasesCountPerUnitOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
		}
	)
	// setup response
	response = &inout.GitHubReleasesCountPerUnitResponse{}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] github releases count per unit adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubRelease](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] github releases count per unit select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github releases count per unit selection failed."))
	} else {
		body.Result = results
	}
	// blank the date format
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

func RegisterGitHubRelases(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + GitHubReleasesSegment + "/list/{start_date}/{end_date}"
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

	uri = "/{version}/" + GitHubReleasesSegment + "/count-per-unit/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubReleasesCountPerUnitOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Count GitHub releases per unit",
		Description:   GitHubReleasesCountPerUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubReleasesTags,
	}, ApiGitHubReleasesCountPerUnitHandler)

}
