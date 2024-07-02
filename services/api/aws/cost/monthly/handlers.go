package monthly

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/response"
	"strings"
)

// Filters
var excludeTax = func(item *cost.Cost) bool {
	return strings.ToLower(item.Service) != "tax"
}

// Helpers used within grouping
var unit = func(i *cost.Cost) (string, string) {
	return "account_unit", strings.ToLower(i.AccountUnit)
}
var account_id = func(i *cost.Cost) (string, string) {
	return "account_id", i.AccountId
}
var account_env = func(i *cost.Cost) (string, string) {
	return "account_environment", strings.ToLower(i.AccountEnvironment)
}
var service = func(i *cost.Cost) (string, string) {
	return "service", strings.ToLower(i.Service)
}

// Group by month
var byUnit = func(item *cost.Cost) string {
	return data.ToIdxF(item, unit)
}
var byUnitEnv = func(item *cost.Cost) string {
	return data.ToIdxF(item, unit, account_env)
}
var byAccountService = func(item *cost.Cost) string {
	return data.ToIdxF(item, account_id, unit, account_env, service)
}

// Index: /aws/costs/v1/monthly
func (a *Api[V, F, C, R]) Index(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	a.End(w, r)

}

// Unit, Env & Service costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}
// Previously "Detailed breakdown" sheet
func (a *Api[V, F, C, R]) UnitEnvironmentServices(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()
	startDate, endDate := a.startAndEndDates(r)

	errs := resp.GetError()

	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(c *cost.Cost) bool {
			return dates.InMonth(c.Date, months)
		}
		withinMonths := store.Filter(inMonthRange)

		// Add table headings
		headingCells := []*response.Cell{
			response.NewCellHeader("Account", ""),
			response.NewCellHeader("Unit", ""),
			response.NewCellHeader("Environment", ""),
			response.NewCellHeader("Service", ""),
		}
		for _, m := range months {
			headingCells = append(headingCells, response.NewCell(m, ""))
		}
		headingCells = append(headingCells, response.NewCellExtra("Totals", ""))

		head := response.NewRow(headingCells...)
		body := []*response.Row[*response.Cell]{}

		for _, g := range withinMonths.Group(byAccountService) {
			rowTotal := 0.0
			first := g.List()[0]
			cells := []*response.Cell{
				response.NewCellHeader(first.AccountId, first.AccountId),
				response.NewCellHeader(first.AccountUnit, first.AccountUnit),
				response.NewCellHeader(first.AccountEnvironment, first.AccountEnvironment),
				response.NewCellHeader(first.Service, first.Service),
			}
			for _, m := range months {
				inM := func(item *cost.Cost) bool {
					return dates.InMonth(item.Date, []string{m})
				}
				values := g.Filter(inM)
				total := cost.Total(values.List())
				rowTotal += total
				cell := response.NewCell(m, total)
				cells = append(cells, cell)
			}
			cells = append(cells, response.NewCellExtra("Totals", rowTotal))

			// Add data times to resp
			for _, i := range g.List() {
				resp.SetDataAge(i.TS())
			}

			row := response.NewRow(cells...)
			body = append(body, row)
		}

		table := response.NewTable(body...)
		table.SetTableHead(head)
		// table.SetTableFoot(foot)

	}
	a.End(w, r)

}

// Unit & Env costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}
// Previously "Service And Environment" sheet
func (a *Api[V, F, C, R]) UnitEnvironments(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()
	startDate, endDate := a.startAndEndDates(r)

	errs := resp.GetError()
	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonthRange)
		// Add table headings
		headingCells := []*response.Cell{
			response.NewCellHeader("Unit", ""),
			response.NewCellHeader("Environment", ""),
		}
		for _, m := range months {
			headingCells = append(headingCells, response.NewCell(m, ""))
		}
		headingCells = append(headingCells, response.NewCellExtra("Totals", ""))

		head := response.NewRow(headingCells...)
		body := []*response.Row[*response.Cell]{}
		// Add table headings

		// loop over month group data to group the other data
		for _, g := range withinMonths.Group(byUnitEnv) {
			first := g.List()[0]
			rowTotal := 0.0
			cells := []*response.Cell{
				response.NewCell(first.AccountUnit, first.AccountUnit),
				response.NewCell(first.AccountEnvironment, first.AccountEnvironment),
			}
			for _, m := range months {
				inM := func(item *cost.Cost) bool {
					return dates.InMonth(item.Date, []string{m})
				}
				values := g.Filter(inM)
				total := cost.Total(values.List())
				cell := response.NewCell(m, total)
				cells = append(cells, cell)
				rowTotal += total
			}
			cells = append(cells, response.NewCellExtra("Totals", rowTotal))

			row := response.NewRow(cells...)
			body = append(body, row)
			// Add data times to resp
			for _, i := range g.List() {
				resp.SetDataAge(i.TS())
			}

		}
		table := response.NewTable(body...)
		table.SetTableHead(head)
		// table.SetTableFoot(foot)

	}
	a.End(w, r)

}

