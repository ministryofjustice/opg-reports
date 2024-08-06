package daily

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/cell"
	"opg-reports/shared/server/resp/row"
)

// DisplayHeadFunctions
func DisplayHeadFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayHeadFunc) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := server.GetStartEndDates(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/uptime/daily] DisplayHeadFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

	hMonths := HeaderMonths(months)
	// -- monthly totals
	// This has a split on in tax is included or not and then
	// a final col for the line total
	var monthly endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/uptime/daily] monthly head func")
		r = row.New()
		r.Add(hMonths...)
		r.Add(cell.New("Average %", "Average %", false, true))
		return
	}

	funcs = map[string]endpoint.DisplayHeadFunc{
		"monthly": monthly,
	}

	return funcs
}

// DisplayRowFunctions
func DisplayRowFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayRowFunc[*uptime.Uptime]) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := server.GetStartEndDates(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/uptime/daily] DisplayHeadFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

	// -- monthly
	var monthly endpoint.DisplayRowFunc[*uptime.Uptime] = func(group string, store data.IStore[*uptime.Uptime], resp *resp.Response) (rows []*row.Row) {
		// row headers
		cells := []*cell.Cell{}
		// get the row months
		rowAvg, monthCells := AvgPerMonth(store, months)
		cells = append(cells, monthCells...)
		// totals
		cells = append(cells, cell.New("Average %", rowAvg, false, true))
		// return the row
		rows = []*row.Row{row.New(cells...)}
		return
	}

	funcs = map[string]endpoint.DisplayRowFunc[*uptime.Uptime]{
		"monthly": monthly,
	}

	return
}

// // DisplayFootFunctions
// func DisplayFootFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayFootFunc) {
// 	// -- get the start & end dates as well as list of all months
// 	startDate, endDate := server.GetStartEndDates(parameters)
// 	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
// 	slog.Debug("[aws/uptime/daily] DisplayFootFunctions",
// 		slog.String("start", startDate.String()),
// 		slog.String("end", endDate.String()),
// 		slog.String("months", fmt.Sprintf("%+v", months)),
// 		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

// 	var perUnit endpoint.DisplayFootFunc = columnTotals
// 	var perUnitEnv endpoint.DisplayFootFunc = columnTotals
// 	var perUnitEnvService endpoint.DisplayFootFunc = columnTotals

// 	funcs = map[string]endpoint.DisplayFootFunc{
// 		"perUnit":           perUnit,
// 		"perUnitEnv":        perUnitEnv,
// 		"perUnitEnvService": perUnitEnvService,
// 	}

// 	return
// }
