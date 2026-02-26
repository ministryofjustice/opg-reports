package costapiteam

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
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// selectStmt is the sql used to fetch data including
// and params (`:name`) that will be replaced by values
// from `Request` (by configuring `Filter`)
const selectStmt string = `
SELECT
	costs.month as month,
	CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
	IIF(accounts.team_name != "", accounts.team_name, "")  as team
FROM costs
LEFT JOIN accounts on accounts.id = costs.account_id
WHERE
	costs.service != 'Tax'
	AND costs.month IN (:months)
GROUP BY
	costs.month,
	accounts.team_name
ORDER BY
	accounts.team_name ASC
;
`

// Request contains the url path / query string values that we will use
// in this handler
type Request struct {
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Team      string `json:"team"`
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
	Version string                        `json:"version"`
	SHA     string                        `json:"sha"`
	Request *Request                      `json:"request"`
	Headers map[tabulate.ColType][]string `json:"headers"` // headers contains details for table headers / rendering
	Data    []map[string]interface{}      `json:"data"`    // the actual data results
	Summary map[string]interface{}        `json:"summary"` // used to contain table totals etc

}

// Filter is with the sql to replace the `:name` named parameters within the
// statement.
// For this endpointm, we only filter by the time period - months
type Filter struct {
	Months []string `json:"months"`
	Team   string   `json:"team"`
}

// Model is the data struct to use when fetching the select
type Model struct {
	Month string  `json:"month"`
	Cost  float64 `json:"cost"`
	Team  string  `json:"team"`
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	return []any{
		&self.Month, &self.Cost, &self.Team,
	}
}

// Responder process the incoming request, queries the database and returns the result as json data.
//
// Data is formatted as a table for easier display.
func Responder(ctx context.Context, conf *apimodels.Args, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter
		months   []string
		in       *Request                      = &Request{}
		bindMap  map[string]interface{}        = map[string]interface{}{}
		all      []*Model                      = []*Model{}
		log      *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "costapiteam", "func", "Responder")
		stmt     string                        = selectStmt
		headings map[tabulate.ColType][]string = map[tabulate.ColType][]string{
			tabulate.KEY:   {"team"},
			tabulate.EXTRA: {"trend"},
			tabulate.END:   {"total"},
		}
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
	// setup months
	headings[tabulate.DATA] = months
	filter = &Filter{Months: months}
	// look for the optional team
	if in.Team != "" {
		log.Info("optional team filter found ...", "team", in.Team)
		filter.Team = in.Team
		stmt = strings.ReplaceAll(stmt, "WHERE", "WHERE accounts.team_name = :team AND")
	}
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
	// get the body
	tableBody := tabulate.TableBody(ctx, all, &tabulate.Args{
		Headers:   headings,
		ColumnKey: "month",
		ValueKey:  "cost"})
	// add rowTotals
	tabulate.RowEnd(tableBody, headings, tabulate.RowTotalF)
	// swap to slice
	tbl := tabulate.TableMapToTable(tableBody)
	// table sort
	tbl = tabulate.SortAscending[string](tbl, "team")
	// do table total
	summary := tabulate.TableEnd(tbl, headings, tabulate.TableTotalF)

	// setup response object
	response = &Response{
		Version: conf.Version,
		SHA:     conf.SHA,
		Request: in,
		Headers: headings,
		Data:    tbl,
		Summary: summary,
	}
	log.Info("complete.")
	respond.AsJSON(ctx, request, writer, response)
}
