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
	"github.com/ministryofjustice/opg-reports/shared/convert"
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

func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {

	archived := query.Get("archived")
	team := query.Get("team")

	// -- actual handler
	mux.HandleFunc("/github/standards/{version}/{$}", func(w http.ResponseWriter, r *http.Request) {
		var err error
		var teams = ""
		var is_archived = ""
		var results []ghs.GithubStandard

		db, err := sqlDB(dbPath)
		if err != nil {
			return
		}
		queries := ghs.New(db)

		// -- fetch the get parameter values
		// team query, add the like logic here
		if qv := query.First(team.Values(r)); qv != "" {
			teams = "%" + qv + "#"
		}
		// archive query
		if qv := query.First(archived.Values(r)); qv != "" {
			is_archived = qv
		}

		// -- run correct query
		if teams != "" && is_archived != "" {
			results, err = queries.ArchivedTeamFilter(ctx, ghs.ArchivedTeamFilterParams{
				IsArchived: strToInt(is_archived), Teams: teams,
			})
		} else if is_archived != "" {
			results, err = queries.ArchivedFilter(ctx, strToInt(is_archived))
		} else if teams != "" {
			results, err = queries.TeamFilter(ctx, teams)
		} else {
			results, err = queries.All(ctx)
		}

		slog.Info("count",
			slog.Int("results", len(results)),
			slog.String("err", fmt.Sprintf("%v", err)))

		// -- different api response to costs / uptime, those are tabular, this is more raw
		// -- add compliance
		c, _ := convert.ListToJson(results)

		w.WriteHeader(200)
		w.Write(c)

	})

	return nil
}
