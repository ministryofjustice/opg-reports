package standards

import (
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiGithubStandardsStatusCode(t *testing.T) {
	logger.LogSetup()
	fs := testhelpers.Fs()

	mux := testhelpers.Mux()
	store := data.NewStore[*std.Repository]()
	store.Add(std.Fake(nil))
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	routes := map[string]int{
		"/github/standards/v1/list/":               http.StatusOK,
		"/github/standards/v1/list/?archived=true": http.StatusOK,
		"/github/standards/v1/list/?teams=foobar":  http.StatusOK,
		"/github/standards/v1/list":                http.StatusMovedPermanently,
		"/github/standards/v1/list?archived=true":  http.StatusMovedPermanently,
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
