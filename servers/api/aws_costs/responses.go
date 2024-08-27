package aws_costs

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apiresponse"
	"github.com/ministryofjustice/opg-reports/shared/convert"
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
		awsc.MonthlyTotalsTaxSplitRow
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
