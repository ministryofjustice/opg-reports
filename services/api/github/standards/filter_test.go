package standards

import (
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiGithubStandardsFiltersForGetParameters(t *testing.T) {
	logger.LogSetup()
	// --- SETUP
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()
	store := data.NewStore[*std.Repository]()
	fakes := 90
	archived := 5
	teams := 3

	for i := 0; i < fakes; i++ {
		c := std.Fake(nil)
		c.Archived = false
		store.Add(c)
	}
	for i := 0; i < archived; i++ {
		c := std.Fake(nil)
		c.Archived = true
		store.Add(c)
	}
	for i := 0; i < teams; i++ {
		c := std.Fake(nil)
		c.Archived = false
		c.Teams = []string{"ABC", "foo"}
		store.Add(c)
	}

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	// --- TEST ARCHIVED ONLY VALUES
	route := "/github/standards/v1/list/?archived=true"
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err := response.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos := data.FromRows[*std.Repository](res.GetData().GetTableBody())

	if len(repos) != archived {
		t.Errorf("archive filter failed")
	}

	// --- TEST TEAM FILTER OR LOGIC
	route = "/github/standards/v1/list/?teams=ABC&team=bar"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = data.FromRows[*std.Repository](res.GetData().GetTableBody())

	if len(repos) != teams {
		t.Errorf("team filter failed")
	}

	// --- TEST TEAM FILTER AND LOGIC
	// this should return empty set as that team does not exist
	route = "/github/standards/v1/list/?teams=ABC,bar"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = data.FromRows[*std.Repository](res.GetData().GetTableBody())

	if len(repos) != 0 {
		t.Errorf("team AND filter failed")
	}

	// --- TEST TEAM FILTER AND OR COMBINED LOGIC
	// this should return all expected teams due to the foo team
	route = "/github/standards/v1/list/?teams=ABC,bar&teams=foo"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = data.FromRows[*std.Repository](res.GetData().GetTableBody())

	if len(repos) != teams {
		t.Errorf("team AND OR filter failed")
	}

}
