package inout

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/models"
)

// OperationBody is a base body that contains operation id string and
// errors that all responses will contain
type OperationBody struct {
	Operation string  `json:"operation,omitempty" doc:"contains the operation id"`
	Errors    []error `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// ResultBody contains the results as a generic list
type ResultBody[T dbs.Row] struct {
	Result []T `json:"result,omitempty" doc:"list of data returned by the api."`
}

func (self *ResultBody[T]) GetResults() []T {
	return self.Result
}

type RequestBody[T any] struct {
	Request T `json:"request,omitempty" doc:"the original request"`
}

// ColumnBody is for api data that is likely shown as a table and provides
// column details to allow maniluplation into tabluar structure
type ColumnBody struct {
	ColumnOrder  []string                 `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{} `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
}

func (self *ColumnBody) GetColumnValues() map[string][]interface{} {
	return self.ColumnValues
}

// DateRangeBody provides list of dates that an api call convers - used in conjuntion
// with start / end date queries
type DateRangeBody struct {
	DateRange []string `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
}

func (self *DateRangeBody) GetDateRange() []string {
	return self.DateRange
}

// DateWideTableBody is for api data that should be converted into table rows
// that should have date range as additional coloumns and will therefore require
// transformation
type DateWideTableBody struct {
	TableRows map[string]map[string]interface{} `json:"-"` // Used for post processing
}

// AwsAccountsListBody is api an response body. Intended to be used for
// returning a list of aws accounts without grouping
//
// Returned from handlers.ApiAwsAccountsListHandler
type AwsAccountsListBody struct {
	OperationBody
	ResultBody[*models.AwsAccount]
	RequestBody[*VersionUnitInput]
}

// AwsAccountsListResponse is the main object returned from a huma handler and
// contains the body
//
// Returned from handlers.ApiAwsAccountsListHandler
type AwsAccountsListResponse struct {
	Body *AwsAccountsListBody
}

// AwsCostsTotalBody is api an response body. Intended to be used for
// returning a single grouped sum of aws costs based on date and unit filters.
//
// Returned from handlers.ApiAwsCostsTotalHandler
type AwsCostsTotalBody struct {
	OperationBody
	ResultBody[*models.AwsCost]
	RequestBody[*DateRangeUnitInput]
}

// AwsCostsTotalResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsTotalHandler
type AwsCostsTotalResponse struct {
	Body *AwsCostsTotalBody
}

// AwsCostsListBody is api an response body. Intended to be used for
// returning a list of aws costs without grouping, but limited to date and unit filters.
//
// Returned from handlers.ApiAwsCostsListHandler
type AwsCostsListBody struct {
	OperationBody
	ResultBody[*models.AwsCost]
	RequestBody[*DateRangeUnitInput]
}

// AwsCostsListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsListHandler
type AwsCostsListResponse struct {
	Body *AwsCostsListBody
}

// AwsCostsTaxesBody is api an response body. Intended to be used for
// returning a list of aws costs grouped by time period and if the sum, contains
// tax or not.
//
// Returned from handlers.ApiAwsCostsTaxesHandler
// Transform: TransformToDateWideTable
type AwsCostsTaxesBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsCost]
	RequestBody[*RequiredGroupedDateRangeUnitInput]
}

// AwsCostsTaxesResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsTaxesHandler
type AwsCostsTaxesResponse struct {
	Body *AwsCostsTaxesBody
}

// AwsCostsSumBody is api an response body. Intended to be used for
// returning a list of total costs per time period (cost per month)
// between the dates and possibly filtered by unit name
//
// Returned from handlers.ApiAwsCostsSumHandler
// Transform: TransformToDateWideTable
type AwsCostsSumBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsCost]
	RequestBody[*RequiredGroupedDateRangeUnitInput]
}

// AwsCostsSumResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumHandler
type AwsCostsSumResponse struct {
	Body *AwsCostsSumBody
}

