package main

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/info"
	"github.com/ministryofjustice/opg-reports/internal/envar"
	"github.com/ministryofjustice/opg-reports/internal/tmplfuncs"
	"github.com/ministryofjustice/opg-reports/servers/front/lib"
)

const assetsDirectory string = "./"
const templateDir string = "./templates"

var apiVersion = "v1"
var mode = info.Fixtures

func main() {
	info.Log()
	Run()
}

func Run() {
	// download the assets
	lib.DownloadGovUKAssets(assetsDirectory)

	svr := lib.NewSvr(
		&lib.Cfg{
			Addr: envar.Get("FRONT_ADDR", info.ServerDefaultFrontAddr),
			Mux:  http.NewServeMux(),
		},
		&lib.Response{
			Organisation: info.Organisation,
			GovUKVersion: info.GovUKFrontendVersion,
			Templates:    lib.TemplateFiles(templateDir),
			Funcs:        tmplfuncs.All,
			Errors:       []error{},
		},
		&lib.Nav{
			Tree: lib.NavigationChoices[mode],
		},
		&lib.Api{
			Version: apiVersion,
			Addr:    envar.Get("API_ADDR", info.ServerDefaultApiAddr),
		},
	)

	svr.Run()

}
