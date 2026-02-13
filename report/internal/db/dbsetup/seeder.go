package dbsetup

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/domain/uptime/uptimemodels"

	"github.com/jmoiron/sqlx"
)

// SeedAll populates the database with generated data
func SeedAll(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	var (
		teams      []*teammodels.Team
		accounts   []*accountmodels.Account
		infracosts []*infracostmodels.Cost
		uptime     []*uptimemodels.Uptime
		codebases  []*codebasemodels.Codebase
		codeowners []*codeownermodels.Codeowner
		importStmt *ImportStatements
		lg         = log.With("func", "dbsetup.SeedAll")
	)

	lg.Debug("starting ...")

	lg.Debug("migrating...")
	err = Migrate(ctx, lg, db)
	if err != nil {
		return
	}

	// seed teams
	lg.Info("seeding teams ...")
	teams = generateTeams(5)
	importStmt = _IMPORTS["teams"]
	_, err = Import[string](ctx, log, db, teams, importStmt)
	if err != nil {
		return
	}
	// seed accounts
	lg.Info("seeding accounts ...")
	accounts = generateAccounts(50, teams)
	importStmt = _IMPORTS["accounts"]
	_, err = Import[string](ctx, log, db, accounts, importStmt)
	if err != nil {
		return
	}
	// costs
	lg.Info("seeding infracosts ...")
	infracosts = generateInfracosts(13000, accounts)
	importStmt = _IMPORTS["infracosts"]
	_, err = Import[int](ctx, log, db, infracosts, importStmt)
	if err != nil {
		return
	}
	// uptime
	lg.Info("seeding uptime ...")
	uptime = generateUptime(8000, accounts)
	importStmt = _IMPORTS["uptime"]
	_, err = Import[int](ctx, log, db, uptime, importStmt)
	if err != nil {
		return
	}

	lg.Info("seeding codebases ...")
	codebases = generateCodebases(50)
	importStmt = _IMPORTS["codebases"]
	_, err = Import[int](ctx, log, db, codebases, importStmt)
	if err != nil {
		return
	}

	lg.Info("seeding codeowners ...")
	codeowners = generateCodeowners(75, teams, codebases)
	importStmt = _IMPORTS["codeowners"]
	_, err = Import[int](ctx, log, db, codeowners, importStmt)
	if err != nil {
		return
	}

	lg.Debug("complete.")
	return
}
