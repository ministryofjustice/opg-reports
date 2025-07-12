package datatable

import "fmt"

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

func New(response ResponseBody) (dt *DataTable, err error) {
	var (
		possibles    []string
		skeleton     map[string]map[string]string
		populated    map[string]map[string]string
		dataRows     []map[string]string = response.DataRows()
		totalCol     string              = response.HasRowTotals()
		trendCol     string              = response.HasTrendColumn()
		dataheaders  []string            = response.DataHeaders()
		identifiers  []string            = response.Identifiers()
		paddedData   []map[string]string = response.PaddedDataRows()
		extraHeaders []string            = []string{}
		extratotals  []string            = identifiers
		sums         []string            = response.SumColumns()
		transform    string              = response.TransformColumn()
		value        string              = response.ValueColumn()
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
