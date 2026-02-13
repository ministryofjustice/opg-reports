package codebaseselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDomainSelectsCodebasesAll(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-domain-selects-codebases.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}

	defer db.Close()
	// insert some dummy selects with seed command
	err = dbsetup.SeedAll(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected seeds error: [%s]", err.Error())
		t.FailNow()
	}
	// select all and compare counts
	_, err = All(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected all error: [%s]", err.Error())
		t.FailNow()
	}

}
