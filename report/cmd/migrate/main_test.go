package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestCMDMigrate(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-cmd-migrate.db")
	)
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	defer db.Close()

	err = runMigrations(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected migration error: %s", err.Error())
	}

}
