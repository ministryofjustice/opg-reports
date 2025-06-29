package team

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Team acts as a top level group matching accounts, github repos etc to be attached
// as owners
type Team struct {
	CreatedAt string `json:"created_at,omitempty" db:"created_at" example:"2019-08-24T14:15:22Z"`
	Name      string `json:"name,omitempty" db:"name" example:"SREs"`

	// Joins
	// AwsAccounts is the one->many join from teams to awsaccounts
	AwsAccounts hasManyAwsAccounts `json:"aws_accounts,omitempty" db:"aws_accounts"`
}

// teamAwsAccount is internal and used for the join on the team model only.
//
// Matches with team.Team fields but copied over to only references values in
// the sql statements and to avoide circular references.
//
// Note: struct name is unique due to how huma generates schema.
type teamAwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"production database"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasManyAwsAccounts is used for one->many join from Team to AwsAccount
//
// Interfaces:
//   - sql.Scanner
type hasManyAwsAccounts []*teamAwsAccount

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

// TeamImport captures the team name under its prior name of
// `billing_unit`and hides the current Name property.
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
//
// We only want `billing_unit` field
type TeamImport struct {
	Name string `json:"billing_unit,omitempty" db:"name"`
}
