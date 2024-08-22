package aws_costs

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

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
			columns["environment"][r.Environment.(string)] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByDetailed && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsDetailed(ctx, awsc.MonthlyCostsDetailedParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}, "account_id": {}, "service": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment.(string)] = true
			columns["account_id"][r.AccountID] = true
			columns["service"][r.Service] = true
		}
		response.Result, _ = convert.Maps(res)
	} else if groupby == gByDetailed && interval == dates.DAY {
		res, _ := queries.DailyCostsDetailed(ctx, awsc.DailyCostsDetailedParams{Start: start, End: end})
		columns = map[string]map[string]bool{"unit": {}, "environment": {}, "account_id": {}, "service": {}}
		for _, r := range res {
			columns["unit"][r.Unit] = true
			columns["environment"][r.Environment.(string)] = true
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
