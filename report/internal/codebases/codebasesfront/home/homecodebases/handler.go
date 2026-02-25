package homecodebases

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/codebases/codebasesapi/codebaseapiall"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/htmlpage"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tmpl"
	"sync"
)

// codebase model, copy from api side
type codebase struct {
	Name                string `json:"name,omitempty"`                  // short name of codebase (without owner)
	FullName            string `json:"full_name,omitempty" `            // full name including the owner
	Url                 string `json:"url,omitempty" `                  // url to access the codebase
	Visibility          string `json:"visibility,omityempty"`           // visibility status
	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
}

type codebaseData struct {
	Codebases []*codebase
}

type PageContent struct {
	htmlpage.HTMLPage
	CodebaseData *codebaseData
}

type dataCallerF func(wg *sync.WaitGroup, page *PageContent)

// Handler deals with the / root page of the reporting site
func Handler(ctx context.Context, args *Args, writer http.ResponseWriter, request *http.Request) {
	var (
		pageName     string         = "OPG Reports"
		pageTitle    string         = "OPG Reports - Active Codebases"
		templateName string         = "home-codebases"
		log          *slog.Logger   = cntxt.GetLogger(ctx).With("package", "homecodebases", "func", "Handler", "url", request.URL.String())
		wg           sync.WaitGroup = sync.WaitGroup{}
		pgArgs       *htmlpage.Args = &htmlpage.Args{Title: pageTitle, Name: pageName, GovUKVersion: args.GovUKVersion, SemVer: args.SemVer}
		page         *PageContent   = &PageContent{HTMLPage: htmlpage.New(request, pgArgs), CodebaseData: &codebaseData{}}
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
	var params = []*rest.Param{}

	return []dataCallerF{
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
			resp, err := rest.FromApi[*codebaseapiall.Response](ctx, args.ApiHost, codebaseapiall.ENDPOINT, request, params...)
			if err == nil {
				codebases := []*codebase{}
				cnv.Convert(resp.Data, &codebases)
				page.CodebaseData.Codebases = codebases
			}
			wg.Done()
		},
	}
}
