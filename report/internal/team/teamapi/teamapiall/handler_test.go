package teamapiall

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/global/seeds"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dump"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
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
)
`

func TestTeamsApiAllHandler(t *testing.T) {
	var (
		err    error
		ctx           = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir           = t.TempDir()
		driver        = "sqlite3"
		dbpath        = filepath.Join(dir, "test-teams-handler.db")
		mfile  string = filepath.Join(dir, "migrate.json")
	)

	// run seeds
	seeds.SeedAll(ctx, &seeds.Args{
		Driver:        driver,
		DB:            dbpath,
		MigrationFile: mfile,
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}
	// setup the server and items
	url := "/v1/teams/"
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
		t.Errorf("error converting ...")
	}
	// test returned data
	if len(rec.Data) < 5 {
		t.Errorf("incorrect number of data rows")
	}
	fmt.Println(dump.Any(rec))
	t.FailNow()

}
