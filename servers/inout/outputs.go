package inout

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/models"
)

// OperationBody is a base body that contains operation id string and
// errors that all responses will contain
type OperationBody struct {
	Operation string  `json:"operation" doc:"contains the operation id"`
	Errors    []error `json:"errors" doc:"list of any errors that occured in the request"`
}

func NewOperationBody() *OperationBody {
	return &OperationBody{}
}

// ColumnBody is for api data that is likely shown as a table and provides
// column details to allow maniluplation into tabluar structure
type ColumnBody struct {
	ColumnOrder  []string                 `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{} `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
}

func NewColumnBody() *ColumnBody {
	return &ColumnBody{}
}

// DateRangeBody provides list of dates that an api call convers - used in conjuntion
// with start / end date queries
type DateRangeBody struct {
	DateRange []string `json:"date_range" db:"-" doc:"all dates within the range requested"`
}

func NewDateRangeBody() *DateRangeBody {
	return &DateRangeBody{}
}

type RequestBody[T any] struct {
	Request T `json:"request" doc:"the original request"`
}

func NewDateRequestBody[T any](t T) *RequestBody[T] {
	return &RequestBody[T]{Request: t}
}

type ResultBody[T dbs.Row] struct {
	Result []T `json:"result" doc:"list of all units returned by the api."`
}

func NewResultBody[T dbs.Row]() *ResultBody[T] {
	var t []T = []T{}
	return &ResultBody[T]{Result: t}
}

// DateWideTableBody is for api data that should be converted into table rows
// that should have date range as additional coloumns and will therefore require
// transformation
type DateWideTableBody struct {
	TableRows map[string]map[string]interface{} `json:"table_rows"` // Used for post processing
}

func NewDateWideTableBody() *DateWideTableBody {
	return &DateWideTableBody{TableRows: map[string]map[string]interface{}{}}
}

// AwsAccountsListBody is api an response body. Intended to be used for
// returning a list of aws accounts without grouping
//
// Returned from handlers.ApiAwsAccountsListHandler
type AwsAccountsListBody struct {
	*OperationBody
	*RequestBody[*VersionUnitInput]
	*ResultBody[*models.AwsAccount]
}

func NewAwsAccountsListBody() *AwsAccountsListBody {
	return &AwsAccountsListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&VersionUnitInput{}),
		ResultBody:    NewResultBody[*models.AwsAccount](),
	}
}

// AwsCostsTotalBody is api an response body. Intended to be used for
// returning a single grouped sum of aws costs based on date and unit filters.
//
// Returned from handlers.ApiAwsCostsTotalHandler
type AwsCostsTotalBody struct {
	*OperationBody
	*RequestBody[*DateRangeUnitInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsTotalBody() *AwsCostsTotalBody {
	return &AwsCostsTotalBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&DateRangeUnitInput{}),
		ResultBody:    NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsListBody is api an response body. Intended to be used for
// returning a list of aws costs without grouping, but limited to date and unit filters.
//
// Returned from handlers.ApiAwsCostsListHandler
type AwsCostsListBody struct {
	*OperationBody
	*RequestBody[*DateRangeUnitInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsListBody() *AwsCostsListBody {
	return &AwsCostsListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&DateRangeUnitInput{}),
		ResultBody:    NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsTaxesBody is api an response body. Intended to be used for
// returning a list of aws costs grouped by time period and if the sum, contains
// tax or not.
//
// Returned from handlers.ApiAwsCostsTaxesHandler
// Transform: TransformToDateWideTable
type AwsCostsTaxesBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeUnitInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsTaxesBody() *AwsCostsTaxesBody {
	return &AwsCostsTaxesBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeUnitInput{}),
		ResultBody:        NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsSumBody is api an response body. Intended to be used for
// returning a list of total costs per time period (cost per month)
// between the dates and possibly filtered by unit name
//
// Returned from handlers.ApiAwsCostsSumHandler
// Transform: TransformToDateWideTable
type AwsCostsSumBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeUnitInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsSumBody() *AwsCostsSumBody {
	return &AwsCostsSumBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeUnitInput{}),
		ResultBody:        NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsSumPerUnitBody is api an response body. Intended to be used for
// returning a list of total costs per time period & unit (cost per month & unit)
// between the dates passed.
//
// Returned from handlers.ApiAwsCostsSumPerUnitHandler
// Transform: TransformToDateWideTable
type AwsCostsSumPerUnitBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsSumPerUnitBody() *AwsCostsSumPerUnitBody {
	return &AwsCostsSumPerUnitBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeInput{}),
		ResultBody:        NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsSumPerUnitEnvBody is api an response body. Intended to be used for
// returning a list of total costs per time period, unit & environment
// (cost per month, for dev / prod & unit) between the dates passed.
//
// Returned from handlers.ApiAwsCostsSumPerUnitEnvHandler
// Transform: TransformToDateWideTable
type AwsCostsSumPerUnitEnvBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsSumPerUnitEnvBody() *AwsCostsSumPerUnitEnvBody {
	return &AwsCostsSumPerUnitEnvBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeInput{}),
		ResultBody:        NewResultBody[*models.AwsCost](),
	}
}

// AwsCostsSumFullDetailsBody is api an response body. Intended to be used for
// returning a list of total costs at a very granual level, grouping by aws
// account, service, region, time period, team name and environment.
// This is then filtered by date range and optional unit name.
//
// Returned from handlers.ApiAwsCostsSumFullDetailsHandler
// Transform: TransformToDateWideTable
type AwsCostsSumFullDetailsBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeUnitInput]
	*ResultBody[*models.AwsCost]
}

