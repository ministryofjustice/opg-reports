package dbsetup

import (
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDBSetupMigrationWorking(t *testing.T) {
	var (
		err     error
		db      *sqlx.DB
		dir     = t.TempDir()
		ctx     = t.Context()
		lg      = logger.New("error")
		driver  = "sqlite3"
		connStr = filepath.Join(dir, "dbsetup-migrations.db")
	)

	db, err = dbconnection.Connection(ctx, lg, driver, connStr)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	defer db.Close()

	err = Migrate(ctx, lg, db)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
}
