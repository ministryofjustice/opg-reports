package awsr

import (
	"context"
	"opg-reports/report/config"
	"opg-reports/report/internal/utils"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	uptimeNamespace string             = "AWS/Route53"
	uptimeMetric    string             = "HealthCheckPercentageHealthy"
	uptimeUnit      types.StandardUnit = types.StandardUnitPercent
	uptimeStat      types.Statistic    = types.StatisticAverage
)

type mockClientCloudwatchMetrics struct{}

func (self *mockClientCloudwatchMetrics) ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.ListMetricsOutput, err error) {
	out = &cloudwatch.ListMetricsOutput{
		NextToken:      nil,
		OwningAccounts: nil,
		Metrics: []types.Metric{
			{
				MetricName: utils.Ptr(uptimeMetric),
				Namespace:  utils.Ptr(uptimeNamespace),
				Dimensions: []types.Dimension{
					{
						Name:  utils.Ptr("HealthCheckId"),
						Value: utils.Ptr("453ca10c-547f-4ee1-a9d6-52081a5bcf33"),
					},
				},
			},
			{
				MetricName: utils.Ptr(uptimeMetric),
				Namespace:  utils.Ptr(uptimeNamespace),
				Dimensions: []types.Dimension{
					{
						Name:  utils.Ptr("HealthCheckId"),
						Value: utils.Ptr("451cb34c-547f-41e1-a9d6-52081a5bcfee"),
					},
				},
			},
		},
	}
	return
}

func (self *mockClientCloudwatchMetrics) GetMetricStatistics(ctx context.Context, params *cloudwatch.GetMetricStatisticsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.GetMetricStatisticsOutput, err error) {
	var now = time.Now().UTC()
	out = &cloudwatch.GetMetricStatisticsOutput{
		Label: utils.Ptr(uptimeMetric),
		Datapoints: []types.Datapoint{
			{Average: utils.Ptr(99.978), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
			{Average: utils.Ptr(98.999), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
			{Average: utils.Ptr(100.00), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
			{Average: utils.Ptr(100.00), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
			{Average: utils.Ptr(100.00), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
			{Average: utils.Ptr(100.00), Unit: types.StandardUnitPercent, Timestamp: utils.Ptr(now)},
		},
	}

	return

}

func TestCWGetUptimeMetrics(t *testing.T) {

	var (
		err     error
		r       *Repository
		client  ClientCloudWatchMetricsLister
		metrics []types.Metric
		ctx     = t.Context()
		conf    = config.NewConfig()
		log     = utils.Logger("ERROR", "TEXT")
	)

	r = Default(ctx, log, conf)
	client = &mockClientCloudwatchMetrics{}
	// // use a real account if token in the env
	// if os.Getenv("AWS_SESSION_TOKEN") != "" {
	// 	client = DefaultClient[*cloudwatch.Client](ctx, "us-east-1")
	// }

	opts := &GetUptimeMetricsOptions{Namespace: "AWS/Route53", MetricName: "HealthCheckPercentageHealthy"}
	metrics, err = r.GetUptimeMetrics(client, opts)
	if err != nil {
		t.Errorf("unexpected error fetching metrics list: %s", err.Error())
		t.FailNow()
	}
	if len(metrics) <= 0 {
		t.Errorf("expected some metrics to be found, nothing returned")
		t.FailNow()
	}
}

func TestCWGetUptimeMetricStats(t *testing.T) {

	var (
		err     error
		r       *Repository
		client  ClientCloudWatchUptime
		metrics []types.Metric
		ctx     = t.Context()
		conf    = config.NewConfig()
		log     = utils.Logger("DEBUG", "JSON")
		now     = time.Now().UTC()
		end     = utils.TimeReset(now, utils.TimeIntervalDay)
		start   = end.AddDate(0, 0, -1)
	)

	r = Default(ctx, log, conf)
	client = &mockClientCloudwatchMetrics{}

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client = DefaultClient[*cloudwatch.Client](ctx, "us-east-1")
	}

	opts := &GetUptimeMetricsOptions{Namespace: "AWS/Route53", MetricName: "HealthCheckPercentageHealthy"}
	metrics, err = r.GetUptimeMetrics(client, opts)
	if err != nil || len(metrics) <= 0 {
		t.Errorf("unexpected error fetching metrics: [err:%v]", err)
	}
	statopts := &GetUptimeStatsOptions{
		Namespace:  uptimeNamespace,
		MetricName: uptimeMetric,
		Unit:       uptimeUnit,
		Statistic:  uptimeStat,
		Period:     60,
		Start:      start,
		End:        end,
	}

	res, err := r.GetUptimeStats(client, metrics, statopts)
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}

	if len(res) <= 0 {
		t.Errorf("no results returned")
	}

}

type tPeriod struct {
	Start    time.Time
	Expected int32
	Err      bool
}

func TestCWGetUptimePeriod(t *testing.T) {
	now := time.Now().UTC()

	tests := []*tPeriod{
		{Start: utils.TimeAdd(now, -1, utils.TimeIntervalHour), Expected: 60},
		{Start: utils.TimeAdd(now, -14, utils.TimeIntervalDay), Expected: 60},
		{Start: utils.TimeAdd(now, -15, utils.TimeIntervalDay), Expected: 300},
		{Start: utils.TimeAdd(now, -62, utils.TimeIntervalDay), Expected: 300},
		{Start: utils.TimeAdd(now, -63, utils.TimeIntervalDay), Expected: 3600},
		{Start: utils.TimeAdd(now, -454, utils.TimeIntervalDay), Expected: 3600},
		{Start: utils.TimeAdd(now, -455, utils.TimeIntervalDay), Err: true},
		{Start: utils.TimeAdd(now, -456, utils.TimeIntervalDay), Err: true},
	}

	for i, test := range tests {

		res, err := GetSuitableUptimePeriod(test.Start)
		if !test.Err && err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if test.Err && err == nil {
			t.Errorf("expected error, but none returned")
		}

		if !test.Err && res != test.Expected {
			t.Errorf("[%d:%v] period incorrect, expected [%v] actual [%v]", i, test.Start, test.Expected, res)
		}

	}

}
