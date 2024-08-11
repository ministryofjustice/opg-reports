package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

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

func Dir() (dir string) {
	td := os.TempDir()
	dir, _ = os.MkdirTemp(td, "test-*")
	return
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
