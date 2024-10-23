package costsapi

import "github.com/ministryofjustice/opg-reports/costs"

type StandardBody struct {
	Type           string         `json:"type" doc:"States what type of data this is for front end handling"`
	Result         []*costs.Cost  `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	OrderedColumns []string       `json:"ordered_columns" doc:"List of columns set in the order they should be rendered for each row"`
	Request        *StandardInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
}
type StandardResult struct {
	Body *StandardBody
}

type TotalBody struct {
	Type    string      `json:"type" doc:"States what type of data this is for front end handling"`
	Request *TotalInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result  float64     `json:"result" doc:"The total sum of all costs as a float without currency." example:"1357.7861"`
}
type TotalResult struct {
	Body *TotalBody
}

type TaxOverviewBody struct {
	Type           string            `json:"type" doc:"States what type of data this is for front end handling"`
	Request        *TaxOverviewInput `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result         []*costs.Cost     `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	OrderedColumns []string          `json:"ordered_columns" doc:"List of columns set in the order they should be rendered for each row"`
}
type TaxOverviewResult struct {
	Body *TaxOverviewBody
}
