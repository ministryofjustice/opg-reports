package main

import (
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

const (
	uptimeNamespace string             = "AWS/Route53"
	uptimeMetric    string             = "HealthCheckPercentageHealthy"
	uptimeUnit      types.StandardUnit = types.StandardUnitPercent
	uptimeStatistic types.Statistic    = types.StatisticAverage
	uptimeRegion    string             = "us-east-1"
)

// awsuptimeCmd imports data from the cost explorer api directly
var awsuptimeCmd = &cobra.Command{
	Use:   "awsuptime",
	Short: "awsuptime fetches data from the r53 health check endpoints",
	Long: `
awsuptime will call the health check api to retrieve data for specified period.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			uptime    []map[string]string                                           // uptime converted to map
			accountID string                                                        // account if from the caller identity
			start     = utils.StringToTimeReset(flagMonth, utils.TimeIntervalMonth) // start of the month
			// clients et al
			stsClient        = awsr.DefaultClient[*sts.Client](ctx, conf.Aws.GetRegion()) // identity client
			awsStore         = awsr.Default(ctx, log, conf)                               // generic aws store
			cloudwatchClient = awsr.DefaultClient[*cloudwatch.Client](ctx, uptimeRegion)  // client, have to fix region to get the correct data
		)

		accountID, err = awsCostsGetAccountID(stsClient, awsStore)
		if err != nil {
			return
		}

		uptime, err = awsUptimeGetData(cloudwatchClient, awsStore, start)
		if err != nil {
			return
		}

		utils.Dump(uptime)
		utils.Dump(accountID)

		return
	},
}

func awsUptimeGetData(
	client awsr.ClientCloudWatchUptime,
	store awsr.RepositoryCloudwatchUptime,
	start time.Time,
) (uptime []map[string]string, err error) {
	var (
		end     = start.AddDate(0, 1, 0)
		options = &cloudwatch.GetMetricStatisticsInput{
			Namespace:  utils.Ptr(uptimeNamespace),
			MetricName: utils.Ptr(uptimeMetric),
			StartTime:  utils.Ptr(start),
			EndTime:    utils.Ptr(end),
			Statistics: []types.Statistic{uptimeStatistic},
			Unit:       uptimeUnit,
		}
	)
	utils.Dump(options)

	uptime, err = store.GetUptimeData(client, options)

	return
}

// awsUptimeGetAccountID returns the account id
func awsUptimeGetAccountID(client awsr.ClientSTSCaller, store awsr.RepositorySTS) (accountID string, err error) {

	caller, err := store.GetCallerIdentity(client)
	if caller != nil {
		accountID = *caller.Account
	}
	return
}
