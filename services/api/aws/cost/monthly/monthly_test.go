package monthly

import (
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiAwsCostMonthlyStatusCode(t *testing.T) {

	fs := testhelpers.Fs()

	mux := testhelpers.Mux()
	min, max, df := testhelpers.Dates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	routes := map[string]int{
		"/aws/costs/v1/monthly/":                    http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/":    http.StatusOK,
		"/aws/costs/v1/monthly/-/-/":                http.StatusOK,
		"/aws/costs/v1/monthly":                     http.StatusMovedPermanently,
		"/aws/costs/v1/monthly/2024-01/2024-02":     http.StatusMovedPermanently,
		"/aws/costs/v1/monthly/not-a-date/-/":       http.StatusConflict,
		"/aws/costs/v1/monthly/not-a-date/2024-02/": http.StatusConflict,
		"/aws/costs/v1/monthly/not-a-date/invalid/": http.StatusConflict,
		"/aws/costs/v1/monthly/2024-04/invalid/":    http.StatusConflict,
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

func TestServicesApiAwsCostMonthlyFSMatch(t *testing.T) {
	fs := testhelpers.Fs()
	store := data.NewStore[*cost.Cost]()
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)

	apiFs := api.FS()
	if apiFs.BaseDir() != fs.BaseDir() {
		t.Errorf("base dir mismatch")
	}
}

func TestServicesApiAwsCostMonthlyStoreMatch(t *testing.T) {
	fs := testhelpers.Fs()
	min, max, df := testhelpers.Dates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)

	apiS := api.Store()

	if apiS.Length() != store.Length() {
		t.Errorf("store data mismatch")
	}
}
