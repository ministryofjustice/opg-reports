package importer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/reset"
	"opg-reports/report/packages/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const InsertStatement string = `
INSERT INTO uptime (
	month,
	average,
	granularity,
	account_id
) VALUES (
	:month,
	:average,
	:granularity,
	:account_id
) ON CONFLICT (account_id,month)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id
;
`

// these are fixed values used for the api calls
const (
	metricRegion     string             = "us-east-1"
	metricsNamespace string             = "AWS/Route53"
	metricsName      string             = "HealthCheckPercentageHealthy"
	metricsStatistic types.Statistic    = types.StatisticAverage
	metricsUnit      types.StandardUnit = types.StandardUnitPercent
)

var (
	ErrIncorrectRegion          = errors.New("metrics must be fetched via us-east-1.")
	ErrFailedGettingMetricsList = errors.New("failed to get metrics list with error.")
	ErrFailedGettingMetricStats = errors.New("failed to get metric statistics with error.")
)

// Model represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type Model struct {
	Month       string `json:"month,omitempty"`
	Average     string `json:"average,omitempty"`
	Granularity string `json:"granularity,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
}

// Map returns a map of all fields on this struct
func (self *Model) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// Get returns the uptime metric stats from the api
func Get(ctx context.Context, client *cloudwatch.Client, opts *args.Import, previous ...types.Datapoint) (found []types.Datapoint, err error) {
	var (
		log       *slog.Logger
		list      *cloudwatch.ListMetricsOutput
		stats     *cloudwatch.GetMetricStatisticsOutput
		start     time.Time = reset.Month(&opts.Filters.Dates.Start)
		end       time.Time = reset.Day(&opts.Filters.Dates.End)
		setRegion string    = client.Options().Region
	)
	ctx, log = logger.Get(ctx)
	log.Info("getting uptime data from aws for timer period...", "start", start, "end", end)
	// check region
	if setRegion != metricRegion {
		err = ErrIncorrectRegion
		log.Error("incorrect region used in client - requires us-east-1", "region", setRegion)
		return
	}
	// fetch the list of all metrics from the api
	log.Info("getting list of health check metrics ...")
	list, err = getHealthCheckMetrics(ctx, client)
	if err != nil {
		return
	}
	// get all the datapoints for each of the metrics
	log.Debug("getting health check metrics stats ...")
	stats, err = getHealthCheckStatistics(ctx, client, list, opts.Filters.Dates)
	if err != nil {
		return
	}
	found = stats.Datapoints
	return

}

// Filter - not filtering on uptime
func Filter(ctx context.Context, items []types.Datapoint, filters *args.Filters) (included []types.Datapoint) {
	included = items
	return
}

// Transform converts the original data into record for local database insertion
func Transform(ctx context.Context, data []types.Datapoint, opts *args.Import) (results []*Model, err error) {
	var (
		log     *slog.Logger
		grouped map[string]float64 = map[string]float64{}
		counter map[string]int     = map[string]int{}
	)
	ctx, log = logger.Get(ctx)
	results = []*Model{}
	log.Info("transforming uptime to local models ...", "count", len(data))

	// create a sum and count of each month uptime to then create the average entries
	for _, point := range data {
		var (
			month time.Time = reset.Day(point.Timestamp)
			key   string    = month.Format(times.YM)
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
			count   int     = counter[key]
			average float64 = (sum / float64(count))
			up      *Model  = &Model{
				Month:     key,
				Average:   fmt.Sprintf("%g", average),
				AccountID: opts.Aws.AccountID,
			}
		)
		results = append(results, up)
	}

	log.Info("uptime transformation completed.", "count", len(results))
	return
}

// getHealthCheckStatistics uses the list of metrics from getHealthCheckMetrics to fetch the uptime datapoints
// from the api and return the api call values
//
// T is *cloudwatch.Client
func getHealthCheckStatistics(ctx context.Context, client *cloudwatch.Client, list *cloudwatch.ListMetricsOutput, dates *args.Dates) (stats *cloudwatch.GetMetricStatisticsOutput, err error) {
	var statsInput *cloudwatch.GetMetricStatisticsInput
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	statsInput = getHeathCheckMetricStatsOptions(list, dates)
	log.Debug("getting health check metrics stats ...")
	// try and get the stats
	stats, err = client.GetMetricStatistics(ctx, statsInput)
	if err != nil {
		log.Error("error getting metric statistics.", "err", err.Error())
		err = errors.Join(ErrFailedGettingMetricStats, err)
		return
	}
	log.With("count", len(stats.Datapoints)).Debug("complete.")
	return
}

// getHealthCheckMetrics returns the list of metrics to use for uptime data. Limit to recently active metrics so
// we aren't picking up aged / dead health checks.
//
// Client is *cloudwatch.Client
func getHealthCheckMetrics(ctx context.Context, client *cloudwatch.Client) (list *cloudwatch.ListMetricsOutput, err error) {
	var log *slog.Logger
	var listOptions *cloudwatch.ListMetricsInput = &cloudwatch.ListMetricsInput{
		Namespace:      convert.Ptr(metricsNamespace),
		MetricName:     convert.Ptr(metricsName),
		RecentlyActive: types.RecentlyActivePt3h,
	}

	ctx, log = logger.Get(ctx)

	log.Debug("fetching metric data for account ...")
	list, err = client.ListMetrics(ctx, listOptions)
	if err != nil {
		log.Error("error getting list of metrics", "err", err.Error())
		err = errors.Join(ErrFailedGettingMetricsList, err)
		return
	}

	log.Debug("found metrics data.", "count", len(list.Metrics))
	return
}

// getHeathCheckMetricStatsOptions is used to generate a suitable GetMetricStatisticsInput struct
// that contains all of the metric dimensions and uses the correct period and units that we need
// to fetch uptime data
func getHeathCheckMetricStatsOptions(list *cloudwatch.ListMetricsOutput, options *args.Dates) (opts *cloudwatch.GetMetricStatisticsInput) {
	var (
		period     int32             = 3600
		dimensions []types.Dimension = []types.Dimension{}
	)

	// merge all of the metric dimensions together
	for _, metric := range list.Metrics {
		dimensions = append(dimensions, metric.Dimensions...)
	}
	// generate the input struct
	opts = &cloudwatch.GetMetricStatisticsInput{
		Namespace:  convert.Ptr(metricsNamespace),
		MetricName: convert.Ptr(metricsName),
		StartTime:  convert.Ptr(options.Start),
		EndTime:    convert.Ptr(options.End),
		Period:     convert.Ptr(period),
		Unit:       metricsUnit,
		Statistics: []types.Statistic{metricsStatistic},
		Dimensions: dimensions,
	}

	return
}
