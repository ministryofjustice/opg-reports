package monthly

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/cell"
	"opg-reports/shared/server/resp/row"
)

// DisplayHeadFunctions
func DisplayHeadFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayHeadFunc) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := startEnd(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/costs/monthly] DisplayHeadFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))
	// -- monthly totals
	// This has a split on in tax is included or not and then
	// a final col for the line total
	var monthlyTotals endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/costs/monthly] monthlyTotals head func")
		r = row.New()
		r.Add(cell.New("Tax", "Tax", true, false))
		for _, m := range months {
			r.Add(cell.New(m, m, true, false))
		}
		r.Add(cell.New("Total", "Total", true, true))
		return
	}

	funcs = map[string]endpoint.DisplayHeadFunc{
		"monthlyTotals": monthlyTotals,
	}

	return funcs
}

// DisplayRowFunctions
func DisplayRowFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayRowFunc[*cost.Cost]) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := startEnd(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/costs/monthly] DisplayHeadFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))
	// -- monthly totals
	var monthlyTotals endpoint.DisplayRowFunc[*cost.Cost] = func(group string, store data.IStore[*cost.Cost], resp *resp.Response) (rows []*row.Row) {
		rows = []*row.Row{}
		// split by tax or not
		withoutTax := store.Filter(excludeTax)
		loop := map[string]data.IStore[*cost.Cost]{
			"Included": store,
			"Excluded": withoutTax,
		}
		for label, s := range loop {
			cells := []*cell.Cell{}
			cells = append(cells, cell.New(label, label, true, false))
			// get the cells for the months
			rowTotal, monthCells := totalPerMonth(s, months)
			cells = append(cells, monthCells...)
			// add extra row total data
			cells = append(cells, cell.New("Totals", rowTotal, false, true))
			r := row.New(cells...)
			rows = append(rows, r)
		}
		return
	}

	funcs = map[string]endpoint.DisplayRowFunc[*cost.Cost]{
		"monthlyTotals": monthlyTotals,
	}

	return
}

// type DisplayFootFunc func(bodyRows []*row.Row) *row.Row
