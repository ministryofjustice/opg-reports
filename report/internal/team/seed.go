package team

import (
	"context"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

// defaultSeeds provides a series of known accounts to be inserted into the database
func defaultSeeds() (seeds []*sqldb.BoundStatement) {
	var now = time.Now().UTC().Format(time.RFC3339)

	seeds = []*sqldb.BoundStatement{
		{Statement: stmtImport, Data: &Team{Name: "TeamA", CreatedAt: now}},
		{Statement: stmtImport, Data: &Team{Name: "TeamB", CreatedAt: now}},
		{Statement: stmtImport, Data: &Team{Name: "TeamC", CreatedAt: now}},
	}

	return
}

// Seed populates the account tables with data passed along. Runs a delete
// on the table before inserting new seeds
//
// If seeds is nil then defaultSeeds are used instead.
func Seed(ctx context.Context, log *slog.Logger, conf *config.Config, seeds []*sqldb.BoundStatement) (inserted []*sqldb.BoundStatement, err error) {
	var (
		store *sqldb.Repository[*Team]
	)

	log = log.With("operation", "Seed", "service", "team")
	// get default seeds
	if seeds == nil || len(seeds) <= 0 {
		seeds = defaultSeeds()
	}

	// create the store for inserting
	store, err = sqldb.New[*Team](ctx, log, conf)
	if err != nil {
		return
	}
	log.Info("deleting all [team] entries ...")
	_, err = store.Exec(stmtDeleteAll)
	if err != nil {
		return
	}
	log.Info("inserting [team] seeds ...")
	err = store.Insert(seeds...)
	// if there is no error then return the data
	if err == nil {
		inserted = seeds
	}
	return
}
