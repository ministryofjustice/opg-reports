// Package costs provides the model and setup functions
// for all cost related data
//
// Currently only capturing AWS cost data.
package costs

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsdb"
)

const RecordsToSeed int = 15000

var insert = costsdb.InsertCosts
var creates = []datastore.CreateStatement{
	costsdb.CreateCostTable,
	costsdb.CreateCostTableIndex,
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup(ctx context.Context, dbFilepath string, seed bool) {
	datastore.Setup[*Cost](ctx, dbFilepath, insert, creates, seed, RecordsToSeed)
}

// CreateNewDB will create a new DB file and then
// try to run table and index creates
func CreateNewDB(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error) {
	return datastore.CreateNewDB(ctx, dbFilepath, creates)
}
