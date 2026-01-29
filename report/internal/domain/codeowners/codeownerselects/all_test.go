package codeownerselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/codeowners/codeownerseeds"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDomainSelectsCodeownersAll(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-domain-selects-codeowners.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()
	// insert some dummy selects with seed command
	seeded, err := codeownerseeds.Seed(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected seeds error: [%s]", err.Error())
		t.FailNow()
	}
	// select all and compare counts
	data, err := All(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected all error: [%s]", err.Error())
		t.FailNow()
	}

	if len(data) != len(seeded) {
		t.Errorf("mismatched row count between seed and select.")
	}

}
