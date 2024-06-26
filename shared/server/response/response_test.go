package response

import (
	"errors"
	"net/http"
	"opg-reports/shared/dates"
	"testing"
	"time"
)

type testI struct {
	Id       string `json:"id"`
	Tag      string `json:"tag"`
	Category string `json:"category"`
}

func (i *testI) UID() string {
	return i.Id
}

func (i *testI) Valid() bool {
	return true
}

func TestSharedServerResultWithResult(t *testing.T) {
	res := NewResult()
	res.Start()
	res.End()
	now := time.Now().UTC()
	if res.Times.Start.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("start time failed")
	}
	if res.Times.End.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("end time failed")
	}
	if res.Times.Duration.String() == "" {
		t.Errorf("duration failed")
	}

	cells := []*Cell{
		{Name: "001", Value: "101"},
		{Name: "002", Value: "102"},
		{Name: "003", Value: "103"},
	}
	rows := NewRows(cells)

	if len(rows) != 1 {
		t.Errorf("incorrect number of rows")
	}

	data := NewData(rows...)

	res.SetResult(data)
	got := res.GetResult()

	if len(got.Rows) != len(rows) || len(got.Rows) != len(data.Rows) {
		t.Errorf("row count mismatch")
	}

}

func TestSharedServerBase(t *testing.T) {

	res := NewSimpleResult()

	// test timings
	now := time.Now().UTC()
	res.Start()
	if res.Times.Start.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("start time failed")
	}
	res.End()
	if res.Times.End.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("end time failed")
	}
	if res.Times.Duration.String() == "" {
		t.Errorf("duration failed")
	}

	// test status
	if res.GetStatus() != http.StatusOK {
		t.Errorf("status error")
	}
	res.SetStatus(http.StatusBadGateway)
	if res.GetStatus() != http.StatusBadGateway || res.Status.Code != http.StatusBadGateway {
		t.Errorf("status error")
	}

	if len(res.GetErrors()) != 0 {
		t.Errorf("errors already set")
	}
	res.AddError(errors.New("test error1"))
	res.AddError(errors.New("test error2"))
	if len(res.GetErrors()) != 2 {
		t.Errorf("errors not set")
	}

	res.SetErrors([]error{})
	if len(res.GetErrors()) != 0 {
		t.Errorf("errors not set properly")
	}

	res.SetErrors([]error{errors.New("test error3")})
	if len(res.GetErrors()) != 1 {
		t.Errorf("errors not set properly")
	}

	res.AddErrorWithStatus(errors.New("test4"), http.StatusBadRequest)
	if len(res.GetErrors()) != 2 {
		t.Errorf("errors not added properly")
	}
	if res.GetStatus() != http.StatusBadRequest || res.Status.Code != http.StatusBadRequest {
		t.Errorf("status not set properly")
	}
}

func TestSharedServerTableDataCell(t *testing.T) {

	c := NewCell("c1", "v1")
	if c.GetValue() != c.Value || c.GetValue() != "v1" {
		t.Errorf("get value error")
	}
	c.SetName("name")
	if c.GetName() != "name" || c.GetName() != c.Name {
		t.Errorf("get name error")
	}
	c.SetValue("val")
	if c.GetValue() != c.Value || c.GetValue() != "val" {
		t.Errorf("get value error")
	}
}

func TestSharedServerTableDataRow(t *testing.T) {

	c := NewCell("c1", "v1")
	r := NewRow[*Cell]()
	r.AddCells(c)

	if len(r.GetCells()) != 1 {
		t.Errorf("get cells mismatch")
	}

	r.SetCells([]*Cell{c, c})
	if len(r.GetCells()) != 2 {
		t.Errorf("get cells mismatch")
	}
}

func TestSharedServerTableDataGetSet(t *testing.T) {

	c := NewCell("c1", "v1")
	r1 := NewRow[*Cell]()
	r1.AddCells(c)
	r2 := NewRow[*Cell]()
	r2.AddCells(c)
	d := NewData(r1, r2)

	if len(d.GetRows()) != 2 {
		t.Errorf("row mismatch")
	}

	d.SetRows([]*Row[*Cell]{r1})
	if len(d.GetRows()) != 1 {
		t.Errorf("row mismatch")
	}

}
