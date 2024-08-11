package github_standards

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/query"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp/cell"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp/row"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp/table"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/exists"
)

const dbMode string = "WAL"
const dbTimeout int = 10000
const dbFk bool = true

func strToInt(s string) int {
	b, err := strconv.ParseBool(s)
	if err == nil && b {
		return 1
	}
	return 0
}

func sqlDB(dbPath string) (db *sql.DB, err error) {
	if exists.FileOrDir(dbPath) != true {
		err = fmt.Errorf("database [%s] does not exist", dbPath)
		return
	}
	// connection string to set modes etc
	conn := fmt.Sprintf("?_journal=%s&_timeout=%d&fk=%t", "WAL", 10000, dbFk)

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

func getResults(
	ctx context.Context,
	queries *ghs.Queries, archived string, team string) (results []ghs.GithubStandard, err error) {

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

	// -- run correct query
	if teamF != "" && archivedF != "" {
		results, err = queries.ArchivedTeamFilter(ctx, ghs.ArchivedTeamFilterParams{
			IsArchived: strToInt(archivedF), Teams: teamF,
		})
	} else if archivedF != "" {
		results, err = queries.ArchivedFilter(ctx, strToInt(archivedF))
	} else if teamF != "" {
		results, err = queries.TeamFilter(ctx, teamF)
	} else {
		results, err = queries.All(ctx)
	}
	return
}

func resultsToRows(results []ghs.GithubStandard, response *resp.Response) (rows []*row.Row) {

	rows = []*row.Row{}
	for _, item := range results {

		response.AddDataAge(dates.Time(item.Ts))
		cells := []*cell.Cell{}
		mapped, _ := convert.ToMap(item)
		for k, v := range mapped {
			cells = append(cells, cell.New(k, v, false, false))
		}
		r := row.New(cells...)
		rows = append(rows, r)
	}
	return
}

func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	// -- allowed filters
	archived := query.Get("archived")
	team := query.Get("team")
	// -- handler functions
	var list = func(w http.ResponseWriter, r *http.Request) {
		var err error
		var response = resp.New()
		var filters = map[string]string{
			"archived": query.First(archived.Values(r)),
			"team":     query.First(team.Values(r)),
		}
		response.Start(w, r)

		db, err := sqlDB(dbPath)
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		queries := ghs.New(db)

		results, err := getResults(ctx, queries, filters["archived"], filters["team"])
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		rows := resultsToRows(results, response)
		tbl := table.New(rows...)
		response.Result = tbl

		response.Metadata["count"] = len(results)
		response.Metadata["filters"] = filters
		response.End(w, r)

		// -- add compliance

	}

	// -- actually register the handler
	mux.HandleFunc("/github/standards/{version}/{$}", list)

	return nil
}
