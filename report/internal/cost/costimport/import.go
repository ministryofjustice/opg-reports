package costimport

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	_ "github.com/mattn/go-sqlite3"
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

// Model represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type Model struct {
	Region    string `json:"region,omitempty"`      // AWS Region
	Service   string `json:"service,omitempty"`     // The AWS service name
	Month     string `json:"month,omitempty"`       // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty"`        // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id,omityempty"` // the actual account id - string as it can have leading zeros. Use in joins as well
}

// Client is used to allow mocking and is a proxy for *costexplorer.Client
type Client interface {
	// GetCostAndUsage method signature from the *costexplorer.Client
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params

	DateStart time.Time `json:"date_start"` // start date, this will be reset to start of the month (and expanded to capture historical data)
	DateEnd   time.Time `json:"date_end"`   // end date
	AccountID string    `json:"account_id"` // AccountID provided by awsid.AccountID
}

func Import(ctx context.Context, client Client, in *Args) (err error) {
	var (
		options *costexplorer.GetCostAndUsageInput
		result  *costexplorer.GetCostAndUsageOutput
		costs   []*Model
		log     *slog.Logger = cntxt.GetLogger(ctx).With("package", "costimport", "func", "Import")
	)
	log.Info("starting ...", "db", in.DB, "date_start", in.DateStart, "date_end", in.DateEnd)
	options = getCostAndUsageInput(
		times.AsYMDString(times.ResetMonth(in.DateStart)),
		times.AsYMDString(in.DateEnd))

	// make the api call
	result, err = client.GetCostAndUsage(ctx, options)
	if err != nil {
		log.Error("error getting cost and usage", "err", err.Error())
		return
	}
	// covnerto models
	costs, err = toModels(ctx, in.AccountID, result)
	if err != nil {
		log.Error("error converting cost and usage", "err", err.Error())
		return
	}

	// fmt.Println(dump.Any(costs))
	// now write to db
	err = dbx.Insert(ctx, InsertStatement, costs, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.Info("complete.")
	return
}

// toModels converts the raw data into a list of models ready to write to the database
func toModels(ctx context.Context, account string, result *costexplorer.GetCostAndUsageOutput) (costs []*Model, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "costimport", "func", "toModels")

	costs = []*Model{}
	log.Debug("starting toModels ... ")

	for _, result := range result.ResultsByTime {
		var day string = *result.TimePeriod.Start
		for _, group := range result.Groups {
			var service string = *&group.Keys[0]
			var region string = *&group.Keys[1]
			for _, cost := range group.Metrics {
				var item = &Model{
					AccountID: account,
					Month:     times.ToYMString(day),
					Service:   service,
					Region:    region,
					Cost:      *cost.Amount,
				}
				costs = append(costs, item)
			}
		}
	}
	log.With("count", len(costs)).Debug("complete.")

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
