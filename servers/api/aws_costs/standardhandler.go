package aws_costs

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/db"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request/get"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
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

// column ordering for each group by
var ordering = map[string][]string{
	string(GroupByUnit):            {"unit"},
	string(GroupByUnitEnvironment): {"unit", "environment"},
	string(GroupByDetailed):        {"account_id", "unit", "environment", "service"},
}

// StandardHandler is configured to deal with `standardUrl` queries and will
// return a ApiResponse. Used by the majority of costs data calls
//
//   - Connects to sqlite db
//   - Uses the group and interval get parameters to determine which db query to run
//   - Adds columns to the apiResponse (driven from group data)
//   - Adds the unique values of each columns to the apiResponse
//   - Adds date range info to the apiResponse
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
func StandardHandler(server *api.ApiServer, w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()

	var (
		err                      error
		awsDB                    *sql.DB
		req                      *request.Request                                                    // details on the incoming request
		now                      time.Time        = time.Now().UTC()                                 // Now
		billingStart, billingEnd time.Time        = dates.BillingDates(now, consts.BILLING_DATE, 12) // Get default start and end dates from billing
		apiResponse              *ApiResponse     = NewResponse()                                    // New response object for return
		allowedIntervals         []string         = []string{string(dates.MONTH), string(dates.DAY)} // limit the values of interval allowed
		allowedGroups            []string         = []string{string(GroupByUnit), string(GroupByUnitEnvironment), string(GroupByDetailed)}
	)
	// -- process request
	response.Start(apiResponse, w, r)
	// -- create the request
	req = request.New(
		get.New("start", billingStart.Format(dates.Format)),
		get.New("end", billingEnd.Format(dates.Format)),
		get.New("unit", ""),
		get.WithChoices("group", allowedGroups),
		get.WithChoices("interval", allowedIntervals),
	)

	// -- connect to the database
	awsDB, err = db.Connect(server.DbPath)
	if err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}
	// -- get the query connection
	queries := awsc.New(awsDB)
	defer awsDB.Close()
	defer queries.Close()

	var (
		interval       string    = req.Param("interval", r)
		intervalFormat string    = dates.IntervalFormat(dates.Interval(interval))
		groupBy        string    = req.Param("group", r)
		start          string    = req.Param("start", r)
		end            string    = req.Param("end", r)
		startDate      time.Time = dates.Time(start)
		endDate        time.Time = dates.Time(end)
		rangeEnd       time.Time = dates.RangeEnd(endDate, dates.Interval(interval))
	)

	slog.Info("running query",
		slog.String("interval", interval),
		slog.String("groupby", groupBy),
		slog.String("format", intervalFormat),
		slog.String("end", end),
		slog.String("start", start))

	// setup apiResponse data
	apiResponse.ColumnOrdering = ordering[groupBy]
	// -- run the query
	apiResponse.Result, _ = StandardQueryResults(server.Ctx, queries, req, r)
	//
	apiResponse.Columns = ColumnPermutations(apiResponse.Result)
	StandardCounters(server.Ctx, queries, apiResponse)
	StandardDates(apiResponse, startDate, endDate, rangeEnd, dates.Interval(interval))

	// --
	apiResponse.QueryFilters = req.Mapped(r)
	response.End(apiResponse, w, r)
	return
}

// --- all the functions that call the correct queries

type queryWrapperF func(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error)

func monthlyPerUnit(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.MonthlyCostsPerUnitParams{Start: start, End: end}
	if res, err := q.MonthlyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyPerUnitEnv(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.MonthlyCostsPerUnitEnvironmentParams{Start: start, End: end}
	if res, err := q.MonthlyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetails(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.MonthlyCostsDetailedParams{Start: start, End: end}
	if res, err := q.MonthlyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetailsForUnit(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		unit  string = req.Param("unit", r)
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)

	params := awsc.MonthlyCostsDetailedForUnitParams{Start: start, End: end, Unit: unit}
	if res, err := q.MonthlyCostsDetailedForUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnit(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.DailyCostsPerUnitParams{Start: start, End: end}
	if res, err := q.DailyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnitEnv(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.DailyCostsPerUnitEnvironmentParams{Start: start, End: end}
	if res, err := q.DailyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetails(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.DailyCostsDetailedParams{Start: start, End: end}
	if res, err := q.DailyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetailsForUnit(ctx context.Context, req *request.Request, q *awsc.Queries, r *http.Request) (results []*CommonResult, err error) {
	var (
		unit  string = req.Param("unit", r)
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsc.DailyCostsDetailedForUnitParams{Start: start, End: end, Unit: unit}
	if res, err := q.DailyCostsDetailedForUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

// StandardQueryResults uses the group and interval values from the request to determine
// which db query to run
//
// The query to use is determined by the following:
// Interval set to `MONTH`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Interval set to `DAY`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Query results are converted to `[]*CommonResult` struct.
func StandardQueryResults(ctx context.Context, q *awsc.Queries, req *request.Request, r *http.Request) (results []*CommonResult, err error) {
	var (
		unit     string = req.Param("unit", r)
		interval string = req.Param("interval", r)
		groupBy  string = req.Param("group", r)
	)
	// define the possible functions and params
	var possibleFuncs = map[string]queryWrapperF{
		string(dates.MONTH) + string(GroupByUnit):            monthlyPerUnit,
		string(dates.MONTH) + string(GroupByUnitEnvironment): monthlyPerUnitEnv,
		string(dates.MONTH) + string(GroupByDetailed):        monthlyDetails,
		string(dates.DAY) + string(GroupByUnit):              dailyPerUnit,
		string(dates.DAY) + string(GroupByUnitEnvironment):   dailyPerUnitEnv,
		string(dates.DAY) + string(GroupByDetailed):          dailyDetails,
		//
	}

	key := interval + groupBy
	// if a filter is set, that over rides the default grouping setup
	if unit != "" && interval == string(dates.MONTH) {
		results, err = monthlyDetailsForUnit(ctx, req, q, r)
	} else if unit != "" && interval == string(dates.DAY) {
		results, err = dailyDetailsForUnit(ctx, req, q, r)
	} else if f, ok := possibleFuncs[key]; ok {
		results, err = f(ctx, req, q, r)
	} else {
		err = fmt.Errorf("error finding query function based on get paremters")
	}
	return
}
