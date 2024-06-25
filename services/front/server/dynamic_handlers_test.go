package server

import (
	"fmt"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
	"os"
	"testing"
)

func TestFrontServerDynamicHandlerMocked(t *testing.T) {
	ms := mockServerAWSCostTotals()
	defer ms.Close()

	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	route := "/costs/aws/totals/"
	conf, _ := cnf.Load([]byte(testRealisticCfg))
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

	str, _ := server.ResponseAsStrings(w.Result())
	fmt.Println(str)
}
