package awsr

import (
	"testing"
	"time"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// mockClientCostExplorerGetter returns fixed cost values for ce calls
type mockClientCostExplorerGetter struct{}

// GetCostAndUsage returns mock / fake cost data so no api call is generated
func (self *mockClientCostExplorerGetter) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (out *costexplorer.GetCostAndUsageOutput, err error) {
	out = &costexplorer.GetCostAndUsageOutput{
		NextPageToken: nil,
		ResultsByTime: []types.ResultByTime{
			{
				TimePeriod: &types.DateInterval{
					Start: params.TimePeriod.Start,
					End:   params.TimePeriod.End,
				},
				Groups: []types.Group{
					{
						Keys: []string{"AWS CloudTrail", "NoRegion"},
						Metrics: map[string]types.MetricValue{
							params.Metrics[0]: {
								Amount: utils.Ptr("-3.61234665"),
								Unit:   utils.Ptr("USD"),
							},
						},
					},
					{
						Keys: []string{"AWS CloudTrail", "eu-west-1"},
						Metrics: map[string]types.MetricValue{
							params.Metrics[0]: {
								Amount: utils.Ptr("10.8865"),
								Unit:   utils.Ptr("USD"),
							},
						},
					},
					{
						Keys: []string{"AWS CloudTrail", "eu-west-2"},
						Metrics: map[string]types.MetricValue{
							params.Metrics[0]: {
								Amount: utils.Ptr("0.1065"),
								Unit:   utils.Ptr("USD"),
							},
						},
					},
					{
						Keys: []string{"Amazon DynamoDB", "eu-west-2"},
						Metrics: map[string]types.MetricValue{
							params.Metrics[0]: {
								Amount: utils.Ptr("0.0050711398"),
								Unit:   utils.Ptr("USD"),
							},
						},
					},
				},
			},
		},
	}
	return
}

func TestCEGetCosts(t *testing.T) {
	var (
		err  error
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)
	client := &mockClientCostExplorerGetter{}
	sv := Default(ctx, log, conf)
	data, err := sv.GetCostData(client, &GetCostDataOptions{
		StartDate:   utils.TimeReset(time.Now().UTC().AddDate(0, -4, 0), "month").Format(utils.DATE_FORMATS.YMD),
		EndDate:     utils.TimeReset(time.Now().UTC().AddDate(0, -3, 0), "month").Format(utils.DATE_FORMATS.YMD),
		Granularity: types.GranularityMonthly,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	if len(data) <= 0 {
		t.Errorf("should return dummy cost values")
	}

}
