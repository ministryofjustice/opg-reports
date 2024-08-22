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
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/convert"
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

// metaExtras adds standard extra db calls to the metadata values
func metaExtras(ctx context.Context, queries *awsc.Queries, response *resp.Response, filters map[string]interface{}) {
	// -- get overall counters
	all, _ := queries.Count(ctx)
	response.Metadata["counters"] = map[string]map[string]int{
		"totals": {
			"count": int(all),
		},
		"this": {
			"count": len(response.Result),
		},
	}
	response.Metadata["filters"] = filters
	// -- add the date min / max values
	min, err := queries.Oldest(ctx)
	max, err := queries.Youngest(ctx)
	if err == nil {
		response.DataAge.Min = min
		response.DataAge.Max = max
	}
}

func runQueries(ctx context.Context, queries *awsc.Queries, response *resp.Response, start string, end string, groupby string, interval dates.Interval) {

	var columns map[string]map[string]bool
	// - per unit, by month
	// - per unit, by day
	// - per unit env, by month
	// - per unit env, by day
	// - per detailed, by month
	if groupby == gByUnit && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsPerUnit(ctx, awsc.MonthlyCostsPerUnitParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByUnit && interval == dates.DAY {
		res, _ := queries.DailyCostsPerUnit(ctx, awsc.DailyCostsPerUnitParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByUnitEnv && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsPerUnitEnvironment(ctx, awsc.MonthlyCostsPerUnitEnvironmentParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment.(string)] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByUnitEnv && interval == dates.DAY {
		res, _ := queries.DailyCostsPerUnitEnvironment(ctx, awsc.DailyCostsPerUnitEnvironmentParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByDetailed && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsDetailed(ctx, awsc.MonthlyCostsDetailedParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}, "account_id": {}, "service": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment] = true
			columns["account_id"][r.AccountID] = true
			columns["service"][r.Service] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByDetailed && interval == dates.DAY {
		res, _ := queries.DailyCostsDetailed(ctx, awsc.DailyCostsDetailedParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}, "account_id": {}, "service": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment] = true
			columns["account_id"][r.AccountID] = true
			columns["service"][r.Service] = true
		}
		response.Result, _ = convert.Maps(res)
	}
	// -- map the columns
	colList := map[string][]string{}
	for col, values := range columns {
		colList[col] = []string{}
		for choice, _ := range values {
			colList[col] = append(colList[col], choice)
		}
	}
	response.Metadata["columns"] = colList
}

// Handlers are all the handler functions, returns a map of funcs
//   - do it this way so multiple routes can share functions and the handler can then
//     get some of the vars set in the scope of this function (like query params etc)
func Handlers(ctx context.Context, mux *http.ServeMux, dbPath string) map[string]func(w http.ResponseWriter, r *http.Request) {
	var (
		startQ           *query.Query = query.Get("start")
		endQ             *query.Query = query.Get("end")
		intervalQ        *query.Query = query.Get("interval")
		groupbyQ         *query.Query = query.Get("group")
		allowedIntervals []string     = []string{
			string(dates.DAY),
			string(dates.MONTH),
		}
		allowedGroups []string = []string{
			gByUnit,
			gByUnitEnv,
			gByDetailed,
		}
	)
	// -- handler functions
	// standard queries, grouping and filtering done by query params
	// 	interval = interval=%Y-%m-%d
	// 	groupby = group=unit&group=environment
	var standard = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		// -- set default values for the get params
		var (
			err      error
			db       *sql.DB
			now      time.Time              = time.Now().UTC()
			s, e                            = dates.BillingDates(now, consts.BILLING_DATE, 9)
			format   string                 = dates.FormatYM
			response *resp.Response         = resp.New()
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
			response.ErrorAndEnd(fmt.Errorf("invalid interval passed"), w, r)
			return
		}
		if !slices.Contains(allowedGroups, groupby) {
			response.ErrorAndEnd(fmt.Errorf("invalid groupby passed"), w, r)
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

	// year to date total
	var ytd = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var db *sql.DB
		var response = resp.New()
		var filters = map[string]interface{}{}
		response.Start(w, r)

		// -- setup db connection
		if db, err = apidb.SqlDB(dbPath); err != nil {
			response.ErrorAndEnd(err, w, r)
			return
		}
		defer db.Close()

		// -- setup the sqlc generated queries
		queries := awsc.New(db)
		defer queries.Close()
		start, end := dates.YearToBillingDate(time.Now(), consts.BILLING_DATE)
		response.Metadata["start_date"] = start.Format(dates.FormatYM)
		response.Metadata["end_date"] = end.Format(dates.FormatYM)

		total, err := queries.Total(ctx, awsc.TotalParams{
			Start: start.Format(dates.FormatYMD),
			End:   end.Format(dates.FormatYMD),
		})
		if err != nil {
			response.ErrorAndEnd(err, w, r)
			return
		}

		response.Result = []map[string]interface{}{
			{"total": total},
		}
		metaExtras(ctx, queries, response, filters)
		response.End(w, r)
		return

	}

	// split by tax
	var taxes = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var db *sql.DB
		var response = resp.New()
		var filters = map[string]interface{}{}
		response.Start(w, r)

		// -- setup db connection
		// -- setup db connection
		if db, err = apidb.SqlDB(dbPath); err != nil {
			response.ErrorAndEnd(err, w, r)
			return
		}
		defer db.Close()
		// -- setup the sqlc generated queries
		queries := awsc.New(db)
		defer queries.Close()

		// get date range
		startDate, endDate := dates.BillingDates(time.Now().UTC(), consts.BILLING_DATE, 12)
		// add the months to the metadata
		response.Metadata["start_date"] = startDate.Format(dates.Format)
		response.Metadata["end_date"] = endDate.Format(dates.Format)
		response.Metadata["date_range"] = dates.Strings(dates.Range(startDate, endDate.AddDate(0, -1, 0), dates.MONTH), dates.FormatYM)

		// -- fetch the raw results
		slog.Info("about to get results, limiting to date range",
			slog.String("end", endDate.Format(dates.FormatYMD)),
			slog.String("start", startDate.Format(dates.FormatYMD)))

		results, err := queries.MonthlyTotalsTaxSplit(ctx, awsc.MonthlyTotalsTaxSplitParams{
			Start: startDate.Format(dates.FormatYMD),
			End:   endDate.Format(dates.FormatYMD),
		})
		if err != nil {
			response.ErrorAndEnd(err, w, r)
			return
		}
		slog.Info("got results")
		// -- add columns
		response.Metadata["columns"] = map[string][]string{
			"service": {"Including Tax", "Excluding Tax"},
		}

		// -- convert results over to output format
		response.Result, _ = convert.Maps(results)
		metaExtras(ctx, queries, response, filters)
		response.End(w, r)
		return
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		"taxes":    taxes,
		"ytd":      ytd,
		"standard": standard,
	}
}

// Register attached the route to the list view
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	funcs := Handlers(ctx, mux, dbPath)
	taxes := funcs["taxes"]
	ytd := funcs["ytd"]
	standard := funcs["standard"]
	// -- actually register the handler
	mux.HandleFunc(taxSplitUrl, mw.Middleware(taxes, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(ytdUrl, mw.Middleware(ytd, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(standardUrl, mw.Middleware(standard, mw.Logging, mw.SecurityHeaders))
	return nil
}
