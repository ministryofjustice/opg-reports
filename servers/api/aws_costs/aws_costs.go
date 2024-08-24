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
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/query"
	"github.com/ministryofjustice/opg-reports/servers/shared/rbase"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const (
	gByUnit     string = "unit"
	gByUnitEnv  string = "unit-env"
	gByDetailed string = "detailed"
)

const (
	ytdUrl      string = "/{version}/aws-costs/ytd/{$}"
	taxSplitUrl string = "/{version}/aws-costs/monthly-tax/{$}"
	standardUrl string = "/{version}/aws-costs/{$}"
)

var ordering = map[string][]string{
	gByUnit:     {"unit"},
	gByUnitEnv:  {"unit", "environment"},
	gByDetailed: {"account_id", "unit", "environment", "service"},
}

var (
	apiCtx    context.Context
	apiDbPath string
)

func StandardHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		now      time.Time       = time.Now().UTC()
		response *CostResponse   = NewResponse()
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
	)
	// -- allowed options for get params
	var (
		allowedIntervals []string = []string{string(dates.DAY), string(dates.MONTH)}
		allowedGroups    []string = []string{gByUnit, gByUnitEnv, gByDetailed}
	)
	// -- get params
	var (
		startQ    *query.Query = query.Get("start")
		endQ      *query.Query = query.Get("end")
		intervalQ *query.Query = query.Get("interval")
		groupbyQ  *query.Query = query.Get("group")
	)
	// -- set default values for the get params
	var (
		s, e                            = dates.BillingDates(now, consts.BILLING_DATE, 12)
		format   string                 = dates.FormatYM
		start    string                 = query.FirstD(startQ.Values(r), s.Format(dates.FormatYM))
		end      string                 = query.FirstD(endQ.Values(r), e.Format(dates.FormatYM))
		interval string                 = query.FirstD(intervalQ.Values(r), string(dates.MONTH))
		inter    dates.Interval         = dates.Interval(interval)
		groupby  string                 = query.FirstD(groupbyQ.Values(r), allowedGroups[0])
		filters  map[string]interface{} = map[string]interface{}{
			"interval": interval,
			"group":    groupby,
		}
	)
	rbase.Start(response, w, r)
	// -- validate incoming params
	if !slices.Contains(allowedIntervals, interval) {
		iErr := fmt.Errorf("invalid interval passed [%s]", interval)
		rbase.ErrorAndEnd(response, iErr, w, r)
		return
	}
	if !slices.Contains(allowedGroups, groupby) {
		gErr := fmt.Errorf("invalid groupby passed [%s]", groupby)
		rbase.ErrorAndEnd(response, gErr, w, r)
		return
	}
	startDate := dates.Time(start)
	// enddate is the first of the month, so reduce month by one for this
	// -- test for day interval
	endDate := dates.Time(end)
	rangeEnd := endDate.AddDate(0, -1, 0)
	// if its day, map the format
	if inter == dates.DAY {
		format = dates.FormatYMD
		rangeEnd = endDate.AddDate(0, 0, -1)
	}
	// setup response data
	response.QueryFilters = filters
	response.ColumnOrdering = ordering[groupby]
	response.StartDate = startDate.Format(format)
	response.EndDate = endDate.Format(format)
	response.DateRange = dates.Strings(dates.Range(startDate, rangeEnd, inter), format)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		rbase.ErrorAndEnd(response, err, w, r)
		return
	}
	defer db.Close()
	queries := awsc.New(db)
	defer queries.Close()

	interval = fmt.Sprintf("'%s'", interval)
	slog.Info("running query",
		slog.String("interval", interval),
		slog.String("groupby", groupby),
		slog.String("format", format),
		slog.String("end", endDate.Format(format)),
		slog.String("start", startDate.Format(format)))

	// -- run the query
	runQueries(ctx, queries, response, startDate.Format(format), endDate.Format(format), groupby, inter)

	extras(ctx, queries, response, startDate, endDate, format, inter)
	response.Counters.This.Count = len(response.Result)
	// end
	rbase.End(response, w, r)
	return
}

func YtdHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
		response *CostResponse   = NewResponse()
	)
	rbase.Start(response, w, r)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		rbase.ErrorAndEnd(response, err, w, r)
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
		rbase.ErrorAndEnd(response, err, w, r)
		return
	}

	response.Result = []*CommonResult{{Total: total.(float64)}}
	// meta data
	extras(ctx, queries, response, start, end, dates.FormatYM, dates.MONTH)
	response.Counters.This.Count = len(response.Result)
	// end
	rbase.End(response, w, r)
	return
}

func MonthlyTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
		response *CostResponse   = NewResponse()
	)
	rbase.Start(response, w, r)
	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		rbase.ErrorAndEnd(response, err, w, r)
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
		rbase.ErrorAndEnd(response, err, w, r)
		return
	}
	slog.Info("got results")
	// -- add columns
	response.Columns = map[string][]string{
		"service": {"Including Tax", "Excluding Tax"},
	}
	response.ColumnOrdering = []string{"service"}
	// add result
	response.Result = Common(results)
	// -- extras
	extras(ctx, queries, response, startDate, endDate, dates.FormatYM, dates.MONTH)
	response.Counters.This.Count = len(response.Result)
	// --
	rbase.End(response, w, r)
	return
}

// Register attached the route to the list view
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
