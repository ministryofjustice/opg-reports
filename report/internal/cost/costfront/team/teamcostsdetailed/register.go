package teamcostsdetailed

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = "/team/{team}/costs/detailed"

type Args struct {
	ApiHost      string `json:"api"`
	GovUKVersion string `json:"govuk_version"`
	SemVer       string `json:"semver"`
	RootDir      string `json:"root_dir"`
	TemplateDir  string `json:"template_dir"`
}

func Register(ctx context.Context, mux *http.ServeMux, args *Args) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "teamcostsdetailed", "func", "Register")

	log.Info("registering handler [`" + ENDPOINT + "/{$}`] ...")
	mux.HandleFunc(fmt.Sprintf("%s/{$}", ENDPOINT), func(writer http.ResponseWriter, request *http.Request) {
		Handler(ctx, args, writer, request)
	})
}
