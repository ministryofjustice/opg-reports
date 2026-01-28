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

	log = log.With("package", "codebases", "func", "Seed")
	log.Debug("starting ...")

	statements, err = codebaseimports.Import(ctx, log, db, seeds)
	if err != nil {
		log.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	log.Debug("complete")
	return

}
