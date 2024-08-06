package standards

import (
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiGithubStandardsFiltersForGetParameters(t *testing.T) {
	logger.LogSetup()
	// --- SETUP
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

	Register(mux, store)
	// --- TEST ARCHIVED ONLY VALUES
	route := "/github/standards/v1/list/?archived=true"
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := resp.New()
	err := resp.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}
	repos := resp.ToEntries[*std.Repository](res.Result.Body)

	if len(repos) != archived {
		t.Errorf("archive filter failed: expected [%v] actual [%v]", archived, len(repos))
	}

	// --- TEST TEAM FILTER OR LOGIC
	route = "/github/standards/v1/list/?team=ABC&team=bar"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = resp.New()
	err = resp.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status code failed: %v", res.StatusCode)
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = resp.ToEntries[*std.Repository](res.Result.Body)

	if len(repos) != teams {
		t.Errorf("team filter failed: actual [%v] expected [%v]", len(repos), teams)
	}

	// --- TEST TEAM FILTER AND LOGIC
	// this should return empty set as that team does not exist
	route = "/github/standards/v1/list/?team=ABC,bar"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = resp.New()
	err = resp.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status code failed: %v", res.StatusCode)
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = resp.ToEntries[*std.Repository](res.Result.Body)

	if len(repos) != 0 {
		t.Errorf("team AND filter failed")
	}

	// --- TEST TEAM FILTER AND OR COMBINED LOGIC
	// this should return all expected teams due to the foo team
	route = "/github/standards/v1/list/?team=ABC,bar&team=foo"
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b = response.Stringify(w.Result())
	res = resp.New()
	err = resp.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status code failed: %v", res.StatusCode)
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos = resp.ToEntries[*std.Repository](res.Result.Body)

	if len(repos) != teams {
		t.Errorf("team AND OR filter failed")
	}

}
