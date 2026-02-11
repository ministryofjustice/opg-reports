package infracostimports

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbimports"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"

	"github.com/jmoiron/sqlx"
)

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
	var lg *slog.Logger = log.With("func", "infracostimports.Import")

	lg.Debug("starting ...")
	statements, err = dbimports.Import[int](ctx, log, db, insertStmt, data)
	if err != nil {
		return
	}
	lg.Debug("complete.")
	return
}
