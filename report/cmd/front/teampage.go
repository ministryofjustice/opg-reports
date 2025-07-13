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

type teampageData struct {
	page.PageContent
	TeamName               string
	CostsByMonthPerAccount *datatable.DataTable
	CostsByMonthDetailed   *datatable.DataTable
}

type conF func(i ...any)

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
		data           *teampageData
		templateName   string            = "team"                                  // teampage uses the index template
		templates      []string          = page.GetTemplateFiles(info.TemplateDir) // all templates in the directory path
		defaultContent page.PageContent  = page.DefaultContent(conf, request)      // fetch the baseline values to render the page
		client         *restr.Repository = restr.Default(ctx, log, conf)
		service        *front.Service    = front.Default(ctx, log, conf)

		// mutex      *sync.Mutex    = &sync.Mutex{}
		wg         sync.WaitGroup = sync.WaitGroup{}
		pageBlocks []conF
	)
	// create the data that will be used in rendering the template
	data = &teampageData{
		PageContent: defaultContent,
		TeamName:    request.PathValue("team")}

	log.Info("processing page", "url", request.URL.String())

	pageBlocks = []conF{
		// get list of all teams
		func(i ...any) {
			data.Teams, _ = service.GetTeamNavigation(client, request)
			wg.Done()
		},
		// get tabular costs grouped by the account name & filtered by the team
		func(i ...any) {
			options := map[string]string{"team": data.TeamName, "account_name": "true"}
			data.CostsByMonthPerAccount, _ = service.GetAwsCostsGrouped(client, request, options)
			wg.Done()
		},
		// get the table of costs broken down in detail
		func(i ...any) {
			options := map[string]string{
				"team":         data.TeamName,
				"account_name": "true",
				"environment":  "true",
				"service":      "true",
				// "region":       "true",
			}
			// adjust := func(p map[string]string) {
			// 	p["start_date"] = utils.DateStringAddMonths(p["start_date"], utils.DATE_FORMATS.YM, 3)
			// }
			data.CostsByMonthDetailed, _ = service.GetAwsCostsGrouped(client, request, options /*, adjust*/)
			wg.Done()
		},
	}

	for _, block := range pageBlocks {
		wg.Add(1)
		block()
	}
	wg.Wait()

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
