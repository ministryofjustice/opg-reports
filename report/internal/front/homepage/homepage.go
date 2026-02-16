package homepage

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/api"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/front/blocks"
	"opg-reports/report/internal/front/page"
	"opg-reports/report/internal/front/respond"
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tmpl"
	"sync"
)

const templateName string = "index"

// Conf contains values needed for this handler, generalyl comes
// from the original front end Info struct
type Conf struct {
	Name         string
	ApiHost      string `json:"api"`           // api host
	TemplateDir  string `json:"template_dir"`  // --template-dir
	GovUKVersion string `json:"govuk_version"` // --govuk_version
	Signature    string `json:"signature"`     // --signature
}

type apiData struct {
	Data      []map[string]interface{}
	Footer    map[string]interface{}
	LabelCols []string
	DataCols  []string
	ExtraCols []string
	EndCols   []string
}

type codeownerData struct {
	Data []*codeownermodels.CodeownerData
}

type PageContent struct {
	page.PageContent
	Teams       []string // list of team names
	Uptime      *apiData
	Infracosts  *apiData
	UnownedCode *codeownerData
}

// blocks
type pageBlockF func(i ...any)

// handleHomepage renders the request for `/` which currently displays:
func handler(ctx context.Context, log *slog.Logger, conf *Conf, writer http.ResponseWriter, request *http.Request) {

	var (
		blockFs []pageBlockF

		wg        sync.WaitGroup = sync.WaitGroup{}
		templates []string       = tmpl.GetTemplateFiles(conf.TemplateDir)
		lg        *slog.Logger   = log.With("func", "homepage.handler")
		data      *PageContent   = &PageContent{PageContent: page.NewContent(request, &page.PageInfo{
			Name:         conf.Name,
			GovUKVersion: conf.GovUKVersion,
			Signature:    conf.Signature,
		})}
	)
	lg.Info("processing page ...", "url", request.URL.String())

	blockFs = []pageBlockF{
		// get list of teams
		func(i ...any) {
			if teams, err := blocks.TeamNavData(ctx, log, conf.ApiHost, request); err == nil {
				data.Teams = teams
			}
			wg.Done()
		},
		// get uptime data
		func(i ...any) {
			// lock in team filter = true
			var uptime, headings, err = blocks.UptimeData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "sort", Value: "team", Locked: true},
			)

			if err == nil && len(uptime) > 0 {
				l := len(uptime)
				data.Uptime = &apiData{
					Data:      uptime[0 : l-1],
					Footer:    uptime[l-1],
					LabelCols: headings[string(headers.KEY)],
					DataCols:  headings[string(headers.DATA)],
					ExtraCols: headings[string(headers.EXTRA)],
					EndCols:   headings[string(headers.END)],
				}
			}
			wg.Done()
		},
		// get cost data
		func(i ...any) {
			// lock in team filter = true
			var costs, headings, err = blocks.InfracostData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "sort", Value: "team", Locked: true},
			)
			if err == nil && len(costs) > 0 {
				l := len(costs)
				data.Infracosts = &apiData{
					Data:      costs[0 : l-1],
					Footer:    costs[l-1],
					LabelCols: headings[string(headers.KEY)],
					DataCols:  headings[string(headers.DATA)],
					ExtraCols: headings[string(headers.EXTRA)],
					EndCols:   headings[string(headers.END)],
				}
			}
			wg.Done()
		},
		// get codeowners
		func(i ...any) {
			var code, err = blocks.CodeownerData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "codeowner", Value: "none", Locked: true},
			)
			if err == nil && len(code) > 0 {
				data.UnownedCode = &codeownerData{Data: code}
			}
			wg.Done()
		},
	}

	for _, blockF := range blockFs {
		wg.Add(1)
		go blockF()
	}
	wg.Wait()
	lg.Info("complete.", "url", request.URL.String())
	respond.Respond(log, writer, request, templateName, templates, data)
}

// registerHomepage is called from rootCmd.RunE for endpoint
// handling delegation
//
// maps `/` to the `handleHomepage` function
func Register(
	ctx context.Context,
	log *slog.Logger,
	mux *http.ServeMux,
	info *Conf,
) {
	log.Info("registering handler [`/{$}`] ...")
	// Homepage
	mux.HandleFunc("/{$}", func(writer http.ResponseWriter, request *http.Request) {
		handler(ctx, log, info, writer, request)
	})
}
