package awscosts

import (
	"github.com/ministryofjustice/opg-reports/report/internal/service/awscost"
)

// AwsCost is an API output for the top20 endpoint
type AwsCost struct {
	awscost.AwsCost
	ID           int    `json:"-"`
	AwsAccountID string `json:"-"` // hide the team id from any output
	CreatedAt    string `json:"-"` // its not in the select, but blank the field incase
}

// AwsGroupedCost is the API output of the grouped cost query and endpoint
type AwsGroupedCost struct {
	Cost         string `json:"cost,omitempty" db:"cost"`
	Date         string `json:"date,omitempty" db:"date"`
	Region       string `json:"region,omitempty" db:"region"`
	Service      string `json:"service,omitempty" db:"service"`
	Team         string `json:"team,omitempty" db:"team_name"`
	Account      string `json:"account_id,omitempty" db:"aws_account_id"`
	AccountName  string `json:"account_name,omitempty" db:"account_name"`
	AccountLabel string `json:"account_label,omitempty" db:"account_label"`
	Environment  string `json:"environment,omitempty" db:"environment"`
}
