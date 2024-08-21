package aws_costs

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/apidb"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const (
	ytdUrl      string = "/{version}/aws-costs/ytd/{$}"
	taxSplitUrl string = "/{version}/aws-costs/monthly-tax/{$}"
)

// metaExtras adds standard extra db calls to the metadata values
func metaExtras(ctx context.Context, queries *awsc.Queries, response *resp.Response, filters map[string]string) {
	// -- get overall counters
	all, _ := queries.Count(ctx)
	response.Metadata["counters"] = map[string]map[string]int{
		"totals": {
			"count": int(all),
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

// Handlers are all the handler functions, returns a map of funcs
//   - do it this way so multiple routes can share functions and the handler can then
//     get some of the vars set in the scope of this function (like query params etc)
func Handlers(ctx context.Context, mux *http.ServeMux, dbPath string) map[string]func(w http.ResponseWriter, r *http.Request) {
	// var start = query.Path("start")
	// var end = query.Path("end")

	// -- handler functions
	// year to date total
	var ytd = func(w http.ResponseWriter, r *http.Request) {
		logger.LogSetup()
		var err error
		var response = resp.New()
		var filters = map[string]string{}
		response.Start(w, r)

		// -- setup db connection
		db, err := apidb.SqlDB(dbPath)
		defer db.Close()

		if err != nil {
			slog.Error("api db error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}
		// -- setup the sqlc generated queries
		queries := awsc.New(db)
		defer queries.Close()
		start, end := dates.YearToBillingDate(time.Now(), consts.BILLING_DATE)
		total, err := queries.Total(ctx, awsc.TotalParams{
			Start: start.Format(dates.FormatYMD),
			End:   end.Format(dates.FormatYMD),
		})
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
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
		var response = resp.New()
		var filters = map[string]string{}
		response.Start(w, r)

		// -- setup db connection
		db, err := apidb.SqlDB(dbPath)
		defer db.Close()

		if err != nil {
			slog.Error("api db error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}
		// -- setup the sqlc generated queries
		queries := awsc.New(db)
		defer queries.Close()

		// -- fetch the raw results
		startDate, endDate := dates.BillingDates(time.Now().UTC(), consts.BILLING_DATE, 11)
		// add the months to the metadata
		response.Metadata["months"] = dates.Strings(dates.Range(startDate, endDate, dates.MONTH), dates.FormatYM)
		slog.Info("about to get results, limiting to date range",
			slog.String("end", endDate.Format(dates.FormatYMD)),
			slog.String("start", startDate.Format(dates.FormatYMD)))

		results, err := queries.MonthlyTotalsTaxSplit(ctx, awsc.MonthlyTotalsTaxSplitParams{
			Start: startDate.Format(dates.FormatYMD),
			End:   endDate.Format(dates.FormatYMD),
		})
		slog.Info("got results")
		if err != nil {
			slog.Error("api error", slog.String("err", err.Error()))
			response.Errors = append(response.Errors, err)
			response.End(w, r)
			return
		}

		// -- convert results over to output format
		response.Result, _ = convert.Maps(results)
		metaExtras(ctx, queries, response, filters)
		response.End(w, r)
		return
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		"taxes": taxes,
		"ytd":   ytd,
	}
}

// Register attached the route to the list view
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	funcs := Handlers(ctx, mux, dbPath)
	taxes := funcs["taxes"]
	ytd := funcs["ytd"]
	// -- actually register the handler
	mux.HandleFunc(taxSplitUrl, mw.Middleware(taxes, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(ytdUrl, mw.Middleware(ytd, mw.Logging, mw.SecurityHeaders))
	return nil
}
