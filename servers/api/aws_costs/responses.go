package aws_costs

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/rbase"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type CountValues struct {
	Count int `json:"count"`
}

type Counters struct {
	Totals *CountValues `json:"totals"`
	This   *CountValues `json:"current"`
}

type PossibleResults interface {
	awsc.MonthlyCostsDetailedRow |
		awsc.MonthlyCostsPerUnitRow |
		awsc.MonthlyCostsPerUnitEnvironmentRow |
		awsc.DailyCostsDetailedRow |
		awsc.DailyCostsPerUnitRow |
		awsc.DailyCostsPerUnitEnvironmentRow |
		awsc.MonthlyTotalsTaxSplitRow
}

// CommonResult used to cast results
type CommonResult struct {
	AccountID   string      `json:"account_id,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	Environment interface{} `json:"environment,omitempty"`
	Service     string      `json:"service,omitempty"`
	Total       interface{} `json:"total,omitempty"`
	Interval    interface{} `json:"interval,omitempty"`
}

// -- Standard
// how to cast result to common type?
type CostResponse struct {
	*rbase.Response
	Counters       *Counters              `json:"counters,omitempty"`
	Columns        map[string][]string    `json:"columns,omitempty"`
	ColumnOrdering []string               `json:"column_ordering,omitempty"`
	QueryFilters   map[string]interface{} `json:"query_filters,omitempty"`
	Result         []*CommonResult        `json:"result"`
}

func Common[T PossibleResults](results []T) (common []*CommonResult) {
	mapList, _ := convert.Maps(results)
	common, _ = convert.Unmaps[*CommonResult](mapList)
	return
}

func NewResponse() *CostResponse {
	return &CostResponse{
		Response: &rbase.Response{
			RequestTimer: &rbase.RequestTimings{},
			DataAge:      &rbase.DataAge{},
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
