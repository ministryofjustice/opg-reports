package main

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/aws/arn"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
	"opg-reports/shared/report"
	"time"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var (
	month         = report.NewMonthArg("month", true, "Month (YYYY-MM) to fetch cost data for", "")
	account_id    = report.NewArg("account_id", true, "AWS Account Id", "")
	account_name  = report.NewArg("account_name", true, "Friendly name for the AWS Account", "")
	account_label = report.NewArg("account_label", true, "Label for the account", "")
	account_unit  = report.NewArg("account_unit", true, "Unit to group the account into (like a team structure)", "")
	account_org   = report.NewArg("account_organisation", true, "Organisation name", "OPG")
	account_env   = report.NewArg("account_environment", true, "Account environment type", "development")
	role          = report.NewArg("role", false, "Role to use to fetch the data", "billing")
	region        = report.NewArg("region", false, "Region to start session from", "eu-west-1")
)

const dir string = "data"

func run(r report.IReport) {
	var m time.Time
	if month.Val() == "-" {
		m = time.Now().UTC().AddDate(0, -1, 0)
	} else {
		m, _ = dates.StringToDate(month.Val())
	}
	id := account_id.Val()
	name := account_name.Val()
	label := account_label.Val()
	unit := account_unit.Val()
	org := account_org.Val()
	env := account_env.Val()
	roleName := role.Val()
	reg := region.Val()

	slog.Info("getting costs",
		slog.String("account name", name),
		slog.String("month", m.Format(dates.FormatYM)),
	)

	startDate := m
	endDate := m.AddDate(0, 1, 0)
	roleArn := arn.RoleArn(id, roleName)
	raw, err := cost.CostAndUsage(roleArn, startDate, endDate, reg, costexplorer.GranularityDaily, dates.FormatYMD)
	if err != nil {
		slog.Error(fmt.Sprintf("error: %v", err.Error()))
		panic(err.Error())
	}

	costs, err := cost.Flatten(raw, id, name, label, unit, org, env)
	content, err := data.ToJsonList[*cost.Cost](costs)
	filename := r.Filename()
	files.WriteFile(dir, filename, content)

}

func main() {
	logger.LogSetup()
	// 1 month ago
	now := time.Now().UTC().AddDate(0, -1, 0)
	month.SetDefault(now.Format(dates.FormatYM))

	costReport := report.New(month, account_id, account_name, account_label, account_unit, account_org, account_env, role, region)
	costReport.SetRunner(run)
	costReport.Run()

}
