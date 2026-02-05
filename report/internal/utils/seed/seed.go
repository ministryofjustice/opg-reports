package seed

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/accounts/accountseeds"
	"opg-reports/report/internal/domain/codebases/codebaseseeds"
	"opg-reports/report/internal/domain/codeowners/codeownerseeds"
	"opg-reports/report/internal/domain/infracosts/infracostseeds"
	"opg-reports/report/internal/domain/teams/teamseeds"
	"opg-reports/report/internal/domain/uptime/uptimeseeds"

	"github.com/jmoiron/sqlx"
)

func SeedDB(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {
	var lg *slog.Logger = log.With("func", "seed.SeedDB")

	lg.Info("starting seed command ...")
	// // migrate the database
	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}

	// seed teams
	_, err = teamseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed accounts
	_, err = accountseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed codebases
	_, err = codebaseseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed infracosts
	_, err = infracostseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed uptime
	_, err = uptimeseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed codeowners
	_, err = codeownerseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}

	lg.Info("complete.")
	return
}
