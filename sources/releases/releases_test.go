package releases_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/releases"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
)

// runs Setup and then checks the
// file exists and the location and that it has
// records in the table
func TestReleasesSetup(t *testing.T) {
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.Background()

	releases.Setup(ctx, dbFile, true)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Errorf("database file was not created in expected locations")
	}
	db, _, err := datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("error connecting to db [%s]", err.Error())
	}
	defer db.Close()

	count, err := datastore.Get[int](ctx, db, releasesdb.ReleaseCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count != releases.RecordsToSeed {
		t.Errorf("incorrect number of rows - expected [%d] actual [%v]", releases.RecordsToSeed, count)
	}

	// // check it made 4 teams
	// count, err = datastore.Get[int](ctx, db, releasesdb.TeamsCount)
	// if err != nil {
	// 	t.Errorf("error counting db rows: [%s]", err.Error())
	// }
	// if count != releases.TeamRecordsToSeed {
	// 	t.Errorf("incorrect number of rows - expected [%d] actual [%v]", releases.TeamRecordsToSeed, count)
	// }

	// // check it made correct number of joins
	// // - for the purpose of seeding, it only links a release to a single team
	// count, err = datastore.Get[int](ctx, db, releasesdb.JoinCount)
	// if err != nil {
	// 	t.Errorf("error counting db rows: [%s]", err.Error())
	// }
	// if count != releases.RecordsToSeed {
	// 	t.Errorf("incorrect number of rows - expected [%d] actual [%v]", releases.TeamRecordsToSeed, count)
	// }

}
