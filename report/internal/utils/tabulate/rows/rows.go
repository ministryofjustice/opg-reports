package rows

import (
	"fmt"
	"opg-reports/report/internal/utils/tabulate/headers"
	"strings"
)

// Empty generates an empty row with all keys present and set to default values
func Empty(headings *headers.Headers) (row map[string]interface{}) {
	row = map[string]interface{}{}
	for _, h := range headings.All() {
		row[h.Field] = h.Default
	}
	return
}

// RowKey returns a unique string to use as a reference for this row. This is
// used as a key within a map or identifier to group rows by
func Key(src map[string]interface{}, headings *headers.Headers) (key string) {
	key = ""
	for _, k := range headings.Keys() {
		if src[k.Field] == nil {
			panic("row key looking for missing field: " + k.Field)
		}
		key += fmt.Sprintf("%s=%s^", k.Field, strings.ToLower(src[k.Field].(string)))
	}
	return
}

type Options struct {
	ColumnKey string
	ValueKey  string
}

// Populate adds fresh data to the `dest` from the `src` map, overwriting value
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
func Populate(src map[string]interface{}, dest map[string]interface{}, headings *headers.Headers, opts *Options) {
	var key = src[opts.ColumnKey].(string)
	var val = src[opts.ValueKey]
	dest[key] = val
	// update the row heaadings with values
	for _, h := range headings.Keys() {
		if v, ok := src[h.Field]; ok && dest[h.Field].(string) == "" {
			dest[h.Field] = v
		}
	}
}

// Average works on the row after its been completely set and adds the average to the end of the row
func AverageF(row map[string]interface{}, headings *headers.Headers) {
	var (
		dataCols []*headers.Header = headings.Data()
		endCol   *headers.Header   = headings.End()
		total    float64           = 0.0
		average  float64           = 0.0
		count    int               = 0
	)
	for _, col := range dataCols {
		var val = headers.Value[float64](col, row) //rowV[float64](row, col)
		total += val
		if val != 0.0 {
			count++
		}
	}
	average = total / float64(count)
	row[endCol.Field] = average
}

// Total works on the row after its been completely set and adds all data columns together
func TotalF(row map[string]interface{}, headings *headers.Headers) {
	var (
		endCol *headers.Header = headings.End()
		total  float64         = 0.0
	)
	for _, col := range headings.Data() {
		var val = headers.Value[float64](col, row)
		total += val
	}
	row[endCol.Field] = total
}

// DiffF only works with 2 columns, used to create the diff between the values in each
func DiffF(row map[string]interface{}, headings *headers.Headers) {
	var (
		last, first, endCol *headers.Header
		dataCols            []*headers.Header = headings.Data()
		count               int               = len(dataCols)
	)
	endCol = headings.End()

	if count == 2 {
		last = dataCols[1]
		first = dataCols[0]
		row[endCol.Field] = headers.Value[float64](last, row) - headers.Value[float64](first, row)
	}
}
