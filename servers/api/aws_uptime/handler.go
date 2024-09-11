package aws_uptime

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_uptime/awsu"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/db"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request/get"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

// Handler
func Handler(server *api.ApiServer, w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()

	var (
		err              error
		awsDB            *sql.DB
		req              *request.Request                                                    // details on the incoming request
		now              time.Time        = time.Now().UTC()                                 // Now
		startD, endD     time.Time        = dates.StartEnd(now, 6)                           // Get default start and end dates
		apiResponse      *ApiResponse     = NewResponse()                                    // New response object for return
		allowedIntervals []string         = []string{string(dates.MONTH), string(dates.DAY)} // limit the values of interval allowed
	)
	response.Start(apiResponse, w, r)
	// -- create the request
	req = request.New(
		get.New("start", startD.Format(dates.Format)),
		get.New("end", endD.Format(dates.Format)),
		get.New("unit", ""),
		get.WithChoices("interval", allowedIntervals),
	)

	// -- connect to the database
	awsDB, err = db.Connect(server.DbPath)
	if err != nil {
		response.ErrorAndEnd(apiResponse, err, w, r)
		return
	}
	// -- get the query connection
	queries := awsu.New(awsDB)
	defer awsDB.Close()
	defer queries.Close()

	var (
		interval       string    = req.Param("interval", r)
		intervalFormat string    = dates.IntervalFormat(dates.Interval(interval))
		start          string    = req.Param("start", r)
		end            string    = req.Param("end", r)
		startDate      time.Time = dates.Time(start)
		endDate        time.Time = dates.Time(end)
		rangeEnd       time.Time = dates.RangeEnd(endDate, dates.Interval(interval))
	)

	slog.Info("running query",
		slog.String("interval", interval),
		slog.String("format", intervalFormat),
		slog.String("end", end),
		slog.String("start", start))

	// setup apiResponse data
	// apiResponse.ColumnOrdering = ordering[groupBy]
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

type queryWrapperF func(ctx context.Context, req *request.Request, q *awsu.Queries, r *http.Request) (results []*Result, err error)

func monthly(ctx context.Context, req *request.Request, q *awsu.Queries, r *http.Request) (results []*Result, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsu.UptimePerMonthParams{Start: start, End: end}
	if res, err := q.UptimePerMonth(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyForUnit(ctx context.Context, req *request.Request, q *awsu.Queries, r *http.Request) (results []*Result, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
		unit  string = req.Param("unit", r)
	)
	params := awsu.UptimePerMonthFilterByUnitParams{Start: start, End: end, Unit: unit}
	if res, err := q.UptimePerMonthFilterByUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func daily(ctx context.Context, req *request.Request, q *awsu.Queries, r *http.Request) (results []*Result, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
	)
	params := awsu.UptimePerDayParams{Start: start, End: end}
	if res, err := q.UptimePerDay(ctx, params); err == nil {
		results = Common(res)
	}
	return
}
func dailyForUnit(ctx context.Context, req *request.Request, q *awsu.Queries, r *http.Request) (results []*Result, err error) {
	var (
		start string = req.Param("start", r)
		end   string = req.Param("end", r)
		unit  string = req.Param("unit", r)
	)
	params := awsu.UptimePerDayFilterByUnitParams{Start: start, End: end, Unit: unit}
	if res, err := q.UptimePerDayFilterByUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

// StandardQueryResults
func StandardQueryResults(ctx context.Context, q *awsu.Queries, req *request.Request, r *http.Request) (results []*Result, err error) {
	var (
		unit     string = req.Param("unit", r)
		interval string = req.Param("interval", r)
	)

	if unit != "" && interval == string(dates.MONTH) {
		results, err = monthlyForUnit(ctx, req, q, r)
	} else if unit != "" && interval == string(dates.DAY) {
		results, err = dailyForUnit(ctx, req, q, r)
	} else if interval == string(dates.MONTH) {
		results, err = monthly(ctx, req, q, r)
	} else if interval == string(dates.DAY) {
		results, err = daily(ctx, req, q, r)
	}

	return
}
