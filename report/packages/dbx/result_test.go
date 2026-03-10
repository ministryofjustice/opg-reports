package dbx

import (
	"database/sql"
	"path/filepath"
	"testing"
)

type testModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (self *testModel) Sequence() []any {
	return []any{
		&self.ID,
		&self.Name,
	}
}

const testTableInsert string = `
CREATE TABLE IF NOT EXISTS test_model (
	id INTEGER PRIMARY KEY,
	name TEXT
) STRICT;

INSERT INTO test_model (name) VALUES ('Z');
INSERT INTO test_model (name) VALUES ('A');
INSERT INTO test_model (name) VALUES ('B');
INSERT INTO test_model (name) VALUES ('C');
INSERT INTO test_model (name) VALUES ('D');
`

const testSelectStmt string = `
SELECT
	id,
	name
FROM test_model
ORDER BY name ASC;
`

func TestPackagesDBXResult(t *testing.T) {
	var (
		err    error
		db     *sql.DB
		rows   *sql.Rows
		ctx    = t.Context()
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-select-result.db")
	)
	// setup a small test db
	db, err = sql.Open(driver, dbpath)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	defer db.Close()

	// run db setup
	_, err = db.ExecContext(ctx, testTableInsert)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	// skip the binding / args stage of a select
	// as there are no filters on this select
	// so hit the query context directly
	rows, err = db.QueryContext(ctx, testSelectStmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	defer rows.Close()
	// now run the row scan and see what we get
	results := Results[*testModel]{}
	for rows.Next() {
		err = results.RowScan(rows)
		if err != nil {
			t.Errorf("unexpected error with row scan: %s", err.Error())
			t.FailNow()
		}
	}

	if len(results.Data()) != 5 {
		t.Errorf("failed to return all data from the db.")
	}

}
