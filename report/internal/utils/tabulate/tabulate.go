package tabulate

type TableHeaders interface {
	Textual() (list []string)
	Numeric() (list []string)
	DataColumns() (list []string)
}

type TableKeyFunc func(r map[string]interface{}) string
type TableRowFunc func(dbRow map[string]interface{}, tableRow map[string]interface{}, headers TableHeaders) (updatedTableRow map[string]interface{})
type TableSortFunc func(table []map[string]interface{}, headers TableHeaders) (updatedTable []map[string]interface{})
type TableSummaryFunc func(table []map[string]interface{}, headers TableHeaders) (updatedTable []map[string]interface{})

type TabulateOptions struct {
	Headers    TableHeaders
	KeyF       TableKeyFunc  // used to generate the key for the table row, from the raw db data
	ColumnF    TableRowFunc  // used to update the table row - generally updates the column values, run on every loop
	LabelF     TableRowFunc  // used to update the table row - generally sets the row labels; run only once per row at the end
	RowEndF    TableRowFunc  // used to handle row totals - run after only once per row after all data is set
	TableSortF TableSortFunc // used to sort the table in a set order
	TableEndF  TableSummaryFunc
}

// Tabulate converts a set of database rows into a table row structure:
func Tabulate(dbRows []map[string]interface{}, opts *TabulateOptions) (table []map[string]interface{}) {
	var done map[string]bool = map[string]bool{}
	var tableRows = map[string]map[string]interface{}{}
	// add all column values into the table row setup
	for _, dbRow := range dbRows {
		var key = opts.KeyF(dbRow)
		var tableRow, ok = tableRows[key]
		if !ok {
			tableRow = SkeletonRow(opts.Headers)
		}
		tableRow = opts.ColumnF(dbRow, tableRow, opts.Headers)
		tableRows[key] = tableRow
	}
	// after adding all values in, loop over and handle labels & totals
	for _, dbRow := range dbRows {
		var key = opts.KeyF(dbRow)
		var tableRow, ok = tableRows[key]

		if _, set := done[key]; !set && ok {
			done[key] = true
			tableRow = opts.LabelF(dbRow, tableRow, opts.Headers)
			tableRow = opts.RowEndF(dbRow, tableRow, opts.Headers)
		}
	}
	// now convert to a slice
	for _, row := range tableRows {
		table = append(table, row)
	}
	// now sort table by team
	table = opts.TableSortF(table, opts.Headers)
	// now process the table totals
	table = opts.TableEndF(table, opts.Headers)
	return

}

// SkeletonRow generates row for a table with empty string / float values
func SkeletonRow(th TableHeaders) (row map[string]interface{}) {
	row = map[string]interface{}{}
	for _, txt := range th.Textual() {
		row[txt] = ""
	}
	for _, num := range th.Numeric() {
		row[num] = 0.0
	}
	return
}
