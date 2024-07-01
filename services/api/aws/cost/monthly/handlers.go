package monthly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/response"
	"strings"
	"time"
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
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := response.NewSimpleResult()
	res.Start()
	res.SetStatus(http.StatusOK)
	res.End()
	content, _ := json.Marshal(res)
	a.Write(w, res.GetStatus(), content)
}

// Unit, Env & Service costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}
// Previously "Detailed breakdown" sheet
func (a *Api[V, F]) UnitEnvironmentServices(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, resp)

	errs := resp.GetErrors()

	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonthRange)
		rows := []*response.Row[*response.Cell]{}
		// Add table headings
		headings := response.NewRow[*response.Cell]()
		headings.AddCells(
			response.NewCell("Account", ""),
			response.NewCell("Unit", ""),
			response.NewCell("Environment", ""),
			response.NewCell("Service", ""))
		for _, m := range months {
			headings.AddCells(response.NewCell(m, ""))
		}
		headings.AddCells(response.NewCell("Totals", ""))

		for _, g := range withinMonths.Group(byAccountService) {
			rowTotal := 0.0
			first := g.List()[0]
			cells := []*response.Cell{
				response.NewCell(first.AccountId, first.AccountId),
				response.NewCell(first.AccountUnit, first.AccountUnit),
				response.NewCell(first.AccountEnvironment, first.AccountEnvironment),
				response.NewCell(first.Service, first.Service),
			}
			for _, m := range months {
				inM := func(item *cost.Cost) bool {
					return dates.InMonth(item.Date, []string{m})
				}
				values := g.Filter(inM)
				total := cost.Total(values.List())
				rowTotal += total
				cell := response.NewCell(m, fmt.Sprintf("%f", total))
				cells = append(cells, cell)
			}

			cells = append(cells, response.NewCell("Totals", fmt.Sprintf("%f", rowTotal)))
			// Add data times to resp
			for _, i := range g.List() {
				resp.AddTimestamp(i.TS())
			}

			row := response.NewRow(cells...)
			rows = append(rows, row)
		}
		result := response.NewData(rows...)
		result.SetHeadingsCounters(4, 1)
		result.SetHeadings(headings)
		resp.SetResult(result)

	}
	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)

}

// Unit & Env costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}
// Previously "Service And Environment" sheet
func (a *Api[V, F]) UnitEnvironments(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, resp)

	errs := resp.GetErrors()

	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonthRange)
		rows := []*response.Row[*response.Cell]{}

		// Add table headings
		headings := response.NewRow[*response.Cell]()
		headings.AddCells(response.NewCell("Unit", ""), response.NewCell("Environment", ""))
		for _, m := range months {
			headings.AddCells(response.NewCell(m, ""))
		}
		headings.AddCells(response.NewCell("Totals", ""))

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
				cell := response.NewCell(m, fmt.Sprintf("%f", total))
				cells = append(cells, cell)
				rowTotal += total
			}
			cells = append(cells, response.NewCell("Totals", fmt.Sprintf("%f", rowTotal)))

			row := response.NewRow(cells...)
			rows = append(rows, row)
			// Add data times to resp
			for _, i := range g.List() {
				resp.AddTimestamp(i.TS())
			}

		}
		result := response.NewData(rows...)
		result.SetHeadingsCounters(2, 1)
		result.SetHeadings(headings)
		resp.SetResult(result)

	}
	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)

}

// Unit costs: /aws/costs/{version}/monthly/{start}/{end}/units/{$}
// Previously "Service" sheet
func (a *Api[V, F]) Units(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, resp)

	errs := resp.GetErrors()
	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}

		withinMonths := store.Filter(inMonthRange)
		// Add table headings
		headings := response.NewRow[*response.Cell]()
		headings.AddCells(response.NewCell("Unit", ""))
		for _, m := range months {
			headings.AddCells(response.NewCell(m, ""))
		}
		headings.AddCells(response.NewCell("Totals", ""))

		rows := []*response.Row[*response.Cell]{}

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
				cell := response.NewCell(m, fmt.Sprintf("%f", total))
				cells = append(cells, cell)
				rowTotal += total
			}
			cells = append(cells, response.NewCell("Totals", fmt.Sprintf("%f", rowTotal)))
			// Add data times to resp
			for _, i := range g.List() {
				resp.AddTimestamp(i.TS())
			}
			row := response.NewRow(cells...)
			rows = append(rows, row)
		}
		result := response.NewData(rows...)
		result.SetHeadingsCounters(1, 1)
		result.SetHeadings(headings)
		resp.SetResult(result)
	}

	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)

}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
// Returns cost data split into with & without tax segments, then grouped by the month
// Previously "Totals" sheet
// Note: if {start} or {end} are "-" it uses current month
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, resp)

	errs := resp.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonthRange := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		// Add table headings
		headings := response.NewRow[*response.Cell]()
		headings.AddCells(response.NewCell("Tax", ""))
		for _, m := range months {
			headings.AddCells(response.NewCell(m, ""))
		}
		// add the row total empty head
		headings.AddCells(response.NewCell("Totals", ""))

		// get everything in range
		withTax := store.Filter(inMonthRange)
		// Add data times to resp
		for _, i := range withTax.List() {
			resp.AddTimestamp(i.TS())
		}
		withTaxRow := withTaxR(withTax, months)
		//  exclude tax from the costs
		withoutTax := withTax.Filter(excludeTax)
		withoutTaxRow := withoutTaxR(withoutTax, months)

		result := response.NewData(withoutTaxRow, withTaxRow)
		result.SetHeadingsCounters(1, 1)
		result.SetHeadings(headings)
		resp.SetResult(result)
	}
	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)
}

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
		cell := response.NewCell(m, fmt.Sprintf("%f", total))
		withoutTaxCells = append(withoutTaxCells, cell)
	}
	withoutTaxCells = append(withoutTaxCells, response.NewCell("Totals", fmt.Sprintf("%f", rowTotal)))
	return response.NewRow(withoutTaxCells...)
}
func withTaxR(withTax data.IStore[*cost.Cost], months []string) *response.Row[*response.Cell] {
	rowTotal := 0.0
	withTaxCells := []*response.Cell{
		response.NewCell("Included", "Included"),
	}
	for _, m := range months {
		inM := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, []string{m})
		}
		values := withTax.Filter(inM)
		total := cost.Total(values.List())
		rowTotal += total
		cell := response.NewCell(m, fmt.Sprintf("%f", total))
		withTaxCells = append(withTaxCells, cell)
	}
	withTaxCells = append(withTaxCells, response.NewCell("Totals", fmt.Sprintf("%f", rowTotal)))
	return response.NewRow(withTaxCells...)
}

func startAndEndDates(r *http.Request, res response.IBase) (startDate time.Time, endDate time.Time) {
	var err error
	now := time.Now().UTC().Format(dates.FormatYM)
	// Get the dates
	startDate, err = dates.StringToDateDefault(r.PathValue("start"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	endDate, err = dates.StringToDateDefault(r.PathValue("end"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	return

}
