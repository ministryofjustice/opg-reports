package monthly

import (
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/response"
)

func withoutTaxR(withoutTax data.IStore[*cost.Cost], months []string) response.IRow[response.ICell] {
	rowTotal := 0.0
	withoutTaxCells := []response.ICell{
		response.NewCellHeader("Excluded", "Excluded"),
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
	withoutTaxCells = append(withoutTaxCells, response.NewCellExtra("Totals", rowTotal))
	return response.NewRow(withoutTaxCells...)
}

func withTaxR(withTax data.IStore[*cost.Cost], months []string) response.IRow[response.ICell] {
	rowTotal := 0.0
	withTaxCells := []response.ICell{
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
