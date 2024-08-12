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

const listRoute string = "/github/standards/{version}/{$}"

func sqlDB(dbPath string) (db *sql.DB, err error) {
	if exists.FileOrDir(dbPath) != true {
		err = fmt.Errorf("database [%s] does not exist", dbPath)
		return
	}
	// connection string to set modes etc
	conn := consts.SQL_CONNECTION_PARAMS
	slog.Info("connecting to db...", slog.String("dbPath", dbPath), slog.String("connection", conn))

	// try to connect to db
	db, err = sql.Open("sqlite3", dbPath+conn)
	if err != nil {
		slog.Error("error connecting to DB", slog.String("err", err.Error()))
		return
	}
	slog.Info("connected to db", slog.String("dbPath", dbPath), slog.String("connection", conn))
	return
}

func getResults(ctx context.Context, queries *ghs.Queries, archived string, team string) (results []ghs.GithubStandard, err error) {

	var teamF = ""
	var archivedF = ""
	results = []ghs.GithubStandard{}

	// -- fetch the get parameter values
	// team query, add the like logic here
	if team != "" {
		teamF = "%" + team + "#"
	}
	// archive query
	if archived != "" {
		archivedF = archived
	}
	// -- run query
	if teamF != "" && archivedF != "" {
		results, err = queries.ArchivedTeamFilter(ctx, ghs.ArchivedTeamFilterParams{
			IsArchived: convert.BoolStringToInt(archivedF), Teams: teamF,
		})
	} else if archivedF != "" {
		results, err = queries.ArchivedFilter(ctx, convert.BoolStringToInt(archivedF))
	} else if teamF != "" {
		results, err = queries.TeamFilter(ctx, teamF)
	} else {
		results, err = queries.All(ctx)
	}
	return
}

func resultsOut(results []ghs.GithubStandard, response *resp.Response) (rows []map[string]interface{}, base int, ext int) {
	base = 0
	ext = 0
	rows = []map[string]interface{}{}
	for _, item := range results {

		// response.AddDataAge(dates.Time(item.Ts))
		if m, err := convert.Map(item); err == nil {

			rows = append(rows, m)
		} else {
			slog.Error("error converting result to map", slog.String("err", err.Error()))
		}
	}
	slog.Info("result out", slog.Int("r", len(response.Errors)))
	return
}

func Handlers(ctx context.Context, mux *http.ServeMux, dbPath string) map[string]func(w http.ResponseWriter, r *http.Request) {
	archived := query.Get("archived")
	team := query.Get("team")
	// -- handler functions
	var list = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var response = resp.New()
		var filters = map[string]string{
			"archived": query.FirstD(archived.Values(r), "true"),
			"team":     query.First(team.Values(r)),
		}
		response.Start(w, r)

		db, err := sqlDB(dbPath)
		if err != nil {
			slog.Error("api db error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		queries := ghs.New(db)
		slog.Info("about to get results")
		results, err := getResults(ctx, queries, filters["archived"], filters["team"])
		slog.Info("got results")
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		// -- work out how many comply and convert over
		res, base, ext := resultsOut(results, response)
		response.Result = res
		response.Metadata["count"] = map[string]int{
			"all":                 len(results),
			"compliance_baseline": base,
			"compliance_extended": ext,
		}
		response.Metadata["filters"] = filters
		response.End(w, r)
	}
	return map[string]func(w http.ResponseWriter, r *http.Request){
		"list": list,
	}
}

func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	funcs := Handlers(ctx, mux, dbPath)
	list := funcs["list"]
	// -- actually register the handler
	mux.HandleFunc(listRoute, mw.Middleware(list, mw.Logging, mw.SecurityHeaders))

	return nil
}
