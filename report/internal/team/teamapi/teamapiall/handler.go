package teamapiall

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
	name
FROM teams
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
	Data    []*Team  `json:"data"` // the actual data results
}

// Filter is with the sql to replace the `:name` named parameters within the
// statement. empty for this endpoint
type Filter struct{}

// Model is the data struct to use when fetching the select
type Team struct {
	Name string `json:"name"`
}

// Sequence is used to return the columns in the order they are selected
func (self *Team) Sequence() []any {
	return []any{&self.Name}
}

// Responder process the incoming request, queries the database and returns the result as json data.
//
// Data is formatted as a table for easier display.
func Responder(ctx context.Context, conf *Config, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter                = &Filter{}
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		all      []*Team                = []*Team{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "teamapiall", "func", "Responder")
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
			var r = &Team{}
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
