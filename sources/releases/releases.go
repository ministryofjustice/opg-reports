package releases

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
)

const RecordsToSeed int = 100

var (
	insert  datastore.InsertStatement //= uptimedb.InsertUptime
	creates []datastore.CreateStatement
)

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup(ctx context.Context, dbFilepath string, seed bool) {
	datastore.Setup[Release](ctx, dbFilepath, insert, creates, seed, RecordsToSeed)
}

// CreateNewDB will create a new DB file and then
// try to run table and index creates
func CreateNewDB(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error) {
	return datastore.CreateNewDB(ctx, dbFilepath, creates)
}
