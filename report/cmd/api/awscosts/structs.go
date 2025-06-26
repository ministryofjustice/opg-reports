package awscosts

import (
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
)

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
