package aws_costs

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

// YtdHandler is configured to handle the `ytdUrl` queries and will return
// a ApiResponse. Returns a single cost value for the entire billing year so far.
// No get parameters are used
//
//   - Connects to sqlite db via `apiDbPath`
//   - Works out the start and end dates (based on billingDate and first of the year)
//   - Gets the single total value for the year to date
//   - Formats apiResponseto have one result with the value
//
// Sample urls
//   - /v1/aws_costs/ytd/
func YtdHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err         error
		db          *sql.DB
		dbPath      string          = apiDbPath
		ctx         context.Context = apiCtx
		apiResponse *ApiResponse    = NewResponse()
	)
	response.Start(apiResponse, w, r)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}
	defer db.Close()
	// -- setup the sqlc generated queries
	queries := awsc.New(db)
	defer queries.Close()
	// get dates
	start, end := dates.YearToBillingDate(time.Now(), consts.BILLING_DATE)

	total, err := queries.Total(ctx, awsc.TotalParams{
		Start: start.Format(dates.FormatYMD),
		End:   end.Format(dates.FormatYMD),
	})
	if err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}

	apiResponse.StartDate = start.Format(dates.FormatYMD)
	apiResponse.EndDate = end.Format(dates.FormatYMD)
	apiResponse.Result = []*CommonResult{{Total: total.(float64)}}
	StandardCounters(ctx, queries, apiResponse)
	StandardDates(apiResponse, start, end, end, dates.MONTH)
	// end
	response.End(apiResponse, w, r)
	return
}
