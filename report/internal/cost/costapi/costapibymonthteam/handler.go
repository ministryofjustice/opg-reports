package costapibymonthteam

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/dump"
	"opg-reports/report/package/queryx"
	"opg-reports/report/package/requested"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/times"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
	costs.month
ORDER BY
	accounts.team_name ASC
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
	Request *Request                 `json:"request"`
	Headers map[string][]string      `json:"headers"` // headers contains details for table headers / rendering
	Data    []map[string]interface{} `json:"data"`    // the actual data results
}

// Filter is with the sql to replace the `:name` named parameters within the
// statement
type Filter struct {
	Months []string `json:"months"`
}

// Model is the data struct to use when fetching the select
type Model struct {
	Month string `json:"month"`
	Cost  string `json:"cost"`
	Team  string `json:"team"`
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	return []any{
		&self.Month, &self.Cost, &self.Team,
	}
}

func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter
		months   []string
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		all      []*Model               = []*Model{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "costapibymonthteam", "func", "Responder")
	)
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// get months between dates
	months = times.AsYMStrings(times.Months(in.Start(), in.End()))
	if len(months) <= 0 {
		log.Error("no months found with date range provided", "err", err.Error())
		return
	}
	filter = &Filter{Months: months}
	// now convert to a map for use in bound statements
	err = cnv.Convert(filter, &bindMap)
	if err != nil {
		log.Error("failed to convert filter into map for binding", "err", err.Error())
		return
	}
	// make the db call via the Select helper that handles row scanning
	// all, err := get(ctx, conf, filter)
	queryx.Select(ctx, selectStmt, &queryx.Input{
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

	fmt.Println(dump.Any(all))

	// setup response object
	response = &Response{
		Request: in,
		Headers: map[string][]string{
			"labels": {"team"},
			"extra":  {"trend"},
			"end":    {"total"},
			"data":   months,
		},
	}
	respond.AsJSON(ctx, request, writer, response)
}
