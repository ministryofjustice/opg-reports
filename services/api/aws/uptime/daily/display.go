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
	// -- monthly
	var monthly endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/uptime/daily] monthly head func")
		r = row.New()
		r.Add(cell.New("", "", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Overall %", "Overall %", false, true))
		return
	}
	// -- monthly by account
	var monthlyByAccountUnit endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/uptime/daily] monthly by account unit head func")
		r = row.New()
		r.Add(cell.New("Unit", "Unit", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Overall %", "Overall %", false, true))
		return
	}

	funcs = map[string]endpoint.DisplayHeadFunc{
		"monthly":              monthly,
		"monthlyByAccountUnit": monthlyByAccountUnit,
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
		cells := []*cell.Cell{
			cell.New("Average", "Average", true, false),
		}
		// get the row months
		rowAvg, monthCells := AvgPerMonth(store, months)
		cells = append(cells, monthCells...)
		// totals
		cells = append(cells, cell.New("Overall %", rowAvg, false, true))
		// return the row
		rows = []*row.Row{row.New(cells...)}
		return
	}

	// -- monthly
	var monthlyByAccountUnit endpoint.DisplayRowFunc[*uptime.Uptime] = func(group string, store data.IStore[*uptime.Uptime], resp *resp.Response) (rows []*row.Row) {
		if store.Length() > 0 {
			first := store.List()[0]
			// row headers
			cells := []*cell.Cell{
				cell.New(first.AccountUnit, first.AccountUnit, true, false),
			}
			// get the row months
			rowAvg, monthCells := AvgPerMonth(store, months)
			cells = append(cells, monthCells...)
			// totals
			cells = append(cells, cell.New("Overall %", rowAvg, false, true))
			// return the row
			rows = []*row.Row{row.New(cells...)}
		}
		return
	}

	funcs = map[string]endpoint.DisplayRowFunc[*uptime.Uptime]{
		"monthly":              monthly,
		"monthlyByAccountUnit": monthlyByAccountUnit,
	}

	return
}
