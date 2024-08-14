package testhelpers

import (
	"fmt"
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

func MockServer(f func(w http.ResponseWriter, r *http.Request), loglevel string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os.Setenv("LOG_LEVEL", loglevel)
		f(w, r)
	}))
}

type Ts struct {
	S time.Time
	E time.Time
}

func (t *Ts) Stop() *Ts {
	t.E = time.Now().UTC()
	return t
}
func (t *Ts) Seconds() string {
	if t.E.Year() == 0 {
		t.Stop()
	}
	dur := t.E.Sub(t.S)
	return fmt.Sprintf("%f", dur.Seconds())
}

func T() *Ts {
	return &Ts{S: time.Now().UTC()}
}
