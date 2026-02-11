package teamimports

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbimports"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teammodels"

	"github.com/jmoiron/sqlx"
)

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO teams (
	name
) VALUES (
	:name
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name
;
`

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*teammodels.Team) (statements []*dbstmts.Insert[*teammodels.Team, string], err error) {
	var lg *slog.Logger = log.With("func", "teamimports.Import")

	lg.Debug("starting ...")
	statements, err = dbimports.Import[string](ctx, log, db, insertStmt, data)
	if err != nil {
		return
	}
	lg.Debug("complete.")
	return
}
