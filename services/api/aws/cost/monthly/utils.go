package monthly

import (
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/resp/cell"
	"opg-reports/shared/server/resp/row"
	"time"
)

func startEnd(parameters map[string][]string) (time.Time, time.Time) {
	now := time.Now().UTC()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startDate := firstDay
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if starts, ok := parameters["start"]; ok && starts[0] != "-" {
		startDate, _ = dates.StringToDateDefault(starts[0], "-", firstDay.Format(dates.FormatYM))
	}
	if ends, ok := parameters["end"]; ok && ends[0] != "-" {
		endDate, _ = dates.StringToDateDefault(ends[0], "-", endDate.Format(dates.FormatYM))
	}
	return startDate, endDate
}

func totalPerMonth(store data.IStore[*cost.Cost], months []string) (rowTotal float64, cells []*cell.Cell) {
	cells = []*cell.Cell{}
	rowTotal = 0.0
	for _, m := range months {
		inMonth := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, []string{m})
		}
		values := store.Filter(inMonth)
		monthTotal := cost.Total(values.List())
		rowTotal += monthTotal
		// add to the set of cells
		cells = append(cells, cell.New(m, monthTotal, false, false))
	}
	return
}

func columnTotals(rows []*row.Row) (totalRow *row.Row) {
	var totals []float64
	totalRow = row.New()
	headingsCount := 0
	if len(rows) > 0 {
		first := rows[0]
		firstCells := first.All()
		l := len(firstCells)
		totals = make([]float64, l)
		for i := 0; i < l; i++ {
			totals[i] = 0.0
		}
		headingsCount = len(first.Headers)
	}

	for _, r := range rows {
		cells := r.All()
		for x := headingsCount; x < len(cells); x++ {
			totals[x] += cells[x].Value.(float64)
		}
	}

	for i, total := range totals {
		var c *cell.Cell
		if i < headingsCount {
			c = cell.New("Total", total, true, false)
		} else {
			c = cell.New("Total", total, false, false)
		}
		totalRow.Add(c)
	}
	return
}