func NewAwsCostsSumFullDetailsBody() *AwsCostsSumFullDetailsBody {
	return &AwsCostsSumFullDetailsBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeUnitInput{}),
		ResultBody:        NewResultBody[*models.AwsCost](),
	}
}

// AwsUptimeListBody is api an response body. Intended to be used for
// returning a list of all uptime checks between two dates with optional
// filter for a unit name.
//
// Returned from handlers.ApiAwsUptimeListHandler
type AwsUptimeListBody struct {
	*OperationBody
	*RequestBody[*DateRangeUnitInput]
	*ResultBody[*models.AwsUptime]
}

func NewAwsUptimeListBody() *AwsUptimeListBody {
	return &AwsUptimeListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&DateRangeUnitInput{}),
		ResultBody:    NewResultBody[*models.AwsUptime](),
	}
}

// AwsUptimeAveragesBody is api an response body. Used to return a
// list of average uptime values between start and end dates passed
// grouped by time period with optional unit filter.
//
// Returned from handlers.ApiAwsUptimeAveragesHandler
type AwsUptimeAveragesBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeUnitInput]
	*ResultBody[*models.AwsUptime]
}

func NewAwsUptimeAveragesBody() *AwsUptimeAveragesBody {
	return &AwsUptimeAveragesBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeUnitInput{}),
		ResultBody:        NewResultBody[*models.AwsUptime](),
	}
}

// AwsUptimeAveragesPerUnitBody is api an response body. Used to return a
// list of average uptime values between start and end dates passed and
// grouped by time period and unit name.
//
// Returned from handlers.ApiAwsUptimeAveragesPerUnitHandler
type AwsUptimeAveragesPerUnitBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeInput]
	*ResultBody[*models.AwsUptime]
}

func NewAwsUptimeAveragesPerUnitBody() *AwsUptimeAveragesPerUnitBody {
	return &AwsUptimeAveragesPerUnitBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeInput{}),
		ResultBody:        NewResultBody[*models.AwsUptime](),
	}
}

// GitHubReleasesListBody is api an response body. Used to return a
// list of all releases between dates passed with and optional filter
// by unit name
//
// Returned from handlers.ApiGitHubReleasesListHandler
type GitHubReleasesListBody struct {
	*OperationBody
	*RequestBody[*DateRangeUnitInput]
	*ResultBody[*models.GitHubRelease]
}

