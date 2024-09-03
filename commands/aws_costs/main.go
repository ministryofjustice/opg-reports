package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/shared/aws"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func Flat(raw *costexplorer.GetCostAndUsageOutput,
	accountId string,
	accountOrg string,
	accountUnit string,
	accountName string,
	accountLabel string,
	accountEnv string,
) (costs []awsc.AwsCost, err error) {
	slog.Debug("Flattening cost data")
	now := time.Now().UTC().Format(dates.Format)

	costs = []awsc.AwsCost{}

	for _, resultByTime := range raw.ResultsByTime {
		day := *resultByTime.TimePeriod.Start

		for _, costGroup := range resultByTime.Groups {
			service := *costGroup.Keys[0]
			region := *costGroup.Keys[1]

			for _, costMetric := range costGroup.Metrics {
				amount := *costMetric.Amount

				c := awsc.AwsCost{
					Service:      service,
					Region:       region,
					Date:         day,
					Cost:         amount,
					AccountID:    accountId,
					Organisation: accountOrg,
					Unit:         accountUnit,
					AccountName:  accountName,
					Label:        accountLabel,
					Environment:  accountEnv,
					Ts:           now,
				}
				costs = append(costs, c)
			}
		}
	}

	return
}

func main() {
	logger.LogSetup()

	var (
		now       = time.Now().UTC()
		lastMonth = dates.ResetMonth(now).AddDate(0, -1, 0)
		group     = flag.NewFlagSet("aws_costs", flag.ExitOnError)
		month     = argument.NewDate(group, "month", lastMonth, dates.FormatYM, "Month (YYYY-MM) to fetch cost data for")
		id        = argument.New(group, "id", "", "AWS Account Id")
		name      = argument.New(group, "name", "", "Friendly name for the AWS Account")
		label     = argument.New(group, "label", "", "Label for the account")
		unit      = argument.New(group, "unit", "", "Unit to group the account into (like a team structure)")
		org       = argument.New(group, "organisation", "OPG", "Organisation name")
		env       = argument.New(group, "environment", "production", "Evironment type")
		out       = argument.New(group, "output", "aws_costs.json", "filename suffix")
		dir       = argument.New(group, "dir", "data", "sub dir")
	)

	group.Parse(os.Args[1:])

	// if any arg is empty, fail
	if *id.Value == "" || *name.Value == "" || *label.Value == "" || *unit.Value == "" || *org.Value == "" || *env.Value == "" {
		slog.Error("required aruments are missing")
		os.Exit(1)
	}
	if month.Value.Format(dates.FormatY) == dates.ErrorYear {
		slog.Error("required aruments are missing")
		os.Exit(1)
	}

	startDate := *month.Value
	endDate := startDate.AddDate(0, 1, 0)

	slog.Info("getting costs",
		slog.String("month", month.Value.Format(dates.FormatYM)),
		slog.String("start", startDate.Format(dates.FormatYMD)),
		slog.String("end", endDate.Format(dates.FormatYMD)),
		slog.String("id", *id.Value),
		slog.String("name", *name.Value),
		slog.String("label", *label.Value),
		slog.String("unit", *unit.Value),
		slog.String("org", *org.Value),
		slog.String("environment", *env.Value),
		slog.String("dir", *dir.Value),
		slog.String("out", *out.Value))

	raw, err := aws.CostAndUsage(startDate, endDate, costexplorer.GranularityDaily, dates.FormatYMD)
	costs, err := Flat(raw, *id.Value, *org.Value, *unit.Value, *name.Value, *label.Value, *env.Value)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	content, err := convert.Marshals(costs)
	if err != nil {
		slog.Error("error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}

	os.MkdirAll(*dir.Value, os.ModePerm)
	filename := filepath.Join(*dir.Value, fmt.Sprintf("%s_%s_%s", *id.Value, startDate.Format(dates.FormatYM), *out.Value))

	os.WriteFile(filename, content, os.ModePerm)

}
