// package uptime handles fetching and returning uptime data via route53 healthchecks.
//
// Fecthes list of metrics from the "AWS/Route53" namespace and then data for those
// metrics to decide uptime values
package uptime

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/utils/ptr"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// AwsClient is used to allow mocking and is a proxy for *cloudwatch.Client
// and the methods the function calls
type AwsClient interface {
	ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error)
	GetMetricStatistics(ctx context.Context, params *cloudwatch.GetMetricStatisticsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricStatisticsOutput, error)
	Options() cloudwatch.Options
}

// GetUptimeDataOptions provides struct to deal with optional settings that can be changed
// for hte api calls
type GetUptimeDataOptions struct {
	Start time.Time
	End   time.Time
}

// these are fixed values used for the api calls
const (
	metricRegion     string             = "us-east-1"
	metricsNamespace string             = "AWS/Route53"
	metricsName      string             = "HealthCheckPercentageHealthy"
	metricsStatistic types.Statistic    = types.StatisticAverage
	metricsUnit      types.StandardUnit = types.StandardUnitPercent
)

// GetUptimeData makes multiple cloudwatch api calls to find uptime data for the months represented within the options.
//
// Firstly, it will find health check metric details and their respective id'. These are then used in the second call as part
// of the dimensions
//
// The route53 metrics require the client region to be set as
func GetUptimeData(ctx context.Context, log *slog.Logger, client *cloudwatch.Client, options *GetUptimeDataOptions) (stats *cloudwatch.GetMetricStatisticsOutput, err error) {
	var (
		list      *cloudwatch.ListMetricsOutput
		setRegion string = client.Options().Region
	)

	log = log.With("package", "uptime", "func", "GetUptimeData")
	log.With("options", options).Debug("starting ...")
	// check region
	if setRegion != metricRegion {
		err = ErrIncorrectRegion
		log.Error("incorrect region used in client - requires us-east-1", "region", setRegion)
		return
	}

	// fetch the list of all metrics from the api
	list, err = getHealthCheckMetrics(ctx, log, client)
	if err != nil {
		return
	}

	// get all the datapoints for each of the metrics
	stats, err = getHealthCheckStatistics(ctx, log, client, list, options)
	if err != nil {
		return
	}

	log.Debug("complete")
	return
}

// getHealthCheckStatistics uses the list of metrics from getHealthCheckMetrics to fetch the uptime datapoints
// from the api and return the api call values
//
// T is *cloudwatch.Client
func getHealthCheckStatistics[T AwsClient](ctx context.Context, log *slog.Logger, client T, list *cloudwatch.ListMetricsOutput, options *GetUptimeDataOptions) (stats *cloudwatch.GetMetricStatisticsOutput, err error) {

	var statsInput *cloudwatch.GetMetricStatisticsInput = getHeathCheckMetricStatsOptions(list, options)

	log = log.With("package", "uptime", "func", "getHealthCheckStatistics")
	log.Debug("starting ...")
	log.With("statsInput", statsInput).Debug("getting metrics statistics ...")

	// try and get the stats
	stats, err = client.GetMetricStatistics(ctx, statsInput)
	if err != nil {
		log.Error("error getting metric statistics", "err", err.Error())
		err = errors.Join(ErrFailedGettingMetricStats, err)
		return
	}

	log.With("count", len(stats.Datapoints)).Debug("complete.")
	return
}

// getHealthCheckMetrics returns the list of metrics to use for uptime data
//
// T is *cloudwatch.Client
func getHealthCheckMetrics[T AwsClient](ctx context.Context, log *slog.Logger, client T) (list *cloudwatch.ListMetricsOutput, err error) {

	var listOptions *cloudwatch.ListMetricsInput = &cloudwatch.ListMetricsInput{
		Namespace:  ptr.Ptr(metricsNamespace),
		MetricName: ptr.Ptr(metricsName),
	}

	log = log.With("package", "uptime", "func", "getHealthCheckMetrics")
	log.Debug("starting ...")
	log.Debug("fetching metric data for account ...")

	list, err = client.ListMetrics(ctx, listOptions)
	if err != nil {
		log.Error("error getting list of metrics", "err", err.Error())
		err = errors.Join(ErrFailedGettingMetricsList, err)
		return
	}

	log.With("count", len(list.Metrics)).Debug("complete.")
	return
}

// getHeathCheckMetricStatsOptions is used to generate a suitable GetMetricStatisticsInput struct
// that contains all of the metric dimensions and uses the correct period and units that we need
// to fetch uptime data
func getHeathCheckMetricStatsOptions(list *cloudwatch.ListMetricsOutput, options *GetUptimeDataOptions) (opts *cloudwatch.GetMetricStatisticsInput) {
	var (
		period     int32             = getPeriod(options.Start)
		dimensions []types.Dimension = []types.Dimension{}
	)

	// merge all of the metric dimensions together
	for _, metric := range list.Metrics {
		dimensions = append(dimensions, metric.Dimensions...)
	}
	// generate the input struct
	opts = &cloudwatch.GetMetricStatisticsInput{
		Namespace:  ptr.Ptr(metricsNamespace),
		MetricName: ptr.Ptr(metricsName),
		StartTime:  ptr.Ptr(options.Start),
		EndTime:    ptr.Ptr(options.End),
		Period:     ptr.Ptr(period),
		Unit:       metricsUnit,
		Statistics: []types.Statistic{metricsStatistic},
		Dimensions: dimensions,
	}

	return
}

// getPeriod works out what period to use based on the api contraints and the start date being requested.
//
// Based on the below details:
//   - Data points with a period of less than 60 seconds are available for 3 hours. These data points are high-resolution metrics and are available only for custom metrics that have been defined with a StorageResolution of 1.
//   - Data points with a period of 60 seconds (1-minute) are available for 15 days.
//   - Data points with a period of 300 seconds (5-minute) are available for 63 days.
//   - Data points with a period of 3600 seconds (1 hour) are available for 455 days (15 months).
//
// Don't use any granularity under 60 seconds.
func getPeriod(start time.Time) (period int32) {
	var (
		now       time.Time     = time.Now().UTC()
		day       time.Duration = (time.Hour * 24)
		days15    time.Duration = (15 * day)
		days63    time.Duration = (63 * day)
		days455   time.Duration = (455 * day)
		hoursDiff float64       = now.Sub(start).Hours()
	)
	period = 3600

	if hoursDiff < days15.Hours() {
		period = 60
	} else if hoursDiff < days63.Hours() {
		period = 300
	} else if hoursDiff < days455.Hours() {
		period = 3600
	}
	return

}
