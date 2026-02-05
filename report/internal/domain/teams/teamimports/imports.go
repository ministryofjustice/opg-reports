package teamimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teammodels"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("team import failed with error.")

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
	var lg *slog.Logger = log.With("func", "domain.teams.teamimports.Import")

	statements = []*dbstmts.Insert[*teammodels.Team, string]{}
	lg.Debug("starting ...")
	lg.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstmts.Insert[*teammodels.Team, string]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// run inserts
	lg.Debug("running import statements via insert ...")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		lg.Error("error with insert.", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return
}
