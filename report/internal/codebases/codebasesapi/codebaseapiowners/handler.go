package codebaseapiowners

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/global/apimodels"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/requested"
	"opg-reports/report/package/respond"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// selectStmt is the sql used to fetch data including
const selectStmt string = `
SELECT
	codebases.full_name,
	codebases.name,
	codebases.url,
	json_group_array(
		DISTINCT json_object(
			'owner', codebase_owners.owner,
			'team_name', codebase_owners.team_name
		)
	) filter ( where codebase_owners.owner is not null) as owners
FROM codebases
LEFT JOIN codebase_owners ON codebase_owners.codebase = codebases.full_name
WHERE
	codebases.archived = 0
GROUP BY
	codebases.full_name
ORDER BY
	codebase_owners.team_name = "none" DESC,
	codebase_owners.team_name ASC,
	codebases.full_name ASC
;
`

// Request contains the url path / query string values that we will use
// in this handler
type Request struct {
	Team string `json:"team"` // option team filter for this handler
}

// Response is the end result thats sent back from the handler via the writter
type Response struct {
	Version string   `json:"version"`
	SHA     string   `json:"sha"`
	Request *Request `json:"request"`
	Data    []*Model `json:"data"` // the actual data results
}

// Filter is with the sql to replace the named parameters
// within the statement.
type Filter struct {
	Team string `json:"team"`
}

// Model is the data struct to use when fetching the select
type Model struct {
	FullName string            `json:"full_name,omitempty" ` // full name including the owner
	Name     string            `json:"name,omitempty"`       // short name of codebase (without owner)
	Url      string            `json:"url,omitempty" `       // url to access the codebase
	Owners   hasManyCodeowners `json:"owners"`               // list of codeowners
}

// Sequence is used to return the columns in the order they are selected
func (self *Model) Sequence() []any {
	self.Owners = hasManyCodeowners{}
	return []any{
		&self.FullName,
		&self.Name,
		&self.Url,
		&self.Owners,
	}
}

// joined teams is the codebase -> codeowners data
type joinedCodeowner struct {
	Owner    string `json:"owner"`
	TeamName string `json:"team_name"`
}

// hasManyCodeowners represents the join
// Interfaces:
//   - sql.Scanner
type hasManyCodeowners []*joinedCodeowner

// Scan handles the processing of the join data
func (self *hasManyCodeowners) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = json.Unmarshal(src.([]byte), self)
	case string:
		err = json.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// Responder process the incoming request, queries the database and returns the result as json data.
func Responder(ctx context.Context, conf *apimodels.Args, request *http.Request, writer http.ResponseWriter) {
	var (
		err      error
		response *Response
		filter   *Filter                = &Filter{}
		in       *Request               = &Request{}
		bindMap  map[string]interface{} = map[string]interface{}{}
		all      []*Model               = []*Model{}
		log      *slog.Logger           = cntxt.GetLogger(ctx).With("package", "codebaseapiowners", "func", "Responder")
		stmt     string                 = selectStmt // localised constant
	)
	log.Info("running http handler ...")
	// convert the http request into Request struct
	requested.Parse(ctx, request, &in)
	// look for the optional team
	if in.Team != "" {
		log.Info("optional team filter found ...", "team", in.Team)
		filter.Team = in.Team
		stmt = strings.ReplaceAll(stmt, "WHERE", "WHERE codebase_owners.team_name = :team AND")
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
