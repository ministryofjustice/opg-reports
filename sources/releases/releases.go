package releases

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
)

// seed counters
const (
	RecordsToSeed     int = 100
	TeamRecordsToSeed int = 4
)

var (
	insert  datastore.InsertStatement   = releasesdb.InsertRelease
	creates []datastore.CreateStatement = []datastore.CreateStatement{
		releasesdb.CreateReleaseTable,
		releasesdb.CreateTeamTable,
		releasesdb.CreateTeamReleaseJoinTable,
		releasesdb.CreateReleaseDateIndex,
		releasesdb.CreateTeamNameIndex,
	}
)

// seedJoins will create the joins between the main release
// table, the team table and the join table to create a set
// of linked data
func seedJoins(ctx context.Context, dbFilepath string) {
	var (
		err       error
		db        *sqlx.DB
		all       []*Release = []*Release{}
		fakeTeams []*Team    = []*Team{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}}
		teamCount int        = 0
	)

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// see if we have any teams - if we dont, then create them
	// and the links to the releases
	teamCount, err = datastore.Get[int](ctx, db, releasesdb.TeamsCount)
	if teamCount <= 0 {
		_, err = datastore.InsertMany(ctx, db, releasesdb.InsertTeam, fakeTeams)
		// panic on error
		if err != nil {
			panic(err)
		}

		// now map releases to teams
		r := &Release{}
		all, err = datastore.List[*Release](ctx, db, releasesdb.AllReleases, r)
		if err != nil {
			panic(err)
		}
		for _, release := range all {
			var team = exfaker.Choice(fakeTeams)
			var join = &Join{TeamID: team.ID, ReleaseID: release.ID}
			_, err = datastore.InsertOne(ctx, db, releasesdb.InsertJoin, join, nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
// Populates the joins as well
func Setup(ctx context.Context, dbFilepath string, seed bool) {

	datastore.Setup[*Release](ctx, dbFilepath, insert, creates, seed, RecordsToSeed)

	if seed {
		seedJoins(ctx, dbFilepath)
	}
}

// CreateNewDB will create a new DB file and then
// try to run table and index creates
func CreateNewDB(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error) {
	return datastore.CreateNewDB(ctx, dbFilepath, creates)
}
