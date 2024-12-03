package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/cols"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/dbs/statements"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/api/inputs"
)

var (
	AwsCostsSegment string   = "aws/costs"
	AwsCostsTags    []string = []string{"aws", "costs"}
)

type AwsCostsTotalBody struct {
	Operation string                     `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.DateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost          `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                    `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsTotalResponse struct {
	Body *AwsCostsTotalBody
}

const AwsCostsTotalOperationID string = "get-aws-costs-ytd"
const AwsCostsTotalDescription string = `Returns total sum of all aws costs.`
const AwsCostsTotalSQL string = `
SELECT
	coalesce(SUM(cost), 0) as cost,
	count(aws_costs.id) as count
FROM aws_costs
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	AND service != 'Tax'
	{WHERE}
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsTotalHandler
//
// Endpoints:
//
//	/version/aws/costs/total/{start_date}/{end_date}?unit=<unit>
func ApiAwsCostsTotalHandler(ctx context.Context, input *inputs.DateRangeUnitInput) (response *AwsCostsTotalResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost  = []*models.AwsCost{}
		dbPath  string             = ctx.Value(dbPathKey).(string)
		where   string             = ""
		replace string             = "{WHERE}"
		sqlStmt string             = AwsCostsTotalSQL
		param   statements.Named   = input
		body    *AwsCostsTotalBody = &AwsCostsTotalBody{
			Request:   input,
			Operation: AwsCostsTotalOperationID,
		}
	)
	// setup response
	response = &AwsCostsTotalResponse{}

	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs total adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs total select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs total selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

type AwsCostsListBody struct {
	Operation string                     `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.DateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost          `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                    `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsListResponse struct {
	Body *AwsCostsListBody
}

const AwsCostsListOperationID string = "get-aws-costs-list"
const AwsCostsListDescription string = `Returns all aws costs between start and end dates`
const awsCostsListSQL string = `
SELECT
	aws_costs.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit,
	json_object(
		'id', aws_accounts.id,
		'number', aws_accounts.number,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment,
		'unit_id', aws_accounts.unit_id
	) as aws_account
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	{WHERE}
GROUP BY aws_costs.id
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsListHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/aws/costs/list/{start_date}/{end_date}?unit=<unit>
func ApiAwsCostsListHandler(ctx context.Context, input *inputs.DateRangeUnitInput) (response *AwsCostsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost = []*models.AwsCost{}
		dbPath  string            = ctx.Value(dbPathKey).(string)
		where   string            = ""
		replace string            = "{WHERE}"
		sqlStmt string            = awsCostsListSQL
		param   statements.Named  = input
		body    *AwsCostsListBody = &AwsCostsListBody{
			Request:   input,
			Operation: AwsCostsListOperationID,
		}
	)
	// setup response
	response = &AwsCostsListResponse{}

	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

type AwsCostsTaxesBody struct {
	Operation    string                                    `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                         `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                                  `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                                  `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}                  `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                                   `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsTaxesResponse struct {
	Body *AwsCostsTaxesBody
}

const AwsCostsTaxesOperationID string = "get-aws-costs-tax"
const AwsCostsTaxesDescription string = `Returns totals per interval, with and without tax`
const AwsCostsTaxesSQL string = `
SELECT
    'Including Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs as incTax
WHERE
    incTax.date >= :start_date
    AND incTax.date < :end_date
	{WHERE}
GROUP BY strftime(:date_format, incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs as excTax
WHERE
    excTax.date >= :start_date
    AND excTax.date < :end_date
	AND excTax.service != 'Tax'
	{WHERE}
GROUP BY strftime(:date_format, date)
ORDER by date ASC
;
`

// ApiAwsCostsTaxesHandler
//
// Endpoints:
//
//	/version/aws/costs/tax/{interval}/{start_date}/{end_date}?unit=<unit>
func ApiAwsCostsTaxesHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeUnitInput) (response *AwsCostsTaxesResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost  = []*models.AwsCost{}
		dbPath  string             = ctx.Value(dbPathKey).(string)
		where   string             = ""
		replace string             = "{WHERE}"
		sqlStmt string             = AwsCostsTaxesSQL
		param   statements.Named   = input
		body    *AwsCostsTaxesBody = &AwsCostsTaxesBody{
			Request:     input,
			Operation:   AwsCostsTaxesOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"service"},
		}
	)
	// setup response
	response = &AwsCostsTaxesResponse{}

	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs tax adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs tax select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs tax selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

