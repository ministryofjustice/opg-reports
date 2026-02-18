package costimport

import (
	"context"
	"opg-reports/report/internal/global"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

// aws-vault exec use-development-operator -- make test name="TestCostImportWithoutMock"
func TestCostImportWithoutMock(t *testing.T) {
	var (
		err       error
		client    *costexplorer.Client
		accountId string
		dir       string          = t.TempDir()
		mfile     string          = filepath.Join(dir, "migrate.json")
		dbpath    string          = filepath.Join(dir, "test-costs-import.db")
		ctx       context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		now       time.Time       = time.Now().UTC()
		start     time.Time       = now.AddDate(0, -4, 0)
		end       time.Time       = now.AddDate(0, -3, 0)
	)

	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = awsclients.New[*costexplorer.Client](ctx, "eu-west-1")
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}

	global.MigrateAll(ctx, &migrations.Args{
		DB:            dbpath,
		Driver:        "sqlite3",
		MigrationFile: mfile,
	})

	// global.MigrateAll(ctx, &globalmodels.MigrationArgs{
	// })

	accountId = awsid.AccountID(ctx, "eu-west-1")
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
	t.FailNow()

}
