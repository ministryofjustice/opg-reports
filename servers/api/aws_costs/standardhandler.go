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
var ordering = map[GroupBy][]string{
	GroupByUnit:            {"unit"},
	GroupByUnitEnvironment: {"unit", "environment"},
	GroupByDetailed:        {"account_id", "unit", "environment", "service"},
}

// StandardHandler is configured to deal with `standardUrl` queries and will
// return a ApiResponse. Used by the majority of costs data calls
//
//   - Connects to sqlite db via `apiDbPath`
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
		// -- request & apiResponse
		apiResponse *ApiResponse = NewResponse()
		req         *ApiRequest  = NewRequest(start, end, dates.MONTH, GroupByUnit)
		// -- validation
		allowedIntervals []string = []string{string(dates.DAY), string(dates.MONTH)}
		allowedGroups    []string = []string{string(GroupByUnit), string(GroupByUnitEnvironment), string(GroupByDetailed)}
	)
	// -- process request
	response.Start(apiResponse, w, r)
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
		response.ErrorAndEnd(apiResponse, iErr, w, r)
		return
	}
	if !slices.Contains(allowedGroups, req.GroupBy) {
		gErr := fmt.Errorf("invalid groupby passed [%s]", req.GroupBy)
		response.ErrorAndEnd(apiResponse, gErr, w, r)
		return
	}
	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
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

	// setup apiResponse data
	apiResponse.QueryFilters = filters
	apiResponse.ColumnOrdering = ordering[req.GroupByT]
	// -- run the query
	apiResponse.Result, _ = StandardQueryResults(ctx, queries, req)
	//
	apiResponse.Columns = ColumnPermutations(apiResponse.Result)
	StandardCounters(ctx, queries, apiResponse)
	StandardDates(apiResponse, req.StartT, req.EndT, req.RangeEnd, req.IntervalT)
	response.End(apiResponse, w, r)
	return
}

// --- all the functions that call the correct queries

type queryWrapperF func(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error)

func monthlyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetailsForUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsDetailedForUnitParams{Start: req.Start, End: req.End, Unit: req.Unit}
	if res, err := q.MonthlyCostsDetailedForUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetailsForUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsDetailedForUnitParams{Start: req.Start, End: req.End, Unit: req.Unit}
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
func StandardQueryResults(ctx context.Context, q *awsc.Queries, req *ApiRequest) (results []*CommonResult, err error) {
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

	key := req.Interval + req.GroupBy
	// if a filter is set, that over rides the default grouping setup
	if req.Unit != "" && req.IntervalT == dates.MONTH {
		results, err = monthlyDetailsForUnit(ctx, req, q)
	} else if req.Unit != "" && req.IntervalT == dates.DAY {
		results, err = dailyDetailsForUnit(ctx, req, q)
	} else if f, ok := possibleFuncs[key]; ok {
		results, err = f(ctx, req, q)
	} else {
		err = fmt.Errorf("error finding query function based on get paremters")
	}
	return
}
