package portfolio

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/codebasereleases/codebasereleasesapi"
	"opg-reports/report/internal/cost/costapi/costapiteam"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteam"
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
	Team        string
	CostData    *frontmodels.TableData
	UptimeData  *frontmodels.TableData
	ReleaseData *frontmodels.ReleaseData
	Dates       *frontmodels.DateRanges
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request, writer http.ResponseWriter) {
	var (
		pageName     string         = "OPG Reports"
		templateName string         = "portfolio"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "portfolio", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{
			Title:        pageName,
			Name:         pageName,
			GovUKVersion: args.GovUKVersion,
			SemVer:       args.SemVer}
		page *PageContent = &PageContent{HTMLPage: htmlpage.New(request, pgArgs)}
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
func dataCallers(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request) []dataCallerF {
	var (
		billingDay = 15
		dateEnd    = times.Add(times.ResetMonth(times.Today()), -1, times.MONTH) // home page uses last complete month
		dateStart  = times.Add(dateEnd, -5, times.MONTH)                         // show 6 months
		params     = []*rest.Param{
			{Type: rest.PATH, Key: "date_end", Value: times.AsYMString(dateEnd)},
			{Type: rest.PATH, Key: "date_start", Value: times.AsYMString(dateStart)},
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
		// get release stats grouped by team
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*codebasereleasesapi.Response](ctx, args.ApiHost, codebasereleasesapi.ENDPOINT_BASE, request, params...)
			if err == nil {
				releases := []*frontmodels.Release{}
				summary := &frontmodels.Release{}
				// covnert to front end version
				cnv.Convert(resp.Data, &releases)
				cnv.Convert(resp.Summary, &summary)
				page.ReleaseData = &frontmodels.ReleaseData{
					Releases: releases,
					Summary:  summary,
				}
			}
			wg.Done()
		},
		// get costs - trigger the same end date as others
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*costapiteam.Response](ctx, args.ApiHost, costapiteam.ENDPOINT_BASE, request, params...)
			if err == nil {
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
		// get uptime breakdown
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*uptimeapiteam.Response](ctx, args.ApiHost, uptimeapiteam.ENDPOINT_BASE, request, params...)
			if err == nil {
				// process the data into local structs
				page.UptimeData = &frontmodels.TableData{
					Data:    resp.Data,
					Summary: resp.Summary,
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
