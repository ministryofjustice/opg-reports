package cost

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// Flatten converts the raw aws api output to just the service and cost data we want in a simple struct
// and attaches the account to that data
func Flatten(raw *costexplorer.GetCostAndUsageOutput, aId string, aName string, aLabel string, aUnit string, aOrg string, aEnv string) (costs []*Cost, err error) {

	slog.Debug(fmt.Sprintf("FlatCosts: flattening cost data"))
	costs = []*Cost{}

	for _, resultByTime := range raw.ResultsByTime {
		day := *resultByTime.TimePeriod.Start

		for _, costGroup := range resultByTime.Groups {
			service := *costGroup.Keys[0]
			region := *costGroup.Keys[1]

			for _, costMetric := range costGroup.Metrics {
				amount := *costMetric.Amount
				c := New(nil)
				c.Service = service
				c.Region = region
				c.Date = day
				c.Cost = amount

				c.AccountId = aId
				c.AccountName = aName
				c.AccountLabel = aLabel
				c.AccountUnit = aUnit
				c.AccountOrganisation = aOrg
				c.AccountEnvironment = aEnv

				costs = append(costs, c)
			}
		}
	}

	return costs, err
}
