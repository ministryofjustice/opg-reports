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
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

var (
	AwsAccountsSegment string   = "aws/accounts"
	AwsAccountsTags    []string = []string{"AWS accounts"}
)

const AwsAccountsListOperationID string = "get-aws-accounts-list"
const AwsAccountsListDescription string = `Returns all aws accounts`
const awsAccountsListSQL string = `
SELECT
	aws_accounts.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit
FROM aws_accounts
LEFT JOIN units on units.id = aws_accounts.unit_id
{WHERE}
GROUP BY aws_accounts.id
ORDER BY aws_accounts.name ASC;
;`

// ApiAwsAccountsListHandler accepts and processes requests to the below endpointutils.
// It will create a new adpator using context details and run sql query using
// crud.Select with the input params being used as named parameters on the query
//
// Endpoints:
//
//	/version/aws/accounts/list
func ApiAwsAccountsListHandler(ctx context.Context, input *inout.VersionUnitInput) (response *inout.AwsAccountsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.AwsAccount       = []*models.AwsAccount{}
		dbPath  string                     = ctx.Value(dbPathKey).(string)
		sqlStmt string                     = awsAccountsListSQL
		where   string                     = ""
		replace string                     = "{WHERE}"
		param   statements.Named           = input
		body    *inout.AwsAccountsListBody = &inout.AwsAccountsListBody{
			Request:   input,
			Operation: AwsAccountsListOperationID,
		}
	)
	// setup response
	response = &inout.AwsAccountsListResponse{}
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
		slog.Error("[api] aws accounts list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.AwsAccount](ctx, adaptor, sqlStmt, param)
	if err != nil {
		slog.Error("[api] aws accounts list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("aws accounts list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

func RegisterAwsAccounts(api huma.API) {
	var uri string = ""

	uri = "/{version}/" + AwsAccountsSegment + "/list"
	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   AwsAccountsListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List AWS accounts",
		Description:   AwsAccountsListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          AwsAccountsTags,
	}, ApiAwsAccountsListHandler)

}
