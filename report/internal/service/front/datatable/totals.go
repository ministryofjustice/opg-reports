package datatable

import (
	"fmt"
	"slices"
	"strconv"
)

func AddColumnsToRows(table map[string]map[string]string, columns ...string) {
	for _, row := range table {
		for _, col := range columns {
			if _, ok := row[col]; !ok {
				row[col] = ""
			}
		}
	}
}

func AddRowTotals(table map[string]map[string]string, identifiers []string, columnName string) {

	for _, row := range table {
		rowTotal := 0.0
		for col, val := range row {
			if slices.Contains(identifiers, col) {
				continue
			}
			if add, e := strconv.ParseFloat(val, 64); e == nil {
				rowTotal += add
			}
		}
		row[columnName] = fmt.Sprintf("%g", rowTotal)
	}

}

func ColumnTotals(table map[string]map[string]string, sumColumns []string, extraCols ...string) (totals map[string]string) {
	totals = map[string]string{}

	sums := map[string]float64{}
	for _, col := range sumColumns {
		sums[col] = 0.0
	}

	for _, row := range table {
		for _, col := range sumColumns {
			if add, e := strconv.ParseFloat(row[col], 64); e == nil {
				sums[col] += add
			}
		}
	}
	// convert to strings
	for k, v := range sums {
		totals[k] = fmt.Sprintf("%g", v)
	}
	// add extra columns
	for _, k := range extraCols {
		totals[k] = ""
	}

	return
}
