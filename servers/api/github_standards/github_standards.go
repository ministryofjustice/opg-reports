package github_standards

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/db"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request/get"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const listRoute string = "/{version}/github-standards/{$}"

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
func ListHandler(server *api.ApiServer, w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err        error
		ghDB       *sql.DB
		ghRequest  *request.Request
		ghResponse *GHSResponse = NewResponse()
		base       int          = 0
		ext        int          = 0
	)
	response.Start(ghResponse, w, r)

	ghRequest = request.New(
		get.New("team", ""),
		get.WithChoices("archived", []string{"false", "true"}))

	// -- connect to the database
	ghDB, err = db.Connect(server.DbPath)
	if err != nil {
		response.ErrorAndEnd(ghResponse, err, w, r)
		return
	}
	// -- get the query connection
	queries := ghs.New(ghDB)
	defer ghDB.Close()
	defer queries.Close()

	// -- get the db results
	results, err := getResults(server.Ctx, queries, ghRequest.Param("archived", r), ghRequest.Param("team", r))
	if err != nil {
		response.ErrorAndEnd(ghResponse, err, w, r)
		return
	}
	// -- set the results
	ghResponse.Result = results
	// -- get the counters
	base, ext = complianceCounters(results)
	all, _ := queries.Count(server.Ctx)
	tBase, _ := queries.TotalCountCompliantBaseline(server.Ctx)
	tExt, _ := queries.TotalCountCompliantExtended(server.Ctx)
	ghResponse.Counters = &Counters{
		Totals: &CountValues{
			Count:             int(all),
			CompliantBaseline: int(tBase),
			CompliantExtended: int(tExt),
		},
		This: &CountValues{
			Count:             len(results),
			CompliantBaseline: base,
			CompliantExtended: ext,
		},
	}
	// -- max ages
	age, err := queries.Age(server.Ctx)
	if err == nil {
		ghResponse.DataAge = &response.DataAge{Min: age, Max: age}
	}
	ghResponse.QueryFilters = ghRequest.Mapped(r)
	// -- end
	response.End(ghResponse, w, r)
	return
}

func Register(mux *http.ServeMux, apiServer *api.ApiServer) {

	handler := api.Wrap(apiServer, ListHandler)
	mux.HandleFunc(listRoute, mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
}
