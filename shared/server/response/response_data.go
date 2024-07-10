package response

// --- CELLS

type ICellName interface {
	SetName(name string)
	GetName() string
}

type ICellValue interface {
	SetValue(value interface{})
	GetValue() interface{}
}

type ICellHeader interface {
	SetIsHeader(is bool)
	GetIsHeader() bool
}
type ICellSupplementary interface {
	SetIsSupplementary(is bool)
	GetIsSupplementary() bool
}

// ICell handles a single peice of information within the table
// structure. It uses Name & Value to represent content and
// some bool flags (IsHeader, IsSupplementary) to give insight
// on what type of cells this may be (header, footer etc)
type ICell interface {
	ICellName
	ICellValue
	ICellHeader
	ICellSupplementary
}

var _ ICell = &Cell{}

// Cell handles a single peice of information within the table
// structure. It uses Name & Value to represent content and
// some bool flags (IsHeader, IsSupplementary) to give insight
// on what type of cells this may be (header, footer etc)
type Cell struct {
	Name            string      `json:"name"`
	Value           interface{} `json:"value"`
	IsHeader        bool        `json:"is_header"`
	IsSupplementary bool        `json:"is_supplementary"`
}

// SetName change the name
func (c *Cell) SetName(name string) {
	c.Name = name
}

// GetName returns the name
func (c *Cell) GetName() string {
	return c.Name
}

// SetValue sets value
func (c *Cell) SetValue(value interface{}) {
	c.Value = value
}

// GetValue returns value
func (c *Cell) GetValue() interface{} {
	return c.Value
}

// SetIsHeader
func (c *Cell) SetIsHeader(h bool) {
	c.IsHeader = h
}

// GetIsHeader
func (c *Cell) GetIsHeader() bool {
	return c.IsHeader
}

// SetIsSupplementary
func (c *Cell) SetIsSupplementary(h bool) {
	c.IsSupplementary = h
}

// GetIsSupplementary
func (c *Cell) GetIsSupplementary() bool {
	return c.IsSupplementary
}

// --- ROWS
// IRowHeader provides interface functions for handling row headers
type IRowHeader[C ICell] interface {
	SetHeaders(cells ...C)
	GetHeaders() []C
	GetHeadersCount() int
}

// IRowData provides interface functions for handling the main body of
// data for this row
type IRowData[C ICell] interface {
	SetData(cells ...C)
	GetData() []C
	GetDataCount() int
}

// IRowSupplementary provides interface functions for handling extra
// columns (generally at the end) in this row - think row totals / averages
type IRowSupplementary[C ICell] interface {
	SetSupplementary(cells ...C)
	GetSupplementary() []C
	GetSupplementaryCount() int
}

// IRow handles a series of table cells and bundles them in a logical
// row structure akin to a spreadsheet or html table.
// To help, the cells within the row are categorised to mark if they
// are a header (like the first cell in a row), supplementary (such
// as an extra column tracking row rotals) or not
type IRow[C ICell] interface {
	IRowHeader[C]
	IRowData[C]
	IRowSupplementary[C]
	SetRaw(cells ...C)
	GetAll() []C
	GetTotalCellCount() int
}

// Row handles a series of table cells and bundles them in a logical
// row structure akin to a spreadsheet or html table.
// To help, the cells within the row are categorised to mark if they
// are a header (like the first cell in a row), supplementary (such
// as an extra column tracking row rotals) or not
type Row[C ICell] struct {
	HeaderCells        []C `json:"headers"`
	DataCells          []C `json:"data"`
	SupplementaryCells []C `json:"supplementary"`
}

// SetHeaders appends the provided cells into the header
// list as long as the cell is marked as being a header
// (via GetIsHeader)
//
// Interface: [IRowHeader]
func (r *Row[C]) SetHeaders(cells ...C) {
	if cells == nil {
		r.HeaderCells = []C{}
	} else {
		for _, h := range cells {
			if h.GetIsHeader() {
				r.HeaderCells = append(r.HeaderCells, h)
			}
		}

	}
}

// GetHeaders returns all the header cells for this row
//
// Interface: [IRowHeader]
func (r *Row[C]) GetHeaders() []C {
	return r.HeaderCells
}

// GetHeadersCount returns a count of how many headers
// there are within this row
//
// Interface: [IRowHeader]
func (r *Row[C]) GetHeadersCount() int {
	if r.HeaderCells == nil {
		return 0
	}
	return len(r.HeaderCells)
}

// SetData appends the cells passed into the current set
// of non-header, non-supplemental cells (think table body)
// as long as that cell is not either
// If the cells passed == nil then then data cells are reset
//
// Interface: [IRowData]
func (r *Row[C]) SetData(cells ...C) {
	if cells == nil {
		r.DataCells = []C{}
	} else {
		for _, d := range cells {
			// not header, not supplementary
			if !d.GetIsHeader() && !d.GetIsSupplementary() {
				r.DataCells = append(r.DataCells, d)
			}
		}
	}
}

// GetData returns all table body data cells
//
// Interface: [IRowData]
func (r *Row[C]) GetData() []C {
	return r.DataCells
}

// GetDataCount returns a count of the table body data cells
//
// Interface: [IRowData]
func (r *Row[C]) GetDataCount() int {
	if r.DataCells == nil {
		return 0
	}
	return len(r.DataCells)
}

