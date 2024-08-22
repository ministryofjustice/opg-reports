package rows

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

const emptyVal float64 = 0.0

// permuteStrings generates all possible combinations between the string slices passed
func permuteStrings(parts ...[]string) (ret []string) {
	{
		var n = 1
		for _, ar := range parts {
			n *= len(ar)
		}
		ret = make([]string, 0, n)
	}
	var at = make([]int, len(parts))
	var buf bytes.Buffer
loop:
	for {
		// increment position counters
		for i := len(parts) - 1; i >= 0; i-- {
			if at[i] > 0 && at[i] >= len(parts[i]) {
				if i == 0 || (i == 1 && at[i-1] == len(parts[0])-1) {
					break loop
				}
				at[i] = 0
				at[i-1]++
			}
		}
		// construct permutated string
		buf.Reset()
		for i, ar := range parts {
			var p = at[i]
			if p >= 0 && p < len(ar) {
				buf.WriteString(ar[p])
			}
		}
		ret = append(ret, buf.String())
		at[len(parts)-1]++
	}
	return ret
}

// K helper to generate a k with known dividers
func K(c string, v interface{}) string {
	return fmt.Sprintf("%s:%v", c, v)
}

func Key(item map[string]interface{}, cols []string) string {
	key := ""
	slices.Sort(cols)
	for _, c := range cols {
		key += K(c, item[c]) + "^"
	}
	return key
}

// flatKeys creates the base map using generated keys
func flatKeys(columns map[string][]interface{}) (ks [][]string) {
	ks = [][]string{}
	cols := []string{}
	for col, _ := range columns {
		cols = append(cols, col)
	}
	slices.Sort(cols)

	for _, col := range cols {
		values := []string{}
		choices := columns[col]
		for _, val := range choices {
			values = append(values, K(col, val)+"^")
		}
		ks = append(ks, values)
	}

	return
}

// splitKey recreates the column & value pairs from the key string
func splitKey(key string) (cols map[string]string) {
	cols = map[string]string{}

	list := strings.Split(key, "^")
	for _, item := range list {
		if again := strings.Split(item, ":"); len(again) == 2 {
			k := again[0]
			v := again[1]
			cols[k] = v
		}
	}
	return
}

// Skeleton creates map of empty values for all possible values based on the settings passed
// eahc key having full set of intervals
func Skeleton(columns map[string][]interface{}, intervals map[string][]string) (skel map[string]map[string]interface{}) {
	skel = map[string]map[string]interface{}{}

	keys := permuteStrings(flatKeys(columns)...)
	for _, key := range keys {
		if _, ok := skel[key]; !ok {
			skel[key] = map[string]interface{}{}
		}
		for intName, values := range intervals {
			i := map[string]interface{}{}
			for _, val := range values {
				i[val] = emptyVal
			}
			skel[key][intName] = i
		}

		skel[key]["columns"] = splitKey(key)
	}
	return

}

// DataToRow converts the api data and using the columns and intervals generates a skeleton struct of all possible records and
// then fills in the interval values where data is found
// This resulting map can be iterated over as if each item is a row in the table
func DataToRows(data []interface{}, columns map[string][]interface{}, intervals map[string][]string, values map[string]string) (rows map[string]map[string]interface{}) {

	rows = Skeleton(columns, intervals)

	cols := []string{}
	for col, _ := range columns {
		cols = append(cols, col)
	}
	slices.Sort(cols)

	keyGroup := map[string][]map[string]interface{}{}
	// -- generate a list of all items from the keys
	for _, entry := range data {
		mapped := entry.(map[string]interface{})
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
func DataRows(data []interface{}, columns map[string][]interface{}, intervals map[string][]string, values map[string]string) (rows map[string]map[string]interface{}) {
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
