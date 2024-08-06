package standards

import (
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/resp"
	"testing"
)

func TestServicesApiGithubStandardsStatusCode(t *testing.T) {
	logger.LogSetup()
	mux := testhelpers.Mux()
	store := data.NewStore[*std.Repository]()
	store.Add(std.Fake(nil))

	Register(mux, store)
	routes := map[string]int{
		"/github/standards/v1/list/":               http.StatusOK,
		"/github/standards/v1/list/?archived=true": http.StatusOK,
		"/github/standards/v1/list/?teams=foobar":  http.StatusOK,
		"/github/standards/v1/list":                http.StatusMovedPermanently,
		"/github/standards/v1/list?archived=true":  http.StatusMovedPermanently,
		// this should return a 200, and the foo param is ignored
		"/github/standards/v1/list/?foo=bar": http.StatusOK,
		// no endpoint for this pattern
		"/github/standards/v1/": http.StatusNotFound,
	}
	for route, status := range routes {
		w, r := testhelpers.WRGet(route)
		mux.ServeHTTP(w, r)
		if w.Result().StatusCode != status {
			r, _ := resp.Stringify(w.Result())
			t.Errorf("http status mismtach [%s] expected [%d], actual [%v]\n---\n%+v\n---\n", route, status, w.Result().StatusCode, r)
		}
	}

}
