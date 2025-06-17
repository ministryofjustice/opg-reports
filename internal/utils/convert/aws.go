package convert

import (
	"time"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/ministryofjustice/opg-reports/internal/services/awscosts"
)

// FromAwsCostExplorerToAwsCosts converts the costexplorer data set into a series of
// AwsCosts
// This will not include any account information, that will need to be added elsewhere
func FromAwsCostExplorerToAwsCosts(response *costexplorer.GetCostAndUsageOutput, account string) (costs []awscosts.Cost) {
	costs = []awscosts.Cost{}

	now := time.Now().UTC().Format(time.RFC3339)

	for _, resultByTime := range response.ResultsByTime {
		day := *resultByTime.TimePeriod.Start
		for _, costGroup := range resultByTime.Groups {
			service := *costGroup.Keys[0]
			region := *costGroup.Keys[1]

			for _, costMetric := range costGroup.Metrics {
				amount := *costMetric.Amount

				cost := awscosts.Cost{
					AccountID: account,
					CreatedAt: now,
					Region:    region,
					Service:   service,
					Date:      day,
					Cost:      amount,
				}
				costs = append(costs, cost)
			}
		}

	}

	return
}
