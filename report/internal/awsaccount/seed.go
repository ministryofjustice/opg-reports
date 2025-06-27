package awsaccount

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

// defaultSeeds provides a series of known accounts to be inserted into the database
func defaultSeeds() (seeds []*sqldb.BoundStatement) {

	seeds = []*sqldb.BoundStatement{
		{Statement: stmtImport, Data: &awsAccountImport{AwsAccount: AwsAccount{ID: "001A", Name: "Acc01A", Label: "A", Environment: "development"}, TeamName: "TeamA"}},
		{Statement: stmtImport, Data: &awsAccountImport{AwsAccount: AwsAccount{ID: "001B", Name: "Acc01B", Label: "B", Environment: "production"}, TeamName: "TeamA"}},
		{Statement: stmtImport, Data: &awsAccountImport{AwsAccount: AwsAccount{ID: "002A", Name: "Acc02A", Label: "A", Environment: "production"}, TeamName: "TeamB"}},
	}

	return
}

// Seed populates the account tables with data passed along. Runs a delete
// on the table before inserting new seeds
//
// If seeds is nil then defaultSeeds are used instead.
func Seed(ctx context.Context, log *slog.Logger, conf *config.Config, seeds []*sqldb.BoundStatement) (inserted []*sqldb.BoundStatement, err error) {
	var (
		store *sqldb.Repository[*AwsAccount]
	)

	log = log.With("operation", "Seed", "service", "awsaccount")
	// get default seeds
	if seeds == nil || len(seeds) <= 0 {
		seeds = defaultSeeds()
	}

	// create the store for inserting
	store, err = sqldb.New[*AwsAccount](ctx, log, conf)
	if err != nil {
		return
	}
	log.Info("deleting all [awsaccount] entries ...")
	_, err = store.Exec(stmtDeleteAll)
	if err != nil {
		return
	}

	log.Info("inserting [awsaccount] seeds ...")
	err = store.Insert(seeds...)
	// if there is no error then return the data
	if err == nil {
		inserted = seeds
	}
	return
}
