package monthly

import (
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiAwsCostMonthlyStatusCode(t *testing.T) {
	logger.LogSetup()
	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))
	Register(mux, store)
	// /units/envs/services/
	routes := map[string]int{
		"/aws/costs/v1/monthly/-/-/":                                 http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/":                     http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/units/":               http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/units/envs/":          http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/units/envs/services/": http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02":                      http.StatusMovedPermanently,
		"/aws/costs/":                                                http.StatusNotFound,
	}
	for route, status := range routes {
		w, r := testhelpers.WRGet(route)
		mux.ServeHTTP(w, r)
		if w.Result().StatusCode != status {
			r, _ := response.Stringify(w.Result())
			t.Errorf("http status mismtach [%s] expected [%d], actual [%v]\n---\n%+v\n---\n", route, status, w.Result().StatusCode, r)
		}
	}

}
