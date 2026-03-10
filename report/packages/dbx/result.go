package dbx

import (
	"database/sql"
	"errors"
	"opg-reports/report/packages/instance"
	"opg-reports/report/packages/types/interfaces"
)

var (
	ErrModelNotPointer = errors.New("model needs to be a pointer.") // returned if the model used for a select result is not a pointer version
	ErrRowScanFailed   = errors.New("rows scan failed with error.") // returned when theres an issue with the `row.Scan` call
)

// Results is used to capture results from a sql select
// into a slice of models and their struct values directly.
type Results[T interfaces.Selectable] struct {
	data []T // stores the results from the `row.Scan` mapping
}

// Data returns the internal results from the row scanning.
//
// Used at the end of a select to return the results.
func (self *Results[T]) Data() []T {
	return self.data
}

// RowScan iis used within a select to process each row from
// a sql query into a model struct.
//
// `T` Model most be a pointer so the `Sequence` func will
// be suitable for the `row.Scan` usage. Will error if
// its not.
//
// Uses `Sequence()` to map the result rows to a struct and
// updaes its internal `data` slice with the results
func (self *Results[T]) RowScan(row *sql.Rows) (err error) {
	var (
		item     T
		sequence []any
	)

	// create a new instance of T
	item = instance.Of[T]()
	// if the Model is not a pointer, then the .Sequence
	// wont work correctly and data wont get mapped, so
	// fail
	if !instance.IsPtr(item) {
		err = ErrModelNotPointer
		return
	}
	// setup default .data value
	if len(self.data) == 0 {
		self.data = []T{}
	}
	// now we get the ordered list of columns for the select
	// and run row scan aginst those
	sequence = item.Sequence()
	if err = row.Scan(sequence...); err == nil {
		self.data = append(self.data, item)
	} else {
		err = errors.Join(ErrRowScanFailed, err)
		return
	}

	return
}
