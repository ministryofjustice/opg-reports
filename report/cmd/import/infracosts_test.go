package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/utils"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/jmoiron/sqlx"
)

// mockInfracostClient returns a positive result with test data
type mockInfracostClient struct{}

func (self *mockInfracostClient) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (out *costexplorer.GetCostAndUsageOutput, err error) {
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

func TestCMDImportsInfracostsWithMock(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		client *mockInfracostClient = &mockInfracostClient{}
		ctx    context.Context      = t.Context()
		log    *slog.Logger         = logger.New("error")
		dir    string               = t.TempDir()
		dbPath string               = filepath.Join(dir, "test-import-mock-infracosts.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	err = importInfracosts(ctx, log, client, db, &InfraOpts{
		AccountID:            "mock-account-id",
		IncludePreviousMonth: true,
		EndDate:              "2025-11-12",
	})
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}

func TestCMDImportsInfracostsWithoutMock(t *testing.T) {

	var (
		err    error
		client *costexplorer.Client
		db     *sqlx.DB
		region string          = "eu-west-1"
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-infracosts.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = awsclients.New[*costexplorer.Client](ctx, log, region)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	err = importInfracosts(ctx, log, client, db, &InfraOpts{
		AccountID:            awsid.AccountID(ctx, log, region),
		IncludePreviousMonth: true,
		EndDate:              "2025-11-12",
	})
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}
