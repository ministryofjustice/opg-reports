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
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
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
		response *resp.Response  = resp.New()
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
		filters  map[string]interface{} = map[string]interface{}{}
	)
	response.Start(w, r)
	// -- validate incoming params
	if !slices.Contains(allowedIntervals, interval) {
		iErr := fmt.Errorf("invalid interval passed [%s]", interval)
		response.ErrorAndEnd(iErr, w, r)
		return
	}
	if !slices.Contains(allowedGroups, groupby) {
		gErr := fmt.Errorf("invalid groupby passed [%s]", groupby)
		response.ErrorAndEnd(gErr, w, r)
		return
	}
	startDate := dates.Time(start)
	endDate := dates.Time(end)
	rangeEnd := endDate.AddDate(0, -1, 0)
	// if its day, map the format
	if inter == dates.DAY {
		format = dates.FormatYMD
		rangeEnd = endDate.AddDate(0, 0, -1)
	}

	filters["group"] = groupby
	filters["interval"] = interval
	response.Metadata["column_ordering"] = ordering[groupby]
	response.Metadata["start_date"] = startDate.Format(format)
	response.Metadata["end_date"] = endDate.Format(format)
	response.Metadata["date_range"] = dates.Strings(dates.Range(startDate, rangeEnd, inter), format)

	// -- setup db connection
	if db, err = apidb.SqlDB(dbPath); err != nil {
		response.ErrorAndEnd(err, w, r)
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
	metaExtras(ctx, queries, response, filters)
	response.End(w, r)
	return
}

func YtdHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string          = apiDbPath
		ctx      context.Context = apiCtx
		response *YtdResponse    = NewYTDResponse()
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
	// meta data
	response.StartDate = start.Format(dates.FormatYM)
	response.EndDate = end.Format(dates.FormatYM)
	response.DateRange = dates.Strings(dates.Range(start, end, dates.MONTH), dates.FormatYM)
	response.Result.Total = total.(float64)
	// -- extras
	all, _ := queries.Count(ctx)
	response.Counters.Totals.Count = int(all)
	response.Counters.This.Count = 1
	if min, err := queries.Oldest(ctx); err == nil {
		response.DataAge.Min = min
	}
	if max, err := queries.Youngest(ctx); err == nil {
		response.DataAge.Max = max
	}
	rbase.End(response, w, r)
	return
}

func MonthlyTaxHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogSetup()
	var (
		err      error
		db       *sql.DB
		dbPath   string              = apiDbPath
		ctx      context.Context     = apiCtx
		response *MonthlyTaxResponse = NewMonthlyTaxResponse()
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
	startDate, endDate := dates.BillingDates(time.Now().UTC(), consts.BILLING_DATE, 9)
	// add the months to the metadata
	response.StartDate = startDate.Format(dates.Format)
	response.EndDate = endDate.Format(dates.Format)
	response.DateRange = dates.Strings(dates.Range(startDate, endDate.AddDate(0, -1, 0), dates.MONTH), dates.FormatYM)

	// -- fetch the raw results
	slog.Info("about to get results, limiting to date range",
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

	// add result
	response.Result = results
	// -- extras
	all, _ := queries.Count(ctx)
	response.Counters.Totals.Count = int(all)
	response.Counters.This.Count = len(results)
	if min, err := queries.Oldest(ctx); err == nil {
		response.DataAge.Min = min
	}
	if max, err := queries.Youngest(ctx); err == nil {
		response.DataAge.Max = max
	}
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
