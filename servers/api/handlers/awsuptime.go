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

}
