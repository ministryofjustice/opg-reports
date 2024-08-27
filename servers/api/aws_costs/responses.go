package aws_costs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apiresponse"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

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

// CommonResult is used instead of the variable versions encapsulated
// by PossibleResults in the CostResponse struct
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

// CostResponse is the response object used and returned by the aws_costs
// api handler
// Based on apiresponse.Response struct as a common ground and then
// add additional fields to the struct that are used for this api
type CostResponse struct {
	*apiresponse.Response
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
func StandardCounters(ctx context.Context, q *awsc.Queries, resp *CostResponse) {

	all, _ := q.Count(ctx)
	min, _ := q.Oldest(ctx)
	max, _ := q.Youngest(ctx)

	resp.Counters = &Counters{
		Totals: &CountValues{Count: int(all)},
		This:   &CountValues{Count: len(resp.Result)},
	}
	resp.DataAge = &apiresponse.DataAge{Min: min, Max: max}
}

func StandardDates(response *CostResponse, start time.Time, end time.Time, rangeEnd time.Time, interval dates.Interval) {
	df := dates.IntervalFormat(interval)
	response.StartDate = start.Format(df)
	response.EndDate = end.Format(df)
	response.DateRange = dates.Strings(dates.Range(start, rangeEnd, interval), df)
}

// --- all the functions that call the correct queries
type queryWrapperF func(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error)

func monthlyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetailsForUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsDetailedForUnitParams{Start: req.Start, End: req.End, Unit: req.Unit}
	if res, err := q.MonthlyCostsDetailedForUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetailsForUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsDetailedForUnitParams{Start: req.Start, End: req.End, Unit: req.Unit}
	if res, err := q.DailyCostsDetailedForUnit(ctx, params); err == nil {
		results = Common(res)
	}
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

// StandardQueryResults uses the group and interval values from the request to determine
// which db query to run
//
// The query to use is determined by the following:
// Interval set to `MONTH`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Interval set to `DAY`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Query results are converted to `[]*CommonResult` struct.
func StandardQueryResults(ctx context.Context, q *awsc.Queries, req *ApiRequest) (results []*CommonResult, err error) {
	// define the possible functions and params
	var possibleFuncs = map[string]queryWrapperF{
		string(dates.MONTH) + string(GroupByUnit):            monthlyPerUnit,
		string(dates.MONTH) + string(GroupByUnitEnvironment): monthlyPerUnitEnv,
		string(dates.MONTH) + string(GroupByDetailed):        monthlyDetails,
		//
		string(dates.DAY) + string(GroupByUnit):            dailyPerUnit,
		string(dates.DAY) + string(GroupByUnitEnvironment): dailyPerUnitEnv,
		string(dates.DAY) + string(GroupByDetailed):        dailyDetails,
		//
	}

	key := req.Interval + req.GroupBy
	// if a filter is set, that over rides the default grouping setup
	if req.Unit != "" && req.IntervalT == dates.MONTH {
		results, err = monthlyDetailsForUnit(ctx, req, q)
	} else if req.Unit != "" && req.IntervalT == dates.DAY {
		results, err = dailyDetailsForUnit(ctx, req, q)
	} else if f, ok := possibleFuncs[key]; ok {
		results, err = f(ctx, req, q)
	} else {
		err = fmt.Errorf("error finding query function based on get paremters")
	}
	return
}

// NewResponse returns a clean response object
func NewResponse() *CostResponse {
	return &CostResponse{
		Response: &apiresponse.Response{
			RequestTimer: &apiresponse.RequestTimings{},
			DataAge:      &apiresponse.DataAge{},
			StatusCode:   http.StatusOK,
			Errors:       []string{},
			DateRange:    []string{},
		},
		Counters: &Counters{
			This:   &CountValues{},
			Totals: &CountValues{},
		},
		Columns:        map[string][]string{},
		ColumnOrdering: []string{},
		Result:         []*CommonResult{},
	}
}