func NewGitHubReleasesListBody() *GitHubReleasesListBody {
	return &GitHubReleasesListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&DateRangeUnitInput{}),
		ResultBody:    NewResultBody[*models.GitHubRelease](),
	}
}

// GitHubReleasesCountBody is api an response body. Used to return a
// count per time period of all releases between dates passed with
// an optional filter by unit name
//
// Returned from handlers.ApiGitHubReleasesCountHandler
type GitHubReleasesCountBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeUnitInput]
	*ResultBody[*models.GitHubRelease]
}

func NewGitHubReleasesCountBody() *GitHubReleasesCountBody {
	return &GitHubReleasesCountBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeUnitInput{}),
		ResultBody:        NewResultBody[*models.GitHubRelease](),
	}
}

// GitHubReleasesCountPerUnitBody is api an response body. Used to return a
// count per time period and unit of all releases between dates passed
//
// Returned from handlers.ApiGitHubReleasesCountPerUnitHandler
type GitHubReleasesCountPerUnitBody struct {
	*OperationBody
	*ColumnBody
	*DateWideTableBody
	*DateRangeBody
	*RequestBody[*RequiredGroupedDateRangeInput]
	*ResultBody[*models.GitHubRelease]
}

func NewGitHubReleasesCountPerUnitBody() *GitHubReleasesCountPerUnitBody {
	return &GitHubReleasesCountPerUnitBody{
		OperationBody:     NewOperationBody(),
		ColumnBody:        NewColumnBody(),
		DateWideTableBody: NewDateWideTableBody(),
		DateRangeBody:     NewDateRangeBody(),
		RequestBody:       NewDateRequestBody(&RequiredGroupedDateRangeInput{}),
		ResultBody:        NewResultBody[*models.GitHubRelease](),
	}
}

// GitHubRepositoriesListBody is api an response body. Used to return a
// list of all repositories with an optional filter by unit name
//
// Returned from handlers.ApiGitHubRepositoriesListHandler
type GitHubRepositoriesListBody struct {
	*OperationBody
	*RequestBody[*VersionUnitInput]
	*ResultBody[*models.GitHubRepository]
}

func NewGitHubRepositoriesListBody() *GitHubRepositoriesListBody {
	return &GitHubRepositoriesListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&VersionUnitInput{}),
		ResultBody:    NewResultBody[*models.GitHubRepository](),
	}
}

// GitHubRepositoryStandardsListBody contains the resposne body to send back
// for a request to the /list endpoint
//
// Returned from handlers.ApiGitHubRepositoryStandardsListHandler
type GitHubRepositoryStandardsListBody struct {
	*OperationBody
	*RequestBody[*VersionUnitInput]
	*ResultBody[*models.GitHubRepositoryStandard]
	BaselineCounters *GitHubRepositoryStandardsCount `json:"baseline_counters" doc:"Compliance counters"`
	ExtendedCounters *GitHubRepositoryStandardsCount `json:"extended_counters" doc:"Compliance counters"`
}

// GitHubRepositoryStandardsCount
type GitHubRepositoryStandardsCount struct {
	Total      int     `json:"total" db:"-" faker:"-" doc:"Total number or records"`
	Compliant  int     `json:"compliant" db:"-" faker:"-" doc:"Number or records that comply."`
	Percentage float64 `json:"percentage" db:"-" faker:"-" doc:"Percentage of record that comply"`
}

func NewGitHubRepositoryStandardsListBody() *GitHubRepositoryStandardsListBody {
	return &GitHubRepositoryStandardsListBody{
		OperationBody:    NewOperationBody(),
		RequestBody:      NewDateRequestBody(&VersionUnitInput{}),
		ResultBody:       NewResultBody[*models.GitHubRepositoryStandard](),
		BaselineCounters: &GitHubRepositoryStandardsCount{},
		ExtendedCounters: &GitHubRepositoryStandardsCount{},
	}
}

// GitHubTeamsListBody is api an response body. Used to return a
// list of all teams with an optional filter by unit name
//
// Returned from handlers.ApiGitHubTeamsListHandler
type GitHubTeamsListBody struct {
	*OperationBody
	*RequestBody[*VersionUnitInput]
	*ResultBody[*models.GitHubTeam]
}

