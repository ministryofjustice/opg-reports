package teamcosts

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapiteamfilter"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"opg-reports/report/package/tmpl"
	"sync"
)

type PageContent struct {
	htmlpage.HTMLPage
	Team     string
	CostData *frontmodels.TableData
	Dates    *frontmodels.DateRanges
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *frontmodels.FrontRegisterArgs, writer http.ResponseWriter, request *http.Request) {
	var (
		team         string         = request.PathValue("team")
		pageName     string         = "OPG Reports"
		pageTitle    string         = fmt.Sprintf("OPG Reports - %s - Costs", cnv.Capitalize(team))
		templateName string         = "teams-costs-by-account"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "teamcosts", "func", "Handler", "url", request.URL.String())
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

	// respond
	respond.AsHTML(ctx, request, writer, page, &respond.Args{
		Template:      templateName,
		TemplateFiles: tmpl.GetTemplateFiles(args.TemplateDir),
		Funcs:         tmpl.TemplateFunctions(),
	})
	log.Info("complete.")
}

// dataCallers provides all the aync / concurrent api calls to fetch and attach data to this page
func dataCallers(ctx context.Context, args *frontmodels.FrontRegisterArgs, request *http.Request) []dataCallerF {
	var (
		billingDay = 15
		dateEnd    = times.ResetMonth(times.Today()) // use this month
		dateStart  = times.Add(dateEnd, -5, times.MONTH)
		params     = []*rest.Param{
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
		// get homepage costs - trigger the same end date as others
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*costapiteamfilter.Response](ctx, args.ApiHost, costapiteamfilter.ENDPOINT, request, params...)
			if err == nil {
				// set date values
				page.Dates = &frontmodels.DateRanges{
					DateStart: resp.Request.DateStart,
					DateEnd:   resp.Request.DateEnd,
					Months: times.AsYMStrings(
						times.Months(times.Add(times.Today(), -12, times.MONTH), times.Today()),
					),
				}
				// process the data into local structs
				page.CostData = &frontmodels.TableData{
					BillingDay: billingDay,
					Data:       resp.Data,
					Summary:    resp.Summary,
					Headers: &frontmodels.TableHeaders{
						Labels: resp.Headers[tabulate.KEY],
						Data:   resp.Headers[tabulate.DATA],
						Extra:  resp.Headers[tabulate.EXTRA],
						End:    resp.Headers[tabulate.END],
					},
				}
			}
			wg.Done()
		},
	}
}
