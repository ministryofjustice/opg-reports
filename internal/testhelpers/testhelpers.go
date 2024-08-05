package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/shared/files"
	"os"
	"time"
)

type TestIEntry struct {
	Id       string    `json:"id"`
	Tags     []string  `json:"tags"`
	Category string    `json:"category"`
	Status   bool      `json:"status"`
	Date     time.Time `json:"date"`
}

func (i *TestIEntry) UID() string {
	return i.Id
}
func (i *TestIEntry) TS() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
func (i *TestIEntry) Valid() bool {
	return true
}

func Fs() *files.WriteFS {
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "test-*")
	dfSys := os.DirFS(tDir).(files.IReadFS)
	return files.NewFS(dfSys, tDir)
}
func Dates() (min time.Time, max time.Time, df string) {
	df = time.RFC3339
	max = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	min = time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	return
}

func Mux() *http.ServeMux {
	return http.NewServeMux()
}
func WRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}

func MockServer(resp string, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(resp))
	}))
	return server
}
