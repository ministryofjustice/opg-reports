package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/times"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jmoiron/sqlx"
)

func TestImportsUptimeWithoutMock(t *testing.T) {

	var (
		err    error
		client *cloudwatch.Client
		db     *sqlx.DB
		region string          = "us-east-1"
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-uptime.db")
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
		client, err = awsclients.New[*cloudwatch.Client](ctx, log, region)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
			t.FailNow()
		}
		err = importUptime(ctx, log, client, db, &UptimeOpts{
			AccountID: awsid.AccountID(ctx, log, region),
			Day:       times.AsYMDString(times.Yesterday()),
		})
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	}
}
