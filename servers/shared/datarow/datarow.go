package datarow

import "slices"

const emptyVal float64 = 0.0

// DataToRow converts the api data and using the columns and intervals generates a skeleton struct of all possible records and
// then fills in the interval values where data is found
// This resulting map can be iterated over as if each item is a row in the table
func DataToRows(data []map[string]interface{}, columns map[string][]string, intervals map[string][]string, values map[string]string) (rows map[string]map[string]interface{}) {

	rows = Skeleton(columns, intervals)

	cols := []string{}
	for col, _ := range columns {
		cols = append(cols, col)
	}
	slices.Sort(cols)

	keyGroup := map[string][]map[string]interface{}{}
	// -- generate a list of all items from the keys
	for _, mapped := range data {
		key := Key(mapped, cols)
		if _, ok := keyGroup[key]; !ok {
			keyGroup[key] = []map[string]interface{}{}
		}
		keyGroup[key] = append(keyGroup[key], mapped)
	}
	// -- loop over the grouped items and update the values in the intervals
	for key, items := range keyGroup {
		row := rows[key]
		for _, item := range items {
			for interval, _ := range intervals {
				rowSet := row[interval]
				rowInterval := rowSet.(map[string]interface{})
				internvalKey := item[interval].(string)
				valueKey := values[interval]
				value := item[valueKey]
				rowInterval[internvalKey] = value

			}
		}
	}

	return

}

// DataRows converts the raw data into row structure (via DataToRows) then removes items
// with empty data, to reduce the size of data sets for display etc
func DataRows(data []map[string]interface{}, columns map[string][]string, intervals map[string][]string, values map[string]string) (rows map[string]map[string]interface{}) {
	rows = DataToRows(data, columns, intervals, values)
	// now we trim off any fully empty versions
	for key, data := range rows {
		empty := true
		for seg, values := range data {
			if seg != "columns" {
				for _, v := range values.(map[string]interface{}) {
					if v != emptyVal {
						empty = false
					}
				}
			}
		}
		if empty {
			delete(rows, key)
		}
	}
	return
}
