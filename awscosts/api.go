package awscosts

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
)

type api struct {
	Register func(api huma.API)
}

// --- inputs

// VersionInput is the version part of the uri path
type VersionInput struct {
	Version string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

// DateRangeInput describes input parameters to capture the start and end of a date range
// The date range is >= start_date < end_date so in the case covering all of Jan 2024
// you would want 2024-01-01 & 2024-02-01
// This makes find last days easier for users
type DateRangeInput struct {
	VersionInput
	StartDate string `json:"start_date" db:"start_date" query:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). Can be YYYY-MM or YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
	EndDate   string `json:"end_date" db:"end_date" query:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). Can be YYYY-MM or YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
}

// GroupedDateRangeInput expands upon DateRangeStatement to include
// interval field thats used for formatting the date column as part
// of the grouping
type GroupedDateRangeInput struct {
	DateRangeInput
	Interval   string `json:"interval" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

// Resolve to convert the interval input param into a date format for the db call
func (self *GroupedDateRangeInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = datastore.Sqlite.YearMonthFormat
	if self.Interval == "year" {
		self.DateFormat = datastore.Sqlite.YearFormat
	} else if self.Interval == "day" {
		self.DateFormat = datastore.Sqlite.YearMonthDayFormat
	}
	return nil
}

// just test as got a custom resolver
var _ huma.Resolver = (*GroupedDateRangeInput)(nil)

// FilterByUnitInput extends GroupedDateRangeInput adding the
// ability to filter the dataset by unit
type FilterByUnitInput struct {
	GroupedDateRangeInput
	Unit string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`
}

type StandardInput struct {
	FilterByUnitInput
}

type StandardBody struct {
	Type           string         `json:"type" doc:"States what type of data this is for front end handling"`
	Result         []*Cost        `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	OrderedColumns []string       `json:"ordered_columns" doc:"List of columns set in the order they should be rendered for each row"`
	Request        *StandardInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
}
type StandardResult struct {
	Body *StandardBody
}

// --- total

var totalDescription string = `Returns a single total for all the AWS costs between the start and end dates - excluding taxes.`

type TotalInput struct {
	DateRangeInput
}
type TotalBody struct {
	Type    string      `json:"type" doc:"States what type of data this is for front end handling"`
	Request *TotalInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result  float64     `json:"result" doc:"The total sum of all costs as a float without currency." example:"1357.7861"`
}
type TotalResult struct {
	Body *TotalBody
}

// apiTotal fetches total sum of all costs within the database
func apiTotal(ctx context.Context, input *TotalInput) (response *TotalResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result float64 = 0.0
	var bdy *TotalBody = &TotalBody{Request: input, Type: "total"}
	response = &TotalResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if result, err = datastore.Get[float64](ctx, db, DB.Total, input.StartDate, input.EndDate); err == nil {
		response.Body.Result = result
	}

	return
}

// --- TaxOverview

var taxOverviewDescription string = `Provides a list of the total costs per interval with and without tax between the start and end date specified.`

type TaxOverviewInput struct {
	GroupedDateRangeInput
}

type TaxOverviewBody struct {
	Type           string            `json:"type" doc:"States what type of data this is for front end handling"`
	Request        *TaxOverviewInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result         []*Cost           `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	OrderedColumns []string          `json:"ordered_columns" doc:"List of columns set in the order they should be rendered for each row"`
}
type TaxOverviewResult struct {
	Body *TaxOverviewBody
}

// apiTotal fetches total sum of all costs within the database
func apiTaxOverview(ctx context.Context, input *TaxOverviewInput) (response *TaxOverviewResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*Cost = []*Cost{}
	var bdy *TaxOverviewBody = &TaxOverviewBody{
		OrderedColumns: []string{"service"},
		Request:        input,
		Type:           "tax-overview",
	}
	response = &TaxOverviewResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if result, err = datastore.Select[[]*Cost](ctx, db, DB.TaxOverview, input); err == nil {
		response.Body.Result = result
	}

	return
}

// --- PerUnit

var perUnitDescription string = `Returns a list of cost data grouped by the unit field as well as the date.

Data is limited to the date range (>= start_date < ) and optional unit filter.`

// apiPerUnit handles getting data grouped by unit
func apiPerUnit(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*Cost = []*Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = DB.Unit

	bdy.OrderedColumns = []string{"unit"}
	bdy.Request = input
	bdy.Type = "unit"
	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if input.Unit != "" {
		stmt = DB.UnitFilter
	}

	if result, err = datastore.Select[[]*Cost](ctx, db, stmt, input); err == nil {
		response.Body.Result = result
	}

	return
}

// --- PerUnitEnv

var perUnitEnvDescription string = `Returns a list of cost data grouped by the unit and environment fields as well as the date.

Data is limited to the date range (>= start_date < ) and optional unit filter.`

// apiPerUnit handles getting data grouped by unit
func apiPerUnitEnv(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*Cost = []*Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = DB.UnitEnv

	bdy.OrderedColumns = []string{"unit", "environment"}
	bdy.Request = input
	bdy.Type = "unit-environment"
	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}
	if input.Unit != "" {
		stmt = DB.UnitEnvFilter
	}

	if result, err = datastore.Select[[]*Cost](ctx, db, stmt, input); err == nil {
		response.Body.Result = result
	}

	return
}

// --- PerUnitEnv

var detailDescription string = `Provides a list of the total costs grouped by a date interval, account_id, environment and service within start (>=) and end (<) dates passed.

Data is limited to the date range (>= start_date < ) and optional unit filter.`

// apiDetailed
func apiDetailed(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*Cost = []*Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = DB.Detailed

	bdy.OrderedColumns = []string{"account_id", "unit", "environment", "service"}
	bdy.Request = input
	bdy.Type = "detail"
	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}
	if input.Unit != "" {
		stmt = DB.DetailedFilter
	}

	if result, err = datastore.Select[[]*Cost](ctx, db, stmt, input); err == nil {
		response.Body.Result = result
	}

	return
}

func register(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-total",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/total",
		Summary:       "Total",
		Description:   totalDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiTotal)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-tax-overview",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/tax-overview/{interval}/",
		Summary:       "Tax overview",
		Description:   taxOverviewDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiTaxOverview)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-costs-per-unit",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/unit/{interval}/",
		Summary:       "Costs per unit",
		Description:   perUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiPerUnit)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-costs-per-unit-environment",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/unit-environment/{interval}/",
		Summary:       "Costs per unit & environment",
		Description:   perUnitEnvDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiPerUnitEnv)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-costs-detailed",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/detailed/{interval}/",
		Summary:       "Detailed Costs",
		Description:   detailDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiDetailed)

}
