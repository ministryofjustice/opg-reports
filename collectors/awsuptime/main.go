/*
awsuptime fetches aws uptime data for the day.

Usage:

	awsuptime [flags]

The flags are:

	-day=<yyyy-mm-dd>
		The month (formated as YYYY-MM-DD) to fetch data for.
		If set to "-", uses the current month.
		Defaults to the current month.
	-unit=<unit>
		Team name for who owns the account.
	-output=<path-pattern>
		Path (with magic values) to the output file
		Default: `./data/{day}_{unit}_aws_uptime.json`

The command presumes an active, autherised session that can connect
to AWS r53 health checks. These are dynamically
fetched from environment variables.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/ministryofjustice/opg-reports/collectors/awsuptime/lib"
	"github.com/ministryofjustice/opg-reports/internal/awsclient"
	"github.com/ministryofjustice/opg-reports/internal/awssession"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/models"
)

var (
	args   = &lib.Arguments{}
	region = "us-east-1"
)

func Run(args *lib.Arguments) (err error) {
	var (
		s          *session.Session
		client     *cloudwatch.CloudWatch
		startDate  time.Time
		endDate    time.Time
		content    []byte
		metrics    []*cloudwatch.Metric
		datapoints []*cloudwatch.Datapoint
		uptimeData []*models.AwsUptime
	)

	if s, err = awssession.NewWithRegion(region); err != nil {
		slog.Error("[awsuptime] aws session failed", slog.String("err", err.Error()))
		return
	}

	if client, err = awsclient.CloudWatch(s); err != nil {
		slog.Error("[awsuptime] aws client failed", slog.String("err", err.Error()))
		return
	}

	if startDate, err = dateutils.Time(args.Day); err != nil {
		slog.Error("[awsuptime] date conversion failed", slog.String("err", err.Error()))
		return
	}

	startDate = dateutils.Reset(startDate, dateintervals.Day)

	// overwrite month with the parsed version
	args.Day = startDate.Format(dateformats.YMD)
	endDate = startDate.AddDate(0, 0, 1)

	now := time.Now().UTC().Format(dateformats.Full)
	// get all of the named metrics for this setup
	if metrics = lib.GetListOfMetrics(client); len(metrics) > 0 {
		datapoints, err = lib.GetMetricsStats(client, metrics, startDate, endDate)
		if err != nil {
			slog.Error("[awsuptime] getting stats failed", slog.String("err", err.Error()))
			return
		}

		unit := &models.Unit{
			Name: args.Unit,
		}
		account := &models.AwsAccount{
			Ts:     now,
			Number: args.AccountID,
		}

		for _, dp := range datapoints {
			ts := time.Now().UTC().Format(dateformats.Full)

			up := &models.AwsUptime{
				Ts:         ts,
				Date:       dp.Timestamp.Format(dateformats.Full),
				Average:    *dp.Average,
				Unit:       (*models.UnitForeignKey)(unit),
				AwsAccount: (*models.AwsAccountForeignKey)(account),
			}

			uptimeData = append(uptimeData, up)
		}
	}

	content, err = json.MarshalIndent(uptimeData, "", "  ")
	if err != nil {
		slog.Error("[awsuptime] error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}
	//
	lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[awsuptime] starting...")
	slog.Debug("[awsuptime]", slog.String("args", fmt.Sprintf("%+v", args)))
	slog.Debug("[awsuptime]", slog.String("region", region))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[awsuptime] done.")

}
