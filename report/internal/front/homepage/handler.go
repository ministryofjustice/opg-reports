package homepage

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapiteam"
	"opg-reports/report/internal/headline/headlineapi/headlineapihome"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteam"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"opg-reports/report/package/tmpl"
	"sync"
)

// headlineData is used to show healdine figures on the page
type headlineData struct {
	TotalCost           float64 `json:"total_cost"`             // total cost result
	AverageCostPerMonth float64 `json:"average_cost_per_month"` // average cost per month
	OverallUptime       float64 `json:"overall_uptime"`         // uptime
	DateStart           string  `json:"date_start"`
	DateEnd             string  `json:"date_end"`
}

type tableHeaders struct {
	Labels []string `json:"labels"`
	Data   []string `json:"data"`
	Extra  []string `json:"extra"`
	End    []string `json:"end"`
}

// tableData is used to handle the cost table data construct
type tableData struct {
	Headers    *tableHeaders            `json:"headers"` // headers contains details for table headers / rendering
	Data       []map[string]interface{} `json:"data"`    // the actual data results
	Summary    map[string]interface{}   `json:"summary"` // used to contain table totals etc
	BillingDay int                      `json:"billing_day"`
}

// datepicker is used for selecting date ranges to show data for
type datePicker struct {
	Months    []string
	DateStart string
	DateEnd   string
}

type PageContent struct {
	htmlpage.HTMLPage
	HeadlineData *headlineData
	CostData     *tableData
	UptimeData   *tableData
	Dates        *datePicker
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		// err  error
		pageName     string         = "OPG Reports"
		templateName string         = "home"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "homepage", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{Title: pageName, Name: pageName, GovUKVersion: args.GovUKVersion, SemVer: args.SemVer}
		page         *PageContent   = &PageContent{HTMLPage: htmlpage.New(request, pgArgs)}
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
		billingDay = 15
		dateEnd    = times.Add(times.ResetMonth(times.Today()), -1, times.MONTH) // home page uses last complete month
		dateStart  = times.Add(dateEnd, -4, times.MONTH)                         // show 3 months
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
		// get homepage stats
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*headlineapihome.Response](ctx, args.ApiHost, headlineapihome.ENDPOINT, request, params...)
			if err == nil {
				// set headlines
				page.HeadlineData = &headlineData{
					TotalCost:           resp.Data.TotalCost,
					AverageCostPerMonth: resp.Data.AverageCostPerMonth,
					OverallUptime:       resp.Data.OverallUptime,
					DateStart:           resp.Request.DateStart,
					DateEnd:             resp.Request.DateEnd,
				}

			}
			wg.Done()
		},
		// get homepage costs - trigger the same end date as others
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*costapiteam.Response](ctx, args.ApiHost, costapiteam.ENDPOINT, request, params...)
			if err == nil {
				// process the data into local structs
				page.CostData = &tableData{
					BillingDay: billingDay,
					Data:       resp.Data,
					Summary:    resp.Summary,
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
		// get homepage uptime breakdown
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*uptimeapiteam.Response](ctx, args.ApiHost, uptimeapiteam.ENDPOINT, request, params...)
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
			}
			wg.Done()
		},
	}
}
