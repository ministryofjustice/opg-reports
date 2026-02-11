package uptimeimports

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbimports"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimemodels"

	"github.com/jmoiron/sqlx"
)

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
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*uptimemodels.Uptime) (statements []*dbstmts.Insert[*uptimemodels.Uptime, int], err error) {
	var lg *slog.Logger = log.With("func", "uptimeimports.Import")

	lg.Debug("starting ...")
	statements, err = dbimports.Import[int](ctx, log, db, insertStmt, data)
	if err != nil {
		return
	}
	lg.Debug("complete.")
	return
}
