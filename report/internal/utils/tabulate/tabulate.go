package tabulate

type TableHeaders interface {
	Textual() (list []string)
	Numeric() (list []string)
	DataColumns() (list []string)
}

type TableRowKeyFunc func(r map[string]interface{}) string
type TableRowUpdateFunc func(dbRow map[string]interface{}, tableRow map[string]interface{}, headers TableHeaders) (updatedTableRow map[string]interface{})

type TabulateOptions struct {
	Headers TableHeaders
	KeyF    TableRowKeyFunc    // used to generate the key for the table row, from the raw db data
	ColumnF TableRowUpdateFunc // used to update the table row - generally updates the column values, run on every loop
	LabelF  TableRowUpdateFunc // used to update the table row - generally sets the row labels; run only once per row at the end
	RowEndF TableRowUpdateFunc // used to handle row totals - run after only once per row after all data is set
}

// Tabulate converts a set of database rows into a table row structure:
func Tabulate(dbRows []map[string]interface{}, opts *TabulateOptions) (tableRows map[string]map[string]interface{}) {
	var done map[string]bool = map[string]bool{}

	tableRows = map[string]map[string]interface{}{}
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
