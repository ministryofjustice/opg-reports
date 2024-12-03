package inout

import "github.com/ministryofjustice/opg-reports/models"

type AwsCostsTotalBody struct {
	Operation string              `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *DateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost   `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error             `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

type AwsCostsTotalResponse struct {
	Body *AwsCostsTotalBody
}

type AwsCostsListBody struct {
	Operation string              `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *DateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost   `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error             `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

type AwsCostsListResponse struct {
	Body *AwsCostsListBody
}

// AwsCostsTaxesBody
// Tabular
type AwsCostsTaxesBody struct {
	Operation    string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                  `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                           `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                           `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}           `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{}  `json:"-"` // Used for post processing
}

type AwsCostsTaxesResponse struct {
	Body *AwsCostsTaxesBody
}

type AwsCostsSumBody struct {
	Operation    string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                  `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                           `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                           `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}           `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{}  `json:"-"` // Used for post processing
}

type AwsCostsSumResponse struct {
	Body *AwsCostsSumBody
}

type AwsCostsSumPerUnitBody struct {
	Operation    string                            `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeInput    `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                 `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                          `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                          `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                           `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{} `json:"-"` // Used for post processing
}

type AwsCostsSumPerUnitResponse struct {
	Body *AwsCostsSumPerUnitBody
}

type AwsCostsSumPerUnitEnvBody struct {
	Operation    string                            `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeInput    `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                 `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                          `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                          `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                           `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{} `json:"-"` // Used for post processing
}

type AwsCostsSumPerUnitEnvResponse struct {
	Body *AwsCostsSumPerUnitEnvBody
}

type AwsCostsSumFullDetailsBody struct {
	Operation    string                             `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                  `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                           `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                           `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}           `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                            `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
	TableRows    map[string]map[string]interface{}  `json:"-"` // Used for post processing
}

type AwsCostsSumFullDetailsResponse struct {
	Body *AwsCostsSumFullDetailsBody
}
