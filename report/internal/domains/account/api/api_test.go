package api

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/config"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/dbx"
	"opg-reports/report/packages/httpx"
	"path/filepath"
	"testing"
)

const testSetup string = `
CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;
CREATE INDEX IF NOT EXISTS accounts_idx ON accounts(id);

INSERT INTO accounts
	(id, name, label, environment, team_name)
VALUES
	('1001', 'account-a', 'test', 'production', 'test-team-a')
;
INSERT INTO accounts
	(id, name, label, environment, team_name)
VALUES
	('1002', 'account-b', 'test', 'production', 'test-team-b')
;
`

func TestDomainsAccountApiHandler(t *testing.T) {
	var (
		ctx    = t.Context()
		cfg    = config.NewApi()
		dir    = t.TempDir()
		dbpath = filepath.Join(dir, "test-handler.db")
	)
	// setup a small test db
	cfg.DBPath = dbpath
	testSeed(ctx, cfg, t)

	// create the mux with the db connection, but not template - so json
	mux := httpx.NewMux()
	// register the handler
	httpx.Register(ctx, mux, cfg, `/accounts/{$}`, nil, Accounts)
	// setup the test call with no filters
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/accounts/?team=test-team-b`, nil)
	// run the request
	mux.ServeHTTP(rr, req)
	resp := rr.Result()

	result := &Result{}
	convert.Between(resp, &result)

	if len(result.Accounts) != 1 {
		t.Errorf("should only be one account for test-team-b")
	}
	if result.Accounts[0].ID != "1002" {
		t.Errorf("incorect account")
	}

}

func testSeed(ctx context.Context, conn dbx.Connector, t *testing.T) {
	var (
		err error
		db  *sql.DB
	)
	// setup a small test db
	db = conn.Connection()
	defer db.Close()
	// run db setup
	_, err = db.ExecContext(ctx, testSetup)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
}