// AwsCostsSumPerUnitBody is api an response body. Intended to be used for
// returning a list of total costs per time period & unit (cost per month & unit)
// between the dates passed.
//
// Returned from handlers.ApiAwsCostsSumPerUnitHandler
// Transform: TransformToDateWideTable
type AwsCostsSumPerUnitBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsCost]
	RequestBody[*RequiredGroupedDateRangeInput]
}

// AwsCostsSumPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumPerUnitHandler
type AwsCostsSumPerUnitResponse struct {
	Body *AwsCostsSumPerUnitBody
}

// AwsCostsSumPerUnitEnvBody is api an response body. Intended to be used for
// returning a list of total costs per time period, unit & environment
// (cost per month, for dev / prod & unit) between the dates passed.
//
// Returned from handlers.ApiAwsCostsSumPerUnitEnvHandler
// Transform: TransformToDateWideTable
type AwsCostsSumPerUnitEnvBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsCost]
	RequestBody[*RequiredGroupedDateRangeInput]
}

// AwsCostsSumPerUnitEnvResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumPerUnitEnvHandler
type AwsCostsSumPerUnitEnvResponse struct {
	Body *AwsCostsSumPerUnitEnvBody
}

// AwsCostsSumFullDetailsBody is api an response body. Intended to be used for
// returning a list of total costs at a very granual level, grouping by aws
// account, service, region, time period, team name and environment.
// This is then filtered by date range and optional unit name.
//
// Returned from handlers.ApiAwsCostsSumFullDetailsHandler
// Transform: TransformToDateWideTable
type AwsCostsSumFullDetailsBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsCost]
	RequestBody[*RequiredGroupedDateRangeUnitInput]
}

// AwsCostsSumFullDetailsResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumFullDetailsHandler
type AwsCostsSumFullDetailsResponse struct {
	Body *AwsCostsSumFullDetailsBody
}

// AwsUptimeListBody is api an response body. Intended to be used for
// returning a list of all uptime checks between two dates with optional
// filter for a unit name.
//
// Returned from handlers.ApiAwsUptimeListHandler
type AwsUptimeListBody struct {
	OperationBody
	ResultBody[*models.AwsUptime]
	RequestBody[*DateRangeUnitInput]
}

// AwsUptimeListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeListHandler
type AwsUptimeListResponse struct {
	Body *AwsUptimeListBody
}

// AwsUptimeAveragesBody is api an response body. Used to return a
// list of average uptime values between start and end dates passed
// grouped by time period with optional unit filter.
//
// Returned from handlers.ApiAwsUptimeAveragesHandler
type AwsUptimeAveragesBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsUptime]
	RequestBody[*RequiredGroupedDateRangeUnitInput]
}

// AwsUptimeAveragesResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeAveragesHandler
type AwsUptimeAveragesResponse struct {
	Body *AwsUptimeAveragesBody
}

// AwsUptimeAveragesPerUnitBody is api an response body. Used to return a
// list of average uptime values between start and end dates passed and
// grouped by time period and unit name.
//
// Returned from handlers.ApiAwsUptimeAveragesPerUnitHandler
type AwsUptimeAveragesPerUnitBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.AwsUptime]
	RequestBody[*RequiredGroupedDateRangeInput]
}

// AwsUptimeAveragesPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeAveragesPerUnitHandler
type AwsUptimeAveragesPerUnitResponse struct {
	Body *AwsUptimeAveragesPerUnitBody
}

// GitHubReleasesListBody is api an response body. Used to return a
// list of all releases between dates passed with and optional filter
// by unit name
//
// Returned from handlers.ApiGitHubReleasesListHandler
type GitHubReleasesListBody struct {
	OperationBody
	ResultBody[*models.GitHubRelease]
	RequestBody[*DateRangeUnitInput]
}

// GitHubReleasesListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesListHandler
type GitHubReleasesListResponse struct {
	Body *GitHubReleasesListBody
}

