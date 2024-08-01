package response

import (
	"fmt"
	"net/http"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"opg-reports/shared/logger"
	"strings"
	"testing"
	"time"
)

func TestSharedServerResponseToJsonFromJson(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	tb := FakeTable(3, 2, 5, 1)
	resp.SetData(tb)

	js, _ := ToJson(resp)
	if !strings.Contains(string(js), `"status": 200`) {
		t.Errorf("json did not contain status")
	}

	resp2 := NewResponse[*Cell, *Row[*Cell]]()
	FromJson(js, resp2)

	if resp2.GetStatus() != resp.GetStatus() {
		t.Errorf("status mismatch")
	}

}
func TestSharedServerResponseRequestData(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()

	tb := FakeTable(3, 2, 5, 1)
	resp.SetData(tb)

	d := resp.GetData()

	if d != tb {
		t.Errorf("table data mismatch")
	}

	h := d.GetTableHead()
	if h.GetHeadersCount() != 2 {
		t.Errorf("incorrect number of headers in row")
		fmt.Printf("%+v\n", h)
	}
	bdy := d.GetTableBody()
	if len(bdy) != 3 {
		t.Errorf("incorrect number of body rows")
		fmt.Printf("%+v\n", h)
	}

	f := d.GetTableFoot()
	if f.GetSupplementaryCount() != 1 {
		t.Errorf("incorrect number of extras")
		fmt.Printf("%+v\n", f)
	}

}
func TestSharedServerResponseRequestErrorAndStatus(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	resp.SetErrorAndStatus(fmt.Errorf("error test!"), http.StatusNotExtended)

	if len(resp.GetError()) != 1 {
		t.Errorf("error set failed")
	}
	if resp.GetStatus() != http.StatusNotExtended {
		t.Errorf("status not set")
	}
}
func TestSharedServerResponseRequestErrors(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	resp.SetError(fmt.Errorf("error test!"))

	if len(resp.GetError()) != 1 {
		t.Errorf("error set failed")
	}
}

func TestSharedServerResponseRequestStatus(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	if resp.GetStatus() != http.StatusOK {
		t.Errorf("default status not set")
	}

	resp.SetStatus(http.StatusBadGateway)
	if resp.GetStatus() != http.StatusBadGateway {
		t.Errorf("status not updated")
	}

}
func TestSharedServerResponseRequestDataAge(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	now := time.Now().UTC()
	min := now.AddDate(-1, 0, 0).UTC()
	ds := []time.Time{
		fake.Date(min, now),
		min,
		fake.Date(min, now),
		now,
		fake.Date(min, now),
	}
	resp.SetDataAge(ds...)

	minT := resp.GetDataAgeMin()
	maxT := resp.GetDataAgeMax()

	if minT.UnixMicro() != min.UnixMicro() {
		t.Errorf(fmt.Sprintf("min failed. actual [%+v] expected [%+v]", minT, min))
	}
	if maxT.UnixMicro() != now.UnixMicro() {
		t.Errorf(fmt.Sprintf("max failed. actual [%+v] expected [%+v]", maxT, now))
		fmt.Printf("%+v\n", maxT)
		fmt.Printf("%+v\n", now)
	}

}

func TestSharedServerResponseRequestDuration(t *testing.T) {
	logger.LogSetup()
	resp := NewResponse[ICell, IRow[ICell]]()
	resp.SetStart()
	resp.SetEnd()
	resp.SetDuration()
	d := resp.GetDuration()

	if d.Microseconds() > 300 {
		t.Errorf("duration failed")
	}

}
func TestSharedServerResponseRequestStart(t *testing.T) {
	logger.LogSetup()
	now := time.Now().UTC()
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetStart()
	s := resp.GetStart()

	if s.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("failed to start resp correctly")
	}
}

func TestSharedServerResponseRequestEnd(t *testing.T) {
	logger.LogSetup()
	now := time.Now().UTC()
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetEnd()
	s := resp.GetEnd()

	if s.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("failed to end resp correctly")
	}
}
