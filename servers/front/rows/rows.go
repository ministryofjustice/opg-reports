package rows

import (
	"bytes"
	"fmt"
)

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
	return fmt.Sprintf("%s_%v", c, v)
}

func Key(item map[string]interface{}, cols []string) string {
	key := ""
	for _, c := range cols {
		key += K(c, item[c]) + "^"
	}
	return key
}

// flatKeys creates the base map using generated keys
func flatKeys(columns map[string][]interface{}) (ks [][]string) {
	ks = [][]string{}

	for col, choices := range columns {
		values := []string{}
		for _, val := range choices {
			values = append(values, K(col, val)+"^")
		}
		ks = append(ks, values)
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
				i[val] = ""
			}
			skel[key][intName] = i
		}
		cols := map[string]string{}
		for c, _ := range columns {
			cols[c] = ""
		}
		skel[key]["columns"] = cols
	}
	return

}

func RowMap(data []interface{}, columns map[string][]interface{}, intervals map[string][]string, values map[string]string) (rows map[string]map[string]interface{}) {

	rows = Skeleton(columns, intervals)
	cols := []string{}
	for col, _ := range columns {
		cols = append(cols, col)
	}

	// loop over each item
	for _, i := range data {
		item := i.(map[string]interface{})
		key := Key(item, cols)
		row := rows[key]

		// cols := row["columns"].(map[string]string)
		// for c, _ := range columns {
		// 	cols[c] = item[c].(string)
		// }

		for intCol, _ := range intervals {
			intervals := row[intCol].(map[string]interface{})
			intervalValue := item[intCol].(string)
			valueKey := values[intCol]
			value := item[valueKey]
			intervals[intervalValue] = value
		}

	}
	return

}
