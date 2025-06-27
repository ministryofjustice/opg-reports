package awscost

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/s3"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Existing generates new enteries for aws_cost by downloading and importing json cost
// data from a location in an s3 bucket.
//
// We've been clollecting cost data for several years into these buckets, so this allows
// the information to be pulled and updated to be used in this version of reporting tools
//
// Using the `conf.Aws.Bucket` settings, calls s3 api to list all files within the bucket that
// match the prefix (typically a subfolder pattern). Then downloads each of those files
// to local storage (to a temp folder).
//
// S3 files are downloaded locally to a temp folder underneath the `conf.Aws.Bucket.Local` path
// which is removed on exit via a defer.
//
// Once downloaded, the each file is converted to a struct (`[]*AwsCostImport`) and merged
// with a sql statement (`stmtImport`) for insertion.
//
// All sql statements are then run in one block, using a `sqldb` repository to handle
// transaction based inserts.
//
// If any of the file->struct conversions fail then the an error is returned and no
// inserts are run.
//
// Example cost entry from a pre-existing json file:
//
//	[{
//	  "id": 0,
//	  "ts": "2024-08-15 18:52:55.055478 +0000 UTC",
//	  "organisation": "ORG",
//	  "account_id": "010256574050",
//	  "account_name": "Dev account",
//	  "unit": "Team A",
//	  "label": "dev",
//	  "environment": "development",
//	  "service": "AWS CloudTrail",
//	  "region": "ap-northeast-2",
//	  "date": "2023-07-01",
//	  "cost": "0.0002485"
//	}]
//
// interface: ImporterExistingCommand
func Existing(ctx context.Context, log *slog.Logger, conf *config.Config, service *s3.Service[*AwsCostImport]) (err error) {
	var (
		store   *sqldb.Repository[*AwsCostImport]
		bucket  string                  = conf.Aws.Buckets.Costs.Name
		prefix  string                  = conf.Aws.Buckets.Costs.Prefix
		costs   []*AwsCostImport        = []*AwsCostImport{}
		inserts []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		sw                              = utils.Stopwatch()
	)
	defer service.Close()
	// timer
	sw.Start()
	// check config values are setup, otherwise we cannot download anything, so error
	if bucket == "" || prefix == "" {
		return fmt.Errorf("required bucket details were not found.")
	}

	// add info to the logger
	log = log.With("operation", "Existing", "service", "awscost")
	log.Info("[awscost] starting existing records import ...")

	// use the service to download and convert data to the structs we want
	costs, err = service.DownloadAndReturnData(bucket, prefix)

	// for each file we need to generate the bounded sql statements
	for _, row := range costs {
		inserts = append(inserts, &sqldb.BoundStatement{Data: row, Statement: stmtImport})
	}

	log.With("count", len(inserts)).Debug("records to insert ...")

	// now insert the cost data
	log.Debug("creating datastore ...")
	store, err = sqldb.New[*AwsCostImport](ctx, log, conf)
	if err != nil {
		return
	}

	log.Debug("running insert ...")
	err = store.Insert(inserts...)
	if err != nil {
		return
	}

	log.With(
		"seconds", sw.Stop().Seconds(),
		"inserted", len(inserts)).
		Info("[awscost] existing records imported.")
	return
}
