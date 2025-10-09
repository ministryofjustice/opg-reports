package seed

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// stmtAwsUptimeSeed is used to insert records into the table
// for seed / fixture data
const stmtAwsUptimeSeed string = `
INSERT INTO aws_uptime (
	date,
	average,
	aws_account_id
) VALUES (
	:date,
	:average,
	:aws_account_id
) ON CONFLICT (aws_account_id,date)
 	DO UPDATE SET average=excluded.average
RETURNING id;`

type awsUptimeSeed struct {
	ID        string `json:"id,omitempty" db:"id"`            // This is the AWS Account ID as a string
	Date      string `json:"date,omitempty" db:"date" `       // The data the uptime value was for
	Average   string `json:"average,omitempty" db:"average" ` // uptime average as a percentage
	AccountID string `json:"account_id" db:"aws_account_id"`
}

var uptimeDate = utils.Month(-1)

var awsUptimeSeeds = []*sqlr.BoundStatement{
	{Statement: stmtAwsUptimeSeed, Data: &awsUptimeSeed{Date: uptimeDate, Average: "99.97", AccountID: "001A"}},
	{Statement: stmtAwsUptimeSeed, Data: &awsUptimeSeed{Date: uptimeDate, Average: "99.96", AccountID: "001B"}},
	{Statement: stmtAwsUptimeSeed, Data: &awsUptimeSeed{Date: uptimeDate, Average: "99.0", AccountID: "002A"}},
	{Statement: stmtAwsUptimeSeed, Data: &awsUptimeSeed{Date: uptimeDate, Average: "90.7", AccountID: "003A"}},
	{Statement: stmtAwsUptimeSeed, Data: &awsUptimeSeed{Date: uptimeDate, Average: "98.99", AccountID: "004A"}},
}

// AwsUptime populates the database (via the sqc var) with standard known enteries
// that can be used for testing and development databases
func (self *Service) AwsUptime(sqc sqlr.RepositoryWriter) (results []*sqlr.BoundStatement, err error) {
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[seed:AwsUptime] seeding finished.")
	}()
	self.log.Info("[seed:AwsUptime] starting seeding ...")
	err = sqc.Insert(awsUptimeSeeds...)
	if err != nil {
		return
	}
	self.log.Info("[seed:AwsUptime] seeding successful")
	results = awsAccountSeeds
	return
}
