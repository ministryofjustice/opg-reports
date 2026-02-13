package dbsetup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/domain/uptime/uptimemodels"

	"github.com/jmoiron/sqlx"
)

var (
	ErrPreFailed       = errors.New("import pre-statement failed with error.")
	ErrImportFailed    = errors.New("import failed with error.")
	ErrUnsupportedType = errors.New("data type unsupported.")
)

// ImportStatements are run combined with the data, the `Pre` attribute
// allows standalone sql (like truncates) to be run before insert
type ImportStatements struct {
	Pre    string
	Insert string
}

// Import writes data into the database connection provided, generating sql statement and returning the result.
//
// If `statements` is nil then we try to identify the correct statement to use from based on `T`
func Import[R any, T any](ctx context.Context, log *slog.Logger, db *sqlx.DB, data []T, statements *ImportStatements) (inserted []*dbstmts.Insert[T, R], err error) {
	var lg *slog.Logger = log.With("func", "dbimports.Import", "types", fmt.Sprintf("%T", inserted))

	lg.Debug("starting ...")
	inserted = []*dbstmts.Insert[T, R]{}

	// if there is no statement, try t ofind it
	if statements == nil {
		statements, err = statementFromT(data)
	}
	if err != nil {
		return
	}
	// run any pre statement - generally a truncate
	if statements.Pre != "" {
		lg.Debug("running pre statement ...")
		_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(statements.Pre))
	}
	// check for errors
	if err != nil {
		lg.Error("pre statement failed.")
		err = errors.Join(ErrPreFailed, err)
		return
	}
	// generate insert statements for the data rows passed
	lg.Debug("generating db insert statements ...")
	for _, row := range data {
		inserted = append(inserted, &dbstmts.Insert[T, R]{
			Statement: statements.Insert,
			Data:      row,
		})
	}
	lg.Debug("running import statements via insert ...")
	err = dbinserts.Insert(ctx, log, db, inserted...)
	if err != nil {
		lg.Error("error with insert.", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return
}

// statementFromT checks the type of data and
func statementFromT[T any](data []T) (stmt *ImportStatements, err error) {
	var t T

	switch any(t).(type) {
	case *accountmodels.Account:
		stmt = _IMPORTS["accounts"]
	case *codebasemodels.Codebase:
		stmt = _IMPORTS["codebases"]
	case *codeownermodels.Codeowner:
		stmt = _IMPORTS["codeowners"]
	case *infracostmodels.Cost:
		stmt = _IMPORTS["infracosts"]
	case *teammodels.Team:
		stmt = _IMPORTS["teams"]
	case *uptimemodels.Uptime:
		stmt = _IMPORTS["uptime"]
	default:
		err = errors.Join(ErrUnsupportedType, fmt.Errorf("data type [%T] is not supported.\n", data))
	}

	return

}
