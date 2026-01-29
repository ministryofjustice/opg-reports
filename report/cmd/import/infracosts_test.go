package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jmoiron/sqlx"
)

func TestImportsInfracostsWithoutMock(t *testing.T) {

	var (
		err    error
		client *costexplorer.Client
		db     *sqlx.DB
		region string          = "eu-west-1"
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-infracosts.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client, err = awsclients.New[*costexplorer.Client](ctx, log, region)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
			t.FailNow()
		}
		err = importInfracosts(ctx, log, client, db, &InfraOpts{
			AccountID:            awsid.AccountID(ctx, log, region),
			IncludePreviousMonth: true,
			EndDate:              "2025-11-12",
		})
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	}

}
