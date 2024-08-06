package daily

import (
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/resp/cell"
)

func HeaderMonths(months []string) (cells []*cell.Cell) {
	cells = []*cell.Cell{}
	for _, m := range months {
		cells = append(cells, cell.New(m, m, false, false))
	}
	return
}

func AvgPerMonth(store data.IStore[*uptime.Uptime], months []string) (rowAvg float64, cells []*cell.Cell) {
	cells = []*cell.Cell{}
	rowCount := 0.0
	for _, m := range months {
		inMonth := func(item *uptime.Uptime) bool {
			dateStr := item.DateTime.Format(dates.FormatYM)
			return dates.InMonth(dateStr, []string{m})
		}
		values := store.Filter(inMonth)
		monthlyAvg := 0.0
		if values.Length() > 0 {
			monthlyAvg = uptime.Average(values.List())
		}
		rowCount += monthlyAvg
		// add to the set of cells
		cells = append(cells, cell.New(m, monthlyAvg, false, false))
	}
	l := len(months)
	rowAvg = rowCount / float64(l)
	return
}
