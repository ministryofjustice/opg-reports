package awscost

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type AwsCost struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Region    string `json:"region,omitempty" db:"region" example:"eu-west-1|eu-west-2"`
	Service   string `json:"service,omitempty" db:"service" example:"AWS service name"` // The AWS service name
	Date      string `json:"date,omitempty" db:"date" example:"2019-08-24"`             // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty" db:"cost" example:"-10.537"`                // The actual cost value as a string - without an currency, but is USD by default
	CreatedAt string `json:"created_at,omitempty" db:"created_at" example:"2019-08-24T14:15:22Z"`

	// Joins
	// AwsAccount joins
	AwsAccountID string            `json:"aws_account_id,omitempty" db:"aws_account_id"`
	AwsAccount   *hasOneAwsAccount `json:"aws_account,omitempty" db:"aws_account"`
}

// awsAccount is used to capture sql join data
type awsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasOneAwsAccount used for the join to awsAccount to handle the scaning into a struct
type hasOneAwsAccount awsAccount

// Scan handles the processing of the join data
func (self *hasOneAwsAccount) Scan(src interface{}) (err error) {
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

// awsCostImport is used for data import / seeding and contains additional data in older formats
type awsCostImport struct {
	AwsCost
	AccountID string `json:"account_id" db:"account_id"`
}
