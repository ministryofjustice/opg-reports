// Package infracost implements methods to fetch data from the AWS API for costexplorer (hence ce).
package infracost

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// AwsClient is used to allow mocking and is a proxy for *costexplorer.Client
type AwsClient interface {
	// GetCostAndUsage method signature from the *costexplorer.Client
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

// GetCostDataOptions options that can change and specify for fetching cost data
type GetCostDataOptions struct {
	Start     time.Time
	End       time.Time
	AccountID string // AccountID provided by awsid.AccountID
}

// GetCostData[T] calls the cost explorer api and returns cost and usage data based on the options that are set.
//
// `T AwsClient` is a proxy for *costexplorer.Client to allow mocking
//
// Expects `options` to resemble the output of `GetCostDataOptions`.
//
// Equilivant cli call:
//
//	aws-vault exec ${profile} -- aws ce get-cost-and-usage \
//		--time-period Start=2025-03-01,End=2025-04-01 \
//		--granularity MONTHLY \
//		--metrics "UnblendedCost" \
//		--group-by Type=DIMENSION,Key=SERVICE Type=DIMENSION,Key=REGION
//
// Note: API limits grouping to 2, so we cant get linked account details at the same time.
func GetCostData[T AwsClient](ctx context.Context, log *slog.Logger, client T, options *GetCostDataOptions) (costs []*infracostmodels.AwsCost, err error) {
	var result *costexplorer.GetCostAndUsageOutput
	var apiOpts *costexplorer.GetCostAndUsageInput = getCostDataOptions(options.Start, options.End)

	log = log.With("package", "infracosts", "func", "GetCostData")
	log.Debug("starting ...")
	// initial call
	result, err = client.GetCostAndUsage(ctx, apiOpts)
	if err != nil {
		err = errors.Join(ErrGettingCostData, err)
		log.Error("error: failed to get cost data", "err", err.Error())
		return
	}

	costs, err = toModels(ctx, log, options.AccountID, result)

	log.Debug("complete.")
	return
}

// toModels converts the raw data into a list of models ready to write to the database
func toModels(ctx context.Context, log *slog.Logger, account string, result *costexplorer.GetCostAndUsageOutput) (costs []*infracostmodels.AwsCost, err error) {

	costs = []*infracostmodels.AwsCost{}
	log = log.With("package", "infracosts", "func", "toModels")
	// convert results into a map
	log.Debug("coverting cost data ... ")
	for _, result := range result.ResultsByTime {
		var day string = *result.TimePeriod.Start
		for _, group := range result.Groups {
			var service string = *&group.Keys[0]
			var region string = *&group.Keys[1]

			for _, cost := range group.Metrics {

				var item = &infracostmodels.AwsCost{
					AccountID: account,
					Date:      day,
					Service:   service,
					Region:    region,
					Cost:      *cost.Amount,
				}
				costs = append(costs, item)

			}
		}
	}
	log.Debug("complete.")

	return
}

// getCostDataOptions returns a CostAndUsageInput struct formatted with expected values
// for monhtly cost data using the start and end date
//
// `start` & `end` dates are reset to the first day of the month so 2026-01-31 => 2026-01-01
func getCostDataOptions(start time.Time, end time.Time) *costexplorer.GetCostAndUsageInput {
	var (
		startDate string   = times.AsString(times.ResetMonth(start), times.YMD)
		endDate   string   = times.AsString(times.ResetMonth(end), times.YMD)
		service   string   = "SERVICE"
		region    string   = "REGION"
		metrics   []string = []string{"UnblendedCost"}
	)
	return &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		TimePeriod: &types.DateInterval{
			Start: &startDate,
			End:   &endDate,
		},
		Metrics: metrics,
		GroupBy: []types.GroupDefinition{
			{Type: types.GroupDefinitionTypeDimension, Key: &service},
			{Type: types.GroupDefinitionTypeDimension, Key: &region},
		},
	}

}
