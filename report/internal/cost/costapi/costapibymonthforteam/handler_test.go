package costapibymonthforteam

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const mockMigrate string = `
CREATE TABLE IF NOT EXISTS costs (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	month TEXT NOT NULL,
	cost TEXT NOT NULL,
	account_id TEXT,
	UNIQUE (account_id,month,region,service)
) STRICT;

CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	vendor TEXT NOT NULL DEFAULT 'aws',
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS idx_costs_date ON costs(month);
CREATE INDEX IF NOT EXISTS idx_costs_date_account ON costs(month, account_id);
CREATE INDEX IF NOT EXISTS idx_costs_unique ON costs(account_id, month, region, service);

INSERT INTO accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	'12300',
	'account01',
	'A01',
	'development',
	'team-a'
);

INSERT INTO accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	'12340',
	'account02',
	'A02',
	'development',
	'team-b'
);

INSERT INTO costs (
	region,
	service,
	month,
	cost,
	account_id
) VALUES(
	'eu-west-1',
	'test',
	'2025-11',
	'123.1',
	'12300'
);

INSERT INTO costs (
	region,
	service,
	month,
	cost,
	account_id
) VALUES(
	'eu-west-1',
	'test',
	'2025-11',
	'123.1',
	'12340'
);
`

func TestCostApiByMonthForTeamHandler(t *testing.T) {
	var (
		err    error
		ctx    = cntxt.AddLogger(t.Context(), logger.New("info"))
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-costs-handler.db")
	)
	os.MkdirAll(filepath.Dir(dbpath), os.ModePerm)
	// setup db base
	db, err := sql.Open(driver, dbpath)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}
	// run migration
	_, err = db.ExecContext(ctx, mockMigrate)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}
	// setup the server and items
	url := "/v1/costs/between/2025-01/2026-01/team/team-a/"
	mux := http.NewServeMux()

	req := httptest.NewRequest(http.MethodGet, url, nil)
	writer := httptest.NewRecorder()

	// setup the bindings to the test handler and call
	Register(ctx, mux, &Config{
		Driver: driver,
		DB:     dbpath,
	})
	mux.ServeHTTP(writer, req)

	// get and parse the result
	resp := writer.Result()
	rec := &Response{}
	err = response.AsT(resp, &rec)
	if err != nil {
		t.Errorf("error converting ... [%s]", err.Error())
	}
	// - test returned data
	// 1 dummy item, 1 table total
	if len(rec.Data) != 2 {
		t.Errorf("incorrect number of data rows")
	}
	if rec.Request.DateEnd != "2026-01" {
		t.Error("data_end failed to return correctly")
	}
	if rec.Request.DateStart != "2025-01" {
		t.Error("data_start failed to return correctly")
	}
	if rec.Request.Team != "team-a" {
		t.Error("team failed to return correctly")
	}
	if len(rec.Headers["labels"]) != 1 {
		t.Error("incorrect number of labels returned")
	}

}
