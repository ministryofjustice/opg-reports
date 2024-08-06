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

	hMonths := HeaderMonths(months)
	// -- monthly totals
	// This has a split on in tax is included or not and then
	// a final col for the line total
	var monthlyTotals endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/costs/monthly] monthlyTotals head func")
		r = row.New()
		r.Add(cell.New("Tax", "Tax", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Totals", "Totals", false, true))
		return
	}
	// -- per unit
	var perUnit endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/costs/monthly] perUnit head func")
		r = row.New()
		r.Add(cell.New("Unit", "Unit", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Totals", "Totals", false, true))
		return
	}
	// -- per unit & env
	var perUnitEnv endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/costs/monthly] perUnitEnv head func")
		r = row.New()
		r.Add(cell.New("Unit", "Unit", true, false))
		r.Add(cell.New("Environment", "Environment", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Totals", "Totals", false, true))
		return
	}
	// -- per unit, env, service
	var perUnitEnvService endpoint.DisplayHeadFunc = func() (r *row.Row) {
		slog.Debug("[aws/costs/monthly] perUnitEnvService head func")
		r = row.New()
		r.Add(cell.New("Account", "Account", true, false))
		r.Add(cell.New("Unit", "Unit", true, false))
		r.Add(cell.New("Environment", "Environment", true, false))
		r.Add(cell.New("Service", "Service", true, false))
		r.Add(hMonths...)
		r.Add(cell.New("Totals", "Totals", false, true))
		return
	}

	funcs = map[string]endpoint.DisplayHeadFunc{
		"monthlyTotals":     monthlyTotals,
		"perUnit":           perUnit,
		"perUnitEnv":        perUnitEnv,
		"perUnitEnvService": perUnitEnvService,
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
			rowTotal, monthCells := TotalPerMonth(s, months)
			cells = append(cells, monthCells...)
			// add extra row total data
			cells = append(cells, cell.New("Totals", rowTotal, false, true))
			r := row.New(cells...)
			rows = append(rows, r)
		}
		return
	}
	// -- per unit
	var perUnit endpoint.DisplayRowFunc[*cost.Cost] = func(group string, store data.IStore[*cost.Cost], resp *resp.Response) (rows []*row.Row) {
		first := store.List()[0]
		// row headers
		cells := []*cell.Cell{
			cell.New(first.AccountUnit, first.AccountUnit, true, false),
		}
		// get the row months
		rowTotal, monthCells := TotalPerMonth(store, months)
		cells = append(cells, monthCells...)
		// totals
		cells = append(cells, cell.New("Totals", rowTotal, false, true))
		// return the row
		rows = []*row.Row{row.New(cells...)}
		return
	}
	// -- per unit & env
	var perUnitEnv endpoint.DisplayRowFunc[*cost.Cost] = func(group string, store data.IStore[*cost.Cost], resp *resp.Response) (rows []*row.Row) {
		first := store.List()[0]
		// row headers
		cells := []*cell.Cell{
			cell.New(first.AccountUnit, first.AccountUnit, true, false),
			cell.New(first.AccountEnvironment, first.AccountEnvironment, true, false),
		}
		// get the row months
		rowTotal, monthCells := TotalPerMonth(store, months)
		cells = append(cells, monthCells...)
		// totals
		cells = append(cells, cell.New("Totals", rowTotal, false, true))
		// return the row
		rows = []*row.Row{row.New(cells...)}
		return
	}
	// perUnitEnvService
	var perUnitEnvService endpoint.DisplayRowFunc[*cost.Cost] = func(group string, store data.IStore[*cost.Cost], resp *resp.Response) (rows []*row.Row) {
		first := store.List()[0]
		// row headers
		cells := []*cell.Cell{
			cell.New(first.AccountId, first.AccountId, true, false),
			cell.New(first.AccountUnit, first.AccountUnit, true, false),
			cell.New(first.AccountEnvironment, first.AccountEnvironment, true, false),
			cell.New(first.Service, first.Service, true, false),
		}
		// get the row months
		rowTotal, monthCells := TotalPerMonth(store, months)
		cells = append(cells, monthCells...)
		// totals
		cells = append(cells, cell.New("Totals", rowTotal, false, true))
		// return the row
		rows = []*row.Row{row.New(cells...)}
		return
	}

	funcs = map[string]endpoint.DisplayRowFunc[*cost.Cost]{
		"monthlyTotals":     monthlyTotals,
		"perUnit":           perUnit,
		"perUnitEnv":        perUnitEnv,
		"perUnitEnvService": perUnitEnvService,
	}

	return
}

// DisplayFootFunctions
func DisplayFootFunctions(parameters map[string][]string) (funcs map[string]endpoint.DisplayFootFunc) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := startEnd(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/costs/monthly] DisplayFootFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

	var perUnit endpoint.DisplayFootFunc = columnTotals
	var perUnitEnv endpoint.DisplayFootFunc = columnTotals
	var perUnitEnvService endpoint.DisplayFootFunc = columnTotals

	funcs = map[string]endpoint.DisplayFootFunc{
		"perUnit":           perUnit,
		"perUnitEnv":        perUnitEnv,
		"perUnitEnvService": perUnitEnvService,
	}

	return
}
