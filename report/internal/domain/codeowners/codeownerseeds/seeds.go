package codeownerseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codeowners/codeownerimports"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed codeowners call failed with an error.")
var seeds []*codeownermodels.Codeowner

func init() {
	seeds = []*codeownermodels.Codeowner{
		{TeamName: "TEAM-B", CodebaseFullName: "mock/codebase-A", Name: "mock-github-team-a"},
		{TeamName: "TEAM-B", CodebaseFullName: "mock/codebase-B", Name: "mock-github-team-a"},
		{TeamName: "TEAM-C", CodebaseFullName: "mock/codebase-B", Name: "mock-github-team-b"},
		{TeamName: "TEAM-D", CodebaseFullName: "mock/codebase-C", Name: "mock-github-team-c"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*codeownermodels.Codeowner, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.codeowners.codeownerseeds.Seed")

	lg.Debug("starting ...")
	statements, err = codeownerimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import.", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.Debug("complete")
	return

}
