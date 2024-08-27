package aws_costs

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/apiresponse"
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
//   - Set the column and column ordering data in response to fixed values
//
// Sample urls:
//   - /v1/aws_costs/monthly-tax/
func MonthlyTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
		response *ApiResponse    = NewResponse()
	)
	apiresponse.Start(response, w, r)
	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		apiresponse.ErrorAndEnd(response, err, w, r)
		return
	}
	defer db.Close()
	// -- setup the sqlc generated queries
	queries := awsc.New(db)
	defer queries.Close()

	// get date range
	startDate, endDate := dates.BillingDates(time.Now().UTC(), consts.BILLING_DATE, 12)
	// -- fetch the raw results
	slog.Info("[MonthlyTaxHandler] about to get results, limiting to date range???",
		slog.String("end", endDate.Format(dates.FormatYMD)),
		slog.String("start", startDate.Format(dates.FormatYMD)))

	results, err := queries.MonthlyTotalsTaxSplit(ctx, awsc.MonthlyTotalsTaxSplitParams{
		Start: startDate.Format(dates.FormatYMD),
		End:   endDate.Format(dates.FormatYMD),
	})
	if err != nil {
		apiresponse.ErrorAndEnd(response, err, w, r)
		return
	}
	slog.Info("got results")
	// -- add columns
	response.Columns = map[string][]string{
		"service": {"Including Tax", "Excluding Tax"},
	}
	response.ColumnOrdering = []string{"service"}
	response.Result = Common(results)
	StandardCounters(ctx, queries, response)
	StandardDates(response, startDate, endDate, endDate.AddDate(0, -1, 0), dates.MONTH)
	// --
	apiresponse.End(response, w, r)
	return
}
