package main

import (
	"log/slog"
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
	"opg-reports/shared/report"
)

var (
	day = report.NewDayArg("day", true, "Day (YYYY-MM-DD) to get uptime data for", "-")
	// As this is only for production, we dont need anything more than the unit
	account_unit = report.NewArg("account_unit", true, "Unit to group the account into (like a team structure)", "")
)

const dir string = "data"

func run(r report.IReport) {
	unit := account_unit.Val()
	start, _ := day.DayValue()
	end := start.AddDate(0, 0, 1)

	slog.Info("getting uptime data",
		slog.String("account unit", unit),
		slog.String("start", start.String()),
		slog.String("end", start.String()))

	//  have to fix the region to us-east-1 to return the correct metrics
	cw, _ := uptime.ClientFromEnv("us-east-1")
	// go fetch the metrics that are setup for uptime
	metrics := uptime.GetListOfMetrics(cw)

	slog.Info("found metrics", slog.Int("count", len(metrics)))
	ups := []*uptime.Uptime{}
	// go get the percentage values per minute for the data of the day
	if len(metrics) > 0 {
		datapoints, _ := uptime.GetMetricsStats(cw, metrics, start, end)
		// convert each data point to an entry
		for _, dp := range datapoints {
			up := uptime.NewFromDatapoint(nil, dp)
			// attach the unit
			up.AccountUnit = unit
			ups = append(ups, up)
		}
	}
	slog.Info("found datapoints", slog.Int("count", len(ups)))
	// write to file
	content, _ := data.ToJsonList[*uptime.Uptime](ups)
	filename := r.Filename()
	files.WriteFile(dir, filename, content)

}

func main() {
	logger.LogSetup()

	costReport := report.New(day, account_unit)
	costReport.SetRunner(run)
	costReport.Run()

}
