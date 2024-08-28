package aws_costs

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/db"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

// MonthlyTaxHandler handles the `taxSplitUrl` requests and returns a CostRepsonse.
// Returns total costs including and excluding tax for the last 12 months. Used to
// make comparing to finace data simpler as that doesnt include tax.
// No get parameters are used
//
//   - Connect to db vai `apiDbPath`
//   - Run query
//   - Set the column and column ordering data in apiResponse to fixed values
//
// Sample urls:
//   - /v1/aws_costs/monthly-tax/
func MonthlyTaxHandler(server *api.ApiServer, w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err         error
		awsDB       *sql.DB
		apiResponse *ApiResponse = NewResponse()
	)
	response.Start(apiResponse, w, r)
	// -- setup db connection
	awsDB, err = db.Connect(server.DbPath)
	if err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}
	queries := awsc.New(awsDB)
	defer awsDB.Close()
	defer queries.Close()

	// get date range
	startDate, endDate := dates.BillingDates(time.Now().UTC(), consts.BILLING_DATE, 12)
	// -- fetch the raw results
	slog.Info("[MonthlyTaxHandler] about to get results, limiting to date range???",
		slog.String("end", endDate.Format(dates.FormatYMD)),
		slog.String("start", startDate.Format(dates.FormatYMD)))

	results, err := queries.MonthlyTotalsTaxSplit(server.Ctx, awsc.MonthlyTotalsTaxSplitParams{
		Start: startDate.Format(dates.FormatYMD),
		End:   endDate.Format(dates.FormatYMD),
	})
	if err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}
	slog.Info("got results")
	// -- add columns
	apiResponse.Columns = map[string][]string{
		"service": {"Including Tax", "Excluding Tax"},
	}
	apiResponse.ColumnOrdering = []string{"service"}
	apiResponse.Result = Common(results)
	StandardCounters(server.Ctx, queries, apiResponse)
	StandardDates(apiResponse, startDate, endDate, endDate.AddDate(0, -1, 0), dates.MONTH)
	// --
	response.End(apiResponse, w, r)
	return
}
