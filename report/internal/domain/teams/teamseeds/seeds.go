package teamseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teamimports"
	"opg-reports/report/internal/domain/teams/teammodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed team call failed with an error.")
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
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*teammodels.Team, string], err error) {
	var lg *slog.Logger = log.With("func", "teamseeds.Seed")

	lg.Debug("starting ...")
	statements, err = teamimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import.", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.Debug("complete")
	return

}
