package uptime

import (
	"fmt"
	"opg-reports/shared/env"
	"testing"
	"time"
)

// TestSharedAWSUptimeRealData is used to test and work on data
// locally using prod values
func TestSharedAWSUptimeRealData(t *testing.T) {
	if env.Get("AWS_SESSION_TOKEN", "") != "" {
		cw, _ := ClientFromEnv("us-east-1")
		metrics := GetListOfMetrics(cw)
		fmt.Printf("found:\n%+v\n", metrics)

		now := time.Now().UTC().AddDate(0, 0, -2)
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		datapoints, _ := GetMetricsStats(cw, metrics, start, end)

		fmt.Printf("found [%v] datapoints", len(datapoints))
	}

}
