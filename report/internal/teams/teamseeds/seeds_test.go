package teamseeds

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/teams/teammigrations"
	"opg-reports/report/internal/teams/teammodels"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestRedoTeamsSeedWorking(t *testing.T) {
	var (
		err        error
		db         *sqlx.DB
		dir        string          = t.TempDir()
		ctx        context.Context = t.Context()
		log        *slog.Logger    = logger.New("error", "text")
		driver     string          = "sqlite3"
		connStr    string          = fmt.Sprintf("%s/%s", dir, "seed-team-working.db")
		statements []*dbstatements.InsertStatement[*teammodels.Team, string]
	)
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()
	// db schema setup
	err = teammigrations.Migrate(ctx, log, db)
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
	if statements[0].Returned == "" {
		t.Errorf("expected name returned for row insert")
	}
}
