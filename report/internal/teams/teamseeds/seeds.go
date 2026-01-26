package teamseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/teams/teamimports"
	"opg-reports/report/internal/teams/teammodels"

	"github.com/jmoiron/sqlx"
)

var (
	costDate string
	seeds    []*teammodels.Team
)

func init() {
	seeds = []*teammodels.Team{
		{Name: "TEAM-A"},
		{Name: "TEAM-B"},
		{Name: "TEAM-C"},
		{Name: "TEAM-D"},
		{Name: "TEAM-E"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstatements.DataStatement[*teammodels.Team, string], err error) {

	log = log.With("package", "teams", "func", "Seed")
	log.Debug("starting ...")

	statements, err = teamimports.Import(ctx, log, db, seeds)
	if err != nil {
		log.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	log.Debug("complete")
	return

}
