package inout

import "github.com/ministryofjustice/opg-reports/models"

// GitHubRepositoryStandardsListBody contains the resposne body to send back
// for a request to the /list endpoint
type GitHubRepositoryStandardsListBody struct {
	Operation        string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request          *VersionUnitInput                  `json:"request,omitempty" doc:"the original request"`
	Result           []*models.GitHubRepositoryStandard `json:"result,omitempty" doc:"list of all units returned by the api."`
	BaselineCounters *GitHubRepositoryStandardsCount    `json:"baseline_counters" doc:"Compliance counters"`
	ExtendedCounters *GitHubRepositoryStandardsCount    `json:"extended_counters" doc:"Compliance counters"`
	Errors           []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

type GitHubRepositoryStandardsCount struct {
	Total      int     `json:"total" db:"-" faker:"-" doc:"Total number or records"`
	Compliant  int     `json:"compliant" db:"-" faker:"-" doc:"Number or records that comply."`
	Percentage float64 `json:"percentage" db:"-" faker:"-" doc:"Percentage of record that comply"`
}

// GitHubRepositoryStandardsListResponse is the main response struct
type GitHubRepositoryStandardsListResponse struct {
	Body *GitHubRepositoryStandardsListBody
}
