package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/aws_uptime/awsu"
	"github.com/ministryofjustice/opg-reports/shared/aws"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const region string = "us-east-1"

func main() {
	logger.LogSetup()

	var (
		now       = time.Now().UTC()
		yesterday = dates.ResetDay(now).AddDate(0, 0, -1)
		group     = flag.NewFlagSet("aws_uptime", flag.ExitOnError)
		day       = argument.NewDate(group, "day", yesterday, dates.FormatYMD, "Day (YYYY-MM-DD) to fetch uptime data for")
		unit      = argument.New(group, "unit", "", "Unit to group the account into (like a team structure)")
		out       = argument.New(group, "output", "aws_uptime.json", "filename suffix")
		dir       = argument.New(group, "dir", "data", "sub dir")
	)

	group.Parse(os.Args[1:])

	if *unit.Value == "" {
		slog.Error("required aruments are missing")
		os.Exit(1)
	}
	if day.Value.Format(dates.FormatY) == dates.ErrorYear {
		slog.Error("required aruments are missing")
		os.Exit(1)
	}
	// reset the inputted value to the start of a day
	startDate := dates.ResetDay(*day.Value)
	// end date is the start of the next day, its a less than only
	endDate := startDate.AddDate(0, 0, 1)

	slog.Info("getting costs",
		slog.String("day", day.String()),
		slog.String("start", startDate.Format(dates.FormatYMD)),
		slog.String("end", endDate.Format(dates.FormatYMD)),
		slog.String("unit", *unit.Value),
		slog.String("dir", *dir.Value),
		slog.String("out", *out.Value))

	cw, err := aws.CWClientFromEnv("us-east-1")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	metrics := aws.GetListOfMetrics(cw)

	slog.Info("found metrics", slog.Int("count", len(metrics)))

	ups := []*awsu.AwsUptime{}
	// go get the percentage values per minute for the data of the day
	if len(metrics) > 0 {
		datapoints, _ := aws.GetMetricsStats(cw, metrics, startDate, endDate)
		// convert each data point to an entry
		for _, dp := range datapoints {
			up := &awsu.AwsUptime{
				Ts:      time.Now().UTC().Format(dates.Format),
				Unit:    *unit.Value,
				Average: *dp.Average,
				Date:    dp.Timestamp.Format(dates.Format),
			}
			ups = append(ups, up)
		}
	}
	slog.Info("found datapoints", slog.Int("count", len(ups)))

	content, err := convert.Marshals(ups)
	if err != nil {
		slog.Error("error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}

	// mkdir and filename
	os.MkdirAll(*dir.Value, os.ModePerm)
	filename := filepath.Join(*dir.Value, fmt.Sprintf("%s_%s_%s", *unit.Value, startDate.Format(dates.FormatYMD), *out.Value))

	os.WriteFile(filename, content, os.ModePerm)
}
