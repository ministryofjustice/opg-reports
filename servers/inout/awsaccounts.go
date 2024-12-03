package inout

import (
	"github.com/ministryofjustice/opg-reports/models"
)

// AwsAccountsListBody contains the resposne body to send back
// for a request to the /list endpoint
type AwsAccountsListBody struct {
	Operation string               `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *VersionUnitInput    `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsAccount `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error              `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsAccountsListResponse struct {
	Body *AwsAccountsListBody
}
