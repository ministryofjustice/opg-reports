package awscost

import (
	"context"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// defaultSeeds provides a series of known costs to be inserted into the database
func defaultSeeds() (seeds []*sqldb.BoundStatement) {
	var now = time.Now().UTC().Format(utils.DATE_FORMATS.Full)
	var date = time.Now().UTC().Format(utils.DATE_FORMATS.YMD)

	seeds = []*sqldb.BoundStatement{
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "-0.01", CreatedAt: now}, AccountID: "001A"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "S3", Date: date, Cost: "10.10", CreatedAt: now}, AccountID: "001A"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "100.57", CreatedAt: now}, AccountID: "001A"}},

		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "-50.21", CreatedAt: now}, AccountID: "001B"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-2", Service: "S3", Date: date, Cost: "603.15", CreatedAt: now}, AccountID: "001B"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "105.7", CreatedAt: now}, AccountID: "001B"}},

		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "1.02", CreatedAt: now}, AccountID: "002A"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-2", Service: "S3", Date: date, Cost: "37.00", CreatedAt: now}, AccountID: "002A"}},
		{Statement: stmtImport, Data: &awsCostImport{AwsCost: AwsCost{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "-300.68", CreatedAt: now}, AccountID: "002A"}},
	}

	return
}

// Seed populates the account tables with data passed along. Runs a delete
// on the table before inserting new seeds
//
// If seeds is nil then defaultSeeds are used instead.
func Seed(ctx context.Context, log *slog.Logger, conf *config.Config, seeds []*sqldb.BoundStatement) (inserted []*sqldb.BoundStatement, err error) {
	var (
		store *sqldb.Repository[*AwsCost]
	)

	log = log.With("operation", "Seed", "service", "awscost")
	// get default seeds
	if seeds == nil || len(seeds) <= 0 {
		seeds = defaultSeeds()
	}

	// create the store for inserting
	store, err = sqldb.New[*AwsCost](ctx, log, conf)
	if err != nil {
		return
	}
	log.Info("deleting all [awscost] entries ...")
	_, err = store.Exec(stmtDeleteAll)
	if err != nil {
		return
	}

	log.Info("inserting [awscost] seeds ...")
	err = store.Insert(seeds...)
	// if there is no error then return the data
	if err == nil {
		inserted = seeds
	}
	return
}
