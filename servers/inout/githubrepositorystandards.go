package inout

import "github.com/ministryofjustice/opg-reports/models"

// GitHubRepositoryStandardsListBody contains the resposne body to send back
// for a request to the /list endpoint
type GitHubRepositoryStandardsListBody struct {
	Operation string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *VersionUnitInput                  `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRepositoryStandard `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// GitHubRepositoryStandardsListResponse is the main response struct
type GitHubRepositoryStandardsListResponse struct {
	Body *GitHubRepositoryStandardsListBody
}
