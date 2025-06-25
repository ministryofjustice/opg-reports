package team

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Team replaces the Unit structure as a grouping of services / accounts
type Team struct {
	ID        int    `json:"id,omitempty" db:"id" example:"1"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at" example:"2019-08-24T14:15:22Z"`
	Name      string `json:"name,omitempty" db:"name" example:"Sirius"`

	// Joins
	// AwsAccounts is the one->many join from teams to awsaccounts
	AwsAccounts hasManyAwsAccounts `json:"aws_accounts,omitempty" db:"aws_accounts"`
}

// account is internal and used for the join on the team model only.
// Matches with team.Team but ccopy over to only have required fields
// and to avoide circular references
type awsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"production database"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasManyAwsAccounts is used for one->many join from Team to AwsAccount
//
// Interfaces:
//   - sql.Scanner
type hasManyAwsAccounts []*awsAccount

// Scan handles the processing of the join data
func (self *hasManyAwsAccounts) Scan(src interface{}) (err error) {
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
