package infracost

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// mockGetter returns a positive result with test data
type mockGetter struct{}

func (self *mockGetter) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (out *costexplorer.GetCostAndUsageOutput, err error) {
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

// mockGetterFailed returns a fake error
type mockGetterFailed struct{}

func (self *mockGetterFailed) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (out *costexplorer.GetCostAndUsageOutput, err error) {
	err = fmt.Errorf("mock error")
	return
}

// check interfaces
var _ AwsClient = &mockGetter{}
var _ AwsClient = &mockGetterFailed{}

// TestInfracostsWithMock uses mock struct above to test logic and values
// without calling AWS api.
func TestDomainInfracostsWithMock(t *testing.T) {
	var (
		err    error
		client *mockGetter
		r      []*infracostmodels.Cost
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		now    time.Time       = time.Now().UTC()
		start  time.Time       = now.AddDate(0, -4, 0)
		end    time.Time       = now.AddDate(0, -3, 0)
	)

	client = &mockGetter{}
	r, err = GetCostData(ctx, log, client, &Options{Start: start, End: end})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(r) <= 0 {
		t.Error("failed to find cost data")
	}
}

// TestInfracostsWithoutMock uses concreate AWS SDK methods rather than mocks
// to directly call the api.
//
// Only runs if there is actually AWS_SESSION_TOKEN env var present.
//
// Run: aws-vault exec use-development-breakglass -- make test name="TestCostsWithoutMock"
func TestDomainInfracostsWithoutMock(t *testing.T) {
	var (
		err       error
		client    *costexplorer.Client
		accountId string
		r         []*infracostmodels.Cost
		ctx       context.Context = t.Context()
		log       *slog.Logger    = logger.New("error")
		now       time.Time       = time.Now().UTC()
		start     time.Time       = now.AddDate(0, -4, 0)
		end       time.Time       = now.AddDate(0, -3, 0)
	)

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client, err = awsclients.New[*costexplorer.Client](ctx, log, "eu-west-1")
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
			t.FailNow()
		}
		accountId = awsid.AccountID(ctx, log, "eu-west-1")
		r, err = GetCostData(ctx, log, client, &Options{Start: start, End: end, AccountID: accountId})
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}

		if len(r) <= 0 {
			t.Error("failed to find cost data")
		}
	} else {
		t.SkipNow()
	}

}

// TestInfracostsWithFailure checks the error returned matches custom error type
func TestDomainInfracostsWithFailure(t *testing.T) {
	var (
		err    error
		client *mockGetterFailed
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		now    time.Time       = time.Now().UTC()
		start  time.Time       = now.AddDate(0, -4, 0)
		end    time.Time       = now.AddDate(0, -3, 0)
	)

	client = &mockGetterFailed{}
	_, err = GetCostData(ctx, log, client, &Options{Start: start, End: end})
	if err == nil {
		t.Errorf("expected error, but nothing returned")
	}
	if !errors.Is(err, ErrGettingCostData) {
		t.Errorf("expected known error type, instead recieved [%v]", err)
	}

}
