package teampage

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

const templateName string = "team"

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
	TeamName           string
	Teams              []string // list of team names
	Uptime             *apiData
	Infracosts         *apiData
	DetailedInfracosts *apiData
	InfracostDiffs     *apiData
	CodeownerData      *codeownerData
}

// blocks
type pageBlockF func(i ...any)

// handleHomepage renders the request for `/` which currently displays:
func handler(ctx context.Context, log *slog.Logger, conf *Conf, writer http.ResponseWriter, request *http.Request) {

	var (
		blockFs []pageBlockF

		wg        sync.WaitGroup = sync.WaitGroup{}
		templates []string       = tmpl.GetTemplateFiles(conf.TemplateDir)
		lg        *slog.Logger   = log.With("func", "teampage.handler")
		data      *PageContent   = &PageContent{
			PageContent: page.NewContent(request, &page.PageInfo{
				Name:         conf.Name,
				GovUKVersion: conf.GovUKVersion,
				Signature:    conf.Signature,
			}),
			TeamName: request.PathValue("team"),
		}
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
		// get codeowners
		func(i ...any) {
			// filte rby this team
			var code, err = blocks.CodeownerData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: data.TeamName, Locked: true},
			)
			if err == nil && len(code) > 0 {
				data.CodeownerData = &codeownerData{Data: code}
			}
			wg.Done()
		},
		// get cost data
		func(i ...any) {
			// lock in account grouping and team filter
			var costs, headings, err = blocks.InfracostData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: data.TeamName, Locked: true},
				&api.Param{Type: api.QUERY, Key: "account", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "sort", Value: "account", Locked: true},
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
		// get detailed cost data
		func(i ...any) {
			// lock in detailed grouping
			var costs, headings, err = blocks.InfracostData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: data.TeamName, Locked: true},
				&api.Param{Type: api.QUERY, Key: "account", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "environment", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "service", Value: "true", Locked: true},
				&api.Param{Type: api.QUERY, Key: "sort", Value: "cost", Locked: true},
			)
			if err == nil && len(costs) > 0 {
				l := len(costs)
				data.DetailedInfracosts = &apiData{
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
		// get differences
		func(i ...any) {
			// lock in detailed grouping
			var costs, headings, err = blocks.InfracostDiffData(ctx, log, conf.ApiHost, request,
				&api.Param{Type: api.QUERY, Key: "team", Value: data.TeamName, Locked: true},
			)
			if err == nil && len(costs) > 0 {
				data.InfracostDiffs = &apiData{
					Data:      costs,
					LabelCols: headings[string(headers.KEY)],
					DataCols:  headings[string(headers.DATA)],
					ExtraCols: headings[string(headers.EXTRA)],
					EndCols:   headings[string(headers.END)],
				}
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

// Register is called from rootCmd.RunE for endpoint
// handling delegation
func Register(
	ctx context.Context,
	log *slog.Logger,
	mux *http.ServeMux,
	info *Conf,
) {
	log.Info("registering handler [`/team/{team}`] ...")
	mux.HandleFunc("/team/{team}/{$}", func(writer http.ResponseWriter, request *http.Request) {
		handler(ctx, log, info, writer, request)
	})
}
