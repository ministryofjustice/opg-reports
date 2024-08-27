package aws_costs

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/apiresponse"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

// GroupBy is used by api / front tohandle group by allowed values
type GroupBy string

const (
	GroupByUnit            GroupBy = "unit"     // Group by the unit only
	GroupByUnitEnvironment GroupBy = "unit-env" // Group by unit and the environment
	GroupByDetailed        GroupBy = "detailed" // Group by more detailed columns (acocunt id, unit, environment, service)
)

// currently supported urls
const (
	ytdUrl      string = "/{version}/aws-costs/ytd/{$}"
	taxSplitUrl string = "/{version}/aws-costs/monthly-tax/{$}"
	standardUrl string = "/{version}/aws-costs/{$}"
)

// column ordering for each group by
var ordering = map[GroupBy][]string{
	GroupByUnit:            {"unit"},
	GroupByUnitEnvironment: {"unit", "environment"},
	GroupByDetailed:        {"account_id", "unit", "environment", "service"},
}

// db and context
var (
	apiCtx    context.Context
	apiDbPath string
)

// StandardHandler is configured to deal with `standardUrl` queries and will
// return a CostResponse. Used by the majority of costs data calls
//
//   - Connects to sqlite db via `apiDbPath`
//   - Uses the group and interval get parameters to determine which db query to run
//   - Adds columns to the response (driven from group data)
//   - Adds the unique values of each columns to the response
//   - Adds date range info to the response
//
// Allows following get parameters:
//   - start: change the start date of the data (default to billingDate - 12)
//   - interval: group the data by DAY or MONTH (default MONTH)
//   - group: how to group the data by other fields (allowed: `unit`, `unit-env`, `detailed` default: unit)
//   - end: change the max date of the data (default to billingDate)
//
// NOTE: the end parameter should be the day after the max you want to capture as a less than (`<`) is used
//
// Sample urls
//   - /v1/aws_costs/?group=unit
//   - /v1/aws_costs/?start=2024-01&end=2024-06
//   - /v1/aws_costs/?start=2024-01-01&end=2024-02-01&interval=DAY
//   - /v1/aws_costs/?start=2024-01-01&end=2024-02-01&interval=DAY&group=detailed
func StandardHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		// -- main
		err     error
		db      *sql.DB
		dbPath  string          = apiDbPath
		ctx     context.Context = apiCtx
		filters map[string]interface{}
		// -- dates
		now        time.Time = time.Now().UTC()
		start, end time.Time = dates.BillingDates(now, consts.BILLING_DATE, 12)
		// -- request & response
		response *CostResponse = NewResponse()
		req      *ApiRequest   = NewRequest(start, end, dates.MONTH, GroupByUnit)
		// -- validation
		allowedIntervals []string = []string{string(dates.DAY), string(dates.MONTH)}
		allowedGroups    []string = []string{string(GroupByUnit), string(GroupByUnitEnvironment), string(GroupByDetailed)}
	)
	// -- process request
	apiresponse.Start(response, w, r)
	req.Update(r)
	// -- the filters being used
	filters = map[string]interface{}{
		"interval": req.Interval,
		"group":    req.GroupBy,
		"unit":     req.Unit,
	}

	// -- validate incoming params
	if !slices.Contains(allowedIntervals, req.Interval) {
		iErr := fmt.Errorf("invalid interval passed [%s]", req.Interval)
		apiresponse.ErrorAndEnd(response, iErr, w, r)
		return
	}
	if !slices.Contains(allowedGroups, req.GroupBy) {
		gErr := fmt.Errorf("invalid groupby passed [%s]", req.GroupBy)
		apiresponse.ErrorAndEnd(response, gErr, w, r)
		return
	}
	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		apiresponse.ErrorAndEnd(response, err, w, r)
		return
	}
	defer db.Close()
	queries := awsc.New(db)
	defer queries.Close()

	slog.Info("running query",
		slog.String("interval", req.Interval),
		slog.String("groupby", req.GroupBy),
		slog.String("format", req.IntervalFormat),
		slog.String("end", req.End),
		slog.String("start", req.Start))

	// setup response data
	response.QueryFilters = filters
	response.ColumnOrdering = ordering[req.GroupByT]
	// -- run the query
	response.Result, _ = StandardQueryResults(ctx, queries, req)
	//
	response.Columns = ColumnPermutations(response.Result)
	StandardCounters(ctx, queries, response)
	StandardDates(response, req.StartT, req.EndT, req.RangeEnd, req.IntervalT)
	apiresponse.End(response, w, r)
	return
}

// YtdHandler is configured to handle the `ytdUrl` queries and will return
// a CostResponse. Returns a single cost value for the entire billing year so far.
// No get parameters are used
//
//   - Connects to sqlite db via `apiDbPath`
//   - Works out the start and end dates (based on billingDate and first of the year)
//   - Gets the single total value for the year to date
//   - Formats responseto have one result with the value
//
// Sample urls
//   - /v1/aws_costs/ytd/
func YtdHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
		response *CostResponse   = NewResponse()
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
	// get dates
	start, end := dates.YearToBillingDate(time.Now(), consts.BILLING_DATE)

	total, err := queries.Total(ctx, awsc.TotalParams{
		Start: start.Format(dates.FormatYMD),
		End:   end.Format(dates.FormatYMD),
	})
	if err != nil {
		apiresponse.ErrorAndEnd(response, err, w, r)
		return
	}

	response.StartDate = start.Format(dates.FormatYMD)
	response.EndDate = end.Format(dates.FormatYMD)
	response.Result = []*CommonResult{{Total: total.(float64)}}
	StandardCounters(ctx, queries, response)
	StandardDates(response, start, end, end, dates.MONTH)
	// end
	apiresponse.End(response, w, r)
	return
}

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
		response *CostResponse   = NewResponse()
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

// Register sets the local context and database paths to the values passed and then
// attaches the local handles to the url patterns supported by aws_costs api
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	SetCtx(ctx)
	SetDBPath(dbPath)
	// -- registers
	mux.HandleFunc(taxSplitUrl, mw.Middleware(MonthlyTaxHandler, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(ytdUrl, mw.Middleware(YtdHandler, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(standardUrl, mw.Middleware(StandardHandler, mw.Logging, mw.SecurityHeaders))
	return nil
}

func SetDBPath(path string) {
	apiDbPath = path
}
func SetCtx(ctx context.Context) {
	apiCtx = ctx
}
