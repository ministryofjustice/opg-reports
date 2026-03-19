package httpx

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"testing"
)

var (
	_ ResponseWriter = &jsonResponse{}
	_ ResponseWriter = &htmlResponse{}
)

type tJSONDataResult struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type tJSONResponseExpected struct {
	Code int
	Type string
	Data *tJSONDataResult
}

type tJSONResponseTest struct {
	SourceData any
	Expected   *tJSONResponseExpected
}

func TestHttpxWriterJSONResponse(t *testing.T) {

	var tests = []*tJSONResponseTest{
		{
			SourceData: &tJSONDataResult{Name: "test", Value: "foobar"},
			Expected: &tJSONResponseExpected{
				Code: http.StatusOK,
				Type: jsonContentType,
				Data: &tJSONDataResult{Name: "test", Value: "foobar"},
			},
		},
	}

	for i, test := range tests {
		var res = &tJSONDataResult{}
		var rr = httptest.NewRecorder()
		var rw = NewJSONResponseWriter(rr)

		actual, _ := rw.Bytes(test.SourceData)
		// check status code
		if rw.Status() != test.Expected.Code {
			t.Errorf("[%d] response code expected to be [%d] but actually [%v]", i, test.Expected.Code, rw.Status())
		}
		// check content type
		_, ct := rw.ContentType()
		if ct != test.Expected.Type {
			t.Errorf("[%d] content type expected to be [%s] but actually [%v]", i, test.Expected.Type, ct)
		}
		// check data
		convert.Between(actual, &res)
		if res.Name != test.Expected.Data.Name {
			t.Errorf("[%d] name expected to be [%s] but actually [%v]", i, test.Expected.Data.Name, res.Name)
		}
		if res.Value != test.Expected.Data.Value {
			t.Errorf("[%d] value expected to be [%s] but actually [%v]", i, test.Expected.Data.Value, res.Value)
		}

	}

}
