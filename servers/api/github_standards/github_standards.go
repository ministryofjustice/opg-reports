package github_standards

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/query"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

// listRoute is the url this handler supports
const listRoute string = "/{version}/github-standards/{$}"

var (
	apiCtx    context.Context
	apiDbPath string
)

// ListHandler is configure to handler the `listRoute` url requests
// and will report a GHSResponse.
//
//   - Connects to sql db via `apiDbPath`
//   - Gets result data dependant on the `archived` and `team` get parameters
//   - Generates compliance counters for the data set and overall
//   - Finds the run date of the report
//
// Sample urls:
//   - /v1/github-standards/
//   - /v1/github-standards/?archived=false
//   - /v1/github-standards/?archived=false&team=<team>
func ListHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err        error
		db         *sql.DB
		dbPath     string            = apiDbPath
		ctx        context.Context   = apiCtx
		ghResponse *GHSResponse      = NewResponse()
		archived   *query.Query      = query.Get("archived")
		team       *query.Query      = query.Get("team")
		filters    map[string]string = map[string]string{
			"archived": query.FirstD(archived.Values(r), "false"),
			"team":     query.First(team.Values(r)),
		}
	)
	response.Start(ghResponse, w, r)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		response.ErrorAndEnd(ghResponse, err, w, r)
		return
	}
	defer db.Close()
	// -- setup the sqlc generated queries
	queries := ghs.New(db)
	defer queries.Close()
	// -- fetch the raw results
	slog.Info("about to get results")
	results, err := getResults(ctx, queries, filters["archived"], filters["team"])
	slog.Info("got results")
	if err != nil {
		response.ErrorAndEnd(ghResponse, err, w, r)
		return
	}
	// -- convert results over to output format
	base, ext := complianceCounters(results)
	ghResponse.Result = results
	// -- get overall counters
	all, _ := queries.Count(ctx)
	tBase, _ := queries.TotalCountCompliantBaseline(ctx)
	tExt, _ := queries.TotalCountCompliantExtended(ctx)

	ghResponse.Counters.Totals.Count = int(all)
	ghResponse.Counters.This.CompliantBaseline = int(tBase)
	ghResponse.Counters.Totals.CompliantExtended = int(tExt)
	ghResponse.Counters.This.Count = len(results)
	ghResponse.Counters.This.CompliantBaseline = base
	ghResponse.Counters.This.CompliantExtended = ext
	ghResponse.QueryFilters = filters

	// -- add the date min / max values
	age, err := queries.Age(ctx)
	if err == nil {
		ghResponse.DataAge.Min = age
		ghResponse.DataAge.Max = age
	}
	response.End(ghResponse, w, r)

	return

}

// Register sets the local context and database paths to the values passed and then
// attaches the local handles to the url patterns supported by the api
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	SetCtx(ctx)
	SetDBPath(dbPath)

	// -- actually register the handler
	mux.HandleFunc(listRoute, mw.Middleware(ListHandler, mw.Logging, mw.SecurityHeaders))

	return nil
}

func SetDBPath(path string) {
	apiDbPath = path
}
func SetCtx(ctx context.Context) {
	apiCtx = ctx
}
