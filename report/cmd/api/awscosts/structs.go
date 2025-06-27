package awscosts

import (
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
)

// AwsCost is an API version of the database model with
// field sanitation
type AwsCost struct {
	awscost.AwsCost
	ID           int    `json:"-"`
	AwsAccountID string `json:"-"` // hide the team id from any output
	CreatedAt    string `json:"-"` // its not in the select, but blank the field incase
}

// GetAwsCostsTop20Response
type GetAwsCostsTop20Response struct {
	Body struct {
		Data []*AwsCost
	}
}
