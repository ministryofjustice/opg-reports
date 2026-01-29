package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
)

func TestImportsCodeownersWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		db     *sqlx.DB
		code   []*codebasemodels.Codebase
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-codeowners.db")
	)
	code = []*codebasemodels.Codebase{
		{FullName: "ministryofjustice/opg-lpa", Name: "opg-lpa"},
		{FullName: "ministryofjustice/opg-use-an-lpa", Name: "opg-use-an-lpa"},
		{FullName: "ministryofjustice/opg-data-lpa-store", Name: "opg-data-lpa-store"},
	}
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	if os.Getenv("GH_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))
		err = importCodeowners(ctx, log, client.Repositories, db, code)
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	} else {
		t.SkipNow()
	}
}
