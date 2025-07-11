package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var emptyCell = "0.00"

// TableHeaderRow allow us to create an ordered row of all headers for the table..
//
// You'd merge group name with the date ranges to get the header
func TableHeaderRow(lists ...[]string) (header []string) {
	header = []string{}
	for _, cells := range lists {
		header = append(header, cells...)
	}
	return
}

func CombinationKey(item map[string]string, identifiers ...string) (key string) {
	slices.Sort(identifiers)
	identifiers = slices.Compact(identifiers)

	key = ""
	for _, k := range identifiers {
		key += fmt.Sprintf("%s:%s^", k, item[k])
	}
	return
}

// uniqueValuesForEachIdentifier generates a slice of slices, with each slice representing
// each identifiers unique values that are within the data (where identifier is the map key)
func uniqueValuesForEachIdentifier(data []map[string]string, identifiers ...string) (combinations [][]string) {
	// sort and remove any duplicate keys
	slices.Sort(identifiers)
	identifiers = slices.Compact(identifiers)

	combinations = [][]string{}

	for _, key := range identifiers {
		var options = []string{}

		for _, item := range data {
			if v, ok := item[key]; ok {
				options = append(options, fmt.Sprintf("%s:%s^", key, v))
			}
		}
		// make unique
		slices.Sort(options)
		options = slices.Compact(options)
		combinations = append(combinations, options)
	}

	return
}

// PossibleCombinationsAsKeys takes slice of maps as raw data and a series of keys (`identifiers`)
// and finds
//
//   - All unique values of each identifier within `data` (returned `uniques`)
//   - A flat list of possible combinations of these unique values (returned as `keys`)
//
// Think of the `identifiers` as column headers, where even if a row is missing a value the
// combination will still be included.
//
// Example:
//
//		Input
//			data = []map[string]string {
//				map[string]string{
//					"account": "A"
//					"region": "2024",
//					"cost": "100"
//				},
//				map[string]string{
//					"account": "B"
//					"region": "2025",
//					"cost": "100"
//				},
//				map[string]string{
//					"account": "A"
//					"region": "2024",
//					"cost": "100"
//				},
//			}
//			identifiers = "account", "region"
//	  Output
//			keys = []string{
//				"account:A^region:2024^",
//				"account:A^region:2025^",
//				"account:B^region:2024^",
//				"account:B^region:2025^"
//			}
//			uniques = [][]string{
//				[]string{"account:A^", "account:B^"}
//				[]string{"region:2024^", "region:2025^"}
//			}
//
// The `keys` returned can be used to reform the data grouped by the identifiers.
func PossibleCombinationsAsKeys(data []map[string]string, identifiers []string) (keys []string, uniques [][]string) {
	uniques = [][]string{}

	slices.Sort(identifiers)
	// get all unique values for each data key
	uniques = uniqueValuesForEachIdentifier(data, identifiers...)
	// generate flat list of possible values - these will be used for row keys
	keys = Permutations(uniques...)
	return

}

// SkeletonTable reates a series of skeleton rows from the known
// column values and date range to form a map of table rows.
//
// `keys` is the output from `PossibleCombinationsAsKeys`
// `cells` is normally the date columns
//
// Output:
//
//	map[string]map[string]string{}{
//		"environment:dev^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:dev^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//	}
func SkeletonTable(keys []string, cells []string) (table map[string]map[string]string) {
	table = map[string]map[string]string{}

	for _, key := range keys {
		row := map[string]string{}
		// recreate the column name and value from the formatted key
		// 	- "environment:backup^account:A" => {"environment":"backup", "account":"A"}
		key = strings.TrimSuffix(key, "^")
		for _, columnAndValue := range strings.Split(key, "^") {
			sp := strings.Split(columnAndValue, ":")
			col, val := sp[0], sp[1]
			row[col] = val
		}
		// append the extra cells on to the table
		for _, cell := range cells {
			row[cell] = emptyCell
		}
		table[key+"^"] = row
	}

	return
}

// PopulateTable loops over each data item and inject its value (from valueColumn)
// into the transformed tables colum (transformColumn) based on the items unique
// key (from identifiers)
//
// This generates a full table, with values in every column merged from the raw
// dataset
func PopulateTable(
	data []map[string]string, table map[string]map[string]string,
	identifiers []string,
	transformColumn string, valueColumn string,
) map[string]map[string]string {

	for _, item := range data {
		key := CombinationKey(item, identifiers...)
		column := item[transformColumn]

		if row, ok := table[key]; ok {
			row[column] = item[valueColumn]
		}
	}

	// remove empty rows...
	for id, row := range table {
		var empty = true
		for k, v := range row {
			if !slices.Contains(identifiers, k) && v != emptyCell {
				empty = false
			}
		}
		if empty {
			delete(table, id)
		}
	}

	return table
}

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

// DummyRows generates fake rows - this is normally added before the
// PossibleCombinationsAsKeys calls to make sure tihngs like all dates
// within a given range are present, when there is a chance that
// they might not be a cost for a particular month.
func DummyRows(extras []string, key string) (dummys []map[string]string) {
	dummys = []map[string]string{}
	for _, d := range extras {
		dummys = append(dummys, map[string]string{key: d})
	}
	return
}
