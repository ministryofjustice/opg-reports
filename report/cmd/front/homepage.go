package main

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/page"
	"opg-reports/report/internal/utils"
	"time"
)

type homepageData struct {
	page.PageContent
	CostsByMonth *dataTable
}

// homepageParams provides the values for placeholders in the api endpoints we
// call on the homepage for costs and others
//
// This is then merged with http.Request.URL.Query() values so they can
// be overwritten by front end (to view other months etc)
func homepageParams(conf *config.Config) (defaults map[string]string) {
	var (
		now       = time.Now().UTC()
		endDate   = utils.BillingMonth(now, conf.Aws.BillingDate)
		startDate = endDate.AddDate(0, -6, 1)
	)

	defaults = map[string]string{
		"start_date":  startDate.Format(utils.DATE_FORMATS.YMD),
		"end_date":    endDate.Format(utils.DATE_FORMATS.YMD),
		"granularity": string(utils.GranularityMonth),
		"team":        "true",
	}

	return
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
		templateName   = "index"                                                       // homepage uses the index template
		templates      = page.GetTemplateFiles(info.TemplateDir)                       // all templates in the directory path
		defaultContent = page.DefaultContent(conf, request)                            // fetch the baseline values to render the page
		data           = &homepageData{PageContent: defaultContent}                    // create the data that will be used in rendering the template
		params         = utils.MergeRequestWithDefaults(request, homepageParams(conf)) // merge the api parameters with one from the current request
	)
	log.Info("processing page", "url", request.URL.String())

	// handle page components
	data.Teams, _ = Components.TeamNavigation.Call(info.RestClient, endpoints.TEAMS_GET_ALL)

	data.CostsByMonth, _ = Components.AwsCostsGroupedByMonth.Call(info.RestClient, endpoints.Parse(endpoints.AWSCOSTS_GROUPED, params))

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
	log.Info("registering homepage handlers ...")
	// Homepage
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleHomepage(ctx, log, conf, info, writer, request)
	})
}
