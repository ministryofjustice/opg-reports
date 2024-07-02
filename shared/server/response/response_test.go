package response

import (
	"fmt"
	"net/http"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"strings"
	"testing"
	"time"
)

func TestSharedServerResponseToJsonFromJson(t *testing.T) {
	resp := NewResponse[*Cell, *Row[*Cell]]()
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
	resp := NewResponse[*Cell, *Row[*Cell]]()

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
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetErrorAndStatus(fmt.Errorf("error test!"), http.StatusNotExtended)

	if len(resp.GetError()) != 1 {
		t.Errorf("error set failed")
	}
	if resp.GetStatus() != http.StatusNotExtended {
		t.Errorf("status not set")
	}
}
func TestSharedServerResponseRequestErrors(t *testing.T) {
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetError(fmt.Errorf("error test!"))

	if len(resp.GetError()) != 1 {
		t.Errorf("error set failed")
	}
}

func TestSharedServerResponseRequestStatus(t *testing.T) {
	resp := NewResponse[*Cell, *Row[*Cell]]()
	if resp.GetStatus() != http.StatusOK {
		t.Errorf("default status not set")
	}

	resp.SetStatus(http.StatusBadGateway)
	if resp.GetStatus() != http.StatusBadGateway {
		t.Errorf("status not updated")
	}

}
func TestSharedServerResponseRequestDataAge(t *testing.T) {
	resp := NewResponse[*Cell, *Row[*Cell]]()
	now := time.Now().UTC()
	min := now.AddDate(-1, 0, 0)
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

	if minT != min {
		t.Errorf("min failed")
		fmt.Printf("%+v\n", minT)
		fmt.Printf("%+v\n", min)
	}
	if maxT != now {
		t.Errorf("max failed")
		fmt.Printf("%+v\n", maxT)
		fmt.Printf("%+v\n", now)
	}

}

func TestSharedServerResponseRequestDuration(t *testing.T) {
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetStart()
	resp.SetEnd()
	resp.SetDuration()
	d := resp.GetDuration()

	if d.Microseconds() > 300 {
		t.Errorf("duration failed")
	}

}
func TestSharedServerResponseRequestStart(t *testing.T) {
	now := time.Now().UTC()
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetStart()
	s := resp.GetStart()

	if s.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("failed to start resp correctly")
	}
}

func TestSharedServerResponseRequestEnd(t *testing.T) {
	now := time.Now().UTC()
	resp := NewResponse[*Cell, *Row[*Cell]]()
	resp.SetEnd()
	s := resp.GetEnd()

	if s.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
		t.Errorf("failed to end resp correctly")
	}
}

// func TestSharedServerResponseTimings(t *testing.T) {
// 	ti := &Timings{RequestTimes: &requestTimes{}, Datatimes: &dataTimings{}}
// 	now := time.Now().UTC()
// 	min := now.AddDate(-5, 0, 0)
// 	ti.AddTimestamp(now)
// 	ti.AddTimestamp(min)

// 	for i := 0; i < 5; i++ {
// 		d := fake.Date(min, now)
// 		ti.AddTimestamp(d)
// 	}

// 	mt, mx := ti.GetMinMax()
// 	if mt != min {
// 		t.Errorf("min failed")
// 		fmt.Println(mt)
// 		fmt.Println(min)
// 	}
// 	if mx != now {
// 		t.Errorf("max failed")
// 		fmt.Println(mx)
// 		fmt.Println(now)
// 	}
// }

// func TestSharedServerResponseResultWithResult(t *testing.T) {
// 	res := NewResponse()
// 	res.Start()
// 	res.End()
// 	now := time.Now().UTC()
// 	if res.RequestTimes.Start.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
// 		t.Errorf("start time failed")
// 	}
// 	if res.GetStart().Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
// 		t.Errorf("start time failed")
// 	}
// 	if res.RequestTimes.End.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
// 		t.Errorf("end time failed")
// 	}
// 	if res.RequestTimes.Duration.String() == "" {
// 		t.Errorf("duration failed")
// 	}

