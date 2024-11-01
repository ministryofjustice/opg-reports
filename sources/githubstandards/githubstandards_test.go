package githubstandards_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/githubstandards"
	"github.com/ministryofjustice/opg-reports/sources/githubstandards/githubstandardsdb"
)

// TempDir runs Setup and then checks the
// file exists and the location and that it has
// records in the table
func TestStandardsSetup(t *testing.T) {
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.Background()

	githubstandards.Setup(ctx, dbFile, true)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Errorf("database file was not created in expected locations")
	}
	db, _, err := datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("error connecting to db [%s]", err.Error())
	}

	count, err := datastore.Get[int](ctx, db, githubstandardsdb.RowCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count != githubstandards.RecordsToSeed {
		t.Errorf("incorrect number of rows - expected [%d] actual [%v]", costs.RecordsToSeed, count)
	}
}
