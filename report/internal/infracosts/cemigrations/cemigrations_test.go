package cemigrations

import (
	"fmt"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInfracostsCeMigrationWorking(t *testing.T) {
	var (
		err error
		db  *sqlx.DB
		dir = t.TempDir()
		// dir     = "./"
		ctx     = t.Context()
		lg      = logger.New("error", "text")
		driver  = "sqlite3"
		connStr = fmt.Sprintf("%s/%s", dir, "migration-working.db")
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
