package codebasecompliance

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/codebases/codebasesapi/codebasesapistats"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tmpl"
	"sync"
)

type PageContent struct {
	htmlpage.HTMLPage
	Team         string
	CodebaseData *frontmodels.CodebaseData
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request, writer http.ResponseWriter) {
	var (
		page         *PageContent
		templateName string
		team         string         = request.PathValue("team")
		wg           sync.WaitGroup = sync.WaitGroup{}
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "codebasecompliance", "func", "Handler", "url", request.URL.String())
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
		Title:        "OPG Reports - Active Codebases",
		GovUKVersion: in.GovUKVersion,
		SemVer:       in.SemVer,
	}
	template = "home-codebase-compliance"
	if team != "" {
		args.Title += " - " + cnv.Capitalize(team)
		template = "team-codebase-compliance"
	}
	page = &PageContent{
		HTMLPage:     htmlpage.New(request, args),
		CodebaseData: &frontmodels.CodebaseData{},
		Team:         team,
	}
	return
}

// dataCallers provides all the aync / concurrent api calls to fetch and attach data to this page
func dataCallers(ctx context.Context, args *frontmodels.RegisterArgs, request *http.Request) (funcs []dataCallerF) {
	var (
		team          = request.PathValue("team")
		statsEndpoint = codebasesapistats.ENDPOINT_BASE
		params        = []*rest.Param{}
	)

	// add team filter values and url
	if team != "" {
		statsEndpoint = codebasesapistats.ENDPOINT_TEAM
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
		// get list of all codebases
		func(wg *sync.WaitGroup, page *PageContent) {
			resp, err := rest.FromApi[*codebasesapistats.Response](ctx, args.ApiHost, statsEndpoint, request, params...)
			if err == nil {
				codebases := []*frontmodels.Codebase{}
				cnv.Convert(resp.Data, &codebases)
				page.CodebaseData.Codebases = codebases
			}
			wg.Done()
		},
	}
	return
}
