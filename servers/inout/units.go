package inout

import "github.com/ministryofjustice/opg-reports/models"

type UnitsListBody struct {
	Operation string         `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *VersionInput  `json:"request,omitempty" doc:"the original request"`
	Result    []*models.Unit `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error        `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}
type UnitsListResponse struct {
	Body *UnitsListBody
}
