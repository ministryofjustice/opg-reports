package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type tValuesFromRequestTest struct {
	Request  *http.Request
	Expected map[string]string
}

func TestHttpxRequestValuesFromRequest(t *testing.T) {
	var tests = []*tValuesFromRequestTest{
		{
			Request:  httptest.NewRequest(http.MethodGet, "/foo/bar/?name=a&age=30&date_start=2025-11", nil),
			Expected: map[string]string{"date_start": "2025-11"},
		},
	}

	for i, test := range tests {
		// pull data
		actual := ValuesFromRequest(test.Request).Map()
		// compare from data we have
		for k, v := range actual {
			if v != test.Expected[k] {
				t.Errorf("[%d] actual field [%s] value [%v] does not match expected [%s]", i, k, v, test.Expected[k])
			}
		}
		// compare from the expected
		for k, v := range test.Expected {
			if v != actual[k] {
				t.Errorf("[%d] expected field [%s] value [%s] does not match actual [%v]", i, k, v, actual[k])
			}
		}
	}

}