type AwsCostsSumBody struct {
	Operation    string                                    `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                         `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                                  `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                                  `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}                  `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                                   `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsSumResponse struct {
	Body *AwsCostsSumBody
}

const AwsCostsSumOperationID string = "get-aws-costs-sum"
const AwsCostsSumDescription string = `Returns sum of aws costs between start and end dates`
const awsCostsSumSQL string = `
SELECT
	'Total' as unit_name,
	coalesce(SUM(cost), 0) as cost,
	count(DISTINCT aws_costs.id) as count,
    strftime(:date_format, aws_costs.date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	AND service != 'Tax'
	{WHERE}
GROUP BY strftime(:date_format, aws_costs.date)
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsSumHandler
//
// Endpoints:
//
//	/version/aws/costs/sum/{interval}/{start_date}/{end_date}?unit=<unit>
func ApiAwsCostsSumHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeUnitInput) (response *AwsCostsSumResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost = []*models.AwsCost{}
		dbPath  string            = ctx.Value(dbPathKey).(string)
		where   string            = ""
		replace string            = "{WHERE}"
		sqlStmt string            = awsCostsSumSQL
		param   statements.Named  = input
		body    *AwsCostsSumBody  = &AwsCostsSumBody{
			Request:     input,
			Operation:   AwsCostsSumOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
		}
	)
	// setup response
	response = &AwsCostsSumResponse{}
	// setup the sql - if unit is set in the input, add where for it
	// otherwise remove it
	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs sum adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs sum select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs sum selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

type AwsCostsSumPerUnitBody struct {
	Operation    string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                     `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                              `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                              `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}              `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsSumPerUnitResponse struct {
	Body *AwsCostsSumPerUnitBody
}

const AwsCostsSumPerUnitOperationID string = "get-aws-costs-sum-per-unit"
const AwsCostsSumPerUnitDescription string = `Returns sum of aws costs between start and end dates grouped by unit.`
const awsCostsSumPerUnitSQL string = `
SELECT
	units.name as unit_name,
	count(DISTINCT aws_costs.id) as count,
	coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, aws_costs.date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	AND service != 'Tax'
GROUP BY units.id, strftime(:date_format, aws_costs.date)
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsSumPerUnitHandler
//
// Endpoints:
//
//	/version/aws/costs/sum-per-unit/{interval}/{start_date}/{end_date}
func ApiAwsCostsSumPerUnitHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeInput) (response *AwsCostsSumPerUnitResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost       = []*models.AwsCost{}
		dbPath  string                  = ctx.Value(dbPathKey).(string)
		sqlStmt string                  = awsCostsSumPerUnitSQL
		param   statements.Named        = input
		body    *AwsCostsSumPerUnitBody = &AwsCostsSumPerUnitBody{
			Request:     input,
			Operation:   AwsCostsSumPerUnitOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
		}
	)
	// setup response
	response = &AwsCostsSumPerUnitResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs sum per unit adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs sum per unit select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs sum per unit selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