// GitHubReleasesCountBody is api an response body. Used to return a
// count per time period of all releases between dates passed with
// an optional filter by unit name
//
// Returned from handlers.ApiGitHubReleasesCountHandler
type GitHubReleasesCountBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.GitHubRelease]
	Request *RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
}

// GitHubReleasesCountResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesCountHandler
type GitHubReleasesCountResponse struct {
	Body *GitHubReleasesCountBody
}

// GitHubReleasesCountPerUnitBody is api an response body. Used to return a
// count per time period and unit of all releases between dates passed
//
// Returned from handlers.ApiGitHubReleasesCountPerUnitHandler
type GitHubReleasesCountPerUnitBody struct {
	OperationBody
	ColumnBody
	DateWideTableBody
	DateRangeBody
	ResultBody[*models.GitHubRelease]
	RequestBody[*RequiredGroupedDateRangeInput]
}

// GitHubReleasesCountPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesCountPerUnitHandler
type GitHubReleasesCountPerUnitResponse struct {
	Body *GitHubReleasesCountPerUnitBody
}

// GitHubRepositoriesListBody is api an response body. Used to return a
// list of all repositories with an optional filter by unit name
//
// Returned from handlers.ApiGitHubRepositoriesListHandler
type GitHubRepositoriesListBody struct {
	OperationBody
	ResultBody[*models.GitHubRepository]
	RequestBody[*VersionUnitInput]
}

// GitHubRepositoriesListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubRepositoriesListHandler
type GitHubRepositoriesListResponse struct {
	Body *GitHubRepositoriesListBody
}

// GitHubRepositoryStandardsListBody contains the resposne body to send back
// for a request to the /list endpoint
//
// Returned from handlers.ApiGitHubRepositoryStandardsListHandler
type GitHubRepositoryStandardsListBody struct {
	OperationBody
	ResultBody[*models.GitHubRepositoryStandard]
	RequestBody[*VersionUnitInput]
	BaselineCounters *GitHubRepositoryStandardsCount `json:"baseline_counters" doc:"Compliance counters"`
	ExtendedCounters *GitHubRepositoryStandardsCount `json:"extended_counters" doc:"Compliance counters"`
}

// GitHubRepositoryStandardsCount
type GitHubRepositoryStandardsCount struct {
	Total      int     `json:"total" db:"-" faker:"-" doc:"Total number or records"`
	Compliant  int     `json:"compliant" db:"-" faker:"-" doc:"Number or records that comply."`
	Percentage float64 `json:"percentage" db:"-" faker:"-" doc:"Percentage of record that comply"`
}

// GitHubRepositoryStandardsListResponse is the main response struct
//
// Returned from handlers.ApiGitHubRepositoryStandardsListHandler
type GitHubRepositoryStandardsListResponse struct {
	Body *GitHubRepositoryStandardsListBody
}

// GitHubTeamsListBody is api an response body. Used to return a
// list of all teams with an optional filter by unit name
//
// Returned from handlers.ApiGitHubTeamsListHandler
type GitHubTeamsListBody struct {
	OperationBody
	ResultBody[*models.GitHubTeam]
	RequestBody[*VersionUnitInput]
}

// GitHubTeamsResponse is a response struct from huma api handler
//
// Returned from handlers.ApiGitHubTeamsListHandler
type GitHubTeamsResponse struct {
	Body *GitHubTeamsListBody
}

// UnitsListBody is api an response body. Used to return a
// list of all units.
//
// Returned from handlers.ApiUnitsListHandler
type UnitsListBody struct {
	OperationBody
	ResultBody[*models.Unit]
	RequestBody[*VersionInput]
}

// UnitsListResponse is a response struct from huma api handler
//
// Returned from handlers.ApiUnitsListHandler
type UnitsListResponse struct {
	Body *UnitsListBody
}
