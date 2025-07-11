package utils

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMergeDefaultsWithQueryStrings(t *testing.T) {
	var testFrontUrl = "https://www.gov.uk/bank-holidays.json?start_date=2025-01"
	var parsedUrl, _ = url.Parse(testFrontUrl)
	var defaults = map[string]string{
		"start_date": "test",
		"end_date":   "-",
	}

	testReq := &http.Request{URL: parsedUrl}

	res := MergeRequestWithDefaults(testReq, defaults)
	if res["end_date"] != "-" {
		t.Errorf("end_date not parsed")
	}
	if res["start_date"] != "2025-01" {
		t.Errorf("start_date failed to be passed")
	}

}
