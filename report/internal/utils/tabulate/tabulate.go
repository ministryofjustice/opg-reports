package tabulate

import (
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tabulate/rows"
	"sort"
)

type RowEndFunc func(tableRow map[string]interface{}, headings *headers.Headers)
type TableEndFunc func(table []map[string]interface{}, headings *headers.Headers) []map[string]interface{}

type Options struct {
	ColumnKey    string // population rows, proxy for rows.PopulateOptions
	ValueKey     string // populating rows, proxy for rows.PopulateOptions
	SortByColumn string // colum to sort the table data on

	RowEndF   RowEndFunc   // run this function against each row of the table at the end
	TableEndF TableEndFunc // runs against the completed table data
}

func Tabulate[T int | float64 | string](databaseRows []map[string]interface{}, headings *headers.Headers, opts *Options) (table []map[string]interface{}) {

	var tableMap = map[string]map[string]interface{}{}
	// generate the table
	for _, src := range databaseRows {
		var rowKey = rows.Key(src, headings)
		var dest, ok = tableMap[rowKey]
		if !ok {
			dest = rows.Empty(headings)
		}
		// populate the destination row with data from the src db record
		rows.Populate(src, dest, headings, &rows.Options{
			ColumnKey: opts.ColumnKey,
			ValueKey:  opts.ValueKey,
		})
		// set table value
		tableMap[rowKey] = dest
	}
	// now add in row end data, if function is set
	if opts.RowEndF != nil {
		for _, row := range tableMap {
			opts.RowEndF(row, headings)
		}
	}
	// convert to slice
	for _, row := range tableMap {
		table = append(table, row)
	}
	// run sort if required
	if opts.SortByColumn != "" {
		SortDescending[T](table, headings, opts.SortByColumn)
	}
	// now add table summary details if set
	if opts.TableEndF != nil {
		table = opts.TableEndF(table, headings)
	}
	return
}

// SortDescending sorts the table slice by the column set
func SortDescending[T int | float64 | string](table []map[string]interface{}, headings *headers.Headers, sortColumn string) {

	sort.Slice(table, func(i, j int) bool {
		var a = table[i][sortColumn].(T)
		var b = table[j][sortColumn].(T)
		return (a > b)
	})
}

// TotalF generates a table summary row contains the totals of each data column combined
func TotalF(table []map[string]interface{}, headings *headers.Headers) []map[string]interface{} {
	var (
		summary    map[string]interface{} = rows.Empty(headings)
		endCol     *headers.Header        = headings.End()
		firstCol   *headers.Header        = headings.First()
		dataCols   []*headers.Header      = headings.Data()
		tableTotal float64                = 0.0
	)

	// if endCol == nil || firstCol == nil {
	// 	return table
	// }

	for _, row := range table {
		for _, col := range dataCols {
			summary[col.Field] = summary[col.Field].(float64) + row[col.Field].(float64)
		}
		tableTotal += row[endCol.Field].(float64)
	}
	// give the first column a name
	summary[firstCol.Field] = endCol.Field
	summary[endCol.Field] = tableTotal
	table = append(table, summary)
	return table
}

// AverageF generates a table summary row contains the totals averages of each data column combined
func AverageF(table []map[string]interface{}, headings *headers.Headers) []map[string]interface{} {
	var (
		summary    map[string]interface{} = rows.Empty(headings)
		endCol     *headers.Header        = headings.End()
		firstCol   *headers.Header        = headings.First()
		dataCols   []*headers.Header      = headings.Data()
		count      int                    = len(table)
		tableTotal float64                = 0.0
	)
	// if endCol == nil || firstCol == nil {
	// 	return table
	// }

	for _, row := range table {
		for _, col := range dataCols {
			summary[col.Field] = summary[col.Field].(float64) + row[col.Field].(float64)
		}
		tableTotal += row[endCol.Field].(float64)
	}
	// now divide to fix create average
	for _, col := range dataCols {
		summary[col.Field] = summary[col.Field].(float64) / float64(count)
	}
	// give the first column a name
	summary[firstCol.Field] = endCol.Field
	// work out overall average
	summary[endCol.Field] = tableTotal / float64(count)
	table = append(table, summary)
	return table
}
