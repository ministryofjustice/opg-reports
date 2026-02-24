package codebaseapiall

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/requested"
	"opg-reports/report/package/respond"

	_ "github.com/mattn/go-sqlite3"
)

// selectStmt is the sql used to fetch data including
const selectStmt string = `
SELECT
	name,
	full_name,
	url,
	compliance_level,
	compliance_report_url,
	compliance_badge
FROM codebases
ORDER BY
	name ASC
;
`

// Request contains the url path / query string values that we will use
// in this handler
type Request struct{}

// Response is the end result thats sent back from the handler via the writter
type Response struct {
	Version string   `json:"version"`
	SHA     string   `json:"sha"`
	Request *Request `json:"request"`
	Data    []*Model `json:"data"` // the actual data results
}

// Filter is with the sql to replace the `:name` named parameters within the
// statement. empty for this endpoint
type Filter struct{}

// Model is the data struct to use when fetching the select
type Model struct {
	Name                string `json:"name,omitempty"`                  // short name of codebase (without owner)
	FullName            string `json:"full_name,omitempty" `            // full name including the owner
	Url                 string `json:"url,omitempty" `                  // url to access the codebase
	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	return []any{
		&self.Name,
		&self.FullName,
		&self.Url,
		&self.ComplianceLevel,
		&self.ComplianceReportUrl,
		&self.ComplianceBadge,
	}
}

// Responder process the incoming request, queries the database and returns the result as json data.
func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter                = &Filter{}
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		all      []*Model               = []*Model{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "codebaseapiall", "func", "Responder")
	)
	log.Info("running http handler ...")
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// now convert to a map for use in bound statements
	err = cnv.Convert(filter, &bindMap)
	if err != nil {
		log.Error("failed to convert filter into map for binding", "err", err.Error())
		return
	}
	// make the db call via the Select helper that handles row scanning.
	// No return value as local values are updates within ScanF lambda
	dbx.Select(ctx, selectStmt, &dbx.SelectArgs{
		DB:      conf.DB,
		Driver:  conf.Driver,
		Params:  conf.Params,
		BindMap: bindMap,
		ScanF: func(rows *sql.Rows) error {
			var r = &Model{}
			var seq = r.Sequence()
			if err = rows.Scan(seq...); err == nil {
				all = append(all, r)
			} else {
				log.Error("row scan failed", "err", err.Error())
			}
			return err
		},
	})

	// setup response object
	response = &Response{
		Version: conf.Version,
		SHA:     conf.SHA,
		Request: in,
		Data:    all,
	}
	log.Info("complete.")
	respond.AsJSON(ctx, request, writer, response)
}
