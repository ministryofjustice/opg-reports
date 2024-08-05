package uptime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const uptimeNamespace string = "AWS/Route53" // Route53 health checks
const uptimeMetric string = "HealthCheckPercentageHealthy"
const uptimePeriod int64 = 60 // 1 minute intervals
const uptimeUnit string = "Percent"
const uptimeStat string = "Average"

// GetListOfMetrics returns all metrics that are tracking uptime percentages within this account
// which can then be used to get the statistics to show uptime amount
func GetListOfMetrics(cw *cloudwatch.CloudWatch) (metrics []*cloudwatch.Metric) {
	slog.Debug("getting uptime metrics list")
	metrics = []*cloudwatch.Metric{}

	in := &cloudwatch.ListMetricsInput{
		Namespace:  aws.String(uptimeNamespace),
		MetricName: aws.String(uptimeMetric),
	}
	page := 0
	cw.ListMetricsPages(in, func(result *cloudwatch.ListMetricsOutput, b bool) bool {
		page += 1
		metrics = append(metrics, result.Metrics...)

		slog.Debug("got list of metrics",
			slog.Int("count", len(result.Metrics)),
			slog.Int("page", page),
			slog.Bool("more pages?", (result.NextToken != nil)))

		// stops looping when we nolonger have next page tokens
		return (result.NextToken != nil)
	})
	slog.Debug("got uptime metrics list", slog.Int("count", len(metrics)))
	return
}

// GetMetricsStats returns all data points for the health check metrics - so allows us to track
// the % uptime of the service
//
// As the namespace and metrics are identical, we can do one call to hte api by merging
// the dimension values together
func GetMetricsStats(cw *cloudwatch.CloudWatch,
	metrics []*cloudwatch.Metric,
	start time.Time,
	end time.Time,
) (datapoints []*cloudwatch.Datapoint, err error) {
	period := uptimePeriod
	unit := uptimeUnit
	stat := uptimeStat
	stats := []*string{&stat}

	slog.Debug("getting metric stats",
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

	slog.Debug("metric stats input", slog.String("input", fmt.Sprintf("%+v", in)))

	results, err := cw.GetMetricStatistics(in)
	datapoints = results.Datapoints

	slog.Debug("metric found count", slog.Int("count", len(datapoints)))

	return

}
