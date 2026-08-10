package sqlx

import (
	"database/sql"
	"opg-reports/app/internal/fmtx"
	"path/filepath"
	"testing"
)

// check the type
var _ RowScanner[*tOrderRow] = RowScan

type tOrderFilter struct{}

type tOrderRow struct {
	ID       int    `json:"id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Username string `json:"username"`
}

func (o *tOrderRow) BindResult() []any {
	return []any{
		&o.ID, &o.Item, &o.Quantity, &o.Username,
	}
}

var (
	tCreateOrderTable string = `CREATE TABLE IF NOT EXISTS test_order_table(
		id INTEGER PRIMARY KEY,
		item TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		username TEXT NOT NULL
	) STRICT;`
	tInsertOrderSQL []string = []string{
		tCreateOrderTable,
		`INSERT INTO test_order_table(item, quantity, username) VALUES ("testItem01", 2, "u01");`,
		`INSERT INTO test_order_table(item, quantity, username) VALUES ("testItem02", 1, "u02"); `,
	}
	tSelectOrderSQL string = `SELECT id, item, quantity, username FROM test_order_table ORDER BY id ASC;`
)

func TestSQLxSelect(t *testing.T) {
	var (
		err error
		sq  *Sqlite
		ctx        = t.Context()
		dir string = t.TempDir()
	)
	// test connecting to the db that works
	sq = NewSQLite(filepath.Join(dir, "test-rowscan.db"), false)

	defer sq.Close()

	// add stuff via exec...
	for _, sql := range tInsertOrderSQL {
		_, err = Exec(ctx, sq, sql)
		if err != nil {
			t.Errorf("unexpected error calling create [%s]", err.Error())
			t.FailNow()
		}

	}
	// now run the select
	results := []*tOrderRow{}

	err = Select(ctx, sq, tSelectOrderSQL, &tOrderFilter{}, &results, nil)
	if err != nil {
		t.Errorf("unexpected error calling Select [%s]", err.Error())
		t.FailNow()
	}

	if len(results) != 2 {
		t.Error("incorrect number of results returned.")
		fmtx.Printj(results)
		t.FailNow()
	}
	if results[0].ID != 1 {
		t.Error("first item is not ID 1.")
	}

}

// TestSQLxRowScan pulls out the based part parsing the SQL results into
// slice of structs
func TestSQLxRowScan(t *testing.T) {
	var (
		err  error
		sq   *Sqlite
		db   *sql.DB
		rows *sql.Rows
		ctx         = t.Context()
		dir  string = t.TempDir()
	)
	// test connecting to the db that works
	sq = NewSQLite(filepath.Join(dir, "test-rowscan.db"), false)
	db, err = sq.Open()
	if err != nil {
		t.Errorf("unexpected error calling open [%s]", err.Error())
		t.FailNow()
	}
	defer sq.Close()

	// add stuff via exec...
	for _, sql := range tInsertOrderSQL {
		_, err = Exec(ctx, sq, sql)
		if err != nil {
			t.Errorf("unexpected error calling create [%s]", err.Error())
			t.FailNow()
		}

	}
	// now select from the table...
	rows, err = db.QueryContext(ctx, tSelectOrderSQL)
	if err != nil {
		t.Errorf("unexpected error calling QueryContext [%s]", err.Error())
		t.FailNow()
	}
	results := []*tOrderRow{}
	// now use row scanning...
	defer rows.Close()
	for rows.Next() {
		err = RowScan(rows, &results)
		if err != nil {
			t.Errorf("unexpected error calling RowScan [%s]", err.Error())
			t.FailNow()
		}

	}

	if len(results) != 2 {
		t.Error("incorrect number of results returned.")
		fmtx.Printj(results)
		t.FailNow()
	}
	if results[0].ID != 1 {
		t.Error("first item is not ID 1.")
	}

	// t.FailNow()

}
