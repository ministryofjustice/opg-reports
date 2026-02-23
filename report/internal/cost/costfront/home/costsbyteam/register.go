package costsbyteam

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = "/home/costs/team"

type Args struct {
	ApiHost      string `json:"api"`
	GovUKVersion string `json:"govuk_version"`
	RootDir      string `json:"root_dir"`
	TemplateDir  string `json:"template_dir"`
}

func Register(ctx context.Context, mux *http.ServeMux, args *Args) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "homepage", "func", "Register")

	// index page
	log.Info("registering handler [`/{$}`] ...")
	mux.HandleFunc(fmt.Sprintf("%s/{$}", ENDPOINT), func(writer http.ResponseWriter, request *http.Request) {
		Handler(ctx, args, writer, request)
	})
}
