// Package costsapi provides the all elements for hum api endpoints relating to cost queries
//
// Contains structs for all query inputs (named as `*Input`), all result bodies (as `*Body`) and
// The actual result structs returned by the end point (named `*Result`).
//
// Contains all api handlers as functions (internal, named `api*`)
//
// Exposes main `Register` method which accepts a huma api and attaches all the
// endpoints to that api with suitable descriptions and details
package costsapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsdb"
)

const Segment string = "costs"
const Tag string = "Cost Data"

// --- total

var totalDescription string = `Returns a single total for all the AWS costs between the start and end dates - excluding taxes.`

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

	if result, err = datastore.Get[float64](ctx, db, costsdb.Total, input.StartDate, input.EndDate); err == nil {
		response.Body.Result = result
	}

	return
}

// --- TaxOverview

var taxOverviewDescription string = `Provides a list of the total costs per interval with and without tax between the start and end date specified.

Each result returned has the following fields:


	service (either 'Including Tax' or 'Excluding Tax')
	date
	cost
`

// apiTotal fetches total sum of all costs within the database
func apiTaxOverview(ctx context.Context, input *TaxOverviewInput) (response *TaxOverviewResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*costs.Cost = []*costs.Cost{}
	var bdy *TaxOverviewBody = &TaxOverviewBody{
		ColumnOrder: []string{"service"},
		Request:     input,
		Type:        "tax-overview",
		DateRange:   convert.DateRange(input.StartTime(), input.EndTime(), input.Interval),
	}
	response = &TaxOverviewResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if result, err = datastore.Select[[]*costs.Cost](ctx, db, costsdb.TaxOverview, input); err == nil {
		response.Body.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
		response.Body.Result = result
	}

	return
}

// --- PerUnit

var perUnitDescription string = `Returns a list of cost data grouped by the unit field as well as the date.

Data is limited to the date range (>= start_date < ) and optional unit filter.

Each result returned has the following fields:


	unit
	date
	cost
`

// apiPerUnit handles getting data grouped by unit
func apiPerUnit(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*costs.Cost = []*costs.Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = costsdb.PerUnit

	bdy.ColumnOrder = []string{"unit"}
	bdy.Request = input
	bdy.Type = "unit"
	bdy.DateRange = convert.DateRange(input.StartTime(), input.EndTime(), input.Interval)

	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if input.Unit != "" {
		stmt = costsdb.PerUnitForUnit
	}

	if result, err = datastore.Select[[]*costs.Cost](ctx, db, stmt, input); err == nil {
		response.Body.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
		response.Body.Result = result
	}

	return
}

// --- PerUnitEnv

var perUnitEnvDescription string = `Returns a list of cost data grouped by the unit and environment fields as well as the date.

Data is limited to the date range (>= start_date < ) and optional unit filter.

Each result returned has the following fields:


	unit
	environment
	date
	cost
`

// apiPerUnit handles getting data grouped by unit
func apiPerUnitEnv(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*costs.Cost = []*costs.Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = costsdb.PerUnitEnvironment

	bdy.ColumnOrder = []string{"unit", "environment"}
	bdy.Request = input
	bdy.Type = "unit-environment"
	bdy.DateRange = convert.DateRange(input.StartTime(), input.EndTime(), input.Interval)

	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}
	if input.Unit != "" {
		stmt = costsdb.PerUnitEnvironmentForUnit
	}

	if result, err = datastore.Select[[]*costs.Cost](ctx, db, stmt, input); err == nil {
		response.Body.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
		response.Body.Result = result
	}

	return
}

var detailDescription string = `Provides a list of the total costs grouped by a date interval, account_id, environment and service within start (>=) and end (<) dates passed.

Data is limited to the date range (>= start_date < ) and optional unit filter.

Each result returned has the following fields:


	account_id
	unit
	label
	environment
	service
	date
	cost
`

