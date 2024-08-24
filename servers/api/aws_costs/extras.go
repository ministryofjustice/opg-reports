package aws_costs

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func extras(
	ctx context.Context,
	queries *awsc.Queries,
	response *CostResponse,
	start time.Time,
	end time.Time,
	df string,
	interval dates.Interval,
) {
	rangeEnd := end.AddDate(0, -1, 0)
	response.StartDate = start.Format(df)
	response.EndDate = end.Format(df)
	response.DateRange = dates.Strings(dates.Range(start, rangeEnd, interval), df)

	// -- extras
	all, _ := queries.Count(ctx)
	response.Counters.Totals.Count = int(all)
	if min, err := queries.Oldest(ctx); err == nil {
		response.DataAge.Min = min
	}
	if max, err := queries.Youngest(ctx); err == nil {
		response.DataAge.Max = max
	}

}
