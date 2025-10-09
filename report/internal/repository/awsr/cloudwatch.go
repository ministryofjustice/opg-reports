package awsr

import (
	"fmt"
	"log/slog"
	"opg-reports/report/internal/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// GetSuitableUptimePeriod works out what period to use based on the api contraints and the dates being requested.
// Based on the below details:
//   - Data points with a period of less than 60 seconds are available for 3 hours. These data points are high-resolution metrics and are available only for custom metrics that have been defined with a StorageResolution of 1.
//   - Data points with a period of 60 seconds (1-minute) are available for 15 days.
//   - Data points with a period of 300 seconds (5-minute) are available for 63 days.
//   - Data points with a period of 3600 seconds (1 hour) are available for 455 days (15 months).
func GetSuitableUptimePeriod(start time.Time) (period int32, err error) {
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

// GetUptimeData fetches both the metrics and the stats for those metrics, effectively
// doing both `GetUptimeMetrics` & `GetUptimeDatapoints`
//
//	options = &cloudwatch.GetMetricStatisticsInput{
//		Namespace:  &Namespace,
//		MetricName: &MetricName,
//		Period:     &Period,
//		StartTime:  &Start,
//		EndTime:    &End,
//		Statistics: []types.Statistic{options.Statistic},
//		Unit:       options.Unit,
//	}
func (self *Repository) GetUptimeData(
	client ClientCloudWatchUptime,
	options *cloudwatch.GetMetricStatisticsInput,
) (stats []map[string]string, err error) {
	var (
		metrics        []types.Metric
		datapoints     []types.Datapoint
		log            *slog.Logger             = self.log.With("operation", "GetUptimeData")
		metricsOptions *GetUptimeMetricsOptions = &GetUptimeMetricsOptions{
			MetricName: *options.MetricName,
			Namespace:  *options.Namespace,
		}
	)
	stats = []map[string]string{}
	// get the metrics
	metrics, err = self.GetUptimeMetrics(client, metricsOptions)
	if err != nil {
		return
	}
	// get the stats for those metrics
	datapoints, err = self.GetUptimeDatapoints(client, metrics, options)
	if err != nil {
		return
	}
	log.Debug("converting datapoints to slice map ... ")

	for _, point := range datapoints {
		stats = append(stats, map[string]string{
			"average": fmt.Sprintf("%g", *point.Average),
			"date":    point.Timestamp.Format(utils.DATE_FORMATS.Full),
		})
	}

	return

}

// GetUptimeDatapoints uses the list of metrics provided to find and return the accumlated average uptime percentage between the start & end date
//
//	options = &cloudwatch.GetMetricStatisticsInput{
//		Namespace:  &Namespace,
//		MetricName: &MetricName,
//		Period:     &Period,
//		StartTime:  &Start,
//		EndTime:    &End,
//		Statistics: []types.Statistic{options.Statistic},
//		Unit:       options.Unit,
//	}
func (self *Repository) GetUptimeDatapoints(client ClientCloudWatchMetricStats, metrics []types.Metric, options *cloudwatch.GetMetricStatisticsInput) (datapoints []types.Datapoint, err error) {
	var (
		output     *cloudwatch.GetMetricStatisticsOutput
		dimensions []types.Dimension = []types.Dimension{}
		log        *slog.Logger      = self.log.With("operation", "GetUptimeDatapoints")
	)

	log.With("options", options).Debug("getting route53 uptime stats for metrics ... ")

	if options.StartTime == nil || options.EndTime == nil {
		err = fmt.Errorf("start or end date missing: \n%v\n", *options)
		return
	}
	// work out period if there isnt one set
	if options.Period == nil {
		period, err := GetSuitableUptimePeriod(*options.StartTime)
		if err != nil {
			options.Period = &period
		}
	}

	if options.Period == nil || *options.Period <= 0 {
		err = fmt.Errorf("time period value missing")
		return
	}

	// merge all metric dimenstions together
	for _, metric := range metrics {
		dimensions = append(dimensions, metric.Dimensions...)
	}
	options.Dimensions = dimensions

	log.With("input", options).Debug("getting metrics stats ... ")

	output, err = client.GetMetricStatistics(self.ctx, options)
	if err != nil {
		return
	}
	log.With("count", len(output.Datapoints)).Debug("found metric datapoints ... ")

	datapoints = output.Datapoints

	return
}

type GetUptimeMetricsOptions struct {
	Namespace  string // "AWS/Route53"
	MetricName string // "HealthCheckPercentageHealthy"
}

// GetUptimeMetrics returns metric details from cloudwatch for the route53 health check
// that can then be used to determine uptime information in other calls
func (self *Repository) GetUptimeMetrics(client ClientCloudWatchMetricsLister, options *GetUptimeMetricsOptions) (metrics []types.Metric, err error) {
	var (
		output *cloudwatch.ListMetricsOutput
		log    *slog.Logger = self.log.With("operation", "GetUptimeMetricsList")
	)

	log.With("options", options).Debug("getting route53 uptime metrics ... ")

	output, err = client.ListMetrics(self.ctx, &cloudwatch.ListMetricsInput{
		Namespace:  &options.Namespace,
		MetricName: &options.MetricName,
	})
	if err != nil {
		return
	}
	metrics = output.Metrics

	return
}
