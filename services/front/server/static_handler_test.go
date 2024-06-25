package server

import (
	"net/http"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/files"
	"os"
	"strings"
	"testing"
)

func TestFrontServerStaticHandler(t *testing.T) {
	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	mux := testMux()
	conf, _ := cnf.Load([]byte(testRealisticCfg))
	s := New(conf, templates)
	s.Register(mux)

	route := "/costs/"
	w, r := testWRGet(route)

	mux.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("should return 200")
	}

	str, _ := strResponse(w.Result())

	if !strings.Contains(str, "OPG Report") {
		t.Errorf("org not rendered")
	}

}
