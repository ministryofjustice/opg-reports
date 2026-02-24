package costapidiff

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
	"opg-reports/report/package/tabulate"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// selectStmt is the sql used to fetch data including
// and params (`:name`) that will be replaced by values
// from `Request` (by configuring `Filter`)
const selectStmt string = `
SELECT
	costs.month as month,
	CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
	costs.service as service,
	IIF(accounts.team_name != "", accounts.team_name, "")  as team,
	IIF(accounts.name != "", accounts.name, "") as account
FROM costs
LEFT JOIN accounts on accounts.id = costs.account_id
WHERE
	costs.service != 'Tax'
	AND costs.month IN (:months)
GROUP BY
	costs.month,
	costs.service,
	accounts.team_name,
	accounts.id
ORDER BY
	accounts.team_name,
	costs.service ASC
;
`

// Request contains the url path / query string values that we will use
// in this handler
type Request struct {
	DateA  string `json:"date_a"`
	DateB  string `json:"date_b"`
	Change string `json:"change"`
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
}

// Model is the data struct to use when fetching the select
type Model struct {
	Month   string  `json:"month"`
	Cost    float64 `json:"cost"`
	Service string  `json:"service"`
	Team    string  `json:"team"`
	Account string  `json:"account"`
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	return []any{
		&self.Month, &self.Cost, &self.Service, &self.Team, &self.Account,
	}
}

// Responder process the incoming request, queries the database and returns the result as json data.
//
// Data is formatted as a table for easier display.
func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter
		months   []string
		change   float64                       = 300
		in       *Request                      = &Request{}
		bindMap  map[string]interface{}        = map[string]interface{}{}
		all      []*Model                      = []*Model{}
		log      *slog.Logger                  = cntxt.GetLogger(ctx).With("package", "costapidiff", "func", "Responder")
		headings map[tabulate.ColType][]string = map[tabulate.ColType][]string{
			tabulate.KEY:   {"team", "account", "service"},
			tabulate.EXTRA: {},
			tabulate.END:   {"diff"},
		}
	)
	log.Info("running http handler ...")
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// get months between dates
	months = []string{in.DateA, in.DateB}
	if len(months) <= 0 {
		log.Error("no months found with date range provided")
		return
	}
	// setup months
	headings[tabulate.DATA] = months
	filter = &Filter{Months: months}
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

	// get the body
	tableBody := tabulate.TableBody(ctx, all, &tabulate.Args{
		Headers:   headings,
		ColumnKey: "month",
		ValueKey:  "cost"})
	// add row diffs
	tabulate.RowEnd(tableBody, headings, tabulate.RowDiffF)
	// swap to slice
	tbl := tabulate.TableMapToTable(tableBody)
	// filter table by min diff
	// parse the change value
	if f, e := strconv.ParseFloat(in.Change, 64); e == nil {
		change = f
	}
	tbl = tabulate.TableFilterByValue(tbl, headings, change)
	// sort
	tbl = tabulate.SortDescending[float64](tbl, "diff")

	// setup response object
	response = &Response{
		Version: conf.Version,
		SHA:     conf.SHA,
		Request: in,
		Headers: headings,
		Data:    tbl,
	}
	log.Info("complete.")
	respond.AsJSON(ctx, request, writer, response)
}
