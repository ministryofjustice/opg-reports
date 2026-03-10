package importer

import (
	"context"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/reset"
	"opg-reports/report/packages/times"
	"time"

	ct "opg-reports/report/internal/domains/cost/types"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

const InsertStatement string = `
INSERT INTO costs (
	region,
	service,
	month,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:month,
	:cost,
	:account_id
) ON CONFLICT (account_id, month, region, service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id
;
`

// Get returns the raw costs stats from the aws api
func Get(ctx context.Context, client *costexplorer.Client, opts *args.Import, previous ...types.ResultByTime) (found []types.ResultByTime, err error) {
	var (
		log     *slog.Logger
		options *costexplorer.GetCostAndUsageInput
		result  *costexplorer.GetCostAndUsageOutput
		start   time.Time = reset.Month(&opts.Filters.Dates.StartCosts)
		end     time.Time = reset.Day(&opts.Filters.Dates.End)
	)
	ctx, log = logger.Get(ctx)

	options = getCostAndUsageInput(start, end)

	log.Info("getting costs from aws for time period...", "start", start, "end", end)
	// make the api call
	result, err = client.GetCostAndUsage(ctx, options)
	if err != nil {
		log.Error("error getting cost and usage", "err", err.Error())
		return
	}
	found = result.ResultsByTime
	log.Info("cost api call completed.", "count", len(found))

	return
}

// Filter - not filtering on costs
func Filter(ctx context.Context, items []types.ResultByTime, filters *args.Filters) (included []types.ResultByTime) {
	included = items
	return
}

// Transform converts the original data into record for local database insertion
func Transform(ctx context.Context, data []types.ResultByTime, opts *args.Import) (results []*ct.ImportCost, err error) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)
	results = []*ct.ImportCost{}

	log.Info("transforming costs to local ct.Costs ...", "count", len(data))
	for _, result := range data {
		var day string = *result.TimePeriod.Start
		for _, group := range result.Groups {
			var service string = group.Keys[0]
			var region string = group.Keys[1]

			for _, cost := range group.Metrics {
				var item = &ct.ImportCost{
					AccountID: opts.Aws.AccountID,
					Month:     day[0:7], // this should be YYYY-MM
					Service:   service,
					Region:    region,
					Cost:      *cost.Amount,
				}
				results = append(results, item)
			}
		}
	}

	log.Info("cost transformation completed.", "count", len(results))
	return
}

// getCostAndUsageInput returns a struct formatted with expected values
// for cost data using the start and end date
//
// `start` is reset the begining of the month to make sure a full month costs are set
func getCostAndUsageInput(start time.Time, end time.Time) *costexplorer.GetCostAndUsageInput {
	var (
		service string   = "SERVICE"
		region  string   = "REGION"
		metrics []string = []string{"UnblendedCost"}
		s       string   = start.Format(times.YMD)
		e       string   = end.Format(times.YMD)
	)
	return &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		TimePeriod: &types.DateInterval{
			Start: &s,
			End:   &e,
		},
		Metrics: metrics,
		GroupBy: []types.GroupDefinition{
			{Type: types.GroupDefinitionTypeDimension, Key: &service},
			{Type: types.GroupDefinitionTypeDimension, Key: &region},
		},
	}

}
