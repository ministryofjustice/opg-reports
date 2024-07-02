package server

import (
	"net/http"
	th "opg-reports/internal/testhelpers"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/files"
	"opg-reports/shared/server/response"
	"os"
	"strings"
	"testing"
	"time"
)

type testEntry struct {
	Id       string `json:"id"`
	Tag      string `json:"tag"`
	Category string `json:"category"`
}

func (i *testEntry) UID() string {
	return i.Id
}
func (i *testEntry) TS() time.Time {
	return time.Now().UTC()
}
func (i *testEntry) Valid() bool {
	return true
}

func TestFrontServerDynamicHandlerMocked(t *testing.T) {

}

func TestFrontServerDynamicHandlerMockedTotals(t *testing.T) {

	ms := th.MockServer(mockAwsCostMonthlyTotals(), http.StatusOK)
	defer ms.Close()

	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	route := "/costs/aws/totals/"
	conf, _ := cnf.Load([]byte(testRealisticServerCnf))
	// create new
	s := New(conf, templates, "", "")
	// point the totals route to look at the test api
	s.Nav.Get(route).Api = map[string]string{"url": ms.URL}
	//
	mux := th.Mux()
	s.Register(mux)
	// now we fetch the local route, which should them call the mocked
	// server
	w, r := th.WRGet(route)
	mux.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected error")
	}

	str, _ := response.Stringify(w.Result())
	if !strings.Contains(str, "AWS Costs") {
		t.Errorf("failed, costs header not found")
	}
	// fmt.Println(str)

}

func TestFrontServerDynamicHandlerMockedUnits(t *testing.T) {
	content := mockAwsCostMonthlyUnits()
	ms := th.MockServer(content, http.StatusOK)
	defer ms.Close()

	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	route := "/costs/aws/units/"
	conf, _ := cnf.Load([]byte(testRealisticServerCnf))
	// create new
	s := New(conf, templates, "", "")
	// point the totals route to look at the test api
	s.Nav.Get(route).Api = map[string]string{"url": ms.URL}
	//
	mux := th.Mux()
	s.Register(mux)
	// now we fetch the local route, which should them call the mocked
	// server
	w, r := th.WRGet(route)
	mux.ServeHTTP(w, r)

	str, _ := response.Stringify(w.Result())
	if !strings.Contains(str, "Costs Per Unit") {
		t.Errorf("failed, costs header not found")
	}
	// fmt.Println(str)
}
