package seed

import (
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// stmtAwsAccountSeed is used to insert records into the team table
// for seed / fixture data
const stmtAwsAccountSeed string = `
INSERT INTO aws_accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:team_name
)
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id;`

type awsAccountSeed struct {
	ID          string `json:"id,omitempty" db:"id"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name"`
	Label       string `json:"label,omitempty" db:"label"`
	Environment string `json:"environment,omitempty" db:"environment"`
	TeamName    string `json:"team_name,omitempty" db:"team_name"`
}

var awsAccountSeeds = []*sqlr.BoundStatement{
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "001A", Name: "Account 1A", Label: "A", Environment: "development", TeamName: "TEAM-A"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "001B", Name: "Account 1B", Label: "B", Environment: "production", TeamName: "TEAM-A"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "002A", Name: "Account 2A", Label: "A", Environment: "production", TeamName: "TEAM-B"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "003A", Name: "Account 3A", Label: "A", Environment: "development", TeamName: "TEAM-C"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "003B", Name: "Account 3B", Label: "B", Environment: "production", TeamName: "TEAM-C"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "004A", Name: "Account 4A", Label: "A", Environment: "production", TeamName: "TEAM-D"}},
	{Statement: stmtAwsAccountSeed, Data: &awsAccountSeed{ID: "004B", Name: "Account 4B", Label: "B", Environment: "production", TeamName: "TEAM-D"}},
}

// AwsAccounts populates the database (via the sqc var) with standard known enteries
// that can be used for testing and development databases
func (self *Service) AwsAccounts(sqc sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[seed:AwsAccounts] seeding finished.")
	}()
	self.log.Info("[seed:AwsAccounts] starting seeding ...")
	err = sqc.Insert(awsAccountSeeds...)
	if err != nil {
		return
	}
	self.log.Info("[seed:AwsAccounts] seeding successful")
	results = awsAccountSeeds
	return
}
