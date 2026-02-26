package costapidetailed

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/global/seeds"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
	"opg-reports/report/package/times"
	"path/filepath"
	"testing"
)

func TestCostApiDetailedHandler(t *testing.T) {
	var (
		err    error
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-handler.db")
		end    = times.AsYMString(times.Today())
		start  = times.AsYMString(times.Add(times.Today(), -3, times.YEAR))
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
	// "/v1/costs/between/{date_start}/{date_end}/detailed"
	url := "/v1/costs/between/" + start + "/" + end + "/detailed/"
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
		t.Errorf("error converting ... [%s]", err.Error())
	}
	// - test returned data
	if len(rec.Data) < 1 {
		t.Errorf("incorrect number of data rows; might be due to seed data using random date")
	}
	if rec.Request.DateEnd != end {
		t.Error("data_end failed to return correctly")
	}
	if rec.Request.DateStart != start {
		t.Error("data_start failed to return correctly")
	}
	if len(rec.Headers["labels"]) < 1 {
		t.Error("incorrect number of labels returned")
	}
}
