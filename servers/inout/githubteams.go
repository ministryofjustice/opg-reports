package inout

import "github.com/ministryofjustice/opg-reports/models"

type GitHubTeamsListBody struct {
	Operation string               `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *VersionUnitInput    `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubTeam `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error              `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}
type GitHubTeamsResponse struct {
	Body *GitHubTeamsListBody
}
