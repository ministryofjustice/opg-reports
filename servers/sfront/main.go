package main

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/pkg/bi"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/envar"
	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/pkg/tmplfuncs"
	"github.com/ministryofjustice/opg-reports/servers/sfront/lib"
)

const assetsDirectory string = "./"
const templateDir string = "./templates"

var apiVersion = bi.ApiVersion
var mode = bi.Mode

func main() {
	bi.Dump()
	Run()
}

func Run() {
	// download the assets
	lib.DownloadGovUKAssets(assetsDirectory)

	svr := lib.NewSvr(
		&lib.Cfg{
			Addr: envar.Get("FRONT_ADDR", consts.ServerDefaultFrontAddr),
			Mux:  http.NewServeMux(),
		},
		&lib.Response{
			Organisation: bi.Organisation,
			GovUKVersion: consts.GovUKFrontendVersion,
			Templates:    lib.TemplateFiles(templateDir),
			Funcs:        tmplfuncs.All,
			Errors:       []error{},
		},
		&lib.Nav{
			Tree: navigation.Configured[mode],
		},
		&lib.Api{
			Version: apiVersion,
			Addr:    envar.Get("API_ADDR", consts.ServerDefaultApiAddr),
		},
	)

	svr.Run()

}
