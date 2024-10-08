package aws_costs

import (
	"context"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

// PossibleResults is used to constrain the type of the value on the Common func
// and simply is interface for all the know result types
type PossibleResults interface {
	awsc.MonthlyCostsDetailedRow |
		awsc.MonthlyCostsPerUnitRow |
		awsc.MonthlyCostsPerUnitEnvironmentRow |
		awsc.DailyCostsDetailedRow |
		awsc.DailyCostsPerUnitRow |
		awsc.DailyCostsPerUnitEnvironmentRow |
		awsc.MonthlyTotalsTaxSplitRow |
		awsc.MonthlyCostsDetailedForUnitRow |
		awsc.DailyCostsDetailedForUnitRow
}

// CountValues tracks a single counter value
// This is normally used to track the total of the data sets
type CountValues struct {
	Count int `json:"count"`
}

// Counters captures multiple count values, in general we
// return the Total (so everyhing in the database) version
// and `This` - which is based on the current query result
type Counters struct {
	Totals *CountValues `json:"totals"`
	This   *CountValues `json:"current"`
}

// CommonResult is used instead of the variable versions encapsulated
// by PossibleResults in the ApiResponse struct
// This is to simplify the parsing on both the api and the consumer
// in front
// To be effective, any empty field is omited in the json
// Converted using `Common` func
type CommonResult struct {
	AccountID   string      `json:"account_id,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	Label       string      `json:"label,omitempty"`
	Environment interface{} `json:"environment,omitempty"`
	Service     string      `json:"service,omitempty"`
	Total       interface{} `json:"total,omitempty"`
	Interval    interface{} `json:"interval,omitempty"`
}

// ApiResponse is the response object used and returned by the aws_costs
// api handler
// Based on response.Response struct as a common ground and then
// add additional fields to the struct that are used for this api
type ApiResponse struct {
	*response.Response
	Counters       *Counters              `json:"counters,omitempty"`
	Columns        map[string][]string    `json:"columns,omitempty"`
	ColumnOrdering []string               `json:"column_ordering,omitempty"`
	QueryFilters   map[string]interface{} `json:"query_filters,omitempty"`
	Result         []*CommonResult        `json:"result"`
}

// Common func converts from the known aws cost structs to the common result
// type via json marshaling
func Common[T PossibleResults](results []T) (common []*CommonResult) {
	mapList, _ := convert.Maps(results)
	common, _ = convert.Unmaps[*CommonResult](mapList)
	return
}

// setIfFound
func setIfFound(r *CommonResult, columns map[string]map[string]bool) {
	mapped, _ := convert.Map(r)

	for k, v := range mapped {
		if k != "total" && k != "interval" {
			if _, ok := columns[k]; !ok {
				columns[k] = map[string]bool{}
			}
			columns[k][v.(string)] = true
		}
	}
}

// ColumnPermutations uses the result set to create a list of columns and
// all of their possible values
// This is normally used to create table headers and the like
func ColumnPermutations(results []*CommonResult) map[string][]string {
	columns := map[string]map[string]bool{}
	colList := map[string][]string{}

	for _, r := range results {
		setIfFound(r, columns)
	}
	for col, values := range columns {
		colList[col] = []string{}
		for choice, _ := range values {
			colList[col] = append(colList[col], choice)
		}
	}
	return colList
}

// StandardCounters adds the standard counter data
func StandardCounters(ctx context.Context, q *awsc.Queries, resp *ApiResponse) {

	all, _ := q.Count(ctx)
	min, _ := q.Oldest(ctx)
	max, _ := q.Youngest(ctx)

	resp.Counters = &Counters{
		Totals: &CountValues{Count: int(all)},
		This:   &CountValues{Count: len(resp.Result)},
	}
	resp.DataAge = &response.DataAge{Min: min, Max: max}
}

func StandardDates(response *ApiResponse, start time.Time, end time.Time, rangeEnd time.Time, interval dates.Interval) {
	df := dates.IntervalFormat(interval)
	response.StartDate = start.Format(df)
	response.EndDate = end.Format(df)
	response.DateRange = dates.Strings(dates.Range(start, rangeEnd, interval), df)
}

// NewResponse returns a clean response object
func NewResponse() *ApiResponse {
	resp := &response.Response{
		RequestTimer: &response.RequestTimings{},
		DataAge:      &response.DataAge{},
		StatusCode:   http.StatusOK,
		Errors:       []string{},
		DateRange:    []string{},
	}
	return &ApiResponse{
		Response: resp,
		Counters: &Counters{
			This:   &CountValues{},
			Totals: &CountValues{},
		},
		Columns:        map[string][]string{},
		ColumnOrdering: []string{},
		Result:         []*CommonResult{},
	}
}
