package costimport

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// Client is used to allow mocking and is a proxy for *costexplorer.Client
type Client interface {
	// GetCostAndUsage method signature from the *costexplorer.Client
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

type Input struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params
	// MigrationFile string    `json:"migration_file"` // migration file
	DateStart time.Time `json:"date_start"` // start date, this will be reset to start of the month (and expanded to capture historical data)
	DateEnd   time.Time `json:"date_end"`   // end date
	AccountID string    `json:"account_id"` // AccountID provided by awsid.AccountID
}

func Import(ctx context.Context, client Client, in *Input) (err error) {
	var (
		options *costexplorer.GetCostAndUsageInput
		// result  *costexplorer.GetCostAndUsageOutput
		log *slog.Logger = cntxt.GetLogger(ctx).WithGroup("costimport")
	)
	log.Info("starting ...", "db", in.DB, "date_start", in.DateStart, "date_end", in.DateEnd)
	options = getCostAndUsageInput(
		times.AsYMDString(times.ResetMonth(in.DateStart)),
		times.AsYMDString(in.DateEnd))

	// make the api call
	_, err = client.GetCostAndUsage(ctx, options)
	if err != nil {
		log.Error("error getting cost and usage", "err", err.Error())
		return
	}

	log.Info("complete.")
	return
}

// Options returns a CostAndUsageInput struct formatted with expected values
// for cost data using the start and end date
//
// `start` is reset the begining of the month to make sure a full month costs are set
func getCostAndUsageInput(start string, end string) *costexplorer.GetCostAndUsageInput {
	var (
		service string   = "SERVICE"
		region  string   = "REGION"
		metrics []string = []string{"UnblendedCost"}
	)
	return &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		TimePeriod: &types.DateInterval{
			Start: &start,
			End:   &end,
		},
		Metrics: metrics,
		GroupBy: []types.GroupDefinition{
			{Type: types.GroupDefinitionTypeDimension, Key: &service},
			{Type: types.GroupDefinitionTypeDimension, Key: &region},
		},
	}

}