func NewGitHubTeamsListBody() *GitHubTeamsListBody {
	return &GitHubTeamsListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&VersionUnitInput{}),
		ResultBody:    NewResultBody[*models.GitHubTeam](),
	}
}

// UnitsListBody is api an response body. Used to return a
// list of all units.
//
// Returned from handlers.ApiUnitsListHandler
type UnitsListBody struct {
	*OperationBody
	*RequestBody[*VersionInput]
	*ResultBody[*models.Unit]
}

func NewUnitsListBody() *UnitsListBody {
	return &UnitsListBody{
		OperationBody: NewOperationBody(),
		RequestBody:   NewDateRequestBody(&VersionInput{}),
		ResultBody:    NewResultBody[*models.Unit](),
	}
}

// AwsAccountsListResponse is the main object returned from a huma handler and
// contains the body
//
// Returned from handlers.ApiAwsAccountsListHandler
type AwsAccountsListResponse struct {
	Body *AwsAccountsListBody
}

// AwsCostsTotalResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsTotalHandler
type AwsCostsTotalResponse struct {
	Body *AwsCostsTotalBody
}

// AwsCostsListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsListHandler
type AwsCostsListResponse struct {
	Body *AwsCostsListBody
}

// AwsCostsTaxesResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsTaxesHandler
type AwsCostsTaxesResponse struct {
	Body *AwsCostsTaxesBody
}

// AwsCostsSumResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumHandler
type AwsCostsSumResponse struct {
	Body *AwsCostsSumBody
}

// AwsCostsSumPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumPerUnitHandler
type AwsCostsSumPerUnitResponse struct {
	Body *AwsCostsSumPerUnitBody
}

// AwsCostsSumPerUnitEnvResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumPerUnitEnvHandler
type AwsCostsSumPerUnitEnvResponse struct {
	Body *AwsCostsSumPerUnitEnvBody
}

// AwsCostsSumFullDetailsResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsCostsSumFullDetailsHandler
type AwsCostsSumFullDetailsResponse struct {
	Body *AwsCostsSumFullDetailsBody
}

// AwsUptimeListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeListHandler
type AwsUptimeListResponse struct {
	Body *AwsUptimeListBody
}

// AwsUptimeAveragesResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeAveragesHandler
type AwsUptimeAveragesResponse struct {
	Body *AwsUptimeAveragesBody
}

// AwsUptimeAveragesPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiAwsUptimeAveragesPerUnitHandler
type AwsUptimeAveragesPerUnitResponse struct {
	Body *AwsUptimeAveragesPerUnitBody
}

// GitHubReleasesListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesListHandler
type GitHubReleasesListResponse struct {
	Body *GitHubReleasesListBody
}

// GitHubReleasesCountResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesCountHandler
type GitHubReleasesCountResponse struct {
	Body *GitHubReleasesCountBody
}

// GitHubReleasesCountPerUnitResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubReleasesCountPerUnitHandler
type GitHubReleasesCountPerUnitResponse struct {
	Body *GitHubReleasesCountPerUnitBody
}

// GitHubRepositoriesListResponse is the main object returned from a huma handler and
// contains the body with more data in.
//
// Returned from handlers.ApiGitHubRepositoriesListHandler
type GitHubRepositoriesListResponse struct {
	Body *GitHubRepositoriesListBody
}

// GitHubRepositoryStandardsListResponse is the main response struct
//
// Returned from handlers.ApiGitHubRepositoryStandardsListHandler
type GitHubRepositoryStandardsListResponse struct {
	Body *GitHubRepositoryStandardsListBody
}

// GitHubTeamsResponse is a response struct from huma api handler
//
// Returned from handlers.ApiGitHubTeamsListHandler
type GitHubTeamsResponse struct {
	Body *GitHubTeamsListBody
}

// UnitsListResponse is a response struct from huma api handler
//
// Returned from handlers.ApiUnitsListHandler
type UnitsListResponse struct {
	Body *UnitsListBody
}
