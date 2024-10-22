package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/datastore/awscosts"
)

// VersionInput is a base input type that we use to as a prefix to version the api - /v1/*
type VersionInput struct {
	Version string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

// StartEndDateInput capture the start and end dates we want to use within an api from the query string.
type StartEndDateInput struct {
	StartDate string `json:"start_date" query:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). Can be YYYY-MM or YYYY-MM-DD." example:"2023-12-01" pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
	EndDate   string `json:"end_date" query:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). Can be YYYY-MM or YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
}

// DateGroupingInput expands on StartEndDateInput by also including the interval period to group the date based data by (month / year)
type DateGroupingInput struct {
	StartEndDateInput
	Interval string `json:"interval" path:"interval" required:"true" doc:"Group the data by the date format" enum:"yearly,monthly,daily" default:"monthly"`
}

// BaseBody is base struct that is used by most results
type BaseBody struct {
	Type        string                   `json:"type" doc:"States what type of data this is for front end handling"`
	Columns     map[string][]interface{} `json:"columns" doc:"List of all the columns and each of their possible values from the dataset. Used for display purposes."`
	ColumnOrder []string                 `json:"column_order" doc:"Ordered list of all of the column names. Used for iterating within a display context to show data correctly."`
}

// TotalInput is used for /v1/awscosts/total requests
type TotalInput struct {
	VersionInput
	StartEndDateInput
}

// TotalResult handles the result for /v1/awscosts/total
type TotalResult struct {
	Body struct {
		Type    string      `json:"type" doc:"States what type of data this is for front end handling"`
		Request *TotalInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
		Result  float64     `json:"result" doc:"The total sum of all costs as a float without currency." example:"1357.7861"`
	}
}

// Total processes the incoming request for /v1/awscosts/total
// Returns a singular float value for the sum of all costs between the date periods passed
func Total(ctx context.Context, input *TotalInput) (response *TotalResult, err error) {
	var databaseFilePath = ctx.Value(Segment).(string)
	var startDate string = input.StartDate
	var endDate string = input.EndDate
	var db *sqlx.DB

	slog.Info("[api.awscosts] TotalInput", slog.String("startDate", startDate), slog.String("endDate", endDate))
	response = &TotalResult{}
	response.Body.Request = input
	response.Body.Type = "TotalWithinDateRange"

	db, _, err = datastore.New(ctx, databaseConfig, databaseFilePath)
	defer db.Close()
	if err != nil {
		return
	}

	result, err := awscosts.GetOne(ctx, db, awscosts.TotalInDateRange, startDate, endDate)
	if err == nil {
		response.Body.Result = result.(float64)
	}

	return
}

// TaxOverviewInput handles the inputs for /v1/awscosts/tax-overview
type TaxOverviewInput struct {
	DateGroupingInput
}

// TaxOverviewResult captures the result data for /v1/awscosts/tax-overview
type TaxOverviewResult struct {
	Body struct {
		BaseBody
		Request *TaxOverviewInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
		Result  []*awscosts.Cost  `json:"result" doc:"List of cost information."`
	}
}

// TaxOverview processes the /v1/awscosts/tax-overview request
// Returns a list of costs which are grouped by YYYY-MM and if the sum includes Tax or not
// Used to provide a monhtly totals with & without tax for comparison
func TaxOverview(ctx context.Context, input *TaxOverviewInput) (response *TaxOverviewResult, err error) {
	var databaseFilePath = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var columns = []string{"service"}
	var params = &awscosts.NamedParameters{
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		DateFormat: databaseConfig.YearMonthFormat,
	}
	slog.Info("[api.awscosts] TotalWithAndWithoutTaxInput", slog.String("startDate", input.StartDate), slog.String("endDate", input.EndDate), slog.String("interval", input.Interval))
	response = &TaxOverviewResult{}
	response.Body.Request = input
	response.Body.Type = "TotalWithAndWithoutTax"

	db, _, err = datastore.New(ctx, databaseConfig, databaseFilePath)
	defer db.Close()
	if err != nil {
		return
	}

	result, err := awscosts.GetMany(ctx, db, awscosts.TotalsWithAndWithoutTax, params)
	if err == nil {
		response.Body.ColumnOrder = columns
		response.Body.Columns = datastore.ColumnValues(result, columns)

		response.Body.Result = result
	}
	return
}

