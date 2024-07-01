package response

type ICell interface {
	SetName(name string)
	GetName() string
	SetValue(v interface{})
	GetValue() interface{}
	SetIsHeader(h bool)
	GetIsHeader() bool
}

// Cell represents a single cell of data - like a spreadsheet
// Used as part of a row
// Impliments [ICell]
type Cell struct {
	Name     string      `json:"name"`
	IsHeader bool        `json:"header"`
	Value    interface{} `json:"value"`
}

// SetName change the name
func (c *Cell) SetName(name string) {
	c.Name = name
}

// GetName returns the name
func (c *Cell) GetName() string {
	return c.Name
}

// SetIsHeader
func (c *Cell) SetIsHeader(h bool) {
	c.IsHeader = h
}

// GetIsHeader
func (c *Cell) GetIsHeader() bool {
	return c.IsHeader
}

// SetValue sets value
func (c *Cell) SetValue(value interface{}) {
	c.Value = value
}

// GetValue returns value
func (c *Cell) GetValue() interface{} {
	return c.Value
}

type IRow[C ICell] interface {
	SetCells(cells []C)
	AddCells(cells ...C)
	GetCells() []C
	Len() int
	// Heading counters for this row
	UpdateCounters()
	GetCounters() (prefix int, suffix int)
}

// Row impliments [IRow]
// Acts like a row in a table / spreadsheet
type Row[C ICell] struct {
	Cells []C `json:"cells"`
	Pre   int `json:"prefix_header_count"`
	Post  int `json:"suffix_header_count"`
}

func (r *Row[C]) UpdateCounters() {
	r.setPreCounter()
	r.setPostCounter()
}
func (r *Row[C]) GetCounters() (prefix int, suffix int) {
	return r.Pre, r.Post
}
func (r *Row[C]) Len() int {
	return len(r.Cells)
}

// SetCells attach cells to this row
func (r *Row[C]) SetCells(cells []C) {
	r.Cells = cells
}

func (r *Row[C]) AddCells(cells ...C) {
	r.Cells = append(r.Cells, cells...)
}

// GetCells return cells of this row
func (r *Row[C]) GetCells() (cells []C) {
	if len(r.Cells) > 0 {
		cells = r.Cells
	}
	return
}

func (r *Row[C]) setPostCounter() {
	r.Post = 0
	cells := r.GetCells()

	for i := len(cells) - 1; i >= 0; i-- {
		c := cells[i]
		if c.GetIsHeader() {
			r.Post += 1
		} else {
			return
		}
	}
}
func (r *Row[C]) setPreCounter() {
	r.Pre = 0
	for _, c := range r.GetCells() {
		if c.GetIsHeader() {
			r.Pre += 1
		} else {
			return
		}
	}
}

type ITableDataWithHeader[C ICell, R IRow[C]] interface {
	SetHeader(row R)
	GetHeader() (row R)
}
type ITableDataWithFooter[C ICell, R IRow[C]] interface {
	SetFooter(row R)
	GetFooter() (row R)
}

type ITableData[C ICell, R IRow[C]] interface {
	ITableDataWithHeader[C, R]
	ITableDataWithFooter[C, R]

	SetRows(rows []R)
	AddRows(rows ...R)
	GetRows() []R
}

// TableData impliments [ITableData]
// Acts like a table
type TableData[C ICell, R IRow[C]] struct {
	Header R   `json:"header"`
	Footer R   `json:"footer"`
	Rows   []R `json:"rows,omitempty"`
}

func (d *TableData[C, R]) SetHeader(row R) {
	row.UpdateCounters()
	d.Header = row
}
func (d *TableData[C, R]) GetHeader() (row R) {
	row = d.Header
	return
}
func (d *TableData[C, R]) SetFooter(row R) {
	row.UpdateCounters()
	d.Footer = row
}
func (d *TableData[C, R]) GetFooter() (row R) {
	row = d.Footer
	return
}

// SetRows sets rows
func (d *TableData[C, R]) SetRows(rows []R) {
	d.Rows = rows
}
func (d *TableData[C, R]) AddRows(rows ...R) {
	d.Rows = append(d.Rows, rows...)
}

// GetRows returns the rows
func (d *TableData[C, R]) GetRows() (rows []R) {
	if len(d.Rows) > 0 {
		rows = d.Rows
	}
	return
}

// NewCell returns an ICell
func NewCell(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value}
}
func NewHeaderCell(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value, IsHeader: true}
}

// NewRow returns a single row
func NewRow[C ICell](cells ...C) (row *Row[C]) {
	row = &Row[C]{}
	row.AddCells(cells...)
	return
}

// NewRows returns multiple rows
func NewRows[C ICell](cellSet ...[]C) (rows []*Row[C]) {
	rows = []*Row[C]{}
	for _, cells := range cellSet {
		rows = append(rows, NewRow(cells...))
	}
	return
}

// NewData returns a data item acts like a table
func NewData[C ICell, R IRow[C]](rows ...R) *TableData[C, R] {
	d := &TableData[C, R]{}
	d.SetRows(rows)
	return d
}
