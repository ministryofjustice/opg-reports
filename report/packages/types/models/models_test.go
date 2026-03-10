package models

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
	"testing"
)

// // check meet the response requirements
// var (
// 	_ interfaces.DataResponse[*tRecord, *tRecord, *Request, *Filter] = &APIResponse[*tRecord, *tRecord, *Request, *Filter]{}
// )

// check meets the filterable interfaces
var (
	_ interfaces.Filterable = &Filter{}
)

// check meets the request related interfaces
var (
	_ interfaces.Populator   = &Request{}
	_ interfaces.FilterMaker = &Request{}
	_ interfaces.ApiRequest  = &Request{}
)

type testReq struct {
	Req      *http.Request
	Expected map[string]interface{}
}

type testReqFilter struct {
	Req      *http.Request
	Expected map[string]interface{}
}

type tRecord struct {
	Name string `json:"name"`
}

func (self *tRecord) Sequence() []any {
	return []any{&self.Name}
}

func TestPackagesTypesModelsRequestFilter(t *testing.T) {
	var tests = []*testReqFilter{
		{
			Req: httptest.NewRequest(http.MethodGet, "/test/?foo=bar&date_a=2025-01&date_b=2025-03&team=team-a", nil),
			Expected: map[string]interface{}{
				"months": []interface{}{"2025-01", "2025-03"},
				"team":   "team-a",
			},
		},
		{
			Req: httptest.NewRequest(http.MethodGet, "/test/?foo=bar&date_start=2025-01&date_end=2025-03&team=team-a", nil),
			Expected: map[string]interface{}{
				"months": []interface{}{"2025-01", "2025-02", "2025-03"},
				"team":   "team-a",
			},
		},
	}

	for _, test := range tests {
		r := &Request{}
		m := map[string]interface{}{}
		f := r.Filter(test.Req)

		convert.Between(f, &m)
		// now check values..
		for k, v := range test.Expected {
			// do something different for months
			if k == "months" {
				exMonths := v.([]interface{})
				for i, m := range m["months"].([]interface{}) {
					if exMonths[i] != m {
						t.Errorf("error with month mismatch")
					}
				}
			} else if m[k] != v {
				t.Errorf("[%s] mismatch, expected [%v] actual [%v]", k, v, m[k])
			}
		}
	}

}

func TestPackagesTypesModelsRequestPopulate(t *testing.T) {
	var tests = []*testReq{
		{
			Req: httptest.NewRequest(http.MethodGet, "/test/?date_a=2025-01&date_b=2025-03", nil),
			Expected: map[string]interface{}{
				"date_a": "2025-01",
				"date_b": "2025-03",
			},
		},
		{
			Req: httptest.NewRequest(http.MethodGet, "/test/?foo=bar&date_start=2025-01&date_end=2025-03", nil),
			Expected: map[string]interface{}{
				"date_start": "2025-01",
				"date_end":   "2025-03",
			},
		},
		{
			Req: httptest.NewRequest(http.MethodGet, "/test/?foo=bar&team=test-team", nil),
			Expected: map[string]interface{}{
				"team": "test-team",
			},
		},
	}

	for _, test := range tests {
		r := &Request{}
		m := map[string]interface{}{}
		r.Populate(test.Req)
		convert.Between(r, &m)
		// now check values..
		for k, v := range test.Expected {
			if m[k] != v {
				t.Errorf("[%s] mismatch, expected [%v] actual [%v]", k, v, m[k])
			}
		}

	}

}

// TestPackagesTypesModelsFilter simple test to check map values are setup
func TestPackagesTypesModelsFilter(t *testing.T) {

	var test = &Filter{Team: "team-a", Months: []string{"2026-01", "2026-02"}}
	var mapped = test.Map()

	if v, ok := mapped["team"]; !ok || v != "team-a" {
		t.Errorf("error with filter to map - team")
	}
	if v, ok := mapped["months"]; !ok || len(v.([]interface{})) != 2 {
		t.Errorf("erroter with filter to map - months")
	}

}