type AwsCostsSumPerUnitEnvBody struct {
	Operation    string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                     `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                              `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                              `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}              `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsSumPerUnitEnvResponse struct {
	Body *AwsCostsSumPerUnitEnvBody
}

const AwsCostsSumPerUnitEnvOperationID string = "get-aws-costs-sum-per-unit-env"
const AwsCostsSumPerUnitEnvDescription string = `Returns sum of aws costs between start and end dates grouped by unit and environment.`
const awsCostsSumPerUnitEnvSQL string = `
SELECT
	units.name as unit_name,
	aws_accounts.environment as aws_account_environment,
	count(DISTINCT aws_costs.id) as count,
	coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, aws_costs.date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	AND service != 'Tax'
GROUP BY units.id, aws_accounts.environment, strftime(:date_format, aws_costs.date)
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsSumPerUnitEnvHandler
//
// Endpoints:
//
//	/version/aws/costs/sum-per-unit-env/{interval}/{start_date}/{end_date}
func ApiAwsCostsSumPerUnitEnvHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeInput) (response *AwsCostsSumPerUnitEnvResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost          = []*models.AwsCost{}
		dbPath  string                     = ctx.Value(dbPathKey).(string)
		sqlStmt string                     = awsCostsSumPerUnitEnvSQL
		param   statements.Named           = input
		body    *AwsCostsSumPerUnitEnvBody = &AwsCostsSumPerUnitEnvBody{
			Request:     input,
			Operation:   AwsCostsSumPerUnitEnvOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name", "aws_account_environment"},
		}
	)
	// setup response
	response = &AwsCostsSumPerUnitEnvResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws costs sum per unit env adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws costs sum per unit env select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws costs sum per unit env selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

type AwsCostsSumFullDetailsBody struct {
	Operation    string                                    `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsCost                         `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                                  `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                                  `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}                  `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                                   `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsSumFullDetailsResponse struct {
	Body *AwsCostsSumFullDetailsBody
}

const AwsCostsSumFullDetailsOperationID string = "get-aws-costs-sum-full-details"
const AwsCostsSumFullDetailsDescription string = `Returns sum of AWS costs between start and end dates detailed to AWS service level.`
const awsCostsSumFullDetailsSQL string = `
SELECT
	aws_costs.service,
	aws_costs.region,
	aws_accounts.environment as aws_account_environment,
	aws_accounts.number as aws_account_number,
	units.name as unit_name,
	count(DISTINCT aws_costs.id) as count,
	coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, aws_costs.date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_costs.date >= :start_date
	AND aws_costs.date < :end_date
	AND service != 'Tax'
	{WHERE}
GROUP BY
	aws_costs.service,
	aws_costs.region,
	aws_accounts.number,
	aws_accounts.environment,
	units.id,
	strftime(:date_format, aws_costs.date)
ORDER BY aws_costs.date ASC;
`

// ApiAwsCostsSumFullDetailsHandler
//
// Endpoints:
//
//	/version/aws/costs/sum-detailed/{interval}/{start_date}/{end_date}
func ApiAwsCostsSumFullDetailsHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeUnitInput) (response *AwsCostsSumFullDetailsResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsCost           = []*models.AwsCost{}
		dbPath  string                      = ctx.Value(dbPathKey).(string)
		sqlStmt string                      = awsCostsSumFullDetailsSQL
		param   statements.Named            = input
		where   string                      = ""
		replace string                      = "{WHERE}"
		body    *AwsCostsSumFullDetailsBody = &AwsCostsSumFullDetailsBody{
			Request:     input,
			Operation:   AwsCostsSumFullDetailsOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name", "aws_account_environment", "aws_account_number", "service", "region"},
		}
	)
	// setup response
	response = &AwsCostsSumFullDetailsResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws cost details adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()

	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}

	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsCost](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws cost details select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws cost details selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

func RegisterAwsCosts(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + AwsCostsSegment + "/total/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsTotalOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Total AWS costs between dates.",
		Description:   AwsCostsTotalDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsTotalHandler)

	uri = "/{version}/" + AwsCostsSegment + "/list/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List AWS costs",
		Description:   AwsCostsListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsListHandler)

	uri = "/{version}/" + AwsCostsSegment + "/tax/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsTaxesOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List AWS costs grouped by period with and without tax",
		Description:   AwsCostsTaxesDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsTaxesHandler)

	uri = "/{version}/" + AwsCostsSegment + "/sum/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsSumOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Sum of AWS costs",
		Description:   AwsCostsSumDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsSumHandler)

	uri = "/{version}/" + AwsCostsSegment + "/sum-per-unit/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsSumPerUnitOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Sum of AWS costs per unit",
		Description:   AwsCostsSumPerUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsSumPerUnitHandler)

	uri = "/{version}/" + AwsCostsSegment + "/sum-per-unit-env/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsSumPerUnitEnvOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Sum of AWS costs per unit and environment",
		Description:   AwsCostsSumPerUnitEnvDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsSumPerUnitEnvHandler)

	uri = "/{version}/" + AwsCostsSegment + "/sum-detailed/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsCostsSumFullDetailsOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Detailed sum of AWS costs",
		Description:   AwsCostsSumFullDetailsDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsSumFullDetailsHandler)

}
