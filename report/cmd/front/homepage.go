package main

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/page"
)

type homepageData struct {
	page.PageContent

	CostsByMonth []map[string]string
}

// handleHomepage
//   - '/'
func handleHomepage(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	writer http.ResponseWriter,
	request *http.Request,
) {
	var (
		templates    = page.GetTemplateFiles(info.TemplateDir)
		templateName = "index"
		data         = page.DefaultContent(conf, request)
		endpoint     = endpoints.AWSCOSTS_GROUPED
	)
	data.Teams = info.Teams

	log.Info("processing page", "url", request.URL.String())
	log.Info("getting data from", "endpoint", endpoint)

	Respond(writer, request, templateName, templates, data)
}

func RegisterHomepageHandlers(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
) {
	log.Info("registering homepage handlers ...")

	// Homepage
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleHomepage(ctx, log, conf, info, writer, request)
	})

}
