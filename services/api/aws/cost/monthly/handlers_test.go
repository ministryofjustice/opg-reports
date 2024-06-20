package monthly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"slices"
	"testing"
	"time"
)

// Index is empty and returns simple api response without a result
// so just check status and errors
func TestServicesApiAwsCostMonthlyHandlerIndex(t *testing.T) {
	fs := testFs()
	mux := testMux()
	store := data.NewStore[*cost.Cost]()
	api := New(store, fs)
	api.Register(mux)

	route := "/aws/costs/v1/monthly/"
	w, r := testWRGet(route)

	mux.ServeHTTP(w, r)

	_, b := strResponse(w.Result())
	res := server.NewApiSimpleResponse()
	json.Unmarshal(b, &res)

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status error")
	}
	if len(res.Errors) != 0 {
		t.Errorf("found error when not expected")
	}
	if res.Times.Duration.String() == "" {
		t.Errorf("duration error")
	}

}

// Generates a series of date in and out of date bounds and then
// triggers the api to get that data.
// Checks the number of items returned matches expectations
func TestServicesApiAwsCostMonthlyHandlerTotals(t *testing.T) {
	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	months := dates.Strings(dates.Months(min, max), dates.FormatYM)
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	l := 500
	x := 100

	for i := 0; i < l; i++ {
		store.Add(cost.Fake(nil, min, max, df))
	}
	for i := 0; i < x; i++ {
		store.Add(cost.Fake(nil, overm, overmx, df))
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	_, b := strResponse(w.Result())
	res := server.NewApiResponseWithResult[*cost.Cost, map[string][]*cost.Cost]()
	json.Unmarshal(b, &res)

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status error")
	}
	if len(res.Errors) != 0 {
		t.Errorf("found error when not expected")
	}
	if res.Times.Duration.String() == "" {
		t.Errorf("duration error")
	}

	if len(res.Result) <= 0 {
		t.Errorf("result not returned correctly")
	}

	total := 0
	for key, list := range res.GetResult() {
		fv := data.FromIdx(key)
		m := fv["month"]
		if !slices.Contains(months, m) {
			t.Errorf("month out of range: %s", m)
		}
		total += len(list)
	}

	if total != l {
		t.Errorf("found extra data!")
	}

}