// 	cells := []*Cell{
// 		{Name: "001", Value: "101"},
// 		{Name: "002", Value: "102"},
// 		{Name: "003", Value: "103"},
// 	}
// 	rows := NewRows(cells)

// 	if len(rows) != 1 {
// 		t.Errorf("incorrect number of rows")
// 	}

// 	data := NewData(rows...)

// 	res.SetResult(data)
// 	got := res.GetResult()

// 	if len(got.Rows) != len(rows) || len(got.Rows) != len(data.Rows) {
// 		t.Errorf("row count mismatch")
// 	}

// }

// func TestSharedServerResponseBase(t *testing.T) {

// 	res := NewSimpleResult()

// 	// test timings
// 	now := time.Now().UTC()
// 	res.Start()
// 	if res.RequestTimes.Start.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
// 		t.Errorf("start time failed")
// 	}
// 	res.End()
// 	if res.RequestTimes.End.Format(dates.FormatYMD) != now.Format(dates.FormatYMD) {
// 		t.Errorf("end time failed")
// 	}
// 	if res.RequestTimes.Duration.String() == "" {
// 		t.Errorf("duration failed")
// 	}

// 	// test status
// 	if res.GetStatus() != http.StatusOK {
// 		t.Errorf("status error")
// 	}
// 	res.SetStatus(http.StatusBadGateway)
// 	if res.GetStatus() != http.StatusBadGateway || res.Status.Code != http.StatusBadGateway {
// 		t.Errorf("status error")
// 	}

// 	if len(res.GetErrors()) != 0 {
// 		t.Errorf("errors already set")
// 	}
// 	res.AddError(errors.New("test error1"))
// 	res.AddError(errors.New("test error2"))
// 	if len(res.GetErrors()) != 2 {
// 		t.Errorf("errors not set")
// 	}

// 	res.SetErrors([]error{})
// 	if len(res.GetErrors()) != 0 {
// 		t.Errorf("errors not set properly")
// 	}

// 	res.SetErrors([]error{errors.New("test error3")})
// 	if len(res.GetErrors()) != 1 {
// 		t.Errorf("errors not set properly")
// 	}

// 	res.AddErrorWithStatus(errors.New("test4"), http.StatusBadRequest)
// 	if len(res.GetErrors()) != 2 {
// 		t.Errorf("errors not added properly")
// 	}
// 	if res.GetStatus() != http.StatusBadRequest || res.Status.Code != http.StatusBadRequest {
// 		t.Errorf("status not set properly")
// 	}
// }

// func TestSharedServerResponseTableDataCell(t *testing.T) {

// 	c := NewCell("c1", "v1")
// 	if c.GetValue() != c.Value || c.GetValue() != "v1" {
// 		t.Errorf("get value error")
// 	}
// 	c.SetName("name")
// 	if c.GetName() != "name" || c.GetName() != c.Name {
// 		t.Errorf("get name error")
// 	}
// 	c.SetValue("val")
// 	if c.GetValue() != c.Value || c.GetValue() != "val" {
// 		t.Errorf("get value error")
// 	}
// }

// func TestSharedServerResponseTableDataRow(t *testing.T) {

// 	c := NewCell("c1", "v1")
// 	r := NewRow[*Cell]()
// 	r.AddCells(c)

// 	if len(r.GetCells()) != 1 {
// 		t.Errorf("get cells mismatch")
// 	}

// 	r.SetCells([]*Cell{c, c})
// 	if len(r.GetCells()) != 2 {
// 		t.Errorf("get cells mismatch")
// 	}
// }

// func TestSharedServerResponseRowHeader(t *testing.T) {
// 	pre1 := NewHeaderCell("h1", "header1")
// 	pre2 := NewHeaderCell("h2", "header2")
// 	pre3 := NewHeaderCell("h3", "header3")
// 	post1 := NewHeaderCell("post1", "suff1")
// 	post2 := NewHeaderCell("post2", "suff2")
// 	c1 := NewCell("c1", "v1")
// 	c2 := NewCell("c1", "v1")

