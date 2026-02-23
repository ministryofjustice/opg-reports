package tabulate

import (
	"context"
	"fmt"
	"math"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/dump"
	"sort"
)

type TableEndFunc func(table []map[string]interface{}, headings map[ColType][]string) []map[string]interface{}

type Args struct {
	ColumnKey string
	ValueKey  string
	Headers   map[ColType][]string
}

// TableBody generates body from T data
func TableBody[T any](ctx context.Context, dbRows []T, in *Args) (tableMap map[string]map[string]interface{}) {
	var rowMap = []map[string]interface{}{}

	cnv.Convert(dbRows, &rowMap)
	tableMap = map[string]map[string]interface{}{}

	for _, src := range rowMap {
		var rowKey = RowKey(src, in.Headers)
		var dest, ok = tableMap[rowKey]
		// if not set, add a default rowq
		if !ok {
			dest = EmptyRow(in.Headers)
		}
		// populate
		PopulateRow(src, dest, in.Headers, in)
		// set table value
		tableMap[rowKey] = dest
	}
	return
}

func TableMapToTable(tableMap map[string]map[string]interface{}) (table []map[string]interface{}) {
	table = []map[string]interface{}{}
	// convert to slice
	for _, row := range tableMap {
		table = append(table, row)
	}
	return
}

// RowEnd runs the row end function (like adding total / average)
func TableEnd(table []map[string]interface{}, headings map[ColType][]string, tableEndF TableEndFunc) []map[string]interface{} {
	if tableEndF == nil {
		return table
	}
	table = tableEndF(table, headings)
	return table
}

// TotalF generates a table summary row contains the totals of each data column combined
func TableTotalF(table []map[string]interface{}, headings map[ColType][]string) []map[string]interface{} {
	var (
		endCol     string   = headings[END][0]
		firstCol   string   = headings[KEY][0]
		dataCols   []string = headings[DATA]
		tableTotal float64  = 0.0
		summary             = EmptyRow(headings)
	)

	if firstCol == "" || endCol == "" {
		panic(fmt.Sprintf("TotalF missing first/end columns in headings:\n[%s]", dump.Any(headings)))
	}
	for _, row := range table {
		for _, col := range dataCols {
			summary[col] = Value[float64](col, 0.0, summary) + Value[float64](col, 0.0, row)
		}
		tableTotal += row[endCol].(float64)
	}

	// give the first column a name
	summary[firstCol] = endCol
	summary[endCol] = tableTotal
	table = append(table, summary)
	return table
}

// AverageF generates a table summary row contains the totals averages of each data column combined
func TableAverageF(table []map[string]interface{}, headings map[ColType][]string) []map[string]interface{} {
	var (
		endCol     string   = headings[END][0]
		firstCol   string   = headings[KEY][0]
		dataCols   []string = headings[DATA]
		count      int      = 0 //len(table)
		tableTotal float64  = 0.0
		summary             = EmptyRow(headings)
	)
	if firstCol == "" || endCol == "" {
		panic(fmt.Sprintf("AverageF missing first/end columns in headings:\n[%s]", dump.Any(headings)))
	}
	for _, row := range table {
		for _, col := range dataCols {
			summary[col] = Value[float64](col, 0.0, summary) + Value[float64](col, 0.0, row)
		}
		rowE := Value[float64](endCol, 0.0, row)
		tableTotal += rowE
		if rowE != 0 {
			count++
		}
	}
	// now divide to fix create average
	for _, col := range dataCols {
		summary[col] = summary[col].(float64) / float64(count)
	}

	// give the first column a name
	summary[firstCol] = endCol
	// work out overall average
	summary[endCol] = tableTotal / float64(count)
	table = append(table, summary)
	return table
}

func TableFilterByValue(table []map[string]interface{}, headings map[ColType][]string, over float64) []map[string]interface{} {
	var endCol = headings[END][0]
	var filtered = []map[string]interface{}{}
	for _, row := range table {
		var value = math.Abs(Value[float64](endCol, 0.0, row))
		if value >= over {
			filtered = append(filtered, row)
		}
	}

	return filtered
}

func SortAscending[T int | float64 | string](table []map[string]interface{}, column string) []map[string]interface{} {
	sort.Slice(table, func(i, j int) bool {
		var a = table[i][column].(T)
		var b = table[j][column].(T)
		return (a < b)
	})
	return table
}

func SortDescending[T int | float64 | string](table []map[string]interface{}, column string) []map[string]interface{} {
	sort.Slice(table, func(i, j int) bool {
		var a = table[i][column].(T)
		var b = table[j][column].(T)
		return (a > b)
	})
	return table
}
