package adaptors

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

var _ dbs.Transactioner = &SqlxTransaction{}

func TestAdaptorsTransactionalFails(t *testing.T) {
	var err error
	var conn = &Connection{Driver: "test", Path: "test.db", Parameters: "?yes=1"}
	var mode = &ReadOnly{}
	var sq = &SqlxDB{}
	var ctx = context.Background()
	var tx = &SqlxTransaction{}

	_, err = tx.Get(ctx, sq, conn, mode)
	if err == nil {
		t.Errorf("tranaction should have returned an error as the db does not exist")
	}
}

func TestAdaptorsTransactionalWorking(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var file = filepath.Join(dir, "test.db")
	var conn = &Connection{Driver: sqliteDriver, Path: file, Parameters: sqliteParams}
	var mode = &ReadOnly{}
	var sq = &SqlxDB{}
	var ctx = context.Background()
	var tx = &SqlxTransaction{}

	// create the database
	err = createDB(file)
	if err != nil {
		t.Errorf("unexpected failure creating database file [%s]", err.Error())
	}

	_, err = tx.Get(ctx, sq, conn, mode)
	if err != nil {
		t.Errorf("GetTransaction returned unexpected error: [%s]", err.Error())
	}

	err = tx.Commit(false)
	if err != nil {
		t.Errorf("commiting empty transaction returned an error: [%s]", err.Error())
	}
}
