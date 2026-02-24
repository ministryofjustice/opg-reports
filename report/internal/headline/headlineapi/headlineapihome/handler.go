// Package headlineapihome is for providing headline figures in one api call.
//
// Used on home page to provide some top line numbers at a glance for ease
package headlineapihome

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
	"opg-reports/report/package/times"
	"time"
)

// total cost within the months requested
const costSelect string = `
SELECT
	CAST(COALESCE(SUM(cost), 0) as REAL) as cost
FROM costs
WHERE
	costs.service != 'Tax'
	AND costs.month IN (:months)
;
`

// average uptime between all services
const uptimeSelect string = `
SELECT
	CAST(COALESCE(AVG(uptime.average), 0) as REAL) as uptime
FROM uptime
WHERE
	uptime.month IN (:months)
`

// Request contains url path / query values
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
	Data    *Result  `json:"data"`
	Months  []string `json:"-"`
}

// Filter is with the sql to replace the `:name` named parameters within the
// statement.
type Filter struct {
	Months []string `json:"months"`
}

type Result struct {
	TotalCost           float64 `json:"total_cost"`             // total cost result
	AverageCostPerMonth float64 `json:"average_cost_per_month"` // average cost per month

	OverallUptime float64 `json:"overall_uptime"` // overall uptime in time period
}

// Responder process the incoming request, queries the database and returns the result as json data.
func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		months   []string
		res      *Result                = &Result{}
		filter   *Filter                = &Filter{}
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "headlineapi", "func", "Responder")
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
	// setup month filter
	filter = &Filter{Months: months}
	// now convert to a map for use in bound statements
	err = cnv.Convert(filter, &bindMap)
	if err != nil {
		log.Error("failed to convert filter into map for binding", "err", err.Error())
		return
	}
	// run cost databases call
	costSelectRun(ctx, conf, filter, bindMap, res)
	// run uptime database call
	uptimeSelectRun(ctx, conf, filter, bindMap, res)

	response = &Response{
		Version: conf.Version,
		SHA:     conf.SHA,
		Request: in,
		Data:    res,
	}
	log.Info("complete.")
	respond.AsJSON(ctx, request, writer, response)

}

// costSelectRun runs the cost select and fetches the val for total cost
func costSelectRun(ctx context.Context, conf *Config, filter *Filter, bindMap map[string]interface{}, res *Result) *Result {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "headlineapi", "func", "costSelectRun")

	dbx.Select(ctx, costSelect, &dbx.SelectArgs{
		DB:      conf.DB,
		Driver:  conf.Driver,
		Params:  conf.Params,
		BindMap: bindMap,
		ScanF: func(rows *sql.Rows) error {
			var err error
			var val float64 = 0.0
			if err = rows.Scan(&val); err == nil {
				res.TotalCost = val
			} else {
				log.Error("row scan failed", "err", err.Error())
			}
			return err
		},
	})
	// work out averages
	res.AverageCostPerMonth = (res.TotalCost / float64(len(filter.Months)))
	return res
}

// uptimeSelectRun runs the uptime select and fetches the val for average uptime
func uptimeSelectRun(ctx context.Context, conf *Config, filter *Filter, bindMap map[string]interface{}, res *Result) *Result {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "headlineapi", "func", "uptimeSelectRun")

	dbx.Select(ctx, uptimeSelect, &dbx.SelectArgs{
		DB:      conf.DB,
		Driver:  conf.Driver,
		Params:  conf.Params,
		BindMap: bindMap,
		ScanF: func(rows *sql.Rows) error {
			var err error
			var val float64 = 0.0
			if err = rows.Scan(&val); err == nil {
				res.OverallUptime = val
			} else {
				log.Error("row scan failed", "err", err.Error())
			}
			return err
		},
	})
	return res
}
