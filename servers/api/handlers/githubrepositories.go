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
	GitHubRepositoriesSegment string   = "github/repository"
	GitHubRepositoriesTags    []string = []string{"GitHub repositories"}
)

const GitHubRepositoriesListOperationID string = "get-github-repository-list"
const GitHubRepositoriesListDescription string = `Returns all github repositories`

const GitHubRepositoriesListSQL string = `
SELECT
	github_repositories.*,
	json_group_array(
		DISTINCT json_object(
			'id', units.id,
			'name', units.name
		)
	) filter ( where units.id is not null) as units,
	json_group_array(
		DISTINCT json_object(
			'id', github_teams.id,
			'slug', github_teams.slug
		)
	) filter ( where github_teams.id is not null) as github_teams
FROM github_repositories
LEFT JOIN github_repositories_github_teams ON github_repositories_github_teams.github_repository_id = github_repositories.id
LEFT JOIN github_teams_units ON github_teams_units.github_team_id = github_repositories_github_teams.github_team_id
LEFT JOIN units on units.id = github_teams_units.unit_id
LEFT JOIN github_teams on github_teams.id = github_teams_units.github_team_id
{WHERE}
GROUP BY github_repositories.id
ORDER BY github_repositories.full_name ASC
;
`

// ApiGitHubRepositoriesListHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/github/repositories/list?unit=<unit>
func ApiGitHubRepositoriesListHandler(ctx context.Context, input *inout.VersionUnitInput) (response *inout.GitHubRepositoriesListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.GitHubRepository        = []*models.GitHubRepository{}
		dbPath  string                            = ctx.Value(dbPathKey).(string)
		sqlStmt string                            = GitHubRepositoriesListSQL
		where   string                            = ""
		replace string                            = "{WHERE}"
		param   statements.Named                  = input
		body    *inout.GitHubRepositoriesListBody = &inout.GitHubRepositoriesListBody{
			Request:   input,
			Operation: GitHubRepositoriesListOperationID,
		}
	)
	// setup response
	response = &inout.GitHubRepositoriesListResponse{}
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
		slog.Error("[api] github repository list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.GitHubRepository](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] github repository list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("github repository list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

func RegisterGitHubRepositories(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + GitHubRepositoriesSegment + "/list"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   GitHubRepositoriesListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List GitHub repositories",
		Description:   GitHubRepositoriesListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          GitHubRepositoriesTags,
	}, ApiGitHubRepositoriesListHandler)

}
