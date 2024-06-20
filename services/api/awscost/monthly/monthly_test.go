package monthly

import (
	"encoding/json"
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
	max = time.Now().UTC()
	min = time.Date(max.Year()-2, max.Month(), 1, 0, 0, 0, 0, time.UTC)
	return
}
func decode[T server.IApiResponse](c io.Reader) T {
	var i T
	json.NewDecoder(c).Decode(&i)
	return i
}

func TestServicesApiAwsCostMonthlyStatusCode(t *testing.T) {

	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	store := data.NewStore[*cost.Cost]()
	store.Add(cost.Fake(nil, min, max, df))
	// New[*cost.Cost, data.IStore[*cost.Cost], files.IWriteFS](store, fs)
	api := New(store, fs)
	api.Register(mux)

	routes := map[string]int{
		"/aws/costs/v1/monthly/":                 http.StatusOK,
		"/aws/costs/v1/monthly/2024-01/2024-02/": http.StatusOK,
		"/aws/costs/v1/monthly":                  http.StatusMovedPermanently,
		"/aws/costs/v1/monthly/2024-01/2024-02":  http.StatusMovedPermanently,
	}
	for route, status := range routes {
		w, r := testWRGet(route)
		mux.ServeHTTP(w, r)
		if w.Result().StatusCode != status {
			t.Errorf("http status mismtach [%s] expected [%d], actual [%v]", route, status, w.Result().StatusCode)
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

func testIsI[V data.IEntry, F files.IReadFS](i server.IApi[V, F]) bool {
	return i == nil
}
func TestServicesApiAwsCostMonthlyInterface(t *testing.T) {
	fs := testFs()
	store := data.NewStore[*cost.Cost]()
	api := New(store, fs)

	if testIsI[*cost.Cost, files.IWriteFS](api) {
		t.Errorf("should not be nil")
	}
}