// 	row := NewRow(pre1, pre2, pre3, c1, c2, post1, post2)
// 	row.UpdateCounters()
// 	pre, post := row.GetCounters()
// 	if pre != 3 {
// 		t.Errorf("pre counter incorrect")
// 	}
// 	if post != 2 {
// 		t.Errorf("post counter incorrect")
// 	}
// }
// func TestSharedServerResponseTableDataGetSet(t *testing.T) {

// 	c := NewCell("c1", "v1")
// 	r1 := NewRow[*Cell]()
// 	r1.AddCells(c)
// 	r2 := NewRow[*Cell]()
// 	r2.AddCells(c)
// 	d := NewData(r1, r2)

// 	if len(d.GetRows()) != 2 {
// 		t.Errorf("row mismatch")
// 	}

// 	d.SetRows([]*Row[*Cell]{r1})
// 	if len(d.GetRows()) != 1 {
// 		t.Errorf("row mismatch")
// 	}

// 	d.AddRows(r2)
// 	if len(d.GetRows()) != 2 {
// 		t.Errorf("adding rows failed")
// 	}

// 	h := NewRow(c)
// 	d.SetHeader(h)
// 	hh := d.GetHeader()
// 	if len(hh.GetCells()) != 1 {
// 		t.Errorf("incorrect heading")
// 	}

// 	h = NewRow(c)
// 	d.SetFooter(h)
// 	ff := d.GetFooter()
// 	if len(ff.GetCells()) != 1 {
// 		t.Errorf("incorrect footer")
// 	}

// 	d2 := NewData(r1, r2)
// 	if d2.GetHeader() != nil {
// 		t.Errorf("should be nil")
// 	}

// }

// func TestSharedServerResponseParseFromJson(t *testing.T) {
// 	resp := NewResponse()
// 	err := ParseFromJson([]byte(sampleRes), resp)
// 	if err != nil {
// 		t.Errorf("failed to parse normal json")
// 	}
// 	if resp.GetStatus() != 404 {
// 		t.Errorf("status mismatch")
// 	}
// }

// func TestSharedServerResponseParseFromHttp(t *testing.T) {
// 	ms := mockServer(sampleRes, http.StatusOK)
// 	defer ms.Close()

// 	httpResp, err := getUrl(ms.URL)
// 	if err != nil {
// 		t.Errorf("unexpected error")
// 	}

// 	resp := NewResponse()
// 	err = ParseFromHttp(httpResp, resp)
// 	if err != nil {
// 		t.Errorf("unexpected error")
// 	}

// 	if resp.GetStatus() != 404 {
// 		t.Errorf("status mismatch")
// 	}
// }

// func mockServer(resp string, status int) *httptest.Server {
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(status)
// 		w.Write([]byte(resp))
// 	}))
// 	return server
// }
// func getUrl(url string) (resp *http.Response, err error) {
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return
// 	}
// 	apiClient := http.Client{Timeout: time.Second * 3}
// 	resp, err = apiClient.Do(req)
// 	return
// }

// var sampleRes string = `{
// 	"timings": {
// 	  "start": "0001-01-01T00:00:00Z",
// 	  "end": "0001-01-01T00:00:00Z",
// 	  "duration": 0
// 	},
// 	"status": 404,
// 	"errors": [],
// 	"result": {
// 	  "headings": {
// 		"cells": [
// 		  {
// 			"name": "h1",
// 			"value": "hv1"
// 		  }
// 		]
// 	  },
// 	  "footer": {
// 		"cells": [
// 		  {
// 			"name": "f1",
// 			"value": "fv1"
// 		  }
// 		]
// 	  },
// 	  "rows": [
// 		{
// 		  "cells": [
// 			{
// 			  "name": "c1",
// 			  "value": "v1"
// 			}
// 		  ]
// 		}
// 	  ]
// 	}
//   }`
