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

type teampageData struct {
	page.PageContent
	TeamName     string
	CostsByMonth *datatable.DataTable
}

// handleTeampage renders the request for `/{team}` which currently displays:
//
//   - Team navigtaion
//   - Last 4 months of costs (considering billing date)
//
// Merge front end query strings with api request values so the front end
// can replace things like start_date
func handleTeampage(
	ctx context.Context, log *slog.Logger, conf *config.Config,
	info *FrontInfo,
	writer http.ResponseWriter, request *http.Request,
) {
	var (
		templateName   string            = "team"                                  // teampage uses the index template
		templates      []string          = page.GetTemplateFiles(info.TemplateDir) // all templates in the directory path
		defaultContent page.PageContent  = page.DefaultContent(conf, request)      // fetch the baseline values to render the page
		client         *restr.Repository = restr.Default(ctx, log, conf)
		service        *front.Service    = front.Default(ctx, log, conf)
		data           *teampageData     = &teampageData{PageContent: defaultContent, TeamName: request.PathValue("team")} // create the data that will be used in rendering the template
		costOptions    map[string]string = map[string]string{"team": data.TeamName, "account_name": "true"}
	)
	log.Info("processing page", "url", request.URL.String())
	// get list of teams
	data.Teams, _ = service.GetTeamNavigation(client, request)
	// get costs grouped by month
	data.CostsByMonth, _ = service.GetAwsCostsGrouped(client, request, costOptions)

	Respond(writer, request, templateName, templates, data)
}

// RegisterTeampageHandlers is called from rootCmd.RunE for endpoint
// handling delegation
//
// maps `/` to the `handleHomepage` function
func RegisterTeampageHandlers(
	ctx context.Context, log *slog.Logger, conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
) {
	log.Info("registering handler [`/team/{team}`] ...")
	// Homepage
	mux.HandleFunc("/team/{team}/", func(writer http.ResponseWriter, request *http.Request) {
		handleTeampage(ctx, log, conf, info, writer, request)
	})
}
