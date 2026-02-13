package dbsetup

import (
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/logger"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDBSetupStatementFromT(t *testing.T) {
	var err error

	d1 := []*accountmodels.Account{}
	_, err = statementFromT(d1)
	if err != nil {
		t.Error("non statement matched")
	}

}

func TestDBSetupImports(t *testing.T) {
	var (
		err     error
		db      *sqlx.DB
		dir     = t.TempDir()
		ctx     = t.Context()
		lg      = logger.New("error")
		driver  = "sqlite3"
		connStr = filepath.Join(dir, "dbsetup-import.db")
	)

	db, err = dbconnection.Connection(ctx, lg, driver, connStr)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	defer db.Close()

	err = Migrate(ctx, lg, db)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	// test team import without a fixed statement
	teams := generateTeams(5)
	s1, err := Import[string](ctx, lg, db, teams, nil)
	if err != nil || len(teams) != len(s1) {
		t.Errorf("unexpected team import error:\n%s", err.Error())
		t.FailNow()
	}

	// test account insert
	accounts := generateAccounts(10, teams)
	importStmt := _IMPORTS["accounts"]
	s2, err := Import[string](ctx, lg, db, accounts, importStmt)
	if err != nil || len(accounts) != len(s2) {
		t.Errorf("unexpected account import error:\n%s", err.Error())
		t.FailNow()
	}

	// test cost insert - the inserted count might differ as dates may overlap
	costs := generateInfracosts(13000, accounts)
	importStmt = _IMPORTS["infracosts"]
	s3, err := Import[int](ctx, lg, db, costs, importStmt)
	if err != nil || len(costs) != len(s3) {
		t.Errorf("unexpected cost import error:\n%s", err.Error())
		t.FailNow()
	}

	// test uptime insert - the inserted count might differ as dates may overlap
	uptimes := generateUptime(6000, accounts)
	importStmt = _IMPORTS["uptime"]
	s4, err := Import[int](ctx, lg, db, uptimes, importStmt)
	if err != nil || len(uptimes) != len(s4) {
		t.Errorf("unexpected uptime import error:\n%s", err.Error())
		t.FailNow()
	}

	codebases := generateCodebases(50)
	importStmt = _IMPORTS["codebases"]
	s5, err := Import[int](ctx, lg, db, codebases, importStmt)
	if err != nil || len(codebases) != len(s5) {
		t.Errorf("unexpected codebases import error:\n%s", err.Error())
		t.FailNow()
	}

	codeowners := generateCodeowners(75, teams, codebases)
	importStmt = _IMPORTS["codeowners"]
	s6, err := Import[int](ctx, lg, db, codeowners, importStmt)
	if err != nil || len(codeowners) != len(s6) {
		t.Errorf("unexpected codeowners import error:\n%s", err.Error())
		t.FailNow()
	}

}
