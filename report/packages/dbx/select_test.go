package dbx

import (
	"database/sql"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/convert"
	"path/filepath"
	"testing"
)

type testFilter struct {
	Months []string `json:"months,omitempty"`
}

// Map returns a map of all fields on this struct
func (self *testFilter) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

type testResult struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Month string `json:"month"`
}

func (self *testResult) Sequence() []any {
	return []any{
		&self.ID,
		&self.Name,
		&self.Month,
	}
}

const testSetup string = `
CREATE TABLE IF NOT EXISTS test_model (
	id INTEGER PRIMARY KEY,
	name TEXT,
	month TEXT
) STRICT;

INSERT INTO test_model (name, month) VALUES ('Z', '2025-12');
INSERT INTO test_model (name, month) VALUES ('A', '2026-01');
INSERT INTO test_model (name, month) VALUES ('B', '2026-02');
INSERT INTO test_model (name, month) VALUES ('C', '2026-03');
INSERT INTO test_model (name, month) VALUES ('D', '2026-03');
`

const testSelStmt string = `
SELECT
	id,
	name,
	month
FROM test_model
WHERE
	month IN(:months)
ORDER BY name ASC;
`

func TestPackagesDBXSelect(t *testing.T) {
	var (
		err    error
		db     *sql.DB
		ctx    = t.Context()
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-select.db")
		filter = &testFilter{Months: []string{
			"2026-01", "2026-02", "2026-03",
		}}
	)

	// setup a small test db
	db, err = sql.Open(driver, dbpath)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	defer db.Close()

	// run db setup
	_, err = db.ExecContext(ctx, testSetup)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	results, err := Select[*testResult](ctx, testSelStmt, filter, &args.DB{
		Driver: driver, DB: dbpath,
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	// results should have 4 items based on months
	if len(results) != 4 {
		t.Errorf("should have 4 results")
	}
	// results should not include 'Z' as its wrong month
	zFound := false
	for _, r := range results {
		if r.Name == "Z" {
			zFound = true
		}
	}
	if zFound {
		t.Errorf("error, should not have found Z")
	}

}
