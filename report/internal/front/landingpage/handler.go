package landingpage

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/codebases/codebasesapi/codebaseapiowners"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/headline/headlineapi/headlineapi"
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
	CodebaseData *frontmodels.CodebaseData
	Dates        *frontmodels.DateRanges
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request, writer http.ResponseWriter) {
	var (
		page         *PageContent
		templateName string
		team         string         = request.PathValue("team")
		wg           sync.WaitGroup = sync.WaitGroup{}
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "landingpage", "func", "Handler", "url", request.URL.String())
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
	page.HeadlineData.Team = team
	page.CodebaseData.Team = team
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
		Title:        "OPG Reports",
		GovUKVersion: in.GovUKVersion,
		SemVer:       in.SemVer,
	}
	template = "landing-page"
	if team != "" {
		args.Title += " - " + cnv.Capitalize(team)
	}
	page = &PageContent{
		HTMLPage:     htmlpage.New(request, args),
		Team:         team,
		CodebaseData: &frontmodels.CodebaseData{},
	}
	return
}

// dataCallers provides all the aync / concurrent api calls to fetch and attach data to this page
func dataCallers(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request) (funcs []dataCallerF) {
	var (
		team          = request.PathValue("team")
		headEndpoint  = headlineapi.ENDPOINT_BASE
		ownerEndpoint = codebaseapiowners.ENDPOINT_TEAM
		dateEnd       = times.ResetMonth(times.Today())
		dateStart     = times.Add(dateEnd, -5, times.MONTH)
		params        = []*rest.Param{
			{Type: rest.PATH, Key: "date_end", Value: times.AsYMString(dateEnd)},
			{Type: rest.PATH, Key: "date_start", Value: times.AsYMString(dateStart)},
		}
	)

	// add team filter values and url
	if team != "" {
		headEndpoint = headlineapi.ENDPOINT_TEAM
		ownerEndpoint = codebaseapiowners.ENDPOINT_TEAM
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
		// get code ownership items
		func(wg *sync.WaitGroup, page *PageContent) {
			// if there is now team, then push a team of `none` to find those without owners
			if team == "" {
				params = append(params, &rest.Param{Type: rest.PATH, Key: "team", Value: "none"})
			}
			resp, err := rest.FromApi[*codebaseapiowners.Response](ctx, args.ApiHost, ownerEndpoint, request, params...)
			if err == nil {
				codeowners := []*frontmodels.Codeowner{}
				cnv.Convert(resp.Data, &codeowners)
				page.CodebaseData.CodeOwners = codeowners
			}
			wg.Done()
		},
		// get landingpage stats
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*headlineapi.Response](ctx, args.ApiHost, headEndpoint, request, params...)
			if err == nil {
				// set headlines
				page.HeadlineData = &frontmodels.HeadlineData{
					DateStart:           resp.Request.DateStart,
					DateEnd:             resp.Request.DateEnd,
					TotalCost:           resp.Data.TotalCost,
					AverageCostPerMonth: resp.Data.AverageCostPerMonth,
					OverallUptime:       resp.Data.OverallUptime,
					CodebaseCount:       resp.Data.CodebaseCount,
					CodebasePassed:      resp.Data.CodebasePassed,
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
	return
}
