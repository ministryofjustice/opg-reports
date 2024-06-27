package server

import (
	"fmt"
	"net/http"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/files"
	"opg-reports/shared/server/response"
	"os"
	"strings"
	"testing"
)

func TestFrontServerDynamicHandlerMockedTotals(t *testing.T) {
	ms := mockServer(mockAwsCostTotalsResponse, http.StatusOK)
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
	s.Nav.Get(route).Api = ms.URL
	//
	mux := testMux()
	s.Register(mux)
	// now we fetch the local route, which should them call the mocked
	// server
	w, r := testWRGet(route)
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
	ms := mockServer(mockAwsCostUnitsResponse, http.StatusOK)
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
	s.Nav.Get(route).Api = ms.URL
	//
	mux := testMux()
	s.Register(mux)
	// now we fetch the local route, which should them call the mocked
	// server
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	str, _ := response.Stringify(w.Result())
	if !strings.Contains(str, "Costs Per Unit") {
		t.Errorf("failed, costs header not found")
	}
	fmt.Println(str)
}