// SetSupplementary appends passed sells into the extra data for the row.
// Think of these as extra columns on the table of a spreadhseet, like row totals
// row average etc.
//
// Interface: [IRowSupplementary]
func (r *Row[C]) SetSupplementary(cells ...C) {
	if cells == nil {
		r.SupplementaryCells = []C{}
	} else {
		for _, s := range cells {
			if s.GetIsSupplementary() {
				r.SupplementaryCells = append(r.SupplementaryCells, s)
			}
		}
	}
}

// GetSupplementary returns the extra columns for this row
//
// Interface: [IRowSupplementary]
func (r *Row[C]) GetSupplementary() []C {
	return r.SupplementaryCells
}

// GetSupplementaryCount returns a count of how many extra columns
// are in the row
//
// Interface: [IRowSupplementary]
func (r *Row[C]) GetSupplementaryCount() int {
	if r.SupplementaryCells == nil {
		return 0
	}
	return len(r.SupplementaryCells)
}

// SetRaw checks each of the cells that are passed in and checks their
// status flags to then assign them in to the correct set (head, body, supplemental)
// directly without having to call the other functions. Helps when attaching
// large number of cols or random sets
//
// Interface: [IRow]
func (r *Row[C]) SetRaw(cells ...C) {
	if cells == nil {
		r.DataCells = []C{}
	} else {
		for _, c := range cells {
			if c.GetIsHeader() {
				r.SetHeaders(c)
			} else if c.GetIsSupplementary() {
				r.SetSupplementary(c)
			} else {
				r.SetData(c)
			}
		}
	}
}

// GetAll returns all cells that have been added to this table.
// If does this in order: GetHeaders(). GetData(), GetSupplementary()
//
// Interface: [IRow]
func (r *Row[C]) GetAll() (all []C) {
	all = []C{}
	for _, h := range r.GetHeaders() {
		all = append(all, h)
	}
	for _, d := range r.GetData() {
		all = append(all, d)
	}
	for _, s := range r.GetSupplementary() {
		all = append(all, s)
	}
	return
}

// GetTotalCellCount returns the overal count in this row
//
// Interface: [IRow]
func (r *Row[C]) GetTotalCellCount() int {
	return r.GetHeadersCount() + r.GetDataCount() + r.GetSupplementaryCount()
}

// ---- TABLE

// ITableHead provides functions to get & set the table header
type ITableHead[C ICell, R IRow[C]] interface {
	SetTableHead(row R)
	GetTableHead() R
}

// ITableBody provides functions to get & set the table body content
type ITableBody[C ICell, R IRow[C]] interface {
	SetTableBody(rows ...R)
	GetTableBody() []R
}

// ITableFoot provides functions to get & set the table foot content
type ITableFoot[C ICell, R IRow[C]] interface {
	SetTableFoot(row R)
	GetTableFoot() R
}

// ITable provides an interface for handling response data in a tabular format.
// Applies the same construct as a typical html table
type ITable[C ICell, R IRow[C]] interface {
	ITableHead[C, R]
	ITableBody[C, R]
	ITableFoot[C, R]
}

// Table organises data into typical head, body and foot structure of
// a HTML table and provides methods to get and set each
type Table[C ICell, R IRow[C]] struct {
	Head R   `json:"head"`
	Body []R `json:"body"`
	Foot R   `json:"foot"`
}

// SetTableHead overwrites the table header
func (t *Table[C, R]) SetTableHead(row R) {
	t.Head = row
}

// GetTableHead returns the table head
func (t *Table[C, R]) GetTableHead() R {
	return t.Head
}

// SetTableBody will append the rows passed into the current table body
// dataset. If the rows == nil then the table body is reset to empty
func (t *Table[C, R]) SetTableBody(rows ...R) {
	if rows == nil {
		t.Body = []R{}
	} else {
		t.Body = append(t.Body, rows...)
	}
}

// GetTableBody returns the table body contents
func (t *Table[C, R]) GetTableBody() []R {
	return t.Body
}

// SetTableFoot overwrites the table foot property
func (t *Table[C, R]) SetTableFoot(row R) {
	t.Foot = row
}

// GetTableFoot returns the table footer
func (t *Table[C, R]) GetTableFoot() R {
	return t.Foot
}

// New helpers
// -- CELLS

// NewCell is a standard cell (think general table data)
func NewCell(name string, value interface{}) ICell {
	return &Cell{Name: name, Value: value, IsHeader: false, IsSupplementary: false}
}

// NewCellHeader is a cell that is also a header for the row it is within
// - the first cells of a table body row
func NewCellHeader(name string, value interface{}) ICell {
	return &Cell{Name: name, Value: value, IsHeader: true, IsSupplementary: false}
}

// NewCellExtra is an extra cell that may sit on the end of a row
// such as a row counter / row average
func NewCellExtra(name string, value interface{}) ICell {
	return &Cell{Name: name, Value: value, IsHeader: false, IsSupplementary: true}
}

// -- ROWS

// NewRow creates a new, empty row. If cells are passed, those
// are added via SetRaw
func NewRow[C ICell](cells ...C) (row IRow[C]) {
	row = &Row[C]{}
	row.SetHeaders()
	row.SetSupplementary()
	row.SetData()
	row.SetRaw(cells...)
	return
}

// -- TABLES

// NewTable generates a new Table and sets defaults. Any rows
// passed are assumed to be table body rows and added as such.
func NewTable[C ICell, R IRow[C]](rows ...R) *Table[C, R] {
	h := NewRow[C]()
	f := NewRow[C]()

	td := &Table[C, R]{
		Head: h.(R),
		Foot: f.(R),
	}
	td.SetTableBody(rows...)
	return td
}
