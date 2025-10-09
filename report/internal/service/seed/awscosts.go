package seed

import (
	"time"

	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// stmtAwsCostSeed is used to insert records into the team table
// for seed / fixture data
const stmtAwsCostSeed string = `
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
	:aws_account_id
) ON CONFLICT (aws_account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`

type awsCostSeed struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Region    string `json:"region,omitempty" db:"region"`
	Service   string `json:"service,omitempty" db:"service"` // The AWS service name
	Date      string `json:"date,omitempty" db:"date" `      // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty" db:"cost" `      // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id" db:"aws_account_id"`
}

var date = time.Now().AddDate(0, -1, 0).UTC().Format(utils.DATE_FORMATS.YMD)
var awsCostSeeds = []*sqlr.BoundStatement{
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "-0.01", AccountID: "001A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "S3", Date: date, Cost: "10.10", AccountID: "001A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "100.57", AccountID: "001A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "SQS", Date: date, Cost: "23.01", AccountID: "001A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-2", Service: "IAM", Date: date, Cost: "0.002", AccountID: "001A"}},

	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "-50.21", AccountID: "001B"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-2", Service: "S3", Date: date, Cost: "603.15", AccountID: "001B"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "105.7", AccountID: "001B"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "R53", Date: date, Cost: "1.7001", AccountID: "001B"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "us-west-1", Service: "EKS", Date: date, Cost: "27501.88", AccountID: "001B"}},

	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "1.02", AccountID: "002A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-2", Service: "S3", Date: date, Cost: "37.00", AccountID: "002A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "RDS", Date: date, Cost: "-300.68", AccountID: "002A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "SNS", Date: date, Cost: "103.51", AccountID: "002A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-2", Service: "RDS", Date: date, Cost: "502.44", AccountID: "002A"}},

	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "ECS", Date: date, Cost: "102.44", AccountID: "003A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-2", Service: "S3", Date: date, Cost: "7.0012", AccountID: "003A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "S3", Date: date, Cost: "96.35", AccountID: "003A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "eu-west-1", Service: "SNS", Date: date, Cost: "18.19", AccountID: "003A"}},
	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "us-west-1", Service: "S3", Date: date, Cost: "2.4474", AccountID: "003A"}},

	{Statement: stmtAwsCostSeed, Data: &awsCostSeed{Region: "us-west-1", Service: "S3", Date: date, Cost: "102.7409", AccountID: "004A"}},
}

// AwsCosts populates the database (via the sqc var) with standard known enteries
// that can be used for testing and development databases
func (self *Service) AwsCosts(sqc sqlr.RepositoryWriter) (results []*sqlr.BoundStatement, err error) {
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[seed:AwsCosts] seeding finished.")
	}()
	self.log.Info("[seed:AwsCosts] starting seeding ...")
	err = sqc.Insert(awsCostSeeds...)
	if err != nil {
		return
	}
	self.log.Info("[seed:AwsCosts] seeding successful")
	results = awsCostSeeds
	return
}
