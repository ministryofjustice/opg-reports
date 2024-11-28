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
	GitHubRepositoryStandardsSegment string   = "github/standards"
	GitHubRepositoryStandardsTags    []string = []string{"github", "repository", "standards"}
)

// -- Release listing

// GitHubRepositoryStandardsListBody contains the resposne body to send back
// for a request to the /list endpoint
type GitHubRepositoryStandardsListBody struct {
	Operation string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.VersionUnitInput           `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRepositoryStandard `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// GitHubRepositoryStandardsListResponse is the main response struct
type GitHubRepositoryStandardsListResponse struct {
	Body *GitHubRepositoryStandardsListBody
}

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
	) as units
FROM github_repository_standards
LEFT JOIN github_repositories on github_repositories.id = github_repository_standards.github_repository_id
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_repository_standards.github_repository_id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN units on units.id = github_teams_units.unit_id
{WHERE}
GROUP BY github_repository_standards.id
ORDER BY github_repository_standards.github_repository_full_name ASC
;
`

// ApiGitHubRepositoryStandardsListHandler accepts and processes requests to the below endpoints.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/github/standards/list?unit=<unit>
func ApiGitHubRepositoryStandardsListHandler(ctx context.Context, input *inputs.VersionUnitInput) (response *GitHubRepositoryStandardsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRepositoryStandard = []*models.GitHubRepositoryStandard{}
		dbPath  string                             = ctx.Value(dbPathKey).(string)
		sqlStmt string                             = gitHubRepositoryStandardsListSQL
		where   string                             = ""
		replace string                             = "{WHERE}"
		param   statements.Named                   = input
		body    *GitHubRepositoryStandardsListBody = &GitHubRepositoryStandardsListBody{
			Request:   input,
			Operation: GitHubRepositoryStandardsListOperationID,
		}
	)
	// setup response
	response = &GitHubRepositoryStandardsListResponse{}
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
