package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/awscosts"
	"github.com/ministryofjustice/opg-reports/datastore"
)

const segment string = "awscosts"

var databaseFilePath *string = nil
var databaseConfig = datastore.Sqlite

// VersionInput is a base input type that we use to as a prefix to version the api - /v1/*
type VersionInput struct {
	Version string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

// StartEndDateInput capture the start and end dates we want to use within an api from the query string.
type StartEndDateInput struct {
	StartDate string `json:"start_date" query:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). Can be YYYY-MM or YYYY-MM-DD." pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
	EndDate   string `json:"end_date" query:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). Can be YYYY-MM or YYYY-MM-DD."  pattern:"([0-9]{4}-[0-9]{2}(-[0-9]{2}){0,1})"`
}

// DateGroupingInput expands on StartEndDateInput by also including the interval period to group the date based data by (month / year)
type DateGroupingInput struct {
	StartEndDateInput
	Interval string `json:"interval" path:"interval" required:"true" doc:"Group the data by the date format" enum:"monthly,daily" default:"monthly"`
}

// BaseBody is base struct that is used by most results
type BaseBody struct {
	Columns     map[string][]interface{} `json:"columns" doc:"List of all the columns and each of their possible values from the dataset. Used for display purposes."`
	ColumnOrder []string                 `json:"column_order" doc:"Ordered list of all of the column names. Used for iterating within a display context to show data correctly."`
}

// TotalWithinDateRangeInput is used for /v1/awscosts/total requests
type TotalWithinDateRangeInput struct {
	VersionInput
	StartEndDateInput
}

// TotalWithinDateRangeResult handles the result for /v1/awscosts/total
type TotalWithinDateRangeResult struct {
	Body struct {
		Request *TotalWithinDateRangeInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
		Result  float64                    `json:"result" doc:"The total sum of all costs as a float without currency." example:"1357.7861"`
	}
}

// TotalWithinDateRangeFunc processes the incoming request for /v1/awscosts/total
func TotalWithinDateRangeFunc(ctx context.Context, input *TotalWithinDateRangeInput) (response *TotalWithinDateRangeResult, err error) {
	var startDate string = input.StartDate
	var endDate string = input.EndDate
	var db *sqlx.DB

	slog.Info("[api.awscosts] TotalWithinDatesInput", slog.String("startDate", startDate), slog.String("endDate", endDate))
	response = &TotalWithinDateRangeResult{}
	response.Body.Request = input

	db, _, err = datastore.New(ctx, databaseConfig, *databaseFilePath)
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

// TotalWithAndWithoutTaxInput handles the inputs for /v1/awscots/tax-overview
type TotalWithAndWithoutTaxInput struct {
	DateGroupingInput
}
type TotalWithAndWithoutTaxResult struct {
	Body struct {
		BaseBody
		Request *TotalWithAndWithoutTaxInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
		Result  []*awscosts.Cost             `json:"result" doc:"List of cost information."`
	}
}

func TotalWithAndWithoutTaxFunc(ctx context.Context, input *TotalWithAndWithoutTaxInput) (response *TotalWithAndWithoutTaxResult, err error) {
	var db *sqlx.DB
	var columns = []string{"service"}
	var params = &awscosts.NamedParameters{
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		DateFormat: databaseConfig.YearMonthFormat,
	}
	slog.Info("[api.awscosts] TotalWithAndWithoutTaxInput", slog.String("startDate", input.StartDate), slog.String("endDate", input.EndDate), slog.String("interval", input.Interval))
	response = &TotalWithAndWithoutTaxResult{}
	response.Body.Request = input
	db, _, err = datastore.New(ctx, databaseConfig, *databaseFilePath)
	defer db.Close()
	if err != nil {
		return
	}

	result, err := awscosts.GetMany(ctx, db, awscosts.TotalsWithAndWithoutTax, params)
	if err == nil {
		response.Body.ColumnOrder = columns
		response.Body.Columns = awscosts.ColumnValues(result, columns)

		response.Body.Result = result
	}
	return
}

func Register(api huma.API, dbFile string) {
	databaseFilePath = &dbFile

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-total",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/total",
		Summary:       "[AWS Costs] Total cost",
		Description:   "Returns a single total for all the AWS costs between the start and end dates - excluding taxes.",
		DefaultStatus: http.StatusOK,
	}, TotalWithinDateRangeFunc)

	huma.Register(api, huma.Operation{
		OperationID:   "get-awscosts-tax-overview",
		Method:        http.MethodGet,
		Path:          "/{version}/awscosts/tax-overview",
		Summary:       "[AWS Costs] Tax overview",
		Description:   "Provides a map of the total costs per interval with and without tax between the start and end date specified.",
		DefaultStatus: http.StatusOK,
	}, TotalWithAndWithoutTaxFunc)

}
