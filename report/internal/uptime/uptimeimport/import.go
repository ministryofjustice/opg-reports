package uptimeimport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/ptr"
	"opg-reports/report/package/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	_ "github.com/mattn/go-sqlite3"
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
	Granularity string `json:"granularity,omityempty"`
	AccountID   string `json:"account_id,omityempty"`
}

// Client is used to allow mocking and is a proxy for *cloudwatch.Client
// and the methods the function calls
type Client interface {
	ListMetrics(ctx context.Context, params *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error)
	GetMetricStatistics(ctx context.Context, params *cloudwatch.GetMetricStatisticsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricStatisticsOutput, error)
	Options() cloudwatch.Options
}

type Args struct {
	DB            string `json:"db"`             // database path
	Driver        string `json:"driver"`         // database driver
	Params        string `json:"params"`         // database connection params
	MigrationFile string `json:"migration_file"` // database migrations

	DateStart time.Time `json:"date_start"` // start date, this will be reset to start of the month (and expanded to capture historical data)
	DateEnd   time.Time `json:"date_end"`   // end date
	AccountID string    `json:"account_id"` // AccountID provided by awsid.AccountID
}

func Import(ctx context.Context, client Client, in *Args) (err error) {
	var (
		list      *cloudwatch.ListMetricsOutput
		statsOpts *cloudwatch.GetMetricStatisticsInput
		stats     *cloudwatch.GetMetricStatisticsOutput
		data      []*Model
		setRegion string       = client.Options().Region
		log       *slog.Logger = cntxt.GetLogger(ctx).With("package", "uptimeimport", "func", "Import")
	)

	log.With("options", in).Info("starting ...")
	// check region
	if setRegion != metricRegion {
		err = ErrIncorrectRegion
		log.Error("incorrect region used in client - requires us-east-1", "region", setRegion)
		return
	}

	// fetch the list of all metrics from the api
	log.Debug("getting lisst of metrics ...")
	list, err = getHealthCheckMetrics(ctx, client)
	if err != nil {
		return
	}

	// get all the datapoints for each of the metrics
	log.Debug("getting metric stats ...")
	stats, statsOpts, err = getHealthCheckStatistics(ctx, client, list, in)
	if err != nil {
		return
	}

	log.Debug("coverting to models ...")
	data, err = toModels(ctx, in.AccountID, *statsOpts.Period, stats)
	if err != nil {
		return
	}

	// now write to db
	err = dbx.Insert(ctx, InsertStatement, data, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.With("count", len(data)).Info("complete.")
	return

}

// toModels converts the raw data into a list of models ready to write to the database
func toModels(ctx context.Context, account string, period int32, result *cloudwatch.GetMetricStatisticsOutput) (data []*Model, err error) {
	var (
		log     *slog.Logger       = cntxt.GetLogger(ctx).With("package", "uptimeimport", "func", "toModels")
		grouped map[string]float64 = map[string]float64{}
		counter map[string]int     = map[string]int{}
	)
	data = []*Model{}
	log.Debug("starting ... ")

	// create a sum and count of each month uptime to then create the average entries
	for _, point := range result.Datapoints {
		var (
			month time.Time = times.ResetDay(*point.Timestamp)
			key   string    = times.AsYMString(month)
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
				Month:       key,
				Average:     fmt.Sprintf("%g", average),
				Granularity: fmt.Sprintf("%d", period),
				AccountID:   account,
			}
		)
		data = append(data, up)
	}
	log.With("count", len(data)).Debug("complete.")
	return
}

// getHealthCheckStatistics uses the list of metrics from getHealthCheckMetrics to fetch the uptime datapoints
// from the api and return the api call values
//
// T is *cloudwatch.Client
func getHealthCheckStatistics[T Client](ctx context.Context, client T, list *cloudwatch.ListMetricsOutput, options *Args) (stats *cloudwatch.GetMetricStatisticsOutput, statsInput *cloudwatch.GetMetricStatisticsInput, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "uptimeimport", "func", "getHealthCheckStatistics")

	statsInput = getHeathCheckMetricStatsOptions(list, options)
	log.Debug("starting ...")
	log.With("period", *statsInput.Period).Debug("getting metrics statistics ...")
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

// getHealthCheckMetrics returns the list of metrics to use for uptime data
//
// T is *cloudwatch.Client
func getHealthCheckMetrics[T Client](ctx context.Context, client T) (list *cloudwatch.ListMetricsOutput, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "uptimeimport", "func", "getHealthCheckMetrics")
	var listOptions *cloudwatch.ListMetricsInput = &cloudwatch.ListMetricsInput{
		Namespace:  ptr.Ptr(metricsNamespace),
		MetricName: ptr.Ptr(metricsName),
	}

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
func getHeathCheckMetricStatsOptions(list *cloudwatch.ListMetricsOutput, options *Args) (opts *cloudwatch.GetMetricStatisticsInput) {
	var (
		period     int32             = getPeriod(options.DateStart)
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
		StartTime:  ptr.Ptr(options.DateStart),
		EndTime:    ptr.Ptr(options.DateEnd),
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
