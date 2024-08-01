package main

import (
	"fmt"
	"opg-reports/shared/logger"
	"opg-reports/shared/report"
)

var (
	day          = report.NewDayArg("day", true, "Day (YYYY-MM-DD) to get uptime data for", "-")
	account_unit = report.NewArg("account_unit", true, "Unit to group the account into (like a team structure)", "")
)

const dir string = "data"

func run(r report.IReport) {
	unit := account_unit.Val()
	start, _ := day.DayValue()
	end := start.AddDate(0, 0, 1)

	fmt.Println(unit)
	fmt.Println(start)
	fmt.Println(end)

}

func main() {
	logger.LogSetup()

	costReport := report.New(day, account_unit)
	costReport.SetRunner(run)
	costReport.Run()

}
