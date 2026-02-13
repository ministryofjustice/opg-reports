package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"opg-reports/report/internal/utils/times"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/jmoiron/sqlx"
)

type mockUptimeClient struct{}

func (self *mockUptimeClient) ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.ListMetricsOutput, err error) {
	var (
		metricsNamespace string = "AWS/Route53"
		metricsName      string = "HealthCheckPercentageHealthy"
	)

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

func (self *mockUptimeClient) GetMetricStatistics(ctx context.Context, params *cloudwatch.GetMetricStatisticsInput, optFns ...func(*cloudwatch.Options)) (out *cloudwatch.GetMetricStatisticsOutput, err error) {
	var metricsUnit types.StandardUnit = types.StandardUnitPercent
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

func (self *mockUptimeClient) Options() cloudwatch.Options {
	return cloudwatch.Options{
		Region: "us-east-1",
	}
}

func TestCMDImportsUptimeWithMock(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		client *mockUptimeClient = &mockUptimeClient{}
		ctx    context.Context   = t.Context()
		log    *slog.Logger      = logger.New("error")
		dir    string            = t.TempDir()
		dbPath string            = filepath.Join(dir, "test-import-mock-uptime.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	err = importUptime(ctx, log, client, db, &UptimeOpts{
		AccountID: "mock-account-a",
		Day:       times.AsYMDString(times.Yesterday()),
	})
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}

func TestCMDImportsUptimeWithoutMock(t *testing.T) {

	var (
		err    error
		client *cloudwatch.Client
		db     *sqlx.DB
		region string          = "us-east-1"
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-uptime.db")
	)
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}

	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	client, err = awsclients.New[*cloudwatch.Client](ctx, log, region)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	err = importUptime(ctx, log, client, db, &UptimeOpts{
		AccountID: awsid.AccountID(ctx, log, region),
		Day:       times.AsYMDString(times.Yesterday()),
	})
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}
