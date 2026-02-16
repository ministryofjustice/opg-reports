package api

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/utils/logger"
	"testing"
	"time"
)

type holiday struct {
	Title string `json:"title"`
}
type holidays struct {
	Division string     `json:"division"`
	Events   []*holiday `json:"events"`
}
type bankHols struct {
	EnglandAndWales *holidays `json:"england-and-wales"`
	Scotland        *holidays `json:"scotland"`
	NorthernIreland *holidays `json:"northern-ireland"`
}

func TestAPIGet(t *testing.T) {
	var (
		ctx = t.Context()
		log = logger.New("error")
	)
	res, _, _ := Get[*bankHols](ctx, log, &Call{
		Host:     `https://www.gov.uk`,
		Endpoint: `/bank-holidays.json`,
		Timeout:  (5 * time.Second),
	})

	if len(res.EnglandAndWales.Events) <= 0 {
		t.Errorf("bank holiday data failed")
	}

}

func TestAPICall(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/?start_date=2025-02", nil)
	c1 := &Call{
		Host:     "localhost:8081",
		Endpoint: "/v1/uptime/between/{start_date}/{end_date}",
		Request:  r,
		Params: []*Param{
			{Type: PATH, Key: "start_date", Value: "2025-01"},
			{Type: PATH, Key: "end_date", Value: "2025-04"},
			{Type: QUERY, Key: "team", Value: "true"},
			{Type: QUERY, Key: "account", Value: "true"},
		},
	}
	expected := "http://localhost:8081/v1/uptime/between/2025-02/2025-04?team=true&account=true"
	u1, _ := c1.URL()
	if expected != u1 {
		t.Errorf("url parsing mismatch")
	}

}
func TestAPIParamValue(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/?start_date=2026-01", nil)
	p1 := &Param{
		Key:   "start_date",
		Value: "2025-01",
	}
	v1 := p1.GetValue(r)
	if v1 != "2026-01" {
		t.Errorf("failed to get param from front")
	}

}

func TestAPIParseURI(t *testing.T) {

	var tests = map[string]string{
		"?test=1":                               "http://localhost/?test=1",
		"/?test":                                "http://localhost/?test",
		"https://www.gov.uk":                    "https://www.gov.uk/",
		"https://www.gov.uk/":                   "https://www.gov.uk/",
		"https://www.gov.uk/bank-holidays.json": "https://www.gov.uk/bank-holidays.json",
		"www.gov.uk/bank-holidays.json?test=1":  "http://www.gov.uk/bank-holidays.json?test=1",
		"localhost:80/test/?test=foo&bar=yes":   "http://localhost:80/test/?test=foo&bar=yes",
	}

	for test, expected := range tests {
		actual, err := parseURI(test)
		if err != nil {
			t.Errorf("unexpected error for [%s]: %s", test, err.Error())
		}
		if expected != actual {
			t.Errorf("mismatch, expected [%s], actual [%s]", expected, actual)
		}
	}

}
