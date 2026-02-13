package dbsetup

import (
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDBSetupSeedsGeneration(t *testing.T) {

	teams := generateTeams(5)
	if len(teams) != 5 {
		t.Errorf("team generation failed, incorrect count")
	}

	accounts := generateAccounts(60, teams)
	if len(accounts) != 60 {
		t.Errorf("account generation failed, incorrect count")
	}

	costs := generateInfracosts(60000, accounts)
	if len(costs) != 60000 {
		t.Errorf("infracosts generation failed, incorrect count")
	}

	uptime := generateUptime(10000, accounts)
	if len(uptime) != 10000 {
		t.Errorf("uptime generation failed, incorrect count")
	}

	codebases := generateCodebases(50)
	if len(codebases) != 50 {
		t.Errorf("codebase generation failed, incorrect count")
	}

	codeowners := generateCodeowners(200, teams, codebases)
	if len(codeowners) != 200 {
		t.Errorf("codeowners generation failed, incorrect count")
	}
}

func TestDBSetupSeedAll(t *testing.T) {
	var (
		err     error
		db      *sqlx.DB
		dir     = t.TempDir()
		ctx     = t.Context()
		lg      = logger.New("error")
		driver  = "sqlite3"
		connStr = filepath.Join(dir, "dbsetup-seed-all.db")
	)

	db, err = dbconnection.Connection(ctx, lg, driver, connStr)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	defer db.Close()

	err = SeedAll(ctx, lg, db)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}

}