// StandardInput handles the inputs for /v1/awscosts/all/{interval}/{grouping}
type StandardInput struct {
	DateGroupingInput
	Grouping string `json:"grouping" path:"grouping" enum:"unit,unit-environment,detailed" default:"unit" doc:"Determines how to group the data in more granular form."`
	Unit     string `json:"unit" query:"unit" doc:"Optional field to filter the data by the value of the unit column."`
}

// StandardResult captures the result data for /v1/awscosts/all/{interval}/{grouping}
type StandardResult struct {
	Body struct {
		BaseBody
		Request *StandardInput   `json:"request" doc:"The public parameters originaly specified in the request to this API."`
		Result  []*awscosts.Cost `json:"result" doc:"List of cost information grouped by interval and unit."`
	}
}

// Standard processes the /v1/awscosts/unit/ request
func Standard(ctx context.Context, input *StandardInput) (response *StandardResult, err error) {
	var databaseFilePath = ctx.Value(Segment).(string)
	var db *sqlx.DB
	var allColumns = map[string][]string{
		"unit":             {"unit"},
		"unit-environment": {"unit", "environment"},
		"detailed":         {"account_id", "unit", "environment", "service"},
	}
	var params = &awscosts.NamedParameters{
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		DateFormat: databaseConfig.YearMonthFormat,
		Unit:       input.Unit,
	}
	var columns = allColumns[input.Grouping]
	var stmt awscosts.ManyStatement

	slog.Info("[api.awscosts] StandardInput",
		slog.String("startDate", input.StartDate),
		slog.String("endDate", input.EndDate),
		slog.String("interval", input.Interval),
		slog.String("unit", input.Unit),
		slog.String("grouping", input.Grouping))

	response = &StandardResult{}

	db, _, err = datastore.New(ctx, databaseConfig, databaseFilePath)
	defer db.Close()
	if err != nil {
		return
	}
	// data interval string format
	if input.Interval == "daily" {
		params.DateFormat = databaseConfig.YearMonthDayFormat
	} else if input.Interval == "yearly" {
		params.DateFormat = databaseConfig.YearFormat
	}

	// work out the sql statment to use
	if input.Grouping == "unit" && input.Unit != "" {
		stmt = awscosts.PerUnitForUnit
	} else if input.Grouping == "unit" {
		stmt = awscosts.PerUnit
	} else if input.Grouping == "unit-environment" && input.Unit != "" {
		stmt = awscosts.PerUnitEnvironmentForUnit
	} else if input.Grouping == "unit-environment" {
		stmt = awscosts.PerUnitEnvironment
	} else if input.Grouping == "detailed" && input.Unit != "" {
		stmt = awscosts.DetailedForUnit
	} else if input.Grouping == "detailed" {
		stmt = awscosts.Detailed
	} else {
		err = huma.Error400BadRequest("requested grouping (" + input.Grouping + ") is invalid.")
		return
	}

	result, err := awscosts.GetMany(ctx, db, stmt, params)
	if err == nil {
		response.Body.ColumnOrder = columns
		response.Body.Columns = datastore.ColumnValues(result, columns)
		response.Body.Result = result
	}
	return
}

// Register handles setting up all routes for awscosts
func Register(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-total",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/total",
		Summary:       "Total",
		Description:   "Returns a single total for all the AWS costs between the start and end dates - excluding taxes.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}, Total)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-tax-overview",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/tax-overview/{interval}/",
		Summary:       "Tax overview",
		Description:   "Provides a list of the total costs per interval with and without tax between the start and end date specified.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}, TaxOverview)

	description := `Provides a list of the total costs grouped by a date interval and the selected grouping within start (>=) and end (<) dates passed.

**Grouping details:**
* *unit*: grouped by the value of the 'unit' field
* *unit-environment*: grouped by 'unit' and 'environment' fields values
* *detailed*: grouped by 'unit', 'environment', 'organisation', 'account_id', and 'service' field values

Allows for optional filtering on unit field.`

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-all",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/all/{interval}/{grouping}",
		Summary:       "Costs (grouped by month or day)",
		Description:   description,
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}, Standard)

}
