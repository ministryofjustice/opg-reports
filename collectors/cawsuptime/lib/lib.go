package lib

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
)

const (
	uptimeNamespace string = "AWS/Route53" // Route53 health checks
	uptimePeriod    int64  = 60            // 1 minute intervals
	uptimeMetric    string = "HealthCheckPercentageHealthy"
	uptimeUnit      string = "Percent"
	uptimeStat      string = "Average"
)

var (
	defDay = convert.DateResetDay(time.Now().UTC()).AddDate(0, 0, -1)
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	Day        string
	Unit       string
	AccountID  string
	OutputFile string
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Day, "day", defDay.Format(consts.DateFormatYearMonthDay), "day to fetch data for.")
	flag.StringVar(&args.Unit, "unit", "", "Unit / team name.")
	flag.StringVar(&args.AccountID, "id", "", "AWS account id")
	flag.StringVar(&args.OutputFile, "output", "./data/{day}_{unit}_aws_uptime.json", "Filepath for the output")

	flag.Parse()
}

// ValidateArgs checks rules and logic for the input arguments
// Make sure some have non empty values and apply default values to others
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"unit":   args.Unit,
		"id":     args.AccountID,
		"output": args.OutputFile,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if args.Day == "-" {
		args.Day = defDay.Format(consts.DateFormat)
	}

	return
}

// WriteToFile writes the content to the file replacing values in
// the filename with values on arg
func WriteToFile(content []byte, args *Arguments) {
	var (
		filename string
		dir      string = filepath.Dir(args.OutputFile)
	)
	os.MkdirAll(dir, os.ModePerm)
	filename = args.OutputFile
	filename = strings.ReplaceAll(filename, "{day}", args.Day)
	filename = strings.ReplaceAll(filename, "{unit}", strings.ToLower(args.Unit))

	os.WriteFile(filename, content, os.ModePerm)

}

// GetListOfMetrics returns all metrics that are tracking uptime percentages within this account
// which can then be used to get the statistics to show uptime amount
func GetListOfMetrics(cw *cloudwatch.CloudWatch) (metrics []*cloudwatch.Metric) {
	slog.Info("getting uptime metrics list")
	metrics = []*cloudwatch.Metric{}

	in := &cloudwatch.ListMetricsInput{
		Namespace:  aws.String(uptimeNamespace),
		MetricName: aws.String(uptimeMetric),
	}
	page := 0
	cw.ListMetricsPages(in, func(result *cloudwatch.ListMetricsOutput, b bool) bool {
		page += 1
		metrics = append(metrics, result.Metrics...)

		slog.Info("got list of metrics",
			slog.Int("count", len(result.Metrics)),
			slog.Int("page", page),
			slog.Bool("more pages?", (result.NextToken != nil)))

		// stops looping when we nolonger have next page tokens
		return (result.NextToken != nil)
	})
	slog.Info("got uptime metrics list", slog.Int("count", len(metrics)))
	return
}

// GetMetricsStats returns all data points for the health check metrics - so allows us to track
// the % uptime of the service
//
// As the namespace and metrics are identical, we can do one call to hte api by merging
// the dimension values together
func GetMetricsStats(cw *cloudwatch.CloudWatch, metrics []*cloudwatch.Metric, start time.Time, end time.Time) (datapoints []*cloudwatch.Datapoint, err error) {
	var (
		period = uptimePeriod
		unit   = uptimeUnit
		stat   = uptimeStat
		stats  = []*string{&stat}
	)
	var results *cloudwatch.GetMetricStatisticsOutput

	slog.Info("getting metric stats",
		slog.String("namespace", uptimeNamespace),
		slog.String("metric", uptimeMetric))

	// merge all dimensions together for one call
	dimensions := []*cloudwatch.Dimension{}
	for _, metric := range metrics {
		dimensions = append(dimensions, metric.Dimensions...)
	}

	in := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(uptimeNamespace),
		MetricName: aws.String(uptimeMetric),
		StartTime:  &start,
		EndTime:    &end,
		Period:     &period,
		Unit:       &unit,
		Statistics: stats,
		Dimensions: dimensions,
	}

	results, err = cw.GetMetricStatistics(in)
	datapoints = results.Datapoints

	slog.Info("metric found count", slog.Int("count", len(datapoints)))

	return

}
