package fmtx

import (
	"testing"
)

type tSprintfNRes struct {
	Str     string
	Ordered []interface{}
}
type tSprintfNSrc struct {
	Str  string
	Sql  bool
	Data map[string]interface{}
}
type testSprintfN struct {
	Src      *tSprintfNSrc
	Expected *tSprintfNRes
}

func TestSharedFmtxSprintfN(t *testing.T) {

	var tests = []*testSprintfN{
		// Test a SQL like statement and a slice
		// 	- named params should become ?
		//  - returned data should be all the months
		{
			Src: &tSprintfNSrc{
				Str:  `SELECT id, name, month FROM test_table WHERE month IN(:months);`,
				Sql:  true,
				Data: map[string]interface{}{"months": []string{"2026-01", "2026-02"}},
			},
			Expected: &tSprintfNRes{
				Str: `SELECT id, name, month FROM test_table WHERE month IN(?,?);`,
				Ordered: []interface{}{
					"2026-01", "2026-02",
				},
			},
		},
		// Test a string with multiple named values
		// 	- one name is reused (:foobar)
		// 	- one key doesnt exist in data and should be skipped
		{
			Src: &tSprintfNSrc{
				Str:  `This [:this] should be foobar and [:that] should be 0 and this is (:this) again. The :end`,
				Sql:  false,
				Data: map[string]interface{}{"this": "foobar", "that": 0},
			},
			Expected: &tSprintfNRes{
				Str: `This [foobar] should be foobar and [0] should be 0 and this is (foobar) again. The `,
				Ordered: []interface{}{
					"foobar", 0, "foobar",
				},
			},
		},
	}

	for _, test := range tests {
		actual, ordered := SprintfNamed(test.Src.Str, test.Src.Data, test.Src.Sql)
		// check string
		if actual != test.Expected.Str {
			t.Errorf("result mismatch, expected '%s' actual '%s'", test.Expected.Str, actual)
		}
		// check ordered values
		for i, v := range ordered {
			var exp = test.Expected.Ordered[i]
			if exp != v {
				t.Errorf("order mismatch, expected [%v] actual [%v]", exp, v)
			}
		}

	}

}
