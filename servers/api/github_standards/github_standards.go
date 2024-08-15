package github_standards

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/query"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const dbMode string = "WAL"
const dbTimeout int = 50000

const listRoute string = "/github/standards/{version}/list/{$}"

// connection to the db
func sqlDB(dbPath string) (db *sql.DB, err error) {
	if exists.FileOrDir(dbPath) != true {
		err = fmt.Errorf("database [%s] does not exist", dbPath)
		return
	}
	// connection string to set modes etc
	conn := consts.SQL_CONNECTION_PARAMS
	slog.Debug("connecting to db...", slog.String("dbPath", dbPath), slog.String("connection", conn))

	// try to connect to db
	db, err = sql.Open("sqlite3", dbPath+conn)
	if err != nil {
		slog.Error("error connecting to DB", slog.String("err", err.Error()))
		return
	}
	slog.Info("connected to db", slog.String("dbPath", dbPath), slog.String("connection", conn))
	return
}

// getResults handles determining which query to call based on the get param values
func getResults(ctx context.Context, queries *ghs.Queries, archived string, team string) (results []ghs.GithubStandard, err error) {

	var teamF = ""
	var archivedF = ""
	results = []ghs.GithubStandard{}
	// -- fetch the get parameter values
	// team query, add the like logic here
	if team != "" {
		teamF = "%#" + team + "#%"
	}
	// archive query
	if archived != "" {
		archivedF = archived
	}
	// -- run queries
	if teamF != "" && archivedF != "" {
		// if both team and archive are set, use joined query
		results, err = queries.FilterByIsArchivedAndTeam(ctx, ghs.FilterByIsArchivedAndTeamParams{
			IsArchived: convert.BoolStringToInt(archivedF), Teams: teamF,
		})
	} else if archivedF != "" {
		// run for just archived - this is defaulted to 1
		results, err = queries.FilterByIsArchived(ctx, convert.BoolStringToInt(archivedF))
	} else if teamF != "" {
		// if only team is set, then return team check
		results, err = queries.FilterByTeam(ctx, teamF)
	} else {
		// table scan - slow!
		results, err = queries.All(ctx)
	}
	return
}

// resultsOut converts the structs to map for output
func resultsOut(results []ghs.GithubStandard, response *resp.Response) (rows []map[string]interface{}, base int, ext int) {
	base = 0
	ext = 0
	rows = []map[string]interface{}{}
	for _, item := range results {
		base += item.CompliantBaseline
		ext += item.CompliantExtended

		if m, err := convert.Map(item); err == nil {
			rows = append(rows, m)
		} else {
			slog.Error("error converting result to map", slog.String("err", err.Error()))
		}
	}
	slog.Info("result out", slog.Int("r", len(response.Errors)))
	return
}

// Handlers are all the handler functions, returns a map of funcs
//   - do it this way so multiple routes can share functions and the handler can then
//     get some of the vars set in the scope of this function (like query params etc)
func Handlers(ctx context.Context, mux *http.ServeMux, dbPath string) map[string]func(w http.ResponseWriter, r *http.Request) {
	archived := query.Get("archived")
	team := query.Get("team")
	// -- handler functions
	var list = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var response = resp.New()
		var filters = map[string]string{
			"archived": query.FirstD(archived.Values(r), "false"),
			"team":     query.First(team.Values(r)),
		}
		response.Start(w, r)

		// -- setup db connection
		db, err := sqlDB(dbPath)
		defer db.Close()

		if err != nil {
			slog.Error("api db error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}
		// -- setup the sqlc generated queries
		queries := ghs.New(db)
		defer queries.Close()

		// -- fetch the raw results
		slog.Info("about to get results")
		results, err := getResults(ctx, queries, filters["archived"], filters["team"])
		slog.Info("got results")
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
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

	return map[string]func(w http.ResponseWriter, r *http.Request){
		"list": list,
	}
}

// Register attached the route to the list view
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	funcs := Handlers(ctx, mux, dbPath)
	list := funcs["list"]
	// -- actually register the handler
	mux.HandleFunc(listRoute, mw.Middleware(list, mw.Logging, mw.SecurityHeaders))

	return nil
}
