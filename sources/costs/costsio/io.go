// Package costsio contains all of the structs used for both api input and output
package costsio

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/costs"
)

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
	StartDate string `json:"start_date" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
}

// GroupedDateRangeInput expands upon DateRangeStatement to include
// interval field thats used for formatting the date column as part
// of the grouping
type GroupedDateRangeInput struct {
	DateRangeInput
	Interval   string `json:"interval" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *GroupedDateRangeInput) StartTime() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *GroupedDateRangeInput) EndTime() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
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

// FilterByUnitInput extends GroupedDateRangeInput adding the
// ability to filter the dataset by unit
type FilterByUnitInput struct {
	GroupedDateRangeInput
	Unit string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`
}

type StandardInput struct {
	FilterByUnitInput
}

type TotalInput struct {
	DateRangeInput
}

type TaxOverviewInput struct {
	GroupedDateRangeInput
}

// -- outputs

type CostsStandardBody struct {
	Type         string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request      *StandardInput                    `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result       []*costs.Cost                     `json:"result" doc:"List of all costs found matching query."`
	ColumnOrder  []string                          `json:"column_order" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	DateRange    []string                          `json:"date_range" doc:"list of string dates between the start and end date"`
	TableRows    map[string]map[string]interface{} `json:"-"`
}

type CostsStandardResult struct {
	Body *CostsStandardBody
}

type CostsTotalBody struct {
	Type    string      `json:"type" doc:"States what type of data this is for front end handling"`
	Request *TotalInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result  float64     `json:"result" doc:"The total sum of all costs as a float without currency." example:"1357.7861"`
}
type CostsTotalResult struct {
	Body *CostsTotalBody
}

type CostsTaxOverviewBody struct {
	Type         string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request      *TaxOverviewInput                 `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result       []*costs.Cost                     `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	ColumnOrder  []string                          `json:"column_order" doc:"List of columns set in the order they should be rendered for each row"`
	ColumnValues map[string][]interface{}          `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	DateRange    []string                          `json:"date_range" doc:"list of string dates between the start and end date"`
	TableRows    map[string]map[string]interface{} `json:"-"`
}
type CostsTaxOverviewResult struct {
	Body *CostsTaxOverviewBody
}
