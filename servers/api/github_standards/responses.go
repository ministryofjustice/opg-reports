package github_standards

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/rbase"
)

type CountValues struct {
	Count             int `json:"count"`
	CompliantBaseline int `json:"compliant_baseline"`
	CompliantExtended int `json:"compliant_extended"`
}

type Counters struct {
	Totals *CountValues `json:"totals"`
	This   *CountValues `json:"current"`
}

// -- Standard
// how to cast result to common type?
type GHSResponse struct {
	*rbase.Response
	Counters     *Counters            `json:"counters,omitempty"`
	QueryFilters map[string]string    `json:"query_filters,omitempty"`
	Result       []ghs.GithubStandard `json:"result"`
}

func NewResponse() *GHSResponse {
	return &GHSResponse{
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
		Result: []ghs.GithubStandard{},
	}
}