// Unit costs: /aws/costs/{version}/monthly/{start}/{end}/units/{$}
// Previously "Service" sheet
func (a *Api[V, F, C, R]) Units(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()
	startDate, endDate := a.startAndEndDates(r)

	errs := resp.GetError()
	if len(errs) == 0 {
		// data range
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonthRange)
		// Add table headings
		headingCells := []*response.Cell{
			response.NewCellHeader("Unit", ""),
		}
		for _, m := range months {
			headingCells = append(headingCells, response.NewCell(m, ""))
		}
		headingCells = append(headingCells, response.NewCellExtra("Totals", ""))

		head := response.NewRow(headingCells...)
		body := []*response.Row[*response.Cell]{}

		for _, g := range withinMonths.Group(byUnit) {
			first := g.List()[0]
			rowTotal := 0.0
			cells := []*response.Cell{
				response.NewCell(first.AccountUnit, first.AccountUnit),
			}
			for _, m := range months {
				inM := func(item *cost.Cost) bool {
					return dates.InMonth(item.Date, []string{m})
				}
				values := g.Filter(inM)
				total := cost.Total(values.List())
				cell := response.NewCell(m, total)
				cells = append(cells, cell)
				rowTotal += total
			}
			cells = append(cells, response.NewCellExtra("Totals", rowTotal))
			/// Add data times to resp
			for _, i := range g.List() {
				resp.SetDataAge(i.TS())
			}
			row := response.NewRow(cells...)
			body = append(body, row)
		}
		table := response.NewTable(body...)
		table.SetTableHead(head)
		// table.SetTableFoot(foot)
	}

	a.End(w, r)

}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
// Returns cost data split into with & without tax segments, then grouped by the month
// Previously "Totals" sheet
// Note: if {start} or {end} are "-" it uses current month
func (a *Api[V, F, C, R]) Totals(w http.ResponseWriter, r *http.Request) {
	slog.Info("[api/aws/costs/monthly] totals", slog.String("uri", r.RequestURI))
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()
	startDate, endDate := a.startAndEndDates(r)

	errs := resp.GetError()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		// Add table headings
		headingCells := []*response.Cell{
			response.NewCellHeader("Tax", ""),
		}
		for _, m := range months {
			headingCells = append(headingCells, response.NewCell(m, ""))
		}
		headingCells = append(headingCells, response.NewCellExtra("Totals", ""))
		head := response.NewRow(headingCells...)

		// get everything in range
		withTax := store.Filter(inMonthRange)
		/// Add data times to resp
		for _, i := range withTax.List() {
			resp.SetDataAge(i.TS())
		}
		withTaxRow := withTaxR(withTax, months)
		//  exclude tax from the costs
		withoutTax := withTax.Filter(excludeTax)
		withoutTaxRow := withoutTaxR(withoutTax, months)

		table := response.NewTable[*response.Cell, *response.Row[*response.Cell]]()
		table.SetTableHead(head)
		table.SetTableBody(withoutTaxRow)
		table.SetTableBody(withTaxRow)

	}
	a.End(w, r)
}

// func columnTotals(rows []*response.Row[*response.Cell], pre int) *response.Row[*response.Cell] {
// 	var totals []float64
// 	footer := response.NewRow[*response.Cell]()

// 	if len(rows) > 0 {
// 		totals = make([]float64, rows[0].Len())
// 		for i := 0; i < rows[0].Len(); i++ {
// 			totals[i] = 0.0
// 		}
// 	}

// 	for _, r := range rows {
// 		cells := r.GetCells()
// 		for x := pre; x < len(cells); x++ {
// 			totals[x] += cells[x].Value.(float64)
// 		}
// 	}
// 	for _, total := range totals {
// 		footer.AddCells(response.NewCell("Total", total))
// 	}
// 	return footer

// }

func withoutTaxR(withoutTax data.IStore[*cost.Cost], months []string) *response.Row[*response.Cell] {
	rowTotal := 0.0
	withoutTaxCells := []*response.Cell{
		response.NewCell("Excluded", "Excluded"),
	}
	for _, m := range months {
		inM := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, []string{m})
		}
		values := withoutTax.Filter(inM)
		total := cost.Total(values.List())
		rowTotal += total
		cell := response.NewCell(m, total)
		withoutTaxCells = append(withoutTaxCells, cell)
	}
	withoutTaxCells = append(withoutTaxCells, response.NewCell("Totals", rowTotal))
	return response.NewRow(withoutTaxCells...)
}
func withTaxR(withTax data.IStore[*cost.Cost], months []string) *response.Row[*response.Cell] {
	rowTotal := 0.0
	withTaxCells := []*response.Cell{
		response.NewCellHeader("Included", "Included"),
	}
	for _, m := range months {
		inM := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, []string{m})
		}
		values := withTax.Filter(inM)
		total := cost.Total(values.List())
		rowTotal += total
		cell := response.NewCell(m, total)
		withTaxCells = append(withTaxCells, cell)
	}
	withTaxCells = append(withTaxCells, response.NewCellExtra("Totals", rowTotal))
	return response.NewRow(withTaxCells...)
}
