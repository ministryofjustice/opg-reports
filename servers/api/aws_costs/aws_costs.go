package aws_costs

import (
	"context"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const overallTotals string = "/aws-costs/{version}/totals/{start}/{end}/{$}"

// Handlers are all the handler functions, returns a map of funcs
//   - do it this way so multiple routes can share functions and the handler can then
//     get some of the vars set in the scope of this function (like query params etc)
func Handlers(ctx context.Context, mux *http.ServeMux, dbPath string) map[string]func(w http.ResponseWriter, r *http.Request) {
	// -- handler functions
	var overallTotals = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var response = resp.New()
		var filters = map[string]string{}
		response.Start(w, r)

		// -- setup db connection
		db, err := apidb.SqlDB(dbPath)
		defer db.Close()

		if err != nil {
			slog.Error("api db error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}
		// -- setup the sqlc generated queries
		queries := awsc.New(db)
		defer queries.Close()

		// -- fetch the raw results
		slog.Info("about to get results")
		results, err := queries.MonthlyTotalsTaxSplit(ctx)
		slog.Info("got results")
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		// -- convert results over to output format
		response.Result, _ = convert.Maps(results)

		// // -- get overall counters
		// all, _ := queries.Count(ctx)
		// tBase, _ := queries.TotalCountCompliantBaseline(ctx)
		// tExt, _ := queries.TotalCountCompliantExtended(ctx)

		// response.Metadata["counters"] = map[string]map[string]int{
		// 	"totals": {
		// 		"count":              int(all),
		// 		"compliant_baseline": int(tBase),
		// 		"compliant_extended": int(tExt),
		// 	},
		// 	"this": {
		// 		"count":              len(res),
		// 		"compliant_baseline": base,
		// 		"compliant_extended": ext,
		// 	},
		// }
		response.Metadata["filters"] = filters
		// -- add the date min / max values
		min, err := queries.Oldest(ctx)
		max, err := queries.Youngest(ctx)
		if err == nil {
			response.DataAge.Min = min
			response.DataAge.Max = max
		}
		response.End(w, r)
		return
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		"overallTotals": overallTotals,
	}
}

// Register attached the route to the list view
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	funcs := Handlers(ctx, mux, dbPath)
	totals := funcs["overallTotals"]
	// -- actually register the handler
	mux.HandleFunc(overallTotals, mw.Middleware(totals, mw.Logging, mw.SecurityHeaders))

	return nil
}
