package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/timer"
)

// CopyFile copys content of in into file name out
func CopyFile(in string, out string) {
	r, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	w, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	w.ReadFrom(r)
}

// Dir generates a temp directory, upto user to delete directory
func Dir() (dir string) {
	td := os.TempDir()
	dir, _ = os.MkdirTemp(td, "test-*")
	return
}

// Dates generate min, max and data formats that we use
func Dates() (min time.Time, max time.Time, df string) {
	df = time.RFC3339
	max = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	min = time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	return
}

// Mux test
func Mux() *http.ServeMux {
	return http.NewServeMux()
}

// WRGet returns test http
func WRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}

// MockServer generates a mockserver with a handler attached and sets the log level
func MockServer(f func(w http.ResponseWriter, r *http.Request), loglevel string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os.Setenv("LOG_LEVEL", loglevel)
		f(w, r)
	}))
}

type Simple struct {
	Name string `json:"name"`
}

// Ts is a test struct
type Ts struct {
	S time.Time `json:"start"`
	E time.Time `json:"end"`
}

// T provides a quick timer method from timer package
func T() *timer.Ts {
	return timer.New()
}
