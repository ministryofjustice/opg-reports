package response

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Cell represents a single cell of data - like a spreadsheet
// Used as part of a row
// Impliments [ICell]
type Cell struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
func (c *Cell) GetValue() string {
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
func (r *Row[C]) GetCells() []C {
	return r.Cells
}

// TableData impliments [ITableData]
// Acts like a table
type TableData[C ICell, R IRow[C]] struct {
	Headings R   `json:"headings"`
	Rows     []R `json:"rows"`
}

// SetRows sets rows
func (d *TableData[C, R]) SetRows(rows []R) {
	d.Rows = rows
}
func (d *TableData[C, R]) AddRows(rows ...R) {
	d.Rows = append(d.Rows, rows...)
}

// GetRows returns the rows
func (d *TableData[C, R]) GetRows() []R {
	return d.Rows
}

func (d *TableData[C, R]) SetHeadings(h R) {
	d.Headings = h
}

func (d *TableData[C, R]) GetHeadings() R {
	return d.Headings
}

// NewCell returns an ICell
func NewCell(name string, value string) *Cell {
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
	d := &TableData[C, R]{}
	d.SetRows(rows)
	return d
}

// Timings impliments [ITimings]
type Timings struct {
	Times struct {
		Start    time.Time     `json:"start"`
		End      time.Time     `json:"end"`
		Duration time.Duration `json:"duration"`
	} `json:"timings"`
}

// Start tracks the start time of this request
func (i *Timings) Start() {
	i.Times.Start = time.Now().UTC()
}

// End tracks the end time and the duration of the request
func (i *Timings) End() {
	i.Times.End = time.Now().UTC()
	i.Times.Duration = i.Times.End.Sub(i.Times.Start)
}

// Status impliments [IStatus]
// Provides http status tracking
type Status struct {
	Code int `json:"status"`
}

// SetStatus updates the status field
func (i *Status) SetStatus(status int) {
	i.Code = status
}

// GetStatus returns the status field
func (i *Status) GetStatus() int {
	return i.Code
}

// Errors handles error tracking and impliments [IErrors]
type Errors struct {
	Errs []error `json:"errors"`
}

// SetErrors replaces errors with those passed
func (r *Errors) SetErrors(errors []error) {
	r.Errs = errors
}

// AddError add a new error to the list
func (r *Errors) AddError(err error) {
	r.Errs = append(r.Errs, err)
}

// GetErrors returns all errors
func (r *Errors) GetErrors() []error {
	return r.Errs
}

// Base impliments IBase
// Would be used for a simple endpoint that doesn't return data,
// such as an api root
type Base struct {
	Timings
	Status
	Errors
}

// AddErrorWithStatus adds an error and updates the status at the same time.
// Helpful when validating fields to do both at once.
func (i *Base) AddErrorWithStatus(err error, status int) {
	i.AddError(err)
	i.SetStatus(status)
}

// Result impliments [IResult].
// It allows a response to return with variable (C) data type. This is currently
// constrained to map[string]R, map[string][]R and []R.
// This means various enpoints can return differing ways collecting the data.IEntry
// so some can group by a field or just list everything that matches
//
// This struct and interface allows you to easily decode a response as long as you know
// its return type
type Result[C ICell, R IRow[C], D ITableData[C, R]] struct {
	Base
	Res D `json:"result"`
}

// SetResult updates the internal result data
func (i *Result[C, R, D]) SetResult(result D) {
	i.Res = result
}

// GetResult returns the result
func (i *Result[C, R, D]) GetResult() D {
	return i.Res
}

// NewSimpleResult returns a fresh Base with
// status set as OK and errors empty
func NewSimpleResult() *Base {
	return &Base{
		Timings: Timings{},
		Status:  Status{Code: http.StatusOK},
		Errors:  Errors{Errs: []error{}},
	}
}

func NewResponse() *Result[*Cell, *Row[*Cell], *TableData[*Cell, *Row[*Cell]]] {
	return &Result[*Cell, *Row[*Cell], *TableData[*Cell, *Row[*Cell]]]{
		Base: *NewSimpleResult(),
	}
}

func NewResponseFromJson[C ICell, R IRow[C], D ITableData[C, R]](content []byte, i *Result[C, R, D]) (response *Result[C, R, D], err error) {
	err = json.Unmarshal(content, i)
	return i, err
}
func NewResponseFromHttp[C ICell, R IRow[C], D ITableData[C, R]](r *http.Response, i *Result[C, R, D]) (response *Result[C, R, D], err error) {
	_, by := Stringify(r)
	return NewResponseFromJson(by, i)
}

// Stringify takes a http.Response and returns string & []byte
// values of the response body
func Stringify(r *http.Response) (string, []byte) {
	b, _ := io.ReadAll(r.Body)
	return string(b), b
}
