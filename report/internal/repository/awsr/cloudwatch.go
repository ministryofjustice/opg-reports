package awsr

import (
	"fmt"
	"log/slog"
	"opg-reports/report/internal/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	uptimeNamespace string             = "AWS/Route53"
	uptimeMetric    string             = "HealthCheckPercentageHealthy"
	uptimeUnit      types.StandardUnit = types.StandardUnitPercent
	uptimeStat      types.Statistic    = types.StatisticAverage
)

// getUptimePeriod works out what period to use based on the api contraints and the dates being requested.
// Based on the below details:
//   - Data points with a period of less than 60 seconds are available for 3 hours. These data points are high-resolution metrics and are available only for custom metrics that have been defined with a StorageResolution of 1.
//   - Data points with a period of 60 seconds (1-minute) are available for 15 days.
//   - Data points with a period of 300 seconds (5-minute) are available for 63 days.
//   - Data points with a period of 3600 seconds (1 hour) are available for 455 days (15 months).
func getUptimePeriod(start time.Time) (period int32, err error) {
	var (
		now       time.Time     = time.Now().UTC()
		day       time.Duration = (time.Hour * 24)
		days15    time.Duration = (15 * day)
		days63    time.Duration = (63 * day)
		days455   time.Duration = (455 * day)
		hoursDiff float64       = now.Sub(start).Hours()
	)

	if hoursDiff < days15.Hours() {
		period = 60
	} else if hoursDiff < days63.Hours() {
		period = 300
	} else if hoursDiff < days455.Hours() {
		period = 3600
	} else {
		err = fmt.Errorf("date range old, max range of 455 days")
	}
	return

}

// GetUptimeStats uses the list of metrics provided to find and return the accumlated average uptime percentage between the start & end date
func (self *Repository) GetUptimeStats(client ClientCloudWatchMetricStats, metrics []types.Metric, start time.Time, end time.Time) (datapoints []types.Datapoint, err error) {
	var (
		input      *cloudwatch.GetMetricStatisticsInput
		output     *cloudwatch.GetMetricStatisticsOutput
		period     int32
		dimensions []types.Dimension = []types.Dimension{}
		log        *slog.Logger      = self.log.With("operation", "GetUptimeStats")
	)

	log.Debug("getting route53 uptime stats for metrics ... ")

	if start.String() == "" || end.String() == "" {
		err = fmt.Errorf("start or end date missing: [start:%v, end:%v]", start, end)
	}
	// try to work out a suitable time period
	period, err = getUptimePeriod(start)
	if err != nil {
		return
	}
	// merge all metric dimenstions together
	for _, metric := range metrics {
		dimensions = append(dimensions, metric.Dimensions...)
	}

	input = &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(uptimeNamespace),
		MetricName: aws.String(uptimeMetric),
		Period:     utils.Ptr(period),
		StartTime:  &start,
		EndTime:    &end,
		Statistics: []types.Statistic{uptimeStat},
		Unit:       uptimeUnit,
		Dimensions: dimensions,
	}
	log.With("input", input).Debug("getting metrics stats ... ")

	output, err = client.GetMetricStatistics(self.ctx, input)
	if err != nil {
		return
	}
	datapoints = output.Datapoints
	log.With("count", len(datapoints)).Debug("found metric datapoints ... ")

	return
}

// GetUptimeMetrics returns metric details from cloudwatch for the route53 health check
// that can then be used to determine uptime information in other calls
func (self *Repository) GetUptimeMetrics(client ClientCloudWatchMetricsLister) (metrics []types.Metric, err error) {
	var (
		output          *cloudwatch.ListMetricsOutput
		log             *slog.Logger                 = self.log.With("operation", "GetUptimeMetricsList")
		metricListInput *cloudwatch.ListMetricsInput = &cloudwatch.ListMetricsInput{
			Namespace:  aws.String(uptimeNamespace),
			MetricName: aws.String(uptimeMetric),
		}
	)

	log.Debug("getting route53 uptime metrics ... ")

	output, err = client.ListMetrics(self.ctx, metricListInput)
	if err != nil {
		return
	}
	metrics = output.Metrics
	return
}
