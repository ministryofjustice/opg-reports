package teamselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/domain/teams/teamseeds"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDomainSelectsTeamsAll(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-domain-selects-teams.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()
	// insert some dummy selects with seed command
	seeded, err := teamseeds.Seed(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected seeds error: [%s]", err.Error())
		t.FailNow()
	}
	// select all and compare counts
	data, err := All(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected all error: [%s]", err.Error())
		t.FailNow()
	}

	if len(data) != len(seeded) {
		t.Errorf("mismatched row count between seed and select.")
	}

}
