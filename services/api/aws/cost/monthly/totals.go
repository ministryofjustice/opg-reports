package monthly

import "opg-reports/shared/server/response"

// columnTotals creates a row whose cell values are the total for that column
func (a *Api[V, F, C, R]) columnTotals(rows []R) (row R) {
	var totals []float64
	row = response.NewRow[C]().(R)
	headingsCount := 0
	if len(rows) > 0 {
		totals = make([]float64, rows[0].GetTotalCellCount())
		for i := 0; i < rows[0].GetTotalCellCount(); i++ {
			totals[i] = 0.0
		}
		headingsCount = rows[0].GetHeadersCount()
	}

	for _, r := range rows {
		cells := r.GetAll()
		for x := headingsCount; x < len(cells); x++ {
			totals[x] += cells[x].GetValue().(float64)
		}
	}

	for i, total := range totals {
		var c C
		if i < headingsCount {
			c = response.NewCellHeader("Total", total).(C)
		} else {
			c = response.NewCell("Total", total).(C)
		}

		row.SetRaw(c)
	}
	return
}
