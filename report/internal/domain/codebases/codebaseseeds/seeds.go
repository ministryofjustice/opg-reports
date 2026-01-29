package codebaseseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codebases/codebaseimports"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed codebases call failed with an error.")
var seeds []*codebasemodels.Codebase

func init() {
	seeds = []*codebasemodels.Codebase{
		{Name: "codebase-A", FullName: "owner/codebase-A"},
		{Name: "codebase-B", FullName: "owner/codebase-B"},
		{Name: "codebase-C", FullName: "owner/codebase-C"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstatements.InsertStatement[*codebasemodels.Codebase, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.codebases.codebaseseeds.Seed")

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
