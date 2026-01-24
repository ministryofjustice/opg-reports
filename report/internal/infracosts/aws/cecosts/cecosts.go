package cecosts

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type ApiClient interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

// GetCostData calls the cost explorer api and returns cost and usage data based on the options that are set.
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
func GetCostData[T ApiClient](ctx context.Context, log *slog.Logger, client T, options *costexplorer.GetCostAndUsageInput) (result *costexplorer.GetCostAndUsageOutput, err error) {
	log = log.With("func", "GetCostData")
	// initial call
	result, err = client.GetCostAndUsage(ctx, options)
	if err != nil {
		err = errors.Join(ErrGettingCostData, err)
		log.Error("error: failed to get cost data", "err", err.Error())
		return
	}

	return
}

// GetCostDataOptions returns a CostAndUsageInput struct formatted with expected values
// to fetch monthly cost data via `Get`
func GetCostDataOptions(start time.Time, end time.Time) *costexplorer.GetCostAndUsageInput {
	var (
		startDate string   = utils.TimeReset(start, utils.TimeIntervalMonth).Format(utils.DATE_FORMATS.YMD)
		endDate   string   = utils.TimeReset(end, utils.TimeIntervalMonth).Format(utils.DATE_FORMATS.YMD)
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
