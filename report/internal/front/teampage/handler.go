package teampage

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/headline/headlineapi/headlineapiteam"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/times"
	"opg-reports/report/package/tmpl"
	"sync"
)

type PageContent struct {
	htmlpage.HTMLPage
	Team         string
	HeadlineData *frontmodels.HeadlineData
	Dates        *frontmodels.DateRanges
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		team         string         = request.PathValue("team")
		pageName     string         = "OPG Reports"
		pageTitle    string         = fmt.Sprintf("OPG Reports - %s Overview", cnv.Capitalize(team))
		templateName string         = "team"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "teampage", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{
			Title:        pageTitle,
			Name:         pageName,
			GovUKVersion: args.GovUKVersion,
			SemVer:       args.SemVer}
		page *PageContent = &PageContent{
			HTMLPage: htmlpage.New(request, pgArgs),
			Team:     team}
	)

	log.Info("starting ...")
	// page data fetched from api via blocks
	for _, blockF := range dataCallers(ctx, args, request) {
		wg.Add(1)
		go blockF(&wg, page)
	}
	wg.Wait()
	page.HeadlineData.Team = team
	// respond
	respond.AsHTML(ctx, request, writer, page, &respond.Args{
		Template:      templateName,
		TemplateFiles: tmpl.GetTemplateFiles(args.TemplateDir),
		Funcs:         tmpl.TemplateFunctions(),
	})
	log.Info("complete.")
}

// dataCallers provides all the aync / concurrent api calls to fetch and attach data to this page
func dataCallers(ctx context.Context, args *Args, request *http.Request) []dataCallerF {
	var (
		dateEnd   = times.ResetMonth(times.Today())
		dateStart = times.Add(dateEnd, -5, times.MONTH)
		params    = []*rest.Param{
			{Type: rest.PATH, Key: "date_end", Value: times.AsYMString(dateEnd)},
			{Type: rest.PATH, Key: "date_start", Value: times.AsYMString(dateStart)},
			{Type: rest.PATH, Key: "team", Value: request.PathValue("team")},
		}
	)

	return []dataCallerF{
		// get teams
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*teamapiall.Response](ctx, args.ApiHost, teamapiall.ENDPOINT, request)
			if err == nil {
				page.Teams = resp.Data
			}
			wg.Done()
		},
		// get homepage stats
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*headlineapiteam.Response](ctx, args.ApiHost, headlineapiteam.ENDPOINT, request, params...)
			if err == nil {
				// set headlines
				page.HeadlineData = &frontmodels.HeadlineData{
					TotalCost:           resp.Data.TotalCost,
					AverageCostPerMonth: resp.Data.AverageCostPerMonth,
					OverallUptime:       resp.Data.OverallUptime,
					DateStart:           resp.Request.DateStart,
					DateEnd:             resp.Request.DateEnd,
				}
				// also set date values
				page.Dates = &frontmodels.DateRanges{
					DateStart: resp.Request.DateStart,
					DateEnd:   resp.Request.DateEnd,
					Months: times.AsYMStrings(
						times.Months(times.Add(times.Today(), -12, times.MONTH), times.Today()),
					),
				}
			}
			wg.Done()
		},
	}
}
