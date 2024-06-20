package server

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

func TestSharedServerApiResponseWithResult(t *testing.T) {

	res := NewApiResponseWithResult[*testI, map[string]*testI]()
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

	mapItems := map[string]*testI{
		"001": {Id: "001", Tag: "t1", Category: "cat1"},
		"002": {Id: "002", Tag: "t1", Category: "cat2"},
		"003": {Id: "003", Tag: "t2", Category: "cat2"},
		"004": {Id: "004", Tag: "t3", Category: "cat1"},
	}

	res.SetResult(mapItems)
	if len(res.GetResult()) != len(mapItems) {
		t.Errorf("set / get failed to match")
	}

	res2 := NewApiResponseWithResult[*testI, []*testI]()
	sliceItems := []*testI{
		{Id: "001", Tag: "t1", Category: "cat1"},
		{Id: "002", Tag: "t1", Category: "cat2"},
		{Id: "003", Tag: "t2", Category: "cat2"},
		{Id: "004", Tag: "t3", Category: "cat1"},
	}
	res2.SetResult(sliceItems)
	if len(res2.GetResult()) != len(sliceItems) {
		t.Errorf("set / get failed to match")
	}

	res3 := NewApiResponseWithResult[*testI, map[string][]*testI]()
	mapSliceItems := map[string][]*testI{
		"tag^t1.": {
			{Id: "001", Tag: "t1", Category: "cat1"},
			{Id: "002", Tag: "t1", Category: "cat2"},
		},
		"tag^t2.": {
			{Id: "003", Tag: "t2", Category: "cat2"},
		},
		"tag^t3.": {
			{Id: "004", Tag: "t3", Category: "cat1"},
		},
	}
	res3.SetResult(mapSliceItems)

	if len(res3.GetResult()) != len(mapSliceItems) {
		t.Errorf("set / get failed to match")
	}

}

func TestSharedServerApiResponseBase(t *testing.T) {

	res := NewApiSimpleResponse()

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
	if res.GetStatus() != http.StatusBadGateway || res.Status != http.StatusBadGateway {
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
	if res.GetStatus() != http.StatusBadRequest || res.Status != http.StatusBadRequest {
		t.Errorf("status not set properly")
	}
}
