package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
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

func MockServer(content string, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(content))
	}))
	return server
}

// create endpoint
func MockEndpoint[T data.IEntry](
	store data.IStore[T],
	allowedParameters []string,
	head endpoint.DisplayHeadFunc,
	row endpoint.DisplayRowFunc[T],
	w http.ResponseWriter,
	r *http.Request) endpoint.IEndpoint[T] {

	qp := endpoint.NewQueryable(allowedParameters)
	parameters := qp.Parse(r)
	response := resp.New()
	response.Metadata["filters"] = parameters
	data := endpoint.NewEndpointData[T](store, nil, nil)
	display := endpoint.NewEndpointDisplay[T](head, row, nil)
	ep := endpoint.New[T]("mock", response, data, display, parameters)
	return ep
}
