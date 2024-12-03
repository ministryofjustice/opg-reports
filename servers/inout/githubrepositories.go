package inout

import "github.com/ministryofjustice/opg-reports/models"

type GitHubRepositoriesListBody struct {
	Operation string                     `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *VersionUnitInput          `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRepository `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                    `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

type GitHubRepositoriesListResponse struct {
	Body *GitHubRepositoriesListBody
}
