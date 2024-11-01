// Package transformers contains standard methods used to transform
// an api result body into a format used within the front end templates.
package transformers

import (
	"fmt"
	"slices"
)

// KVPair generates a string from the key and value in the form `name:value`.
//
// Used when flatterning a map into a slice so the key of the map
// is not lost and can be parsed later if needed. In particular when
// creating flat table rows from the api data
//
// Allows optional suffixes to be passed which are then directly added
// to the result string.
//
//	KVPair("A", 1)        // "A:1"
//	KVPair("A", 1, "s")   // "A:1s"
func KVPair(key string, value interface{}, suffixes ...string) (s string) {
	s = fmt.Sprintf("%s:%v", key, value)

	for _, suffix := range suffixes {
		s += suffix
	}
	return
}

func RowKV(sortedColumns []string, values map[string]interface{}) (kv string) {
	kv = ""
	for _, column := range sortedColumns {
		var ok bool
		var value interface{}
		if value, ok = values[column]; !ok {
			value = ""
		}
		kv += KVPair(column, value, "^")
	}
	return
}

// SortedColumnNames takes the columnValues (from `.column_values`) and
// returns a list of the column names, sorted
// Used to ensure generation of values using the column values have
// a predictable & consistent order
func SortedColumnNames(columnValues map[string][]string) (sorted []string) {
	sorted = []string{}
	for columnName := range columnValues {
		sorted = append(sorted, columnName)
	}
	slices.Sort(sorted)
	return
}

// ColumnValuesList takes the existing set of columnValues and flattens
// the map into a slice of slices.
// Each top level item from the map becomes its own slice with the returned
// values.
//
// Within each slice the key name of the map is included in the new value
// so the field name is not lost from that data
//
// This is used to create all of the possible table rows that will be needed
// to display the api data.
//
// Input:
//
//	map[string][]string {
//		"A": []string{ "1", "2"},
//		"Z": []string{ "One", "Two"},
//	}
//
// Output:
//
//	[][]string {
//		[]string{ "A:1^", "A:2^"},
//		[]string{ "Z:One^", "Z:Two^"},
//	}
//
// Normally the output is then passed to Permutations
func ColumnValuesList(columnValues map[string][]string) (values [][]string) {
	values = [][]string{}

	for _, columnName := range SortedColumnNames(columnValues) {
		var choices = []string{}
		for _, value := range columnValues[columnName] {
			choices = append(choices, KVPair(columnName, value, "^"))
		}
		values = append(values, choices)
	}
	return
}
