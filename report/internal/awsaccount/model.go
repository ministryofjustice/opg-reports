package awsaccount

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// AwsAccount is db model representing the fields we use
// directly.
//
// Imports use AwsAccountImport and API returns awsaccounts.AwsAccount
type AwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at" example:"2019-08-24T14:15:22Z"`

	// Joins
	// TeamName is the raw db column that stores the association, should not
	TeamName string      `json:"team_name,omitempty" db:"team_name"`
	Team     *hasOneTeam `json:"team,omitempty" db:"team"`
}

// team is internal and used for handling the account->team join
type accountTeam struct {
	Name string `json:"name,omitempty" db:"name" example:"SRE"`
}

type hasOneTeam accountTeam

// Scan handles the processing of the join data
func (self *hasOneTeam) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = utils.Unmarshal(src.([]byte), self)
	case string:
		err = utils.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// AwsAccountImport captures an extra field from the metadata which
// is used in the stmtInsert to create the initial join to team based
// on the billing_unit name
//
// Example account from the opg-metadata source file:
//
//	{
//		"id": "500000067891",
//		"name": "My production",
//		"billing_unit": "Team A",
//		"label": "prod",
//		"environment": "production",
//		"type": "aws",
//		"uptime_tracking": true
//	}
type AwsAccountImport struct {
	AwsAccount
	TeamName string `json:"billing_unit,omitempty" db:"team_name"`
}
