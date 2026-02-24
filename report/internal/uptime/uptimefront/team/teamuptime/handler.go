package teamuptime

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteamfilter"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"opg-reports/report/package/tmpl"
	"sync"
)

type tableHeaders struct {
	Labels []string `json:"labels"`
	Data   []string `json:"data"`
	Extra  []string `json:"extra"`
	End    []string `json:"end"`
}

// tableData is used to handle the cost table data construct
type tableData struct {
	Headers *tableHeaders            `json:"headers"` // headers contains details for table headers / rendering
	Data    []map[string]interface{} `json:"data"`    // the actual data results
	Summary map[string]interface{}   `json:"summary"` // used to contain table totals etc
}

// datepicker is used for selecting date ranges to show data for
type datePicker struct {
	Months    []string
	DateStart string
	DateEnd   string
}

type PageContent struct {
	htmlpage.HTMLPage
	Team       string
	UptimeData *tableData
	Dates      *datePicker
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		// err  error
		pageName     string         = "OPG Reports"
		templateName string         = "team-uptime-by-team"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "teamuptime", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		page         *PageContent   = &PageContent{
			HTMLPage: htmlpage.New(request, &htmlpage.Args{Name: pageName, GovUKVersion: args.GovUKVersion}),
			Team:     request.PathValue("team"),
		}
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
func dataCallers(ctx context.Context, args *Args, request *http.Request) []dataCallerF {
	var (
		dateEnd   = times.ResetMonth(times.Today()) // use this month
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
		// get team uptime breakdown
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*uptimeapiteamfilter.Response](ctx, args.ApiHost, uptimeapiteamfilter.ENDPOINT, request, params...)
			if err == nil {
				// process the data into local structs
				page.UptimeData = &tableData{
					Data:    resp.Data,
					Summary: resp.Summary,
					Headers: &tableHeaders{
						Labels: resp.Headers[tabulate.KEY],
						Data:   resp.Headers[tabulate.DATA],
						Extra:  resp.Headers[tabulate.EXTRA],
						End:    resp.Headers[tabulate.END],
					},
				}
				// also set date values
				page.Dates = &datePicker{
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
