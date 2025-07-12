package main

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/page"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front"
	"opg-reports/report/internal/service/front/datatable"
)

type homepageData struct {
	page.PageContent
	CostsByMonth *datatable.DataTable
}

// handleHomepage renders the request for `/` which currently displays:
//
//   - Team navigtaion
//   - Last 4 months of costs (considering billing date)
//
// Merge front end query strings with api request values so the front end
// can replace things like start_date
//
// Uses multiple `Components` to generate all the data displayed on this
// page
func handleHomepage(
	ctx context.Context, log *slog.Logger, conf *config.Config,
	info *FrontInfo,
	writer http.ResponseWriter, request *http.Request,
) {
	var (
		templateName   string            = "index"                                    // homepage uses the index template
		templates      []string          = page.GetTemplateFiles(info.TemplateDir)    // all templates in the directory path
		defaultContent page.PageContent  = page.DefaultContent(conf, request)         // fetch the baseline values to render the page
		data           *homepageData     = &homepageData{PageContent: defaultContent} // create the data that will be used in rendering the template
		client         *restr.Repository = restr.Default(ctx, log, conf)
		service        *front.Service    = front.Default(ctx, log, conf)
	)
	log.Info("processing page", "url", request.URL.String())
	// get list of teams
	data.Teams, _ = service.GetTeamNavigation(client, request)
	// get costs grouped by month
	data.CostsByMonth, _ = service.GetAwsCostsGrouped(client, request, map[string]string{"team": "true"})

	Respond(writer, request, templateName, templates, data)
}

// RegisterHomepageHandlers is called from rootCmd.RunE for endpoint
// handling delegation
//
// maps `/` to the `handleHomepage` function
func RegisterHomepageHandlers(
	ctx context.Context, log *slog.Logger, conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
) {
	log.Info("registering handler [`/`] ...")
	// Homepage
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleHomepage(ctx, log, conf, info, writer, request)
	})
}
