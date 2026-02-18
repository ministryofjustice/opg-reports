package accountapiall

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

func TestAccountApiAllHandler(t *testing.T) {
	var (
		err    error
		ctx           = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir           = t.TempDir()
		driver        = "sqlite3"
		dbpath        = filepath.Join(dir, "test-account-handler.db")
		mfile  string = filepath.Join(dir, "migrate.json")
	)

	// run seeds
	_, err = seeds.SeedAll(ctx, &seeds.Args{
		Driver:        driver,
		DB:            dbpath,
		MigrationFile: mfile,
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}
	// setup the server and items
	url := "/v1/accounts/"
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
	err = response.AsT(resp, &rec)
	if err != nil {
		t.Errorf("error converting ...")
	}
	// test returned data
	if len(rec.Data) < 5 {
		t.Errorf("incorrect number of data rows")
	}

}
