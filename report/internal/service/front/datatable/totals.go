package datatable

import (
	"fmt"
	"slices"
	"strconv"
)

// ColumnTotalsAveraged assumes the totals are summed values so uses the table values to determine the number of rows present for each column
// and generates the average from that
func ColumnTotalsAveraged(table map[string]map[string]string, identifiers []string, extraCols []string, totals map[string]string, totalCol string) {
	var allCols = []string{}
	var totalCols = []string{}
	var columnCounters = map[string]int{}

	allCols = append(identifiers, extraCols...)
	// find all the column keys that could have values used in total
	for key, _ := range totals {
		if !slices.Contains(allCols, key) {
			totalCols = append(totalCols, key)
			columnCounters[key] = 0
		}
	}
	// now loop over the table and work count how many valid entries for each colum total there are
	for _, row := range table {
		for _, key := range totalCols {
			if _, ok := row[key]; ok {
				if val, e := strconv.ParseFloat(row[key], 64); e == nil && val > 0.0 {
					columnCounters[key]++
				}
			}
		}
	}

	// now adjust the totals to represent the averages
	for key, count := range columnCounters {
		if count > 0 {
			floatCount := float64(count)
			if current, e := strconv.ParseFloat(totals[key], 64); e == nil {
				newTotal := current / floatCount
				totals[key] = fmt.Sprintf("%g", newTotal)
			}
		}

	}

}

// ColumnTotalsSummed is empty, as data is assumed to be sum of totals already
func ColumnTotalsSummed(table map[string]map[string]string, identifiers []string, extraCols []string, totals map[string]string, totalCol string) {
}

// RowTotalsSummed does nothing, data remains the same as the raw data is assumed to data totals
func RowTotalsSummed(table map[string]map[string]string, identifiers []string, columnName string) {}

// RowTotalsAveraged converts the total current in each row to be an average, based on count of how many
// data columns are presents
func RowTotalsAveraged(table map[string]map[string]string, identifiers []string, columnName string) {

	for _, row := range table {
		var columnCount float64 = float64(rowDataColumnCount(row, identifiers, columnName))
		var currentTotal string = row[columnName]
		// if its a float, and greater than 0, work out its average
		if total, e := strconv.ParseFloat(currentTotal, 64); e == nil && total > 0.0 {
			total = total / columnCount
			row[columnName] = fmt.Sprintf("%g", total)
		}
	}
}

// rowDataColumnCount worke out how many data columns are in the row
func rowDataColumnCount(row map[string]string, identifiers []string, totalCol string) (count int) {
	count = 0
	// find only the number of non 0 row values
	for col, v := range row {
		if !slices.Contains(identifiers, col) && col != totalCol {
			if value, e := strconv.ParseFloat(v, 64); e == nil && value > 0.0 {
				count++
			}
		}
	}
	return
}

// AddColumnsToRows injects a every `column` into every row contains in table.
//
// Used to ensure every table row has every column contains, avoiding errors /
// missing keys durning rendering
func AddColumnsToRows(table map[string]map[string]string, columns ...string) {
	for _, row := range table {
		for _, col := range columns {
			if _, ok := row[col]; !ok {
				row[col] = ""
			}
		}
	}
}

// AddRowTotals works our the total values of the numeric entries in each row, assigning the
// value to the `columnName` key.
//
// It knows which keys to use in its calculation by ignoring the `identifiers` slice values, which
// will be the textual data used to group the results on the api.
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

// ColumnTotals is similar to `AddRowTables`, but operated on each column rather than row of the table. It
// loops over every row of the table and creates a total value for each column.
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
