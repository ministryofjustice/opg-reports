package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
)

func TestImportsTeamsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-teams.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))

		err = teamsImport(ctx, log, client.Repositories, db)
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	} else {
		t.SkipNow()
	}

}
