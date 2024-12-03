package inout

import "github.com/ministryofjustice/opg-reports/models"

type GitHubReleasesListBody struct {
	Operation string                  `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *OptionalDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.GitHubRelease `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                 `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

type GitHubReleasesListResponse struct {
	Body *GitHubReleasesListBody
}

// GitHubReleasesCountBody contains the resposne details for a request to the /count
// endpoint
// Tabular
type GitHubReleasesCountBody struct {
	Operation    string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.GitHubRelease            `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                           `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                           `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}           `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{}  `json:"-"` // Used for post processing
}

// GitHubReleasesCountResponse is used by the /count endpoint
type GitHubReleasesCountResponse struct {
	Body *GitHubReleasesCountBody
}

// GitHubReleasesCountPerUnitBody contains the resposne details for a request to the /count-per-unit
// endpoint
type GitHubReleasesCountPerUnitBody struct {
	Operation    string                            `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeInput    `json:"request,omitempty" doc:"the original request"`
	Result       []*models.GitHubRelease           `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                          `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                          `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                           `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{} `json:"-"` // Used for post processing
}

// GitHubReleasesCountResponse is used by the /count endpoint
type GitHubReleasesCountPerUnitResponse struct {
	Body *GitHubReleasesCountPerUnitBody
}
