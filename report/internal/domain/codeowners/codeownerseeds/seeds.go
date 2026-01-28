package codeownerseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codeowners/codeownerimports"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed codeowners call failed with an error.")
var seeds []*codeownermodels.Codeowner

func init() {
	seeds = []*codeownermodels.Codeowner{
		{TeamName: "TEAM-A", CodebaseFullName: "owner/codebase-A", Name: "mock-owner-a"},
		{TeamName: "TEAM-A", CodebaseFullName: "owner/codebase-B", Name: "mock-owner-a"},
		{TeamName: "TEAM-B", CodebaseFullName: "owner/codebase-B", Name: "mock-owner-b"},
		{TeamName: "TEAM-C", CodebaseFullName: "owner/codebase-C", Name: "mock-owner-c"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstatements.InsertStatement[*codeownermodels.Codeowner, int], err error) {

	log = log.With("package", "codeowners", "func", "Seed")
	log.Debug("starting ...")

	statements, err = codeownerimports.Import(ctx, log, db, seeds)
	if err != nil {
		log.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	log.Debug("complete")
	return

}
