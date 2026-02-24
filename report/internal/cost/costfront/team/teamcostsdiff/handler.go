package teamcostsdiff

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapidiffteamfilter"
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
	Months  []string
	Changes []string
	Change  string
	DateA   string
	DateB   string
}

type PageContent struct {
	htmlpage.HTMLPage
	Team     string
	CostData *tableData
	Dates    *datePicker
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		team         string         = request.PathValue("team")
		pageName     string         = "OPG Reports"
		pageTitle    string         = fmt.Sprintf("OPG Reports - %s - Cost Differences", cnv.Capitalize(team))
		templateName string         = "team-costs-differences"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "teamcostsdiff", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{Title: pageTitle, Name: pageName, GovUKVersion: args.GovUKVersion, SemVer: args.SemVer}
		page         *PageContent   = &PageContent{HTMLPage: htmlpage.New(request, pgArgs), Team: team}
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
		dateB      = times.ResetMonth(times.Today()) // use this month
		dateA      = times.Add(dateB, -1, times.MONTH)
		params     = []*rest.Param{
			{Type: rest.PATH, Key: "date_a", Value: times.AsYMString(dateA)},
			{Type: rest.PATH, Key: "date_b", Value: times.AsYMString(dateB)},
			{Type: rest.PATH, Key: "team", Value: request.PathValue("team")},
			{Type: rest.QUERY, Key: "change", Value: "100"},
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

			resp, err := rest.FromApi[*costapidiffteamfilter.Response](ctx, args.ApiHost, costapidiffteamfilter.ENDPOINT, request, params...)
			if err == nil {
				// set date values
				page.Dates = &datePicker{
					DateA:   resp.Request.DateA,
					DateB:   resp.Request.DateB,
					Change:  resp.Request.Change,
					Changes: changes,
					Months: times.AsYMStrings(
						times.Months(times.Add(times.Today(), -12, times.MONTH), times.Today()),
					),
				}
				// process the data into local structs
				page.CostData = &tableData{
					BillingDay: billingDay,
					Data:       resp.Data,
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
