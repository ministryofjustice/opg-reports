package uptime

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type mockGetter struct{}

func (self *mockGetter) ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.ListMetricsOutput, err error) {
	out = &cloudwatch.ListMetricsOutput{
		Metrics: []types.Metric{
			{
				MetricName: ptr.Ptr(metricsName),
				Namespace:  ptr.Ptr(metricsNamespace),
				Dimensions: []types.Dimension{
					{
						Name:  ptr.Ptr("metric-A"),
						Value: ptr.Ptr("value-A"),
					},
				},
			},
		},
	}
	return
}

func (self *mockGetter) GetMetricStatistics(ctx context.Context, params *cloudwatch.GetMetricStatisticsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.GetMetricStatisticsOutput, err error) {
	var ts = time.Now()
	out = &cloudwatch.GetMetricStatisticsOutput{
		Datapoints: []types.Datapoint{
			{
				Unit:      metricsUnit,
				Timestamp: ptr.Ptr(ts),
				Average:   ptr.Ptr(99.999),
			},
		},
	}
	return
}

func (self *mockGetter) Options() cloudwatch.Options {
	return cloudwatch.Options{
		Region: "us-east-1",
	}
}

func TestRedoUptimeWithMock(t *testing.T) {
	var (
		err    error
		client *mockGetter = &mockGetter{}
		r      *cloudwatch.GetMetricStatisticsOutput
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		now    time.Time       = time.Now().UTC()
		start  time.Time       = now.AddDate(0, -4, 0)
		end    time.Time       = now.AddDate(0, -3, 0)
	)

	r, err = GetUptimeData(ctx, log, client, &GetUptimeDataOptions{Start: start, End: end})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(r.Datapoints) <= 0 {
		t.Error("failed to find uptime data")
	}

}

func TestRedoUptimeWithoutMock(t *testing.T) {
	var (
		err    error
		client *cloudwatch.Client
		r      *cloudwatch.GetMetricStatisticsOutput
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		now    time.Time       = time.Now().UTC()
		start  time.Time       = now.AddDate(0, -4, 0)
		end    time.Time       = now.AddDate(0, -3, 0)
	)

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client, err = awsclients.New[*cloudwatch.Client](ctx, log, "us-east-1")
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
			t.FailNow()
		}

		r, err = GetUptimeData(ctx, log, client, &GetUptimeDataOptions{Start: start, End: end})
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(r.Datapoints) <= 0 {
			t.Error("failed to find uptime data")
		}
	} else {
		t.Skip()
	}
}
