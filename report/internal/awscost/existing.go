package awscost

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/s3bucket"
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
// Once downloaded, the each file is converted to a struct (`[]*awsCostImport`) and merged
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
func Existing(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var (
		s3b        *s3bucket.Repository
		store      *sqldb.Repository[*awsCostImport]
		localDir   string
		bucket     string                  = conf.Aws.Buckets.Costs.Name
		prefix     string                  = conf.Aws.Buckets.Costs.Prefix
		list       []string                = []string{}
		downloaded []string                = []string{}
		inserts    []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		sw                                 = utils.Stopwatch()
	)
	// timer
	sw.Start()
	// check config values are setup, otherwise we cannot download anything, so error
	if bucket == "" || prefix == "" {
		return fmt.Errorf("required bucket details were not found.")
	}

	// setup a temp directory
	os.MkdirAll(conf.Aws.Buckets.Local, os.ModePerm)
	localDir, _ = os.MkdirTemp(conf.Aws.Buckets.Local, "aws_costs-*")
	// clean up tmp directory, but we leave the parent
	defer os.RemoveAll(localDir)

	// add info to the logger
	log = log.With("operation", "Existing", "service", "awscost", "dir", localDir)
	log.Info("[awscost] starting existing records import ...")

	// create the s3 repository
	s3b, err = s3bucket.New(ctx, log, conf)
	if err != nil {
		return
	}
	// get list of all files from the bucket
	list, err = s3b.ListBucket(bucket, prefix)
	if err != nil {
		return
	}
	log.With("count", len(list)).Debug("found files to download ...")
	// download all the listed files
	downloaded, err = s3b.Download(bucket, list, localDir)
	if err != nil {
		return
	}
	log.With("count", len(downloaded)).Debug("files downloaded ...")

	// for each file we need to generate the bounded sql statements
	for _, file := range downloaded {
		var stmts []*sqldb.BoundStatement

		stmts, err = fileToInsertStmts(log, file)
		if err != nil {
			log.Error("failed stmt fetch", "error", err.Error())
			return
		}
		inserts = append(inserts, stmts...)
	}

	log.With("count", len(inserts)).Debug("records to insert ...")

	// now insert the cost data
	log.Debug("creating datastore ...")
	store, err = sqldb.New[*awsCostImport](ctx, log, conf)
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
		"inserted", len(inserts),
		"files", len(list),
		"downloaded", len(downloaded)).
		Info("[awscost] existing records imported.")
	return
}

// fileToInsertStmts loads the file into a slice (`[]*awsCostImport{}`) and merges each entry
// with the import sql statement, returning the resuling list.
//
// If the file->struct conversion fails then an error is returned instead
func fileToInsertStmts(log *slog.Logger, filename string) (inserts []*sqldb.BoundStatement, err error) {
	var importCosts []*awsCostImport = []*awsCostImport{}

	inserts = []*sqldb.BoundStatement{}
	log = log.With("file", filename)

	log.Debug("loading json file into struct")
	err = utils.StructFromJsonFile(filename, &importCosts)
	if err != nil {
		log.Error("failed to load", "error", err.Error())
		return
	}

	for _, row := range importCosts {
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}

	return
}
