package uptimeimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/uptime/uptimemodels"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("uptime import failed with error")

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO uptime (
	date,
	average,
	granularity,
	account_id
) VALUES (
	:date,
	:average,
	:granularity,
	:account_id
) ON CONFLICT (account_id,date)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id;
`

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*uptimemodels.Uptime) (statements []*dbstatements.InsertStatement[*uptimemodels.Uptime, int], err error) {

	statements = []*dbstatements.InsertStatement[*uptimemodels.Uptime, int]{}
	log = log.With("package", "uptime", "func", "Import")

	log.Debug("starting ...")
	log.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstatements.InsertStatement[*uptimemodels.Uptime, int]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// run inserts
	log.Debug("running import statements via insert")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		log.Error("error with insert", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	log.Debug("complete.")
	return
}
