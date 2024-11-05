// Package standards provides info about required standards
//
// Currently, this only github repository settings
package standards

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsdb"
)

const RecordsToSeed int = 100

var insert = standardsdb.InsertStandard
var creates = []datastore.CreateStatement{
	standardsdb.CreateStandardsTable,
	standardsdb.CreateStandardsIndexIsArchived,
	standardsdb.CreateStandardsIndexIsArchivedTeams,
	standardsdb.CreateStandardsIndexTeams,
	standardsdb.CreateStandardsIndexBaseline,
	standardsdb.CreateStandardsIndexExtended,
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup(ctx context.Context, dbFilepath string, seed bool) {

	datastore.Setup[Standard](ctx, dbFilepath, insert, creates, seed, RecordsToSeed)

}

// CreateNewDB will create a new DB file and then
// try to run table and index creates
func CreateNewDB(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error) {

	return datastore.CreateNewDB(ctx, dbFilepath, creates)
}
