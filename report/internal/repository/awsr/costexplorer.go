package awsr

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

// GetCostData calls the cost explorer api and returns cost and usage data in the bellow format:
//
//	[{
//		"cost": "100.2335589669",
//		"date": "2025-01-01",
//		"region": "eu-west-1",
//		"service": "AmazonCloudWatch"
//	  },
//	  {
//		"cost": "10.6836594846",
//		"date": "2025-03-01",
//		"region": "eu-west-2",
//		"service": "AmazonCloudWatch"
//	  }]
//
// Equilivant cli call:
//
//	aws-vault exec ${profile} -- aws ce get-cost-and-usage \
//		--time-period Start=2025-03-01,End=2025-04-01 \
//		--granularity MONTHLY \
//		--metrics "UnblendedCost" \
//		--group-by Type=DIMENSION,Key=SERVICE Type=DIMENSION,Key=REGION
//
// Note: API limits grouping to 2, so we cant get linked account details at the same time
func (self *Repository) GetCostData(client ClientCostExplorerGetter, options *costexplorer.GetCostAndUsageInput) (values []map[string]string, err error) {
	var log *slog.Logger = self.log.With("operation", "GetCostData", "options", options)
	log.Debug("getting cost data ... ")
	values = []map[string]string{}

	// validate the options passed along for correct settings
	if options.TimePeriod == nil || options.TimePeriod.Start == nil || options.TimePeriod.End == nil {
		err = fmt.Errorf("time period is not set correctly:\n %v", options.TimePeriod)
		return
	}
	if options.Granularity == "" {
		err = fmt.Errorf("granularity is not set")
		return
	}
	if len(options.Metrics) <= 0 {
		err = fmt.Errorf("metrics are not set")
		return
	}
	if len(options.GroupBy) <= 0 {
		err = fmt.Errorf("grouping is not set")
		return
	}

	// call the api
	out, err := client.GetCostAndUsage(self.ctx, options)
	if err != nil {
		return
	}
	// convert results into a map
	log.Debug("coverting cost data ... ")
	for _, result := range out.ResultsByTime {
		var day string = *result.TimePeriod.Start

		for _, group := range result.Groups {
			var service string = *&group.Keys[0]
			var region string = *&group.Keys[1]

			for _, cost := range group.Metrics {
				var cost string = *cost.Amount

				values = append(values, map[string]string{
					"date":    day,
					"service": service,
					"region":  region,
					"cost":    cost,
				})
			}
		}
	}

	return
}
