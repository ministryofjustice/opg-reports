package awsaccounts

import (
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
)

// AwsAccount is an API version of the database record - so
// some fields are removed / ignored for sanitisation
type AwsAccount struct {
	awsaccount.AwsAccount
	TeamID    int    `json:"-"` // hide the team id from any output
	CreatedAt string `json:"-"` // its not in the select, but blank the field incase
}

// GetAwsAccountsAllResponse
type GetAwsAccountsAllResponse struct {
	Body struct {
		Count int           `json:"count,omitempty"`
		Data  []*AwsAccount `json:"data"`
	}
}
