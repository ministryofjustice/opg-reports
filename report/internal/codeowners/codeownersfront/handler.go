package codeownersfront

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/codeowners/codeownersapi"
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
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "codeownersfront", "func", "Handler", "url", request.URL.String())
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
		Title:        "OPG Reports - Code Owners",
		GovUKVersion: in.GovUKVersion,
		SemVer:       in.SemVer,
	}
	template = "codebase-owners"
	if team != "" {
		args.Title += " - " + cnv.Capitalize(team)
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
		ownerEndpoint = codeownersapi.ENDPOINT_BASE
		params        = []*rest.Param{}
	)

	// add team filter values and url
	if team != "" {
		ownerEndpoint = codeownersapi.ENDPOINT_TEAM
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
			resp, err := rest.FromApi[*codeownersapi.Response](ctx, args.ApiHost, ownerEndpoint, request, params...)
			if err == nil {
				codeowners := []*frontmodels.Codeowner{}
				cnv.Convert(resp.Data, &codeowners)
				page.CodebaseData.CodeOwners = codeowners
			}
			wg.Done()
		},
	}
	return
}
