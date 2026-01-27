package accountseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/accounts/accountimports"
	"opg-reports/report/internal/accounts/accountmodels"
	"opg-reports/report/internal/db/dbstatements"

	"github.com/jmoiron/sqlx"
)

var seeds []*accountmodels.AwsAccount

func init() {
	seeds = []*accountmodels.AwsAccount{
		{ID: "001A", Name: "Account 1A", Label: "A", Environment: "development", TeamName: "TEAM-A"},
		{ID: "001B", Name: "Account 1B", Label: "B", Environment: "production", TeamName: "TEAM-A"},
		{ID: "002A", Name: "Account 2A", Label: "A", Environment: "production", TeamName: "TEAM-B"},
		{ID: "003A", Name: "Account 3A", Label: "A", Environment: "development", TeamName: "TEAM-C"},
		{ID: "003B", Name: "Account 3B", Label: "B", Environment: "production", TeamName: "TEAM-C"},
		{ID: "004A", Name: "Account 4A", Label: "A", Environment: "production", TeamName: "TEAM-D"},
		{ID: "004B", Name: "Account 4B", Label: "B", Environment: "production", TeamName: "TEAM-D"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstatements.DataStatement[*accountmodels.AwsAccount, string], err error) {

	log = log.With("package", "accounts", "func", "Seed")
	log.Debug("starting ...")

	statements, err = accountimports.Import(ctx, log, db, seeds)
	if err != nil {
		log.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	log.Debug("complete")
	return

}
