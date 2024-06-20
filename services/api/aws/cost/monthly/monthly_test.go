package monthly

import (
	"io"
	"net/http"
	"net/http/httptest"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
	"os"
	"testing"
	"time"
)

func TestServicesApiAwsCostMonthlyStatusCode(t *testing.T) {

	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))

	api := New(store, fs)
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
		w, r := testWRGet(route)
		mux.ServeHTTP(w, r)
		if w.Result().StatusCode != status {
			r, _ := strResponse(w.Result())
			t.Errorf("http status mismtach [%s] expected [%d], actual [%v]\n---\n%+v\n---\n", route, status, w.Result().StatusCode, r)
		}
	}

}

func TestServicesApiAwsCostMonthlyFSMatch(t *testing.T) {
	fs := testFs()
	store := data.NewStore[*cost.Cost]()
	api := New(store, fs)

	apiFs := api.FS()
	if apiFs.BaseDir() != fs.BaseDir() {
		t.Errorf("base dir mismatch")
	}
}

func TestServicesApiAwsCostMonthlyStoreMatch(t *testing.T) {
	fs := testFs()
	min, max, df := testDates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))

	api := New(store, fs)
	apiS := api.Store()

	if apiS.Length() != store.Length() {
		t.Errorf("store data mismatch")
	}
}

func TestServicesApiAwsCostMonthlyInterface(t *testing.T) {
	fs := testFs()
	store := data.NewStore[*cost.Cost]()
	api := New(store, fs)

	if testIsI[*cost.Cost, files.IWriteFS](api) {
		t.Errorf("should not be nil")
	}
}

func testIsI[V data.IEntry, F files.IReadFS](i server.IApi[V, F]) bool {
	return i == nil
}

func testFs() *files.WriteFS {
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "files-all-*")
	dfSys := os.DirFS(tDir).(files.IReadFS)
	return files.NewFS(dfSys, tDir)
}

func testMux() *http.ServeMux {
	return http.NewServeMux()
}
func testWRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}
func testDates() (min time.Time, max time.Time, df string) {
	df = time.RFC3339
	max = time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	min = time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	return
}

func strResponse(r *http.Response) (string, []byte) {
	b, _ := io.ReadAll(r.Body)
	return string(b), b
}
