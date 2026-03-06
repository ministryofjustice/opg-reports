package codebasereleasesapi

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/global/apimodels"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/requested"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/times"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// selectStmt is the sql used to fetch data including
const selectStmt string = `
SELECT
	codebase_metrics.month,
	COALESCE(SUM(codebase_metrics.releases),0) as releases,
	COALESCE(SUM(codebase_metrics.releases_securityish),0) as releases_securityish,
	CAST(COALESCE(AVG(codebase_metrics.releases_average_time), 0) as REAL) as releases_average_time
FROM codebase_metrics
LEFT JOIN codebases on codebases.full_name = codebase_metrics.codebase
WHERE
	codebases.archived = 0
	AND codebase_metrics.month IN (:months)
GROUP BY
	codebase_metrics.month
ORDER BY
	codebase_metrics.month DESC
;
`

// Request contains the url path / query string values that we will use
// in this handler
type Request struct {
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
}

func (self *Request) Start() (t time.Time) {
	t = times.MustFromString(self.DateStart)
	return
}
func (self *Request) End() (t time.Time) {
	t = times.MustFromString(self.DateEnd)
	return
}

// Response is the end result thats sent back from the handler via the writter
type Response struct {
	Version string   `json:"version"`
	SHA     string   `json:"sha"`
	Request *Request `json:"request"`
	Data    []*Model `json:"data"`    // the actual data results
	Summary *Model   `json:"summary"` // overall totals of each
}

// Filter is with the sql to replace the named parameters
// within the statement.
type Filter struct {
	Months []string `json:"months"`
}

// Model is the data struct to use when fetching the select
type Model struct {
	Month               string  `json:"month"`                 // month as YYYY-MM string
	Releases            int     `json:"releases"`              // count of releases for this month
	ReleasesSecurityish int     `json:"releases_securityish"`  // count of releases for this month that seem to be security related
	ReleasesAverageTime float64 `json:"releases_average_time"` // average time path to live workflow took (in milliseconds)
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	return []any{
		&self.Month,
		&self.Releases,
		&self.ReleasesSecurityish,
		&self.ReleasesAverageTime,
	}
}

// Responder process the incoming request, queries the database and returns the result as json data.
func Responder(ctx context.Context, conf *apimodels.Args, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		months   []string
		filter   *Filter                = &Filter{}
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		all      []*Model               = []*Model{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "codebasereleasesapi", "func", "Responder")
		stmt     string                 = selectStmt // localised constant

	)
	log.Info("running http handler ...")
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// get months between dates
	months = times.AsYMStrings(times.Months(in.Start(), in.End()))
	if len(months) <= 0 {
		log.Error("no months found with date range provided")
		return
	}
	filter.Months = months
	// now convert to a map for use in bound statements
	err = cnv.Convert(filter, &bindMap)
	if err != nil {
		log.Error("failed to convert filter into map for binding", "err", err.Error())
		return
	}
	// make the db call via the Select helper that handles row scanning.
	// No return value as local values are updates within ScanF lambda
	dbx.Select(ctx, stmt, &dbx.SelectArgs{
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
	summary := &Model{
		Releases:            0,
		ReleasesSecurityish: 0,
		ReleasesAverageTime: 0.0,
	}
	for _, r := range all {
		summary.Releases += r.Releases
		summary.ReleasesSecurityish += r.ReleasesSecurityish
		summary.ReleasesAverageTime += r.ReleasesAverageTime
	}

	// setup response object
	response = &Response{
		Version: conf.Version,
		SHA:     conf.SHA,
		Request: in,
		Data:    all,
		Summary: summary,
	}
	log.Info("complete.")
	respond.AsJSON(ctx, request, writer, response)
}
