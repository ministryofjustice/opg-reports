package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/exists"
)

func Connect(dbPath string) (db *sql.DB, err error) {
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
