package dbx

import (
	"database/sql"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/convert"
	"path/filepath"
	"testing"
)

type tModel struct {
	Name  string `json:"name"`
	Month string `json:"month"`
}

func (self *tModel) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

const testTable string = `
CREATE TABLE IF NOT EXISTS test_model (
	id INTEGER PRIMARY KEY,
	name TEXT,
	month TEXT
) STRICT;
`
const testInsert string = `INSERT INTO test_model (name, month) VALUES (:name, :month);`

func TestPackageDBXInsert(t *testing.T) {

	var (
		err    error
		db     *sql.DB
		ctx    = t.Context()
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-insert.db")
	)
	// setup a small test db
	db, err = sql.Open(driver, dbpath)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	defer db.Close()

	// run db setup
	_, err = db.ExecContext(ctx, testTable)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	data := []*tModel{
		{Name: "X", Month: "2026-01"},
		{Name: "Y", Month: "2026-02"},
		{Name: "Z", Month: "2026-03"},
		{Name: "A", Month: "2026-03"},
	}

	err = Insert(ctx, testInsert, data, &args.DB{
		Driver: driver, DB: dbpath,
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

}
