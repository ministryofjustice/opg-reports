package uptimeimport

import (
	"context"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/times"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// aws-vault exec use-production-operator -- make test name="TestUptimeImportWithoutMock"
func TestUptimeImportWithoutMock(t *testing.T) {
	var (
		err       error
		client    *cloudwatch.Client
		accountId string
		dir       string          = t.TempDir()
		mfile     string          = filepath.Join(dir, "migrate.json")
		dbpath    string          = filepath.Join(dir, "test-costs-import.db")
		ctx       context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		end       time.Time       = times.ResetMonth(time.Now().UTC())
		start     time.Time       = times.Add(end, -1, times.MONTH)
	)

	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = awsclients.New[*cloudwatch.Client](ctx, "us-east-1")
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}

	migrations.MigrateAll(ctx, &migrations.Args{
		DB:            dbpath,
		Driver:        "sqlite3",
		MigrationFile: mfile,
	})

	accountId = awsid.AccountID(ctx, "us-east-1")
	err = Import(ctx, client, &Args{
		DB:        dbpath,
		Driver:    "sqlite3",
		DateStart: start,
		DateEnd:   end,
		AccountID: accountId,
	})

	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

}
