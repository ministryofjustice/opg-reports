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
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

var (
	GitHubRepositoryStandardsSegment string   = "github/standard"
	GitHubRepositoryStandardsTags    []string = []string{"GitHub repository standards"}
)

const GitHubRepositoryStandardsListOperationID string = "get-github-standards-list"
const GitHubRepositoryStandardsListDescription string = `Returns all github standards data`

const gitHubRepositoryStandardsListSQL string = `
SELECT
	github_repository_standards.*,
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
	) as units,
	json_group_array(
		DISTINCT json_object(
			'id', github_teams.id,
			'slug', github_teams.slug
		)
	) as github_teams
FROM github_repository_standards
LEFT JOIN github_repositories on github_repositories.id = github_repository_standards.github_repository_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_repository_standards.github_repository_id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN github_teams ON github_teams.id = github_repositories_github_teams.github_team_id
LEFT JOIN units on units.id = github_teams_units.unit_id
WHERE
	github_repository_standards.is_archived = 0
	{WHERE}
GROUP BY github_repository_standards.id
ORDER BY github_repository_standards.github_repository_full_name ASC
;
`

// ApiGitHubRepositoryStandardsListHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/github/standards/list?unit=<unit>
func ApiGitHubRepositoryStandardsListHandler(ctx context.Context, input *inout.VersionUnitInput) (response *inout.GitHubRepositoryStandardsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		base    *inout.GitHubRepositoryStandardsCount
		ext     *inout.GitHubRepositoryStandardsCount
		results []*models.GitHubRepositoryStandard       = []*models.GitHubRepositoryStandard{}
		dbPath  string                                   = ctx.Value(dbPathKey).(string)
		sqlStmt string                                   = gitHubRepositoryStandardsListSQL
		where   string                                   = ""
		replace string                                   = "{WHERE}"
		param   statements.Named                         = input
		body    *inout.GitHubRepositoryStandardsListBody = &inout.GitHubRepositoryStandardsListBody{}
	)
	body.Request = input
	body.Operation = GitHubRepositoryStandardsListOperationID

	// setup response
	response = &inout.GitHubRepositoryStandardsListResponse{}
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
		slog.Error("[api] github repository standards list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubRepositoryStandard](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] github repository standards list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github repository standards list selection failed."))
	} else {
		body.Result = results
	}

	base = &inout.GitHubRepositoryStandardsCount{Total: len(results), Compliant: 0, Percentage: 0.0}
	ext = &inout.GitHubRepositoryStandardsCount{Total: len(results), Compliant: 0, Percentage: 0.0}
	// update counters
	for _, row := range results {
		if row.IsCompliantBaseline() {
			base.Compliant += 1
		}
		if row.IsCompliantExtended() {
			ext.Compliant += 1
		}
	}
	base.Percentage = (float64(base.Compliant) / float64(base.Total)) * 100
	ext.Percentage = (float64(ext.Compliant) / float64(ext.Total)) * 100
	body.BaselineCounters = base
	body.ExtendedCounters = ext
	response.Body = body
	return
}

func RegisterGitHubRepositoryStandards(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + GitHubRepositoryStandardsSegment + "/list"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubRepositoryStandardsListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List GitHub repository standards",
		Description:   GitHubRepositoryStandardsListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubRepositoryStandardsTags,
	}, ApiGitHubRepositoryStandardsListHandler)

}
