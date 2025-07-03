package existing

import (
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
)

const stmtAwsCostImport string = `
INSERT INTO aws_costs (
	region,
	service,
	date,
	cost,
	aws_account_id
) SELECT
	:region,
	:service,
	:date,
	:cost,
	id
FROM aws_accounts WHERE aws_accounts.id = :account_id
ON CONFLICT (aws_account_id,date,region,service)
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
//		"account_id": "050116572273",
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

// InsertAwsCosts handles the inserting otf team data from opgmetadata reository
// into the local database service.
// Example cost entry:
//
//	{
//		"id": 0,
//		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
//		"organisation": "OPG",
//		"account_id": "050116572273",
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
func (self *Service) InsertAwsCosts(client awsr.ClientS3, sq sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
	return
}
