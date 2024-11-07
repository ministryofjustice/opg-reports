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

	// run the main setup, we want to make sure
	// database is setup and seeded correctly with random data
	releases.Setup(ctx, dbFile, true)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Errorf("database file was not created in expected locations")
	}

	// grab db connection
	db, _, err := datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("error connecting to db [%s]", err.Error())
	}
	defer db.Close()

	// check the number of releases created is the same as the seed count
	count, err := datastore.Get[int](ctx, db, releasesdb.ReleaseCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count != releases.RecordsToSeed {
		t.Errorf("incorrect number of rows - expected [%d] actual [%v]", releases.RecordsToSeed, count)
	}

	// Based on the faker config on the release model every release should
	// generate a single
	// Those teams may not be unique, but we can check the number joins
	// which should be 1 x numberof releases
	expectedJoins := 1 * releases.RecordsToSeed
	count, err = datastore.Get[int](ctx, db, releasesdb.JoinCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count != expectedJoins {
		t.Errorf("incorrect number of joins - expected [%d] actual [%v]", expectedJoins, count)
	}

	// Now check that a record fetched will return the correct team details
	rand := &releases.Release{}
	err = datastore.GetRecord(ctx, db, releasesdb.GetRandomRelease, rand)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	// check that we get a single a team from the seeding!
	teams, err := rand.Teams(ctx, db)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}

	if len(teams) != 1 {
		t.Errorf("failed to get team for record")
	}
}
