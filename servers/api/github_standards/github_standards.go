package github_standards

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/query"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const listRoute string = "/{version}/github-standards/{$}"

var (
	apiCtx    context.Context
	apiDbPath string
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string            = apiDbPath
		ctx      context.Context   = apiCtx
		response *resp.Response    = resp.New()
		archived *query.Query      = query.Get("archived")
		team     *query.Query      = query.Get("team")
		filters  map[string]string = map[string]string{
			"archived": query.FirstD(archived.Values(r), "false"),
			"team":     query.First(team.Values(r)),
		}
	)

	response.Start(w, r)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		response.ErrorAndEnd(err, w, r)
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
		response.ErrorAndEnd(err, w, r)
		return
	}
	// -- convert results over to output format
	res, base, ext := resultsOut(results, response)
	response.Result = res
	// -- get overall counters
	all, _ := queries.Count(ctx)
	tBase, _ := queries.TotalCountCompliantBaseline(ctx)
	tExt, _ := queries.TotalCountCompliantExtended(ctx)

	response.Metadata["counters"] = map[string]map[string]int{
		"totals": {
			"count":              int(all),
			"compliant_baseline": int(tBase),
			"compliant_extended": int(tExt),
		},
		"this": {
			"count":              len(res),
			"compliant_baseline": base,
			"compliant_extended": ext,
		},
	}
	response.Metadata["filters"] = filters
	// -- add the date min / max values
	age, err := queries.Age(ctx)
	if err == nil {
		response.DataAge.Min = age
		response.DataAge.Max = age
	}
	response.End(w, r)
	return

}

// Register attached the route to the list view
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
