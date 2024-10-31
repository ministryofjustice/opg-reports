package costsapi

import "github.com/ministryofjustice/opg-reports/sources/costs"

type StandardBody struct {
	Type         string                   `json:"type" doc:"States what type of data this is for front end handling"`
	Result       []*costs.Cost            `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	ColumnOrder  []string                 `json:"column_order" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{} `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Request      *StandardInput           `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	DateRange    []string                 `json:"date_range" doc:"list of string dates between the start and end date"`
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
	Type         string                   `json:"type" doc:"States what type of data this is for front end handling"`
	Request      *TaxOverviewInput        `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result       []*costs.Cost            `json:"result" doc:"List of call costs grouped by interval for with and without tax costs."`
	ColumnOrder  []string                 `json:"column_order" doc:"List of columns set in the order they should be rendered for each row"`
	ColumnValues map[string][]interface{} `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	DateRange    []string                 `json:"date_range" doc:"list of string dates between the start and end date"`
}
type TaxOverviewResult struct {
	Body *TaxOverviewBody
}
