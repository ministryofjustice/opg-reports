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
		monthlyAvg := uptime.Average(values.List())
		rowCount += monthlyAvg
		// add to the set of cells
		cells = append(cells, cell.New(m, monthlyAvg, false, false))
	}
	l := len(months)
	rowAvg = rowCount / float64(l)
	return
}

// func columnTotals(rows []*row.Row) (totalRow *row.Row) {
// 	var totals []float64
// 	totalRow = row.New()
// 	headingsCount := 0
// 	if len(rows) > 0 {
// 		first := rows[0]
// 		firstCells := first.All()
// 		l := len(firstCells)
// 		totals = make([]float64, l)
// 		for i := 0; i < l; i++ {
// 			totals[i] = 0.0
// 		}
// 		headingsCount = len(first.Headers)
// 	}

// 	for _, r := range rows {
// 		cells := r.All()
// 		for x := headingsCount; x < len(cells); x++ {
// 			totals[x] += cells[x].Value.(float64)
// 		}
// 	}

// 	for i, total := range totals {
// 		var c *cell.Cell
// 		if i < headingsCount {
// 			c = cell.New("Total", total, true, false)
// 		} else {
// 			c = cell.New("Total", total, false, true)
// 		}
// 		totalRow.Add(c)
// 	}
// 	return
// }
