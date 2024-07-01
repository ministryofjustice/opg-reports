package response

type ICell interface {
	SetName(name string)
	GetName() string
}
type IRow[C ICell] interface {
	SetCells(cells []C)
	AddCells(cells ...C)
	GetCells() []C
}

type IHeadings[C ICell, R IRow[C]] interface {
	SetHeadings(h R)
	GetHeadings() R
	SetHeadingsCounters(pre int, post int)
	GetHeadingsCounters() (pre int, post int)
}

type IFooter[C ICell, R IRow[C]] interface {
	SetFooter(f R)
	GetFooter() R
	SetFooterCounters(pre int, post int)
	GetFooterCounters() (pre int, post int)
}

type ITableData[C ICell, R IRow[C]] interface {
	IHeadings[C, R]
	IFooter[C, R]
	SetRows(rows []R)
	AddRows(rows ...R)
	GetRows() []R
}

// Cell represents a single cell of data - like a spreadsheet
// Used as part of a row
// Impliments [ICell]
type Cell struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
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
func (c *Cell) SetValue(value string) {
	c.Value = value
}

// GetValue returns value
func (c *Cell) GetValue() interface{} {
	return c.Value
}

// Row impliments [IRow]
// Acts like a row in a table / spreadsheet
type Row[C ICell] struct {
	Cells []C `json:"cells"`
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

// TableHeadings represents specifc headers for the tabular data
type TableHeadings[C ICell, R IRow[C]] struct {
	PreHeadings  int `json:"pre_headings"`
	PostHeadings int `json:"post_headings"`
	Headings     R   `json:"headings,omitempty"`
}

// SetHeadings sets the header
func (d *TableHeadings[C, R]) SetHeadings(h R) {
	d.Headings = h
}

// GetHeadings returns them
func (d *TableHeadings[C, R]) GetHeadings() (head R) {
	if len(d.Headings.GetCells()) > 0 {
		head = d.Headings
	}
	return
}
func (d *TableHeadings[C, R]) SetHeadingsCounters(pre int, post int) {
	d.PreHeadings = pre
	d.PostHeadings = post
}
func (d *TableHeadings[C, R]) GetHeadingsCounters() (pre int, post int) {
	pre = d.PreHeadings
	post = d.PostHeadings
	return
}

// TableFooter represents footer of a table, used for totals etc
type TableFooter[C ICell, R IRow[C]] struct {
	Footer     R   `json:"footer,omitempty"`
	PreFooter  int `json:"pre_footer"`
	PostFooter int `json:"post_footer"`
}

func (d *TableFooter[C, R]) SetFooter(h R) {
	d.Footer = h
}

func (d *TableFooter[C, R]) GetFooter() R {
	return d.Footer
}
func (d *TableFooter[C, R]) SetFooterCounters(pre int, post int) {
	d.PreFooter = pre
	d.PostFooter = post
}
func (d *TableFooter[C, R]) GetFooterCounters() (pre int, post int) {
	pre = d.PreFooter
	post = d.PostFooter
	return
}

// TableData impliments [ITableData]
// Acts like a table
type TableData[C ICell, R IRow[C]] struct {
	*TableHeadings[C, R]
	*TableFooter[C, R]
	Rows []R `json:"rows,omitempty"`
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

// Overwrite the get headings to check the pointer
func (d *TableData[C, R]) GetHeadings() (head R) {
	if d.TableHeadings != nil {
		head = d.TableHeadings.GetHeadings()
	}
	return
}
func (d *TableData[C, R]) GetFooter() (f R) {
	if d.TableFooter != nil {
		f = d.TableFooter.GetFooter()
	}
	return
}

// NewCell returns an ICell
func NewCell(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value}
}

// NewRow returns a single row
func NewRow[C ICell](cells ...C) (row *Row[C]) {
	row = &Row[C]{}
	for _, cell := range cells {
		row.AddCells(cell)
	}
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
	d := &TableData[C, R]{
		TableHeadings: &TableHeadings[C, R]{},
		TableFooter:   &TableFooter[C, R]{},
	}
	d.SetRows(rows)
	return d
}
