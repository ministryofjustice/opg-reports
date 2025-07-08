package existing

import (
	"fmt"
	"os"

	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

const stmtAwsCostImport string = `
INSERT INTO aws_costs (
	region,
	service,
	date,
	cost,
	aws_account_id
) VALUES (
	:region,
	:service,
	:date,
	:cost,
	:account_id
) ON CONFLICT (aws_account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`

// awsCost is used for data import / seeding and contains additional data in older formats
//
// Example cost entry:
//
//	{
//		"id": 0,
//		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
//		"organisation": "OPG",
//		"account_id": "050116572003",
//		"account_name": "development",
//		"unit": "TeamA",
//		"label": "A",
//		"environment": "development",
//		"service": "Amazon Simple Storage Service",
//		"region": "eu-west-1",
//		"date": "2023-07-31",
//		"cost": "0.2309542206"
//	}
//
// We use the old account_id field for the join information
type awsCost struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Region    string `json:"region,omitempty" db:"region"`
	Service   string `json:"service,omitempty" db:"service"` // The AWS service name
	Date      string `json:"date,omitempty" db:"date" `      // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty" db:"cost" `      // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id" db:"account_id"`
}

// InsertAwsCosts handles the inserting of existing cost data from files
// in s3 buckets into the local database service.
//
// Example cost entry:
//
//	{
//		"id": 0,
//		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
//		"organisation": "OPG",
//		"account_id": "050116570073",
//		"account_name": "development",
//		"unit": "TeamA",
//		"label": "A",
//		"environment": "development",
//		"service": "Amazon Simple Storage Service",
//		"region": "eu-west-1",
//		"date": "2023-07-31",
//		"cost": "0.2309542206"
//	}
//
// We use the old account_id field for the join information
func (self *Service) InsertAwsCosts(client awsr.ClientS3ListAndGetter, source awsr.RepositoryS3BucketDownloader, sq sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
	var dir string
	var downloaded []string
	var totalInserted = 0
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds()).
			Info("[existing:AwsCosts] existing func finished.")
	}()
	self.log.Info("[existing:AwsCosts] starting existing records import ...")

	// check if source or sql service are empty
	if source == nil {
		err = fmt.Errorf("source was nil")
		return
	}
	if sq == nil {
		err = fmt.Errorf("sq was nil")
		return
	}
	// make temp directory locally to sync files into
	dir, err = os.MkdirTemp("./", "__download-s3-*")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	// download all the files
	downloaded, err = source.DownloadBucket(client,
		self.conf.Aws.Buckets.Costs.Name,
		self.conf.Aws.Buckets.Costs.Prefix,
		dir,
	)
	if err != nil {
		return
	}

	self.log.With("count", len(downloaded)).
		Debug("[existing:AwsCosts] downloaded files.")
	// loop over all the downloaded files, read them into a local struct
	// and then run insert against the database
	for _, file := range downloaded {
		var (
			costs   []*awsCost             = []*awsCost{}
			inserts []*sqlr.BoundStatement = []*sqlr.BoundStatement{}
			lg                             = self.log.With("file", file)
		)
		if e := utils.StructFromJsonFile(file, &costs); e != nil {
			lg.Error("error converting file", "err", e.Error())
			return nil, e
		}
		// for each file we need to generate the bounded sql statements
		for _, row := range costs {
			inserts = append(inserts, &sqlr.BoundStatement{Data: row, Statement: stmtAwsCostImport})
		}
		lg.With("count", len(inserts)).Debug("[existing:AwsCosts] inserting records from file ...")

		// run inserts
		if e := sq.Insert(inserts...); e != nil {
			return
		}
		// only merge in the items with return values
		for _, in := range inserts {
			if in.Returned != nil {
				totalInserted++
				results = append(results, in)
			}
		}
	}

	self.log.With("inserted", totalInserted).Info("[existing:AwsCosts] existing records successful")

	return
}
