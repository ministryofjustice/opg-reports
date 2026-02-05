package accountimports

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbimports"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountmodels"

	"github.com/jmoiron/sqlx"
)

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:team_name
)
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id
;
`

// Import uses combines the data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*accountmodels.Account) (statements []*dbstmts.Insert[*accountmodels.Account, string], err error) {
	var lg *slog.Logger = log.With("func", "domain.accounts.accountimports.Import")

	lg.Debug("starting ...")
	statements, err = dbimports.Import[string](ctx, log, db, insertStmt, data)
	if err != nil {
		return
	}
	log.Debug("complete.")
	return
}
