// package uptime handles fetching and returning uptime data via route53 healthchecks.
//
// Fecthes list of metrics from the "AWS/Route53" namespace and then data for those
// metrics to decide uptime values
package uptime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/ptr"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var (
	ErrIncorrectRegion          = errors.New("metrics must be fetched via us-east-1.")
	ErrFailedGettingMetricsList = errors.New("failed to get metrics list with error.")
	ErrFailedGettingMetricStats = errors.New("failed to get metric statistics with error.")
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
	Start     time.Time
	End       time.Time
	AccountID string
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
// of the dimensions.
//
// The route53 metrics require the client region to be set as us-east-1.
//
// T is `*cloudwatch.Client`.
func GetUptimeData[T AwsClient](ctx context.Context, log *slog.Logger, client T, options *GetUptimeDataOptions) (data []*uptimemodels.Uptime, err error) {
	var (
		list      *cloudwatch.ListMetricsOutput
		statsOpts *cloudwatch.GetMetricStatisticsInput
		stats     *cloudwatch.GetMetricStatisticsOutput
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
	log.Debug("getting lisst of metrics ...")
	list, err = getHealthCheckMetrics(ctx, log, client)
	if err != nil {
		return
	}

	// get all the datapoints for each of the metrics
	log.Debug("getting metric stats ...")
	stats, statsOpts, err = getHealthCheckStatistics(ctx, log, client, list, options)
	if err != nil {
		return
	}

	log.Debug("coverting to models ...")
	data, err = toModels(ctx, log, options.AccountID, *statsOpts.Period, stats)
	if err != nil {
		return
	}

	log.With("count", len(data)).Debug("complete")
	return
}

// toModels converts the raw data into a list of models ready to write to the database
func toModels(ctx context.Context, log *slog.Logger, account string, period int32, result *cloudwatch.GetMetricStatisticsOutput) (data []*uptimemodels.Uptime, err error) {
	var (
		grouped map[string]float64 = map[string]float64{}
		counter map[string]int     = map[string]int{}
	)
	data = []*uptimemodels.Uptime{}
	log = log.With("package", "uptime", "func", "toModels")
	log.Debug("starting ... ")

	// create a sum and count of each month uptime to then create the average entries
	for _, point := range result.Datapoints {
		var (
			month time.Time = times.ResetDay(*point.Timestamp)
			key   string    = times.AsYMDString(month)
		)
		// find or update the value in
		if _, ok := grouped[key]; !ok {
			grouped[key] = 0.0
			counter[key] = 0
		}
		grouped[key] += *point.Average
		counter[key]++
	}
	// now generate the average values
	for key, sum := range grouped {
		var (
			count   int                  = counter[key]
			average float64              = (sum / float64(count))
			up      *uptimemodels.Uptime = &uptimemodels.Uptime{
				Date:        key,
				Average:     fmt.Sprintf("%g", average),
				Granularity: fmt.Sprintf("%d", period),
				AccountID:   account,
			}
		)
		data = append(data, up)
	}
	log.Debug("complete.")
	return
}

// getHealthCheckStatistics uses the list of metrics from getHealthCheckMetrics to fetch the uptime datapoints
// from the api and return the api call values
//
// T is *cloudwatch.Client
func getHealthCheckStatistics[T AwsClient](ctx context.Context, log *slog.Logger, client T, list *cloudwatch.ListMetricsOutput, options *GetUptimeDataOptions) (stats *cloudwatch.GetMetricStatisticsOutput, statsInput *cloudwatch.GetMetricStatisticsInput, err error) {

	log = log.With("package", "uptime", "func", "getHealthCheckStatistics")
	statsInput = getHeathCheckMetricStatsOptions(list, options)

	log.Debug("starting ...")
	log.With("period", *statsInput.Period).Debug("getting metrics statistics ...")

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
// If time period is more than 15 days ago, use hourly.
//
// Don't use any granularity under 60 seconds.
func getPeriod(start time.Time) (period int32) {
	var (
		now       time.Time     = time.Now().UTC()
		day       time.Duration = (time.Hour * 24)
		days15    time.Duration = (15 * day)
		hoursDiff float64       = now.Sub(start).Hours()
	)
	period = 3600

	if hoursDiff < days15.Hours() {
		period = 60
	}
	return

}
