package aws_costs

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/rbase"
)

type CountValues struct {
	Count int `json:"count"`
}

type Counters struct {
	Totals CountValues `json:"totals"`
	This   CountValues `json:"this"`
}

// -- YTD
type YtdResult struct {
	Total float64 `json:"total"`
}

type YtdResponse struct {
	*rbase.Response
	Counters Counters   `json:"counters"`
	Result   *YtdResult `json:"result"`
}

type MonthlyTaxResponse struct {
	*rbase.Response
	Counters Counters                        `json:"counters"`
	Result   []awsc.MonthlyTotalsTaxSplitRow `json:"result"`
	Columns  map[string][]string             `json:"columns"`
}

func NewMonthlyTaxResponse() *MonthlyTaxResponse {
	return &MonthlyTaxResponse{
		Response: &rbase.Response{
			RequestTimer: &rbase.RequestTimings{},
			DataAge:      &rbase.DataAge{},
			StatusCode:   http.StatusOK,
			Errors:       []string{},
			DateRange:    []string{},
		},
		Columns: map[string][]string{},
		Result:  []awsc.MonthlyTotalsTaxSplitRow{},
	}
}

func NewYTDResponse() *YtdResponse {
	return &YtdResponse{
		Response: &rbase.Response{
			RequestTimer: &rbase.RequestTimings{},
			DataAge:      &rbase.DataAge{},
			StatusCode:   http.StatusOK,
			Errors:       []string{},
		},
		Result: &YtdResult{},
	}
}
