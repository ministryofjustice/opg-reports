package datatable

var emptyCell = "0.00"

type ResponseBody interface {
	DataHeaders() (dh []string)
	DataRows() (data []map[string]string)
	PaddedDataRows() (all []map[string]string)
	Identifiers() (identifiers []string)
	Cells() (cells []string)
	TransformColumn() string
	ValueColumn() string
	SumColumns() (cols []string)
	HasRowTotals() string
	HasTrendColumn() string
}

// DataTable is the end result of transforming a list into a set of rows
// based on grouping and ordering data. This is then used within templates
// to render content in a table
type DataTable struct {
	Body         map[string]map[string]string
	RowHeaders   []string
	DataHeaders  []string
	ExtraHeaders []string
	Header       []string
	Footer       map[string]string
}

func New(response ResponseBody) (dt *DataTable) {
	var (
		totalCol     = response.HasRowTotals()
		trendCol     = response.HasTrendColumn()
		dataheaders  = response.DataHeaders()
		identifiers  = response.Identifiers()
		paddedData   = response.PaddedDataRows()
		possibles, _ = PossibleCombinationsAsKeys(paddedData, identifiers)
		skeleton     = SkeletonTable(possibles, response.Cells())
		populated    = PopulateTable(response.DataRows(),
			skeleton, identifiers, response.TransformColumn(), response.ValueColumn())
		extraHeaders = []string{}
		extratotals  = identifiers
		sums         = response.SumColumns()
	)
	// split logic as it matters where in the order trend and total are injected
	if trendCol != "" {
		extraHeaders = append(extraHeaders, trendCol)
	}
	if totalCol != "" {
		AddRowTotals(populated, identifiers, totalCol)
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
	}
	if len(sums) >= 0 {
		dt.Footer = ColumnTotals(populated, sums, extratotals...)
	}

	return
}
