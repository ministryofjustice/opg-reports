package adaptors

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

var _ dbs.DBer = &SqlxDB{}

func TestAdaptorsDBerGetFailsAsNotExists(t *testing.T) {
	var conn = &Connection{Driver: "test", Path: "test.db", Parameters: "?yes=1"}
	var sq = &SqlxDB{}
	var ctx = context.Background()

	if _, err := sq.Get(ctx, conn); err == nil {
		t.Errorf("did not get the expected a failure - the database does not exist.")
	}

	// should be fine to call close after
	if err := sq.Close(); err != nil {
		t.Errorf("close should not have failed, the pointer should be nil after error and ignored")
	}
}

func TestAdaptorsDBerGet(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var file = filepath.Join(dir, "test.db")
	// create the database
	err = createDB(file)
	if err != nil {
		t.Errorf("unexpected failure creating database file [%s]", err.Error())
	}

	var conn = &Connection{Driver: sqliteDriver, Path: file, Parameters: sqliteParams}
	var sq = &SqlxDB{}
	var ctx = context.Background()

	_, err = sq.Get(ctx, conn)
	if err != nil {
		t.Errorf("unexpected error getting database: [%s]", err.Error())
	}

	err = sq.Close()
	if err != nil {
		t.Errorf("close should not have failed: [%s]", err.Error())
	}

}
