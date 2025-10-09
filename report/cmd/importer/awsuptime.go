package main

import (
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
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

const uptimeLongDesc string = `
awsuptime will call the health check api to retrieve data for yesterday only.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used

`

var (
	uptimeDayFlag string         = "" // represents --day="YYYY-MM-DD"
	awsuptimeCmd  *cobra.Command = &cobra.Command{
		Use:   "awsuptime",
		Short: "awsuptime fetches data from the r53 health check endpoints for yesterday only",
		Long:  uptimeLongDesc,
		RunE:  awsUptimeRunner,
	} // awsuptimeCmd imports data from the aws api directly
)

// awsUptimeRunner used by the cobra command (awsuptimeCmd) to process the cli request to fetch data from
// the aws api and import to local database
func awsUptimeRunner(cmd *cobra.Command, args []string) (err error) {
	var (
		uptime    []map[string]string                                             // uptime converted to map
		accountID string                                                          // account if from the caller identity
		start     = utils.StringToTimeReset(uptimeDayFlag, utils.TimeIntervalDay) // start of yesterday
		// clients et al
		stsClient        = awsr.DefaultClient[*sts.Client](ctx, conf.Aws.GetRegion()) // identity client
		awsStore         = awsr.Default(ctx, log, conf)                               // generic aws store
		cloudwatchClient = awsr.DefaultClient[*cloudwatch.Client](ctx, uptimeRegion)  // client, have to fix region to get the correct data
		// inserts
		sqClient   = sqlr.DefaultWithSelect[*api.AwsUptime](ctx, log, conf)
		apiService = api.Default[*api.AwsUptime](ctx, log, conf)
	)

	accountID, err = awsAccountID(stsClient, awsStore)
	if err != nil {
		return
	}

	uptime, err = awsUptimeGetData(cloudwatchClient, awsStore, start)
	if err != nil {
		return
	}

	utils.Dump(uptime)
	utils.Dump(accountID)

	err = awsUptimeInsert(sqClient, apiService, accountID, uptime)

	return
}

// awsCostsInsert adds new data into the existing database for aws costs
func awsUptimeInsert(
	client sqlr.RepositoryWriter,
	service *api.Service[*api.AwsUptime],
	accountID string,
	apiData []map[string]string,
) (err error) {
	var dbData = []*api.AwsUptime{}

	// convert to AwsCosts struct
	err = utils.Convert(apiData, &dbData)
	if err != nil {
		log.Error("error converting", "err", err.Error())
		return
	}

	// add account id into each row
	for _, c := range dbData {
		c.AwsAccountID = accountID
	}

	// insert
	_, err = service.PutAwsUptime(client, dbData)
	if err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}

	return
}

// awsUptimeGetData calls the api to fetch metrics and the uptime percentages from the aws api
func awsUptimeGetData(
	client awsr.ClientCloudWatchUptime,
	store awsr.RepositoryCloudwatchUptime,
	start time.Time,
) (uptime []map[string]string, err error) {
	var (
		end     = start.AddDate(0, 0, 1)
		options = &cloudwatch.GetMetricStatisticsInput{
			Namespace:  utils.Ptr(uptimeNamespace),
			MetricName: utils.Ptr(uptimeMetric),
			StartTime:  utils.Ptr(start),
			EndTime:    utils.Ptr(end),
			Statistics: []types.Statistic{uptimeStatistic},
			Unit:       uptimeUnit,
		}
	)

	uptime, err = store.GetUptimeData(client, options)

	return
}

func init() {
	awsuptimeCmd.Flags().StringVar(&uptimeDayFlag, "day", utils.Yesterday().Format(utils.DATE_FORMATS.YMD), "The day to get uptime data for. (YYYY-MM-DD)")
}
