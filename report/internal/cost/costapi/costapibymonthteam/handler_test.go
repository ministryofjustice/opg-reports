package costapibymonthteam

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dump"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
	"testing"
)

func TestCostApiByMonthTeamHandler(t *testing.T) {

	// resp := writer.Result()
	// response.AsT(resp, recieved)
	// fmt.Println(dump.Any(recieved))
	mux := http.NewServeMux()
	ctx := cntxt.AddLogger(t.Context(), logger.New("error"))
	url := "/v1/costs/between/2025-01/2026-01/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	writer := httptest.NewRecorder()

	// setup the bindings to the test handler and call
	Register(ctx, mux, &Config{})
	mux.ServeHTTP(writer, req)

	// get and parse the result
	resp := writer.Result()
	rec := &Response{}
	err := response.AsT(resp, &rec)
	if err != nil {
		t.Errorf("error converting ...")
	}

	fmt.Println(dump.Any(rec.Headers))
	t.FailNow()
}
