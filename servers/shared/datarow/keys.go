package datarow

import (
	"fmt"
	"slices"
	"strings"
)

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
func flatKeys(columns map[string][]string) (ks [][]string) {
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
