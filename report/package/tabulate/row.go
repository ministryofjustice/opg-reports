package tabulate

import (
	"fmt"
	"strings"
)

// ColType used as enum constraint
type ColType string

// types of headers we'd use in a table
const (
	KEY   ColType = "labels"
	DATA  ColType = "data"
	EXTRA ColType = "extra"
	END   ColType = "end"
)

type RowEndFunc func(tableRow map[string]interface{}, headings map[ColType][]string)

// PopulateRow adds fresh data to the `dest` from the `src` map, overwriting value
// with the version found in `src`. It will set `Key` fields as well as `Data`
// fields.
//
// Assumes all columns names are strings and values are floats.
//
// Example:
//
//	key := src[ColumnKey].(string)
//	dest[key] = src[ValueKey]
//
// Notes:
//
//	`ColumnKey` should be the field which contains the name of the column to write to in `dest`
//	`ValueKey` is the `src` field that contains th value to write to `dest`
func PopulateRow(src map[string]interface{}, dest map[string]interface{}, headings map[ColType][]string, opts *Args) {
	var key = src[opts.ColumnKey].(string)
	// set this key
	if val, ok := src[opts.ValueKey]; ok && val != nil {
		dest[key] = val
	}
	// update the row heaadings with values
	for _, h := range headings[KEY] {
		if v, ok := src[h]; ok && dest[h].(string) == "" {
			dest[h] = v
		}
	}
}

// RowKey returns a unique string to use as a reference for this row. This is
// used as a key within a map or identifier to group rows by
func RowKey(src map[string]interface{}, headings map[ColType][]string) (key string) {
	key = ""
	for _, k := range headings[KEY] {
		if src[k] == nil {
			panic("row key looking for missing field: " + k)
		}
		key += fmt.Sprintf("%s=%s^", k, strings.ToLower(src[k].(string)))
	}
	return
}

// Empty generates an empty row with all keys present and set to default values
func EmptyRow(headings map[ColType][]string) (row map[string]interface{}) {
	var defaultF float64 = 0.0
	var defaultS string = ""

	row = map[string]interface{}{}

	for t, set := range headings {
		var def interface{} = defaultF
		if t == KEY || t == EXTRA {
			def = defaultS
		}
		for _, col := range set {
			row[col] = def
		}
	}

	return
}

// RowEnd runs the row end function (like adding total / average)
func RowEnd(tableMap map[string]map[string]interface{}, headings map[ColType][]string, rowF RowEndFunc) {
	if rowF == nil {
		return
	}
	for _, row := range tableMap {
		rowF(row, headings)
	}
}

// Average works on the row after its been completely set and adds the average to the end of the row
func RowAverageF(row map[string]interface{}, headings map[ColType][]string) {
	var (
		dataCols []string = headings[DATA]
		endCol   []string = headings[END]
		total    float64  = 0.0
		average  float64  = 0.0
		count    int      = 0
	)
	for _, col := range dataCols {
		var val = Value[float64](col, 0.0, row) //rowV[float64](row, col)
		total += val
		if val != 0.0 {
			count++
		}
	}
	average = total / float64(count)
	if len(endCol) > 0 {
		row[endCol[0]] = average
	}
}

// Total works on the row after its been completely set and adds all data columns together
func RowTotalF(row map[string]interface{}, headings map[ColType][]string) {
	var (
		endCol []string = headings[END]
		total  float64  = 0.0
	)
	for _, col := range headings[DATA] {
		var val = Value[float64](col, 0.0, row)
		total += val
	}
	if len(endCol) > 0 {
		row[endCol[0]] = total
	}
}

// DiffF only works with 2 columns, used to create the diff between the values in each
func RowDiffF(row map[string]interface{}, headings map[ColType][]string) {
	var (
		last, first string
		endCol      []string = headings[END]
		dataCols    []string = headings[DATA]
		count       int      = len(dataCols)
	)
	if count == 2 {
		last = dataCols[1]
		first = dataCols[0]
		if len(endCol) > 0 {
			row[endCol[0]] = Value[float64](last, 0.0, row) - Value[float64](first, 0.0, row)
		}
	}
}

// Value will return the value from the row, or the default value of the header
func Value[T any](h string, def T, src map[string]interface{}) (val T) {
	val = def
	if v, ok := src[h]; ok && v != nil {
		val = v.(T)
	}
	return
}
