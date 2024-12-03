package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
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

// -- Aws costs listing

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
	AND service != 'Tax'
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

type AwsCostsSumBody struct {
	Operation string                                    `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost                         `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                                   `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsCostsSumResponse struct {
	Body *AwsCostsSumBody
}

const AwsCostsSumOperationID string = "get-aws-costs-sum"
const AwsCostsSumDescription string = `Returns sum of aws costs between start and end dates`
const awsCostsSumSQL string = `
SELECT
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
			Request:   input,
			Operation: AwsCostsSumOperationID,
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
	response.Body = body
	return
}

type AwsCostsSumPerUnitBody struct {
	Operation string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.RequiredGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsCost                     `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
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
			Request:   input,
			Operation: AwsCostsSumPerUnitOperationID,
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
	response.Body = body
	return
}

func RegisterAwsCosts(api huma.API) {
	var uri string = ""

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
		Summary:       "Sum of AWS costs grouped per unit",
		Description:   AwsCostsSumPerUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsCostsTags,
	}, ApiAwsCostsSumPerUnitHandler)

}
