package infracostimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("infracost import failed with error")

const insertStmt string = `
INSERT INTO infracosts (
	region,
	service,
	date,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:date,
	:cost,
	:account_id
) ON CONFLICT (account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*infracostmodels.Cost) (statements []*dbstmts.Insert[*infracostmodels.Cost, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.infracosts.infracostimports.Import")
	statements = []*dbstmts.Insert[*infracostmodels.Cost, int]{}

	lg.Debug("starting ...")
	lg.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstmts.Insert[*infracostmodels.Cost, int]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// run inserts
	lg.Debug("running import statements via insert ...")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		log.Error("error with insert", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return
}
