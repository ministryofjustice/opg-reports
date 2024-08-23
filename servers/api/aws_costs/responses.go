package aws_costs

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/rbase"
)

type CountValues struct {
	Count int `json:"count"`
}

type Counters struct {
	Totals CountValues `json:"totals"`
	This   CountValues `json:"this"`
}
type YtdResult struct {
	Total float64 `json:"total"`
}

type YtdResponse struct {
	*rbase.Response
	StartDate string     `json:"start_date"`
	EndDate   string     `json:"end_date"`
	DateRange []string   `json:"date_range"`
	Counters  Counters   `json:"counters"`
	Result    *YtdResult `json:"result"`
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
