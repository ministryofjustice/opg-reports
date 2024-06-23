package monthly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
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
	services := []string{"ec2", "ecs", "tax", "rds", "r53"}
	l := 900
	x := 100

	for i := 0; i < l; i++ {
		c := cost.Fake(nil, min, max, df)
		c.Service = fake.Choice(services)
		store.Add(c)
	}
	for i := 0; i < x; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.Service = fake.Choice(services)
		store.Add(c)
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	_, b := strResponse(w.Result())
	res := server.NewApiResponseWithResult[*cost.Cost, map[string]map[string][]*cost.Cost]()
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
	for _, dataSet := range res.GetResult() {
		for key, list := range dataSet {

			fv := data.FromIdx(key)
			m := fv["month"]
			if !slices.Contains(months, m) {
				t.Errorf("month out of range: %s", m)
			}
			total += len(list)
		}
	}
	// as items may appear in both top level segements if they are
	// not tax, the total should always be at least l, but can be higher
	if total <= l {
		t.Errorf("data mis match")
	}

}

func TestServicesApiAwsCostMonthlyHandlerUnits(t *testing.T) {
	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	months := dates.Strings(dates.Months(min, max), dates.FormatYM)
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	units := []string{"teamOne", "teamTwo", "teamThree"}
	l := 90
	x := 5

	for i := 0; i < l; i++ {
		c := cost.Fake(nil, min, max, df)
		c.AccountUnit = fake.Choice(units)
		store.Add(c)
	}
	for i := 0; i < x; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.AccountUnit = fake.Choice(units)
		store.Add(c)
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	_, b := strResponse(w.Result())
	res := server.NewApiResponseWithResult[*cost.Cost, map[string][]*cost.Cost]()
	json.Unmarshal(b, &res)

	total := 0
	for key, list := range res.GetResult() {
		fv := data.FromIdx(key)
		m := fv["month"]
		if !slices.Contains(months, m) {
			t.Errorf("month out of range: %s", m)
		}
		total += len(list)
	}
	// as items may appear in both top level segements if they are
	// not tax, the total should always be at least l, but can be higher
	if total != l {
		t.Errorf("data mis match")
	}

}

func TestServicesApiAwsCostMonthlyHandlerUnitEnvs(t *testing.T) {
	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	months := dates.Strings(dates.Months(min, max), dates.FormatYM)
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	units := []string{"teamOne", "teamTwo", "teamThree"}
	envs := []string{"dev", "preprod", "prod"}
	l := 9
	x := 5

	for i := 0; i < l; i++ {
		c := cost.Fake(nil, min, max, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
	}
	for i := 0; i < x; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	str, b := strResponse(w.Result())
	fmt.Println(str)
	res := server.NewApiResponseWithResult[*cost.Cost, map[string][]*cost.Cost]()
	json.Unmarshal(b, &res)

	total := 0
	for key, list := range res.GetResult() {
		fv := data.FromIdx(key)
		m := fv["month"]
		if !slices.Contains(months, m) {
			t.Errorf("month out of range: %s", m)
		}
		total += len(list)
	}
	// as items may appear in both top level segements if they are
	// not tax, the total should always be at least l, but can be higher
	if total != l {
		t.Errorf("data mis match")
	}

}

func TestServicesApiAwsCostMonthlyHandlerUnitEnvServices(t *testing.T) {
	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	months := dates.Strings(dates.Months(min, max), dates.FormatYM)
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	units := []string{"teamOne", "teamTwo", "teamThree"}
	envs := []string{"dev", "preprod", "prod"}
	l := 9
	x := 5

	for i := 0; i < l; i++ {
		c := cost.Fake(nil, min, max, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
	}
	for i := 0; i < x; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/services/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	str, b := strResponse(w.Result())
	fmt.Println(str)
	res := server.NewApiResponseWithResult[*cost.Cost, map[string][]*cost.Cost]()
	json.Unmarshal(b, &res)

	total := 0
	for key, list := range res.GetResult() {
		fv := data.FromIdx(key)
		m := fv["month"]
		if !slices.Contains(months, m) {
			t.Errorf("month out of range: %s", m)
		}
		total += len(list)
	}
	// as items may appear in both top level segements if they are
	// not tax, the total should always be at least l, but can be higher
	if total != l {
		t.Errorf("data mis match")
	}

}
