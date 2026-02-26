package accountsapiforteam

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/global/seeds"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
	"path/filepath"
	"testing"
)

func TestAccountForTeamApiAllHandler(t *testing.T) {
	var (
		err    error
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-handler.db")
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
	// setup the server and items
	url := "/v1/accounts/team/team-a/"
	mux := http.NewServeMux()

	req := httptest.NewRequest(http.MethodGet, url, nil)
	writer := httptest.NewRecorder()

	// setup the bindings to the test handler and call
	Register(ctx, mux, &Config{
		Driver: driver,
		DB:     dbpath,
	})
	mux.ServeHTTP(writer, req)

	// get and parse the result
	resp := writer.Result()
	rec := &Response{}
	err = response.As(resp, &rec)
	if err != nil {
		t.Errorf("error converting ...")
	}
	// test returned data
	if len(rec.Data) <= 0 {
		t.Errorf("incorrect number of data rows' this may be due to a seeding being randomised.")
	}
	for _, row := range rec.Data {
		if row.Team != "team-a" {
			t.Errorf("incorrect team returned in data.")
		}
	}

}
