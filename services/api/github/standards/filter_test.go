package standards

import (
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server/response"
	"testing"
)

func TestServicesApiGithubStandardsFilter(t *testing.T) {
	fs := testhelpers.Fs()
	mux := testhelpers.Mux()

	store := data.NewStore[*std.Repository]()
	l := 90
	x := 5

	for i := 0; i < l; i++ {
		c := std.FakeCompliant(nil, std.DefaultBaselineCompliance)
		c.Archived = false
		store.Add(c)
	}
	for i := 0; i < x; i++ {
		c := std.Fake(nil)
		c.Archived = true
		store.Add(c)
	}

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := New(store, fs, resp)
	api.Register(mux)

	route := "/github/standards/v1/list/?archived=true"
	w, r := testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	str, b := response.Stringify(w.Result())
	res := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err := response.FromJson(b, res)

	if err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if res.GetStatus() != http.StatusOK {
		t.Errorf("status code failed")
		fmt.Println(str)
	}

	// convert the row back to a repo
	repos := []*std.Repository{}
	for _, row := range res.GetData().GetTableBody() {
		rep := data.FromRow[*std.Repository](row)
		repos = append(repos, rep)
	}

	if len(repos) != x {
		t.Errorf("archive filter failed")
	}

}
