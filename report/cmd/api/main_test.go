package main

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/global/seeds"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"path/filepath"
	"testing"
)

// TestAPIEndpointsRespond sets up the server and then makes sure all endpoints
// get a 200 response
func TestAPIEndpointsRespond(t *testing.T) {

	var (
		err    error
		mux    *http.ServeMux
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-teams-handler.db")
	)

	// run seeds
	_, err = seeds.SeedAll(ctx, &seeds.Args{
		Driver: driver,
		DB:     dbpath,
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}

	mux = http.NewServeMux()

	// register all endpoints
	registerEndpoints(ctx, mux, &cli{
		Driver: driver,
		DB:     dbpath,
	})

	// all endpoints
	endpoints := []string{
		"/",
		"/ping/",
		"/v1/teams/",
		"/v1/accounts/",
		"/v1/accounts/team/team-a/",
		"/v1/costs/teams/between/2026-01/2026-02/",
		"/v1/costs/accounts/between/2026-01/2026-02/team/team-a/",
	}
	for _, url := range endpoints {
		writer := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)

		mux.ServeHTTP(writer, req)
		resp := writer.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("failed to get [%s]; received error code [%d]", url, resp.StatusCode)
		}

	}

}
