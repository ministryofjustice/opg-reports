package server

import (
	"fmt"
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
	conf, _ := cnf.Load(testRealServerCnf())
	s := New(conf, templates, "", "")
	s.Register(mux)

	// --- TEST KNOWN TOP LEVEL STATIC ROUTES
	routes := []string{}
	for _, sect := range conf.Sections {
		if len(sect.Api) == 0 {
			routes = append(routes, sect.Href)
		}
		for _, s := range sect.Sections {
			if len(s.Api) == 0 {
				routes = append(routes, s.Href)
			}
		}
	}

	for _, route := range routes {
		w, r := th.WRGet(route)
		mux.ServeHTTP(w, r)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("should return 200: %v", w.Result().StatusCode)
		}

		str, _ := response.Stringify(w.Result())
		if !strings.Contains(str, "OPG") {
			fmt.Println(str)
			t.Errorf("org not rendered for route [%s]", route)

		}
	}

}
