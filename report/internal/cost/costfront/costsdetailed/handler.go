package costsdetailed

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapidetailed"
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
func Handler(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request, writer http.ResponseWriter) {
	var (
		page         *PageContent
		templateName string
		team         string         = request.PathValue("team")
		wg           sync.WaitGroup = sync.WaitGroup{}
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "costsdetailed", "func", "Handler", "url", request.URL.String())
	)

	log.Info("starting ...")
	page, templateName = getPage(team, args, request)
	if team != "" {
		log.Info("found team parameter ... ", "team", team)
	}
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

func getPage(team string, in *frontmodels.RegisterArgs, request *http.Request) (page *PageContent, template string) {
	var args *htmlpage.Args = &htmlpage.Args{
		Name:         "OPG Reports",
		Title:        "OPG Reports - Detailed AWS Costs",
		GovUKVersion: in.GovUKVersion,
		SemVer:       in.SemVer,
	}
	template = "home-costs-detailed"
	if team != "" {
		args.Title += " - " + cnv.Capitalize(team)
		template = "team-costs-detailed"
	}
	page = &PageContent{
		HTMLPage: htmlpage.New(request, args),
		Team:     team,
	}
	return
}

// dataCallers provides all the aync / concurrent api calls to fetch and attach data to this page
//
// Will add team filter into the calling endpoint if required
func dataCallers(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request) (funcs []dataCallerF) {
	var (
		team         = request.PathValue("team")
		billingDay   = 15
		costEndpoint = costapidetailed.ENDPOINT_BASE
		dateEnd      = times.ResetMonth(times.Today())
		dateStart    = times.Add(dateEnd, -5, times.MONTH)
		params       = []*rest.Param{
			{Type: rest.PATH, Key: "date_end", Value: times.AsYMString(dateEnd)},
			{Type: rest.PATH, Key: "date_start", Value: times.AsYMString(dateStart)},
		}
	)
	// add team filter values and url
	if team != "" {
		costEndpoint = costapidetailed.ENDPOINT_TEAM
		params = append(params, &rest.Param{Type: rest.PATH, Key: "team", Value: team})
	}

	funcs = []dataCallerF{
		// get teams
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*teamapiall.Response](ctx, args.ApiHost, teamapiall.ENDPOINT, request)
			if err == nil {
				page.Teams = resp.Data
			}
			wg.Done()
		},
		// get detailed costs
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*costapidetailed.Response](ctx, args.ApiHost, costEndpoint, request, params...)
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
	return
}
