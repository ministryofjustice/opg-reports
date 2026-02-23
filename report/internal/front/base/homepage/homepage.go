package homepage

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/team/teamget"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/tmpl"
	"sync"
)

type PageContent struct {
	htmlpage.HTMLPage
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		// err  error
		pageName     string         = "OPG Reports"
		templateName string         = "home"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "home", "func", "handlerIndex", "url", request.URL.String())
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
			if teams, err := teamget.NavigationData(ctx, args.ApiHost, request); err == nil {
				page.Teams = teams
			}
			wg.Done()
		},
	}
}
