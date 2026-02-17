package costapibymonthteam

import (
	"context"
	"net/http"
	"opg-reports/report/package/requested"
	"opg-reports/report/package/respond"
	"opg-reports/report/package/times"
	"time"
)

const selectStmt string = `
SELECT
	costs.month as month,
	CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
	accounts.team_name as team
FROM costs
LEFT JOIN accounts on accounts.id = costs.account_id
WHERE
	infracosts.service != 'Tax'
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

// Sequence provides the columns matching the order of the select
// for use with rows.Scan()
func (self *Model) Sequence() (columns []*string) {
	return []*string{
		&self.Month, &self.Cost, &self.Team,
	}
}

func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		response *Response
		months   []string
		in       *Request = &Request{}
	)
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// get months between dates
	months = times.AsYMStrings(times.Months(in.Start(), in.End()))

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

func get(ctx context.Context, conf *Config) {

}
