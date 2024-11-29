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
	AwsUptimeSegment string   = "aws/uptime"
	AwsUptimeTags    []string = []string{"aws", "uptime"}
)

// AwsUptimeListBody contains the resposne body to send back
// for a request to the /list endpoint
type AwsUptimeListBody struct {
	Operation string                     `json:"operation,omitempty" doc:"contains the operation id"`
	Request   *inputs.DateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result    []*models.AwsUptime        `json:"result,omitempty" doc:"list of all units returned by the api."`
	Errors    []error                    `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsUptimeListResponse struct {
	Body *AwsUptimeListBody
}

const AwsUptimeListOperationID string = "get-aws-uptime-list"
const AwsUptimeListDescription string = `Returns all uptime data between start and end dates.`
const awsUptimeListSQL string = `
SELECT
	aws_uptime.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit,
	 json_object(
		'id', aws_accounts.id,
		'number', aws_accounts.number,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account
FROM aws_uptime
LEFT JOIN aws_accounts on aws_accounts.id = aws_uptime.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_uptime.date >= :start_date
	AND aws_uptime.date < :end_date
	{WHERE}
GROUP BY aws_uptime.id
ORDER BY aws_uptime.date ASC;
;`

// ApiAwsUptimeListHandler accepts and processes requests to the below endpoints.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/aws/uptime/list?unit=<unit>
func ApiAwsUptimeListHandler(ctx context.Context, input *inputs.DateRangeUnitInput) (response *AwsUptimeListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsUptime = []*models.AwsUptime{}
		dbPath  string              = ctx.Value(dbPathKey).(string)
		sqlStmt string              = awsUptimeListSQL
		where   string              = ""
		replace string              = "{WHERE}"
		param   statements.Named    = input
		body    *AwsUptimeListBody  = &AwsUptimeListBody{
			Request:   input,
			Operation: AwsUptimeListOperationID,
		}
	)
	// setup response
	response = &AwsUptimeListResponse{}
	// check for unit
	if input.Unit != "" {
		where = "WHERE units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws uptime list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsUptime](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws uptime list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws uptime list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

// AwsUptimeAveragesBody contains the resposne body to send back
// for a request to the /list endpoint
type AwsUptimeAveragesBody struct {
	Operation    string                                    `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeUnitInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsUptime                       `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                                  `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                                  `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}                  `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                                   `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsUptimeAveragesResponse struct {
	Body *AwsUptimeAveragesBody
}

const AwsUptimeAveragesOperationID string = "get-aws-uptime-averages"
const AwsUptimeAveragesDescription string = `Returns average uptime data group by time period.`
const awsUptimeAveragesSQL string = `
SELECT
	'Average' as unit_name,
	count(DISTINCT aws_uptime.id) as count,
    (coalesce(SUM(average), 0) / count(DISTINCT aws_uptime.id) ) as average,
    strftime(:date_format, aws_uptime.date) as date
FROM aws_uptime
LEFT JOIN aws_accounts on aws_accounts.id = aws_uptime.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_uptime.date >= :start_date
	AND aws_uptime.date < :end_date
	{WHERE}
GROUP BY strftime(:date_format, aws_uptime.date)
ORDER BY aws_uptime.date ASC;
;`

// ApiAwsUptimeAveragesHandler
// Endpoints:
//
//	/version/aws/uptime/averages/{interval}/{start_date}/{end_date}?unit=<unit>
func ApiAwsUptimeAveragesHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeUnitInput) (response *AwsUptimeAveragesResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsUptime    = []*models.AwsUptime{}
		dbPath  string                 = ctx.Value(dbPathKey).(string)
		sqlStmt string                 = awsUptimeAveragesSQL
		where   string                 = ""
		replace string                 = "{WHERE}"
		param   statements.Named       = input
		body    *AwsUptimeAveragesBody = &AwsUptimeAveragesBody{
			Request:     input,
			Operation:   AwsUptimeAveragesOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
		}
	)
	// setup response
	response = &AwsUptimeAveragesResponse{}
	// check for unit
	if input.Unit != "" {
		where = "AND units.Name = :unit "
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	} else {
		sqlStmt = strings.ReplaceAll(sqlStmt, replace, where)
	}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws uptime averages adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsUptime](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws uptime averages select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws uptime averages selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

// AwsUptimeAveragesPerUnitBody contains the resposne body to send back
// for a request to the /list endpoint
type AwsUptimeAveragesPerUnitBody struct {
	Operation    string                                `json:"operation,omitempty" doc:"contains the operation id"`
	Request      *inputs.RequiredGroupedDateRangeInput `json:"request,omitempty" doc:"the original request"`
	Result       []*models.AwsUptime                   `json:"result,omitempty" doc:"list of all units returned by the api."`
	DateRange    []string                              `json:"date_range,omitempty" db:"-" doc:"all dates within the range requested"`
	ColumnOrder  []string                              `json:"column_order" db:"-" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}              `json:"column_values" db:"-" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	Errors       []error                               `json:"errors,omitempty" doc:"list of any errors that occured in the request"`
}

// the main response struct
type AwsUptimeAveragesPerUnitResponse struct {
	Body *AwsUptimeAveragesPerUnitBody
}

const AwsUptimeAveragesPerUnitOperationID string = "get-aws-uptime-averages-per-unit"
const AwsUptimeAveragesPerUnitDescription string = `Returns average uptime data grouped by time period and unit.`
const AwsUptimeAveragesPerUnitSQL string = `
SELECT
	units.name as unit_name,
	count(DISTINCT aws_uptime.id) as count,
    (coalesce(SUM(average), 0) / count(DISTINCT aws_uptime.id) ) as average,
    strftime(:date_format, aws_uptime.date) as date
FROM aws_uptime
LEFT JOIN aws_accounts on aws_accounts.id = aws_uptime.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
WHERE
	aws_uptime.date >= :start_date
	AND aws_uptime.date < :end_date
GROUP BY units.id, strftime(:date_format, aws_uptime.date)
ORDER BY aws_uptime.date ASC;
;`

// ApiAwsUptimeAveragesPerUnitHandler
// Endpoints:
//
//	/version/aws/uptime/averages-per-unit/{interval}/{start_date}/{end_date}
func ApiAwsUptimeAveragesPerUnitHandler(ctx context.Context, input *inputs.RequiredGroupedDateRangeInput) (response *AwsUptimeAveragesPerUnitResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsUptime           = []*models.AwsUptime{}
		dbPath  string                        = ctx.Value(dbPathKey).(string)
		sqlStmt string                        = AwsUptimeAveragesPerUnitSQL
		param   statements.Named              = input
		body    *AwsUptimeAveragesPerUnitBody = &AwsUptimeAveragesPerUnitBody{
			Request:     input,
			Operation:   AwsUptimeAveragesPerUnitOperationID,
			DateRange:   dateutils.Dates(input.Start(), input.End(), input.GetInterval()),
			ColumnOrder: []string{"unit_name"},
		}
	)
	// setup response
	response = &AwsUptimeAveragesPerUnitResponse{}

	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] aws uptime averages per unit adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsUptime](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws uptime averages per unit select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws uptime averages per unit selection failed."))
	} else {
		body.Result = results
	}
	body.Request.DateFormat = ""
	body.ColumnValues = cols.Values(body.Result, body.ColumnOrder)
	response.Body = body
	return
}

func RegisterAwsUptime(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + AwsUptimeSegment + "/list"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsUptimeListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List AWS uptime",
		Description:   AwsUptimeListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsUptimeTags,
	}, ApiAwsUptimeListHandler)

	uri = "/{version}/" + AwsUptimeSegment + "/averages/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsUptimeAveragesOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Average AWS uptime",
		Description:   AwsUptimeAveragesDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsUptimeTags,
	}, ApiAwsUptimeAveragesHandler)

	uri = "/{version}/" + AwsUptimeSegment + "/averages-per-unit/{interval}/{start_date}/{end_date}"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsUptimeAveragesPerUnitOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "Average AWS uptime per unit",
		Description:   AwsUptimeAveragesPerUnitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsUptimeTags,
	}, ApiAwsUptimeAveragesPerUnitHandler)

}
