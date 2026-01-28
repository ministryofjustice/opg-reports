package accountimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/accounts/accountmodels"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("account import failed with error")

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

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
//
// `data` is presumed to come from the account.GetAwsAccountData
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*accountmodels.AwsAccount) (statements []*dbstatements.InsertStatement[*accountmodels.AwsAccount, string], err error) {

	statements = []*dbstatements.InsertStatement[*accountmodels.AwsAccount, string]{}
	log = log.With("package", "accounts", "func", "Import")

	log.Debug("starting ...")
	log.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstatements.InsertStatement[*accountmodels.AwsAccount, string]{
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
