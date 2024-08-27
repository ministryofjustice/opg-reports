package aws_costs

import (
	"context"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

// runQueries uses the groupby and interval values to determines which query to run and store result.
// Using the result of the query it then generates the column data
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
// Query results are converted to `[]*CommonResult` struct, which inclues all of the possible columns from all queries.
// This conversion is done via json marshaling
func runQueries(ctx context.Context, queries *awsc.Queries, response *CostResponse, start string, end string, groupby string, interval dates.Interval) {

	// - per unit, by month
	// - per unit, by day
	// - per unit env, by month
	// - per unit env, by day
	// - per detailed, by month
	if groupby == gByUnit && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsPerUnit(ctx, awsc.MonthlyCostsPerUnitParams{Start: start, End: end})
		response.Result = Common(res)
	} else if groupby == gByUnit && interval == dates.DAY {
		res, _ := queries.DailyCostsPerUnit(ctx, awsc.DailyCostsPerUnitParams{Start: start, End: end})
		response.Result = Common(res)
	} else if groupby == gByUnitEnv && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsPerUnitEnvironment(ctx, awsc.MonthlyCostsPerUnitEnvironmentParams{Start: start, End: end})
		response.Result = Common(res)
	} else if groupby == gByUnitEnv && interval == dates.DAY {
		res, _ := queries.DailyCostsPerUnitEnvironment(ctx, awsc.DailyCostsPerUnitEnvironmentParams{Start: start, End: end})

		response.Result = Common(res)
	} else if groupby == gByDetailed && interval == dates.MONTH {
		res, _ := queries.MonthlyCostsDetailed(ctx, awsc.MonthlyCostsDetailedParams{Start: start, End: end})
		response.Result = Common(res)
	} else if groupby == gByDetailed && interval == dates.DAY {
		res, _ := queries.DailyCostsDetailed(ctx, awsc.DailyCostsDetailedParams{Start: start, End: end})
		response.Result = Common(res)
	}
	// -- generate all the unique column values for
	// put to as map of maps first so we dont get dups
	columns := map[string]map[string]bool{}
	for _, r := range response.Result {
		if r.Unit != "" {
			if _, ok := columns["unit"]; !ok {
				columns["unit"] = map[string]bool{}
			}
			columns["unit"][r.Unit] = true
		}
		if r.Environment != nil && r.Environment.(string) != "" {
			if _, ok := columns["environment"]; !ok {
				columns["environment"] = map[string]bool{}
			}
			columns["environment"][r.Environment.(string)] = true
		}
		if r.AccountID != "" {
			if _, ok := columns["account_id"]; !ok {
				columns["account_id"] = map[string]bool{}
			}
			columns["account_id"][r.AccountID] = true
		}
		if r.Service != "" {
			if _, ok := columns["service"]; !ok {
				columns["service"] = map[string]bool{}
			}
			columns["service"][r.Service] = true
		}
	}
	// now convert the mapinto format for output
	colList := map[string][]string{}
	for col, values := range columns {
		colList[col] = []string{}
		for choice, _ := range values {
			colList[col] = append(colList[col], choice)
		}
	}
	response.Columns = colList
}
