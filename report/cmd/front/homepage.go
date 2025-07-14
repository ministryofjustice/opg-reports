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
	"sync"
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
		wg             sync.WaitGroup    = sync.WaitGroup{} // used for concurrency
		blocks         []conF
	)
	log.Info("processing page", "url", request.URL.String())

	blocks = []conF{
		func(i ...any) {
			// get list of teams
			data.Teams, _ = service.GetTeamNavigation(client, request)
			wg.Done()
		},
		func(i ...any) {
			// get costs grouped by month
			opts := map[string]string{"team": "true"}
			data.CostsByMonth, _ = service.GetAwsCostsGrouped(client, request, opts)
			wg.Done()
		},
	}
	for _, block := range blocks {
		wg.Add(1)
		go block()
	}
	wg.Wait()
	log.Info("procesed page", "url", request.URL.String())
	Respond(log, writer, request, templateName, templates, data)
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
	log.Info("registering handler [`/{$}`] ...")
	// Homepage
	mux.HandleFunc("/{$}", func(writer http.ResponseWriter, request *http.Request) {
		handleHomepage(ctx, log, conf, info, writer, request)
	})
}
