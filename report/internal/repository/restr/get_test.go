package restr

import (
	"context"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/utils"
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

func TestRestGet(t *testing.T) {
	var (
		url    = `https://www.gov.uk/bank-holidays.json`
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		client = http.Client{Timeout: (2 * time.Second)}
		store  = Default(ctx, log, conf)
		result = &bankHols{}
	)

	code, err := store.Get(client, url, result)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if code != http.StatusOK {
		t.Errorf("incorrect status code: %v", code)
	}
	if len(result.EnglandAndWales.Events) <= 0 {
		t.Errorf("missing data from call")
	}

}

func TestRestParseURI(t *testing.T) {

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
