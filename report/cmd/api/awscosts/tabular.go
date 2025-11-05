package awscosts

import (
	"fmt"
	"log/slog"
	"opg-reports/report/cmd/api/tables"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"
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
		trendCol     string                       = "trend"                        // append to each row for front end rendering
		totalsCol    string                       = "total"                        // name of column to add row totals into
		transformCol string                       = "date"                         // the column name in each row to generate headers and merge data on
		valueCol     string                       = "cost"                         // the value column used to merge rows with
		emptyVal     string                       = "0.00"                         // empty / place holder string
		records      []map[string]string          = []map[string]string{}          // converted version of records into a generic slice map
		rowKeys      []string                     = []string{}                     // contains all the possible row keys based on the values of the group columns
		table        map[string]map[string]string = map[string]map[string]string{} // all of the data converted into a table format, firstly via skeleton, then via populated
		lastDate     string                                                        // last date, used for sorting
	)
	log = log.With("operation", "TabulateGroupedCosts")
	// init
	sortedTable = []map[string]string{}
	table = map[string]map[string]string{}
	footer = map[string]string{}
	// tidy dates
	sort.Strings(dates)
	lastDate = dates[len(dates)-1]

	log.Debug("converting []T to slice of maps")
	// convert to slice map from the records
	err = utils.Convert(dbRecords, &records)
	if err != nil {
		log.Error("error converting T[] to slice of maps")
		return
	}

	log.With("columns", columns).Debug("generating row keys from data and columns")
	// generate a set of possible key combinations to identify each row grouping
	if len(columns) > 0 {
		rowKeys, _ = tables.PossibleCombinationsAsKeys(records, columns)
	}
	log.Debug("generating skeleton table structure")
	// now create a skeleton table from the rowKeys & date values
	table = tables.Skeleton(rowKeys, dates, emptyVal)
	log.Debug("populating table structure")
	// now populate the table
	table = tables.Populate(records, table, columns, transformCol, valueCol, emptyVal)
	// inject trend column
	log.Debug("adding trend column")
	addTrendCol(table, trendCol)
	// now inject rowTotals
	log.Debug("adding row totals")
	rowTotals(table, dates, totalsCol)
	// now create the footer
	log.Debug("creating column totals for the footer")
	footer = columnTotals(table, append(columns, trendCol), append(dates, totalsCol))
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

func addTrendCol(table map[string]map[string]string, col string) {
	for _, row := range table {
		row[col] = ""
	}
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
