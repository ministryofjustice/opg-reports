package github_standards

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/query"
	"github.com/ministryofjustice/opg-reports/shared/exists"
)

const dbMode string = "WAL"
const dbTimeout int = 10000
const dbFk bool = true

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

	// -- actual handler
	mux.HandleFunc("/github/standards/{version}/{$}", func(w http.ResponseWriter, r *http.Request) {
		db, err := sqlDB(dbPath)
		if err != nil {
			return
		}
		queries := ghs.New(db)

		is_archived := 0
		if q := query.First(archived.Values(r)); q == "true" {
			is_archived = 1
		}

		all, err := queries.Archived(ctx, is_archived)

		slog.Info("count",
			slog.Int("count", len(all)),
			slog.String("err", fmt.Sprintf("%v", err)))

		// logic of getting what data

	})

	return nil
}
