package monthly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/response"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Index is empty and returns simple api response without a result
// so just check status and errors
func TestServicesApiAwsCostMonthlyHandlerIndex(t *testing.T) {
	logger.LogSetup()
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	store := data.NewStore[*cost.Cost]()
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)

	api.Register(mux)

	route := "/aws/costs/v1/monthly/"
	w, r := testhelpers.WRGet(route)

	mux.ServeHTTP(w, r)

	_, b := response.Stringify(w.Result())
	res := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	json.Unmarshal(b, &res)

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status error")
	}
	if len(res.GetError()) != 0 {
		t.Errorf("found error when not expected")
	}

	if res.GetDuration().String() == "" {
		t.Errorf("duration error")
	}

	costs := data.FromRows[*cost.Cost](res.GetData().GetTableBody())
	if len(costs) != 0 {
		t.Errorf("unexpected data returned from empty data store")
	}
}

// Generates a series of date in and out of date bounds and then
// triggers the api to get that data.
// Checks the number of items returned matches expectations
// URLS:
//   - /aws/costs/v1/monthly/%s/%s/
func TestServicesApiAwsCostMonthlyHandlerTotals(t *testing.T) {
	logger.LogSetup()
	p := message.NewPrinter(language.English)
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	services := []string{"ec2", "ecs", "tax", "rds", "r53"}

	inrange := 100
	outofrange := 30
	inrangeunit := 20

	// within range and limited to know service
	allCost := []*cost.Cost{}
	for i := 0; i < inrange; i++ {
		c := cost.Fake(nil, min, max, df)
		c.Service = fake.Choice(services)
		store.Add(c)
		allCost = append(allCost, c)
	}
	// out of range, so dont add to cost totals
	for i := 0; i < outofrange; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.Service = fake.Choice(services)
		store.Add(c)

	}
	// in range with a unit to be filtered on
	unitFoo := []*cost.Cost{}
	for i := 0; i < inrangeunit; i++ {
		c := cost.Fake(nil, min, max, df)
		c.Service = fake.Choice(services)
		c.AccountUnit = "foobar"
		store.Add(c)
		unitFoo = append(unitFoo, c)
		allCost = append(allCost, c)
	}
	unitTotal := cost.Total(unitFoo)
	allTotal := cost.Total(allCost)

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	// --- TEST WITH FILTER
	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/?unit=foobar", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err := response.FromJson(b, res)

	if err != nil {
		fmt.Println(str)
		t.Errorf("error: %+v", err)
	}

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	apiTotal := 0.0
	for _, row := range res.GetData().GetTableBody() {
		if row.HeaderCells[0].Name == "Included" {
			apiTotal = row.SupplementaryCells[0].Value.(float64)
		}
	}
	// round the values down
	apiTotalS := p.Sprintf("%.4f", apiTotal)
	unitTotalS := p.Sprintf("%.4f", unitTotal)
	if apiTotalS != unitTotalS {
		t.Errorf("filtering by unit failed: expected [%v] actual [%v]", unitTotal, apiTotal)
	}

	// --- TEST WITHOUT FILTER
	route = fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromJson(b, res)

	if err != nil {
		fmt.Println(str)
		t.Errorf("error: %+v", err)
	}

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	apiTotal = 0.0
	for _, row := range res.GetData().GetTableBody() {
		if row.HeaderCells[0].Name == "Included" {
			apiTotal = row.SupplementaryCells[0].Value.(float64)
		}
	}
	// round the values down
	apiTotalS = p.Sprintf("%.4f", apiTotal)
	allTotalS := p.Sprintf("%.4f", allTotal)
	if apiTotalS != allTotalS {
		t.Errorf("filtering by unit failed: expected [%v] actual [%v]", allTotal, apiTotal)
	}

}

// URLS:
//   - /aws/costs/v1/monthly/%s/%s/units/
func TestServicesApiAwsCostMonthlyHandlerUnits(t *testing.T) {
	logger.LogSetup()
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	units := []string{"teamOne", "teamTwo", "teamThree"}
	l := 900
	x := 50

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

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	response.FromJson(b, res)
	// fmt.Println(str)

	if resp.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

}

// URLS:
//   - /aws/costs/v1/monthly/%s/%s/units/envs
func TestServicesApiAwsCostMonthlyHandlerUnitEnvs(t *testing.T) {
	logger.LogSetup()
	pr := message.NewPrinter(language.English)
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)
	store := data.NewStore[*cost.Cost]()
	units := []string{"teamOne", "teamTwo", "teamThree"}
	envs := []string{"dev", "preprod", "prod"}
	l := 9
	x := 5

	// force a prod version in
	p := cost.Fake(nil, min, max, df)
	p.AccountUnit = fake.Choice(units)
	p.AccountEnvironment = "prod"
	prods := []*cost.Cost{p}
	store.Add(p)
	for i := 0; i < (l - 1); i++ {
		c := cost.Fake(nil, min, max, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
		if c.AccountEnvironment == "prod" {
			prods = append(prods, c)
		}
	}
	for i := 0; i < x; i++ {
		c := cost.Fake(nil, overm, overmx, df)
		c.AccountUnit = fake.Choice(units)
		c.AccountEnvironment = fake.Choice(envs)
		store.Add(c)
	}

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	// -- TEST WITH FILTER

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/?environment=prod", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err := response.FromJson(b, res)
	if err != nil {
		t.Errorf("error parsing data: %v", err)
	}

	if resp.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}
	prodTotal := cost.Total(prods)
	// multiple rows can have prod info due to splitting by unit as well
	apiTotal := 0.0
	for _, row := range res.GetData().GetTableBody() {
		if row.HeaderCells[1].Name == "prod" {
			apiTotal += row.SupplementaryCells[0].Value.(float64)
		}
	}
	// round the values down
	apiTotalS := pr.Sprintf("%.4f", apiTotal)
	prodTotalS := pr.Sprintf("%.4f", prodTotal)
	if apiTotalS != prodTotalS {
		fmt.Println(str)
		t.Errorf("filtering by unit failed: expected [%v] actual [%v]", prodTotal, apiTotal)
	}

	// -- TEST WITHOUT FILTER
	route = fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromJson(b, res)
	if err != nil {
		t.Errorf("error parsing data: %v", err)
	}
	if resp.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

}

func TestServicesApiAwsCostMonthlyHandlerUnitEnvServices(t *testing.T) {
	logger.LogSetup()
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
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

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/services/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	response.FromJson(b, res)

	if resp.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

}
