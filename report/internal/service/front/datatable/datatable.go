package datatable

// DataTable is the end result of transforming a list into a set of rows
// based on grouping and ordering data. This is then used within templates
// to render content in a table
type DataTable struct {
	Body         []map[string]string
	RowHeaders   []string
	DataHeaders  []string
	ExtraHeaders []string
	Footer       map[string]string
	Others       map[string]interface{} // extra information that handles can add to
}
