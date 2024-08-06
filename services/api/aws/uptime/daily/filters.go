package daily

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"opg-reports/shared/server/resp"
)

func FilterFunctions(parameters map[string][]string, response *resp.Response) (funcs map[string]data.IStoreFilterFunc[*uptime.Uptime]) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := server.GetStartEndDates(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/uptime/daily] FilterFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

	response.Metadata["startDate"] = startDate
	response.Metadata["endDate"] = endDate
	response.Metadata["filters"] = parameters

	// -- month range filters
	var inMonth = func(item *uptime.Uptime) bool {
		dateStr := item.DateTime.Format(dates.FormatYM)
		return dates.InMonth(dateStr, months)
	}

	funcs = map[string]data.IStoreFilterFunc[*uptime.Uptime]{
		"inMonth": inMonth,
	}

	return
}
