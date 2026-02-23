package costsdetailed

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapidetailed"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/times"
	"opg-reports/report/package/tmpl"
	"sync"
	"time"
)

type PageContent struct {
	htmlpage.HTMLPage
	CostData *costapidetailed.Response
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		// err  error
		pageName     string         = "OPG Reports"
		templateName string         = "home-costs-detailed"
		endDate      time.Time      = times.Today()
		startDate    time.Time      = times.Add(endDate, -12, times.MONTH)
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "home", "func", "costsdetailed", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		page         *PageContent   = &PageContent{
			HTMLPage: htmlpage.New(request, &htmlpage.Args{Name: pageName, GovUKVersion: args.GovUKVersion}),
		}
	)
	log.Info("starting ...")
	// page data fetched from api via blocks
	for _, blockF := range dataCallers(ctx, args, request) {
		wg.Add(1)
		go blockF(&wg, page)
	}
	wg.Wait()
	page.CostData.Months = times.AsYMStrings(times.Months(startDate, endDate))

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
	return []dataCallerF{
		// get teams
		func(wg *sync.WaitGroup, page *PageContent) {
			if teams, err := teamapiall.Get(ctx, args.ApiHost, request); err == nil {
				page.Teams = teams
			}
			wg.Done()
		},
		// get homepage stats
		func(wg *sync.WaitGroup, page *PageContent) {
			if costs, err := costapidetailed.Get(ctx, args.ApiHost, request); err == nil {
				page.CostData = costs
			}
			wg.Done()
		},
	}
}
