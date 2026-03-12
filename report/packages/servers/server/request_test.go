package server

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types"
	"slices"
	"testing"
)

// interface checks
var (
	_ types.Parameters         = &parameters{}
	_ types.HttpRequester      = &Request{}
	_ types.ParameterRequester = &Request{}
	_ types.ServerRequester    = &Request{}
)

type tScenario struct {
	Url      string            // the source url used with the request
	Expected map[string]string // expected processed values from the url

}

func TestPackagesServersServerRequest(t *testing.T) {

	var tests = []*tScenario{
		// should return the first value of date_start
		{
			Url:      "/test/?date_start=2026-01&date_start=2025-01",
			Expected: map[string]string{"date_start": "2026-01"},
		},
		// should return multiple fields but ignore the foo param
		{
			Url:      "/test/?date_end=2026-01&date_start=2025-01&foo=bar",
			Expected: map[string]string{"date_start": "2025-01", "date_end": "2026-01"},
		},
	}

	for i, test := range tests {
		var r = &Request{}
		req := httptest.NewRequest(http.MethodGet, test.Url, nil)
		// set the request
		r.SetRequest(req)
		// get the data
		data := r.Parameters().Data()
		for k, actual := range data {
			var expected = test.Expected[k]
			if expected != actual {
				t.Errorf("[%d] expected [%s] to be [%v] but actual [%v]", i, k, expected, actual)
			}
		}
	}

}

func TestPackagesServersServerParameterKeys(t *testing.T) {
	var p = &parameters{}
	var compare = map[string]string{}
	var keys = p.Keys()
	// create a map via json marshaling to compare against the
	// keys list generated
	convert.Between(p, &compare)

	for k, _ := range compare {
		if !slices.Contains(keys, k) {
			t.Errorf("key list is missing an expected field: [%s]", k)
		}
	}

}
