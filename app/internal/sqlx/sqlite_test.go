package sqlx

import (
	"path/filepath"
	"testing"
)

func TestSQLxSqlite(t *testing.T) {

	var (
		sq  *Sqlite
		err error
		dir string = t.TempDir()
	)
	// test connecting to the db that works
	sq = NewSQLite(filepath.Join(dir, "test-connection.sql"))
	_, err = sq.Open()

	if err != nil {
		t.Errorf("unexpected error calling open [%s]", err.Error())
		t.FailNow()
	}
	sq.Close()
}
