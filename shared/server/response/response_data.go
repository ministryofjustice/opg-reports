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
type ICell interface {
	ICellName
	ICellValue
	ICellHeader
	ICellSupplementary
}

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
type IRowHeader[C ICell] interface {
	SetHeaders(cells ...C)
	GetHeaders() []C
	GetHeadersCount() int
}
type IRowData[C ICell] interface {
	SetData(cells ...C)
	GetData() []C
	GetDataCount() int
}
type IRowSupplementary[C ICell] interface {
	SetSupplementary(cells ...C)
	GetSupplementary() []C
	GetSupplementaryCount() int
}

type IRow[C ICell] interface {
	IRowHeader[C]
	IRowData[C]
	IRowSupplementary[C]
	SetRaw(cells ...C)
	GetRaw() []C
	GetTotalCellCount() int
}

type Row[C ICell] struct {
	HeaderCells        []C `json:"headers"`
	DataCells          []C `json:"data"`
	SupplementaryCells []C `json:"supplementary"`
}

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

func (r *Row[C]) GetHeaders() []C {
	return r.HeaderCells
}

func (r *Row[C]) GetHeadersCount() int {
	if r.HeaderCells == nil {
		return 0
	}
	return len(r.HeaderCells)
}

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

func (r *Row[C]) GetData() []C {
	return r.DataCells
}

func (r *Row[C]) GetDataCount() int {
	if r.DataCells == nil {
		return 0
	}
	return len(r.DataCells)
}

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

func (r *Row[C]) GetSupplementary() []C {
	return r.SupplementaryCells
}

func (r *Row[C]) GetSupplementaryCount() int {
	if r.SupplementaryCells == nil {
		return 0
	}
	return len(r.SupplementaryCells)
}

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

func (r *Row[C]) GetRaw() (all []C) {
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

func (r *Row[C]) GetTotalCellCount() int {
	return r.GetHeadersCount() + r.GetDataCount() + r.GetSupplementaryCount()
}

// ---- TABLE

type ITableHead[C ICell, R IRow[C]] interface {
	SetTableHead(row R)
	GetTableHead() R
}

type ITableBody[C ICell, R IRow[C]] interface {
	SetTableBody(rows ...R)
	GetTableBody() []R
}

type ITableFoot[C ICell, R IRow[C]] interface {
	SetTableFoot(row R)
	GetTableFoot() R
}

type ITable[C ICell, R IRow[C]] interface {
	ITableHead[C, R]
	ITableBody[C, R]
	ITableFoot[C, R]
}

type Table[C ICell, R IRow[C]] struct {
	Head R   `json:"head"`
	Body []R `json:"body"`
	Foot R   `json:"foot"`
}

func (t *Table[C, R]) SetTableHead(row R) {
	t.Head = row
}

func (t *Table[C, R]) GetTableHead() R {
	return t.Head
}

func (t *Table[C, R]) SetTableBody(rows ...R) {
	if rows == nil {
		t.Body = []R{}
	} else {
		t.Body = append(t.Body, rows...)
	}
}

func (t *Table[C, R]) GetTableBody() []R {
	return t.Body
}

func (t *Table[C, R]) SetTableFoot(row R) {
	t.Foot = row
}

func (t *Table[C, R]) GetTableFoot() R {
	return t.Foot
}

// New helpers
// -- CELLS
func NewCell(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value, IsHeader: false, IsSupplementary: false}
}
func NewCellHeader(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value, IsHeader: true, IsSupplementary: false}
}
func NewCellExtra(name string, value interface{}) *Cell {
	return &Cell{Name: name, Value: value, IsHeader: false, IsSupplementary: true}
}

// -- ROWS
func NewRow[C ICell](cells ...C) (row *Row[C]) {
	row = &Row[C]{}
	row.SetHeaders()
	row.SetSupplementary()
	row.SetData()
	row.SetRaw(cells...)
	return
}

// -- TABLES
func NewTable[C ICell, R IRow[C]](rows ...R) *Table[C, R] {
	td := &Table[C, R]{}
	td.SetTableBody(rows...)
	return td
}
