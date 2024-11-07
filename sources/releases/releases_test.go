package releases_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
	"github.com/ministryofjustice/opg-reports/sources/releases"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
)

// check interfaces
var _ record.Record = &releases.Release{}
var _ record.JoinInserter = &releases.Release{}
var _ record.JoinSelector = &releases.Release{}

// Runs Setup to create a seeded db and then checks records are inserted correctly
// for the primary record type (Release) and checks join insert and selects.
//
//   - Adds dummy / generated data via .Setup and checks success with row count
//   - Checks the number of joins created between addresses (should be 1 each)
//   - Checks calling .Teams will return the correct team data
//   - Checks that calling a db select operation will in .TeamList correctly
//
// Tests the JoinInserter & JoinSelector interfaces
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

	// every release should have at least 1 team, maybe more
	// make sure the number of joins is at least the number of releases
	count, err = datastore.Get[int](ctx, db, releasesdb.JoinCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count < releases.RecordsToSeed {
		t.Errorf("incorrect number of joins - actual [%v]", count)
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

	if len(teams) <= 0 {
		t.Errorf("failed to get team for record")
	}
	// now compare fetched team to selected version
	if len(rand.TeamList) != len(teams) {
		t.Errorf("automatic team fetching via JoinSelector interface failed")
	}
	// check content on each
	for _, og := range rand.TeamList {
		found := false
		for _, pulled := range teams {
			if pulled.Name == og.Name {
				found = true
			}
		}
		if !found {
			t.Errorf("failed to match teams between interface join selector and direct fetch.")
		}
	}
}

// If we generate several thousand records, whats the performance of the join
// select and insert interfaces like.
// This test will time an insert and full fetch
func TestReleasesPerformance(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "perf.db")
	var ctx = context.Background()
	var n = 10000

	inTick := timer.New()
	// make sure extra providers are enabled
	exfaker.AddProviders()
	// create db
	db, _, err = releases.CreateNewDB(ctx, dbFile)
	if err != nil {
		t.Errorf("error [%s]", err.Error())
	}
	defer db.Close()

	// -- generate large fake dataset and then insert, timing them
	faked := exfaker.Many[*releases.Release](n)
	ids, err := datastore.InsertMany(ctx, db, releasesdb.InsertRelease, faked)
	if err != nil {
		t.Errorf("error inserting: [%s]", err.Error())
	}
	if len(ids) != n {
		t.Errorf("incorrect number of row ids returned - expected [%d] actual [%v]", n, len(ids))
	}
	inTick.Stop()

	// -- now fetch all of those records back and time that
	outTick := timer.New()
	all, err := datastore.SelectMany[*releases.Release](ctx, db, "SELECT * FROM releases;", &releases.Release{})
	if err != nil {
		t.Errorf("error inserting: [%s]", err.Error())
	}

	outTick.Stop()

	if len(all) != n {
		t.Errorf("incorrect number of rows returned - expected [%d] actual [%v]", n, len(all))
	}

	fmt.Println("insert duration")
	fmt.Println(inTick.Duration())
	fmt.Println("select duration")
	fmt.Println(outTick.Duration())

}
