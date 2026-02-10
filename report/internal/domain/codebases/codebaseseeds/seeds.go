package codebaseseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebaseimports"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed codebases call failed with an error.")
var seeds []*codebasemodels.Codebase

func init() {
	seeds = []*codebasemodels.Codebase{
		{Name: "codebase-A", FullName: "mock/codebase-A"},
		{Name: "codebase-B", FullName: "mock/codebase-B"},
		{Name: "codebase-C", FullName: "mock/codebase-C"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*codebasemodels.Codebase, int], err error) {
	var lg *slog.Logger = log.With("func", "codebaseseeds.Seed")

	lg.Debug("starting ...")
	statements, err = codebaseimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import.", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.Debug("complete.")
	return

}
