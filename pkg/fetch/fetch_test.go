package fetch_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
	"github.com/ministryofjustice/opg-reports/pkg/fetch"
)

type testO struct{}
type testResponse struct {
	Body struct {
		Message string `json:"message" example:"Successful connection."`
	}
}

func mockServer(f func(w http.ResponseWriter, r *http.Request), loglevel string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os.Setenv("LOG_LEVEL", loglevel)
		f(w, r)
	}))
}

func serverOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TestFetchSimple calls a local mock server to make
// sure the response is correct
func TestFetchSimple(t *testing.T) {
	var ok = "OK"
	var mock = mockServer(serverOk, "info")
	defer mock.Close()
	// call that api endpoint
	var ep endpoints.ApiEndpoint = "/"
	var full, _ = url.Parse(mock.URL)

	content, code, err := fetch.Fetch("http://"+full.Host, ep)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	if code != http.StatusOK {
		t.Errorf("response code mismatch - expected [%d] actual [%v]", http.StatusOK, code)
	}

	if string(content) != ok {
		t.Errorf("response body mismatch - expected [%s] actual [%s]", ok, content)
	}

}
