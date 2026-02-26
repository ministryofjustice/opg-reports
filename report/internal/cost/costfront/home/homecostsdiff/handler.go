package homecostsdiff

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapidiff"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
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
	CostData *frontmodels.TableData
	Dates    *frontmodels.DateComparision
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *frontmodels.FrontRegisterArgs, writer http.ResponseWriter, request *http.Request) {
	var (
		pageName     string         = "OPG Reports"
		pageTitle    string         = "OPG Reports - Cost Differences"
		templateName string         = "home-costs-differences"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "homecostsdiff", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{
			Title:        pageTitle,
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
func dataCallers(ctx context.Context, args *frontmodels.FrontRegisterArgs, request *http.Request) []dataCallerF {
	var (
		billingDay = 15
		dateB      = times.ResetMonth(times.Today()) // use this month
		dateA      = times.Add(dateB, -1, times.MONTH)
		params     = []*rest.Param{
			{Type: rest.PATH, Key: "date_a", Value: times.AsYMString(dateA)},
			{Type: rest.PATH, Key: "date_b", Value: times.AsYMString(dateB)},
			{Type: rest.QUERY, Key: "change", Value: "300"},
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
		// get cost differences
		func(wg *sync.WaitGroup, page *PageContent) {
			var changes = []string{}
			for i := 100; i <= 3000; i += 100 {
				changes = append(changes, fmt.Sprintf("%d", i))
			}

			resp, err := rest.FromApi[*costapidiff.Response](ctx, args.ApiHost, costapidiff.ENDPOINT_BASE, request, params...)
			if err == nil {
				// set date values
				page.Dates = &frontmodels.DateComparision{
					DateA:   resp.Request.DateA,
					DateB:   resp.Request.DateB,
					Change:  resp.Request.Change,
					Changes: changes,
					Months: times.AsYMStrings(
						times.Months(times.Add(times.Today(), -12, times.MONTH), times.Today()),
					),
				}
				// process the data into local structs
				page.CostData = &frontmodels.TableData{
					BillingDay: billingDay,
					Data:       resp.Data,
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