// apiDetailed
func apiDetailed(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	// grab the db from the context
	var dbFilepath string = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var result []*costs.Cost = []*costs.Cost{}
	var bdy *StandardBody = &StandardBody{}
	var stmt = costsdb.Detailed

	bdy.ColumnOrder = []string{"account_id", "unit", "environment", "service", "label"}
	bdy.Request = input
	bdy.Type = "detail"
	bdy.DateRange = convert.DateRange(input.StartTime(), input.EndTime(), input.Interval)

	response = &StandardResult{Body: bdy}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}
	if input.Unit != "" {
		stmt = costsdb.DetailedForUnit
	}

	if result, err = datastore.Select[[]*costs.Cost](ctx, db, stmt, input); err == nil {
		response.Body.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
		response.Body.Result = result
	}

	return
}

// List endpoints the costapi will handle ready for use in the navigation structs
const (
	UriTotal                  endpoints.ApiEndpoint = "/{version}/costs/aws/total/{billing_date:-11}/{billing_date:0}"
	UriMonthlyTax             endpoints.ApiEndpoint = "/{version}/costs/aws/tax-overview/{billing_date:-11}/{billing_date:0}/month"
	UriMonthlyUnit            endpoints.ApiEndpoint = "/{version}/costs/aws/unit/{billing_date:-9}/{billing_date:0}/month"
	UriMonthlyUnitEnvironment endpoints.ApiEndpoint = "/{version}/costs/aws/unit-environment/{billing_date:-9}/{billing_date:0}/month"
	UriMonthlyDetailed        endpoints.ApiEndpoint = "/{version}/costs/aws/detailed/{billing_date:-11}/{billing_date:0}/month"
	UriDailyTax               endpoints.ApiEndpoint = "/{version}/costs/aws/tax-overview/{billing_date:-11}/{billing_date:0}/day"
	UriDailyUnit              endpoints.ApiEndpoint = "/{version}/costs/aws/unit/{billing_date:-11}/{billing_date:0}/day"
	UriDailyUnitEnvironment   endpoints.ApiEndpoint = "/{version}/costs/aws/unit-environment/{billing_date:-11}/{billing_date:0}/day"
	UriDailyDetailed          endpoints.ApiEndpoint = "/{version}/costs/aws/detailed/{billing_date:-11}/{billing_date:0}/day"
)

// Register attaches all the endpoints this module handles on the passed huma api
//
// Currently supports the following endpoints:
//   - /{version}/costs/aws/total
//   - /{version}/costs/aws/tax-overview/{interval}
//   - /{version}/costs/aws/unit/{interval}
//   - /{version}/costs/aws/unit-environment/{internval}
//   - /{version}/costs/aws/detailed/{interval}
func Register(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID:   "get-costs-aws-total",
		Method:        http.MethodGet,
		Path:          "/{version}/costs/aws/total/{start_date}/{end_date}",
		Summary:       "Total",
		Description:   totalDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiTotal)

	huma.Register(api, huma.Operation{
		OperationID:   "get-costs-aws-tax-overview",
		Method:        http.MethodGet,
		Path:          "/{version}/costs/aws/tax-overview/{start_date}/{end_date}/{interval}",
		Summary:       "Tax overview",
		Description:   taxOverviewDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiTaxOverview)

	huma.Register(api, huma.Operation{
		OperationID:   "get-costs-aws-costs-per-unit",
		Method:        http.MethodGet,
		Path:          "/{version}/costs/aws/unit/{start_date}/{end_date}/{interval}",
		Summary:       "Costs per unit",
		Description:   perUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiPerUnit)

	huma.Register(api, huma.Operation{
		OperationID:   "get-costs-aws-costs-per-unit-environment",
		Method:        http.MethodGet,
		Path:          "/{version}/costs/aws/unit-environment/{start_date}/{end_date}/{interval}",
		Summary:       "Costs per unit & environment",
		Description:   perUnitEnvDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiPerUnitEnv)

	huma.Register(api, huma.Operation{
		OperationID:   "get-costs-aws-costs-detailed",
		Method:        http.MethodGet,
		Path:          "/{version}/costs/aws/detailed/{start_date}/{end_date}/{interval}",
		Summary:       "Detailed Costs",
		Description:   detailDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiDetailed)

}
