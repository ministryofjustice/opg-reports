package awscosts

import (
	"fmt"
	"log/slog"
	"opg-reports/report/cmd/api/tables"
	"opg-reports/report/internal/service/api"
	"sort"
	"strconv"
)

// TabulateGroupedCosts converts the data into structgure table values
// that will then be displayed directly by the front end
//
// Uses all possible combinations of the columns to generate a skeleton
// table which is then populated with data from the dbRecords
//
// Table is then sorted in descending order using the last date item
func TabulateGroupedCosts[T api.Model](
	log *slog.Logger,
	columns []string,
	dates []string,
	dbRecords []T,
) (sortedTable []map[string]string, footer map[string]string, err error) {

	var (
		lastDate    string
		table       map[string]map[string]string
		trendColumn string                    = "trend"
		totalColumn string                    = "total"
		tableCfg    *tables.ListToTableConfig = &tables.ListToTableConfig{
			TextColumns:     columns,
			DataColumns:     dates,
			ValueField:      "cost",
			DataSourceField: "date",
			DefaultValue:    "0.00",
		}
	)
	// inits
	sortedTable = []map[string]string{}
	footer = map[string]string{}
	sort.Strings(dates)
	lastDate = dates[len(dates)-1]

	log = log.With("tableCfg", tableCfg).With("operation", "TabulateGroupedCosts")

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
	log.Debug("add in row totals")
	rowTotals(table, dates, totalColumn)
	// now create the footer with totals
	log.Debug("creating column totals for the footer")
	footer = columnTotals(table,
		append(columns, trendColumn),
		append(dates, totalColumn))

	// now copy over the map to a slice for sorting
	for _, v := range table {
		sortedTable = append(sortedTable, v)
	}
	log.With("lastDate", lastDate).Debug("sorted table decending order by last column value")
	// now sort the data by the last date column, highest value first (descending)
	sort.SliceStable(sortedTable, func(i, j int) bool {
		var (
			a    float64           = 0.0
			b    float64           = 0.0
			aRow map[string]string = sortedTable[i]
			bRow map[string]string = sortedTable[j]
			res  bool
		)
		if v, err := strconv.ParseFloat(aRow[lastDate], 64); err == nil {
			a = v
		}
		if v, err := strconv.ParseFloat(bRow[lastDate], 64); err == nil {
			b = v
		}
		res = (a > b)
		return res

	})

	return
}

func columnTotals(table map[string]map[string]string, cols []string, colsToSum []string) (totals map[string]string) {
	var sums map[string]float64 = map[string]float64{}
	totals = map[string]string{}
	// blank out the columns
	for _, k := range cols {
		totals[k] = ""
	}
	// blank out the map to 0
	for _, col := range colsToSum {
		sums[col] = 0.0
	}

	for _, row := range table {
		for _, col := range colsToSum {
			if add, e := strconv.ParseFloat(row[col], 64); e == nil {
				sums[col] += add
			}
		}
	}
	// convert to strings
	for k, v := range sums {
		totals[k] = fmt.Sprintf("%g", v)
	}
	return
}

func rowTotals(table map[string]map[string]string, dates []string, totalCol string) {
	// now inject row totals
	for _, row := range table {
		var rowTotal float64 = 0.0
		for _, date := range dates {
			var dateVal float64 = 0.0
			if val, e := strconv.ParseFloat(row[date], 64); e == nil && val != 0.0 {
				dateVal = val
			}
			rowTotal += dateVal
		}
		row[totalCol] = fmt.Sprintf("%g", rowTotal)

	}

}
