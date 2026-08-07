package sqlx

import (
	"path/filepath"
	"testing"
)

func TestSQLxExec(t *testing.T) {
	var (
		sq  *Sqlite
		err error
		ctx        = t.Context()
		dir string = t.TempDir()
	)
	// test connecting to the db that works
	sq = NewSQLite(filepath.Join(dir, "test-exec.sql"))
	_, err = sq.Open()

	if err != nil {
		t.Errorf("unexpected error calling open [%s]", err.Error())
		t.FailNow()
	}
	defer sq.Close()

	// test creating a table and inserting a row
	sql := "CREATE TABLE IF NOT EXISTS test_table(id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;"
	_, err = Exec(ctx, sq, sql)
	if err != nil {
		t.Errorf("unexpected error calling exec [%s]", err.Error())
		t.FailNow()
	}

	insert := "INSERT INTO test_table(name) VALUES ('testA');"
	r, err := Exec(ctx, sq, insert)
	if err != nil {
		t.Errorf("unexpected error calling insert exec [%s]", err.Error())
		t.FailNow()
	}

	id, err := r.LastInsertId()
	if err != nil {
		t.Errorf("unexpected error with last insert id [%s]", err.Error())
		t.FailNow()
	}

	if id <= 0 {
		t.Errorf("invalid last insert id returned (%v)", id)
	}
}
