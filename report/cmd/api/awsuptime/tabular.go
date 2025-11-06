package awsuptime

import (
	"fmt"
	"log/slog"
	"opg-reports/report/cmd/api/tables"
	"opg-reports/report/internal/service/api"
	"sort"
	"strconv"
)

// TabulateGroupedUptime converts the data into structured table values
// that will then be displayed directly by the front end
//
// Uses all possible combinations of the columns to generate a skeleton
// table which is then populated with data from the dbRecords
func TabulateGroupedUptime[T api.Model](
	log *slog.Logger,
	columns []string,
	dates []string,
	dbRecords []T,
) (sortedTable []map[string]string, footer map[string]string, err error) {

	var (
		table       map[string]map[string]string
		trendColumn string                    = "trend"
		totalColumn string                    = "total"
		tableCfg    *tables.ListToTableConfig = &tables.ListToTableConfig{
			TextColumns:     columns,
			DataColumns:     dates,
			ValueField:      "average",
			DataSourceField: "date",
			DefaultValue:    "0.00",
		}
	)
	log = log.With("tableCfg", tableCfg).With("operation", "TabulateGroupedUptime")
	// inits
	sortedTable = []map[string]string{}
	footer = map[string]string{}
	sort.Strings(dates)

	// convert the raw database records into table structure
	log.Debug("converting db records to table")
	table, err = tables.ListToTable(log, tableCfg, dbRecords)
	if err != nil {
		return
	}

	// add in trend
	log.Debug("add the trend column into each row on the table")
	tables.AddColumnToEachRow(table, trendColumn)
	// add in row totals
	log.Debug("add in row averages for total")
	rowAverages(table, dates, totalColumn)
	// now create the footer with totals
	log.Debug("creating column averages for the footer")
	footer = columnAverages(table,
		append(columns, trendColumn),
		append(dates, totalColumn))

	// now copy over the map to a slice for sorting
	for _, v := range table {
		sortedTable = append(sortedTable, v)
	}
	sort.SliceStable(sortedTable, func(i, j int) bool {
		var a = sortedTable[i]["team"]
		var b = sortedTable[j]["team"]
		return (a < b)
	})

	return
}

func columnAverages(table map[string]map[string]string, cols []string, colsToSum []string) (totals map[string]string) {
	var sums map[string]float64 = map[string]float64{}
	var counts map[string]float64 = map[string]float64{}
	totals = map[string]string{}
	// blank out the columns
	for _, k := range cols {
		totals[k] = ""
	}
	// blank out the map to 0
	for _, col := range colsToSum {
		sums[col] = 0.0
		counts[col] = 0.0
	}

	for _, row := range table {
		for _, col := range colsToSum {
			if add, e := strconv.ParseFloat(row[col], 64); e == nil {
				sums[col] += add
				counts[col] += 1.0
			}
		}
	}
	// convert to strings
	for k, v := range sums {
		totals[k] = fmt.Sprintf("%g", (v / counts[k]))
	}
	return
}

func rowAverages(table map[string]map[string]string, dates []string, totalCol string) {

	// now inject row totals
	for _, row := range table {
		var count float64 = 0.0
		var total float64 = 0.0

		for _, date := range dates {
			var dateVal float64 = 0.0
			if val, e := strconv.ParseFloat(row[date], 64); e == nil && val != 0.0 {
				dateVal = val
				count += 1.0
			}
			total += dateVal
		}
		row[totalCol] = fmt.Sprintf("%g", (total / count))

	}

}
