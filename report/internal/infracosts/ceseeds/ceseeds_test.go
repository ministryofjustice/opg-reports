package ceseeds

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/infracosts/cemigrations"
	"opg-reports/report/internal/infracosts/cemodels"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInfracostsCeSeedWorking(t *testing.T) {
	var (
		err        error
		db         *sqlx.DB
		dir        string          = t.TempDir()
		ctx        context.Context = t.Context()
		log        *slog.Logger    = logger.New("error", "text")
		driver     string          = "sqlite3"
		connStr    string          = fmt.Sprintf("%s/%s", dir, "seed-working.db")
		statements []*dbstatements.DataStatement[*cemodels.AwsCost, int]
	)
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()
	// db schema setup
	err = cemigrations.Migrate(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected migration issue:\n%v", err.Error())
	}

	statements, err = Seed(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected seed issue:\n%v", err.Error())
	}
	if len(statements) < 1 {
		t.Errorf("expected multiple results to be returned")
	}
	if statements[0].Returned <= 0 {
		t.Errorf("expected positive id for row insert")
	}
}
