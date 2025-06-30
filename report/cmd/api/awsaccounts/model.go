package awsaccounts

import (
	"github.com/ministryofjustice/opg-reports/report/internal/service/awsaccount"
)

// AwsAccount is an API version of the database record - so
// some fields are removed / ignored for sanitisation
type AwsAccount struct {
	awsaccount.AwsAccount
	TeamName  int    `json:"-"` // hide the team entry from any output
	CreatedAt string `json:"-"` // its not in the select, but blank the field incase
}
