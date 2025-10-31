package datatable

import "fmt"

var emptyCell = "0.00"

type RowTotalCleaner func(table map[string]map[string]string, identifiers []string, columnName string)

type ResponseBody interface {
	// DataHeaders returns the column headings used for the core data - generally Dates
	DataHeaders() (dh []string)
	// DataRows returns the data from within the api response (.Data from the API)
	//
	// This is used as a base to find all possible combination values and then populate
	// the final table
	DataRows() (data []map[string]string)
	// PaddedRows is used to inject fake rows into the data set to ensure all values
	// of used fields are present.
	// Typically used to add months / dates that might not exist in the data set by default
	PaddedDataRows() (all []map[string]string)
	// Idenfiers returns set of fields that were used to group this data together on the
	// api - these also form the first columns (RowHeaders).
	// These are the non numeric table columns, like Account Name, that we want to be
	// present at the start of the table row, but are not used in calculations.
	Identifiers() (identifiers []string)
	// Cells returns a list of column names that should be added to each row, this is
	// generally used to insert date keys into each row of the table
	Cells() (cells []string)
	// TransformColumn with ValueColumn are used to insert a value from a different field
	// name into a new column - so fetching the value of .Date, using that as the column
	// to then insert .Cost into:
	//
	// 	col := r.TransformColumn()
	//  val := t.ValueColumn()
	// 	row[col] =  row[val]
	TransformColumn() string
	// TransformColumn with ValueColumn are used to insert a value from a different field
	// name into a new column - so fetching the value of .Date, using that as the column
	// to then insert .Cost into:
	//
	// 	col := r.TransformColumn()
	//  val := t.ValueColumn()
	// 	row[col] =  row[val]
	ValueColumn() string
	// RowTotalKeyName returns the name of a key that should be used to
	// store the total value of each row.
	//
	// If it returns "" a row total should not be added
	RowTotalKeyName() string
	// TrendKeyName is like RowTotalKeyName, it returns the name of the
	// field/key to add to a row that will be rendered for displaying
	// trends.
	// If empty / "", then no trend should be added
	TrendKeyName() string
	// SumColumns returns a list of columns/key values that should be used
	// to generate a sum for each column in the table - effectively the
	// overall table total within the tfoot.
	//
	// If empty, dont generate a table footer totals row
	SumColumns() (cols []string)
	// RowTotalCleanup is called to allow data to be modified afterwards; typically
	// used to convert a sum to an average for things like uptime / success rates
	RowTotalCleanup() RowTotalCleaner
}

// DataTable is the end result of transforming a list into a set of rows
// based on grouping and ordering data. This is then used within templates
// to render content in a table
type DataTable struct {
	Body         map[string]map[string]string
	RowHeaders   []string
	DataHeaders  []string
	ExtraHeaders []string
	Footer       map[string]string
	Others       map[string]interface{} // extra information that handles can add to
}

func New(response ResponseBody) (dt *DataTable, err error) {
	var (
		possibles    []string
		skeleton     map[string]map[string]string
		populated    map[string]map[string]string
		dataRows     []map[string]string = response.DataRows()
		totalCol     string              = response.RowTotalKeyName()
		trendCol     string              = response.TrendKeyName()
		dataheaders  []string            = response.DataHeaders()
		identifiers  []string            = response.Identifiers()
		paddedData   []map[string]string = response.PaddedDataRows()
		extraHeaders []string            = []string{}
		extratotals  []string            = identifiers
		sums         []string            = response.SumColumns()
		transform    string              = response.TransformColumn()
		value        string              = response.ValueColumn()
		rowTotalF    RowTotalCleaner     = response.RowTotalCleanup()
	)
	// if there are no idenfifiers, this fails!
	if len(identifiers) <= 0 {
		err = fmt.Errorf("no grouping data found in api result, so cannot create datatable")
		return
	}

	possibles, _ = PossibleCombinationsAsKeys(paddedData, identifiers)
	skeleton = SkeletonTable(possibles, response.Cells())
	populated = PopulateTable(dataRows, skeleton, identifiers, transform, value)

	// split logic as it matters where in the order trend and total are injected
	if trendCol != "" {
		extraHeaders = append(extraHeaders, trendCol)
	}
	if totalCol != "" {
		AddRowTotals(populated, identifiers, totalCol)
		rowTotalF(populated, identifiers, totalCol)
		extraHeaders = append(extraHeaders, totalCol)

	}
	if trendCol != "" {
		AddColumnsToRows(populated, trendCol)
		extratotals = append(extratotals, trendCol)
	}

	dt = &DataTable{
		Body:         populated,
		RowHeaders:   identifiers,
		DataHeaders:  dataheaders,
		ExtraHeaders: extraHeaders,
		Others:       map[string]interface{}{},
	}
	if len(sums) >= 0 {
		dt.Footer = ColumnTotals(populated, sums, extratotals...)
	}

	return
}
