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
)

func TestFrontServerStaticHandler(t *testing.T) {
	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	mux := th.Mux()
	conf, _ := cnf.Load([]byte(testRealisticServerCnf))
	s := New(conf, templates, "", "")
	s.Register(mux)

	route := "/costs/"
	w, r := th.WRGet(route)

	mux.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("should return 200")
	}

	str, _ := response.Stringify(w.Result())

	if !strings.Contains(str, "OPG Report") {
		t.Errorf("org not rendered")
	}

}
