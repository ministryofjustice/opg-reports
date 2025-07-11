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
		// allTeamsSrv    = front.Default[*apiAllTeams, []string](ctx, log, conf)
		templates      = page.GetTemplateFiles(info.TemplateDir)
		templateName   = "index"
		defaultContent = page.DefaultContent(conf, request)
		data           = &homepageData{PageContent: defaultContent}
	)

	log.Info("processing page", "url", request.URL.String())
	// handle page components
	data.Teams, _ = Components.TeamNavigation.Call(info.RestClient, endpoints.TEAMS_GET_ALL)
	data.CostsByMonth, _ = Components.AwsCostsGroupedByMonth.Call(
		info.RestClient,
		endpoints.Parse(
			endpoints.AWSCOSTS_GROUPED,
			map[string]string{
				"granularity": "monthly",
				"start_date":  "2025-01-01",
				"end_date":    "2025-02-01",
			},
		))

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
