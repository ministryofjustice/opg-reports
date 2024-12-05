package transformers_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/transformers"
)

// TestTransformersKVPair checks that the KVPair is generated
// as expected
func TestTransformersKVPair(t *testing.T) {
	var expected = "foobar:1AB"
	var actual = transformers.KVPair("foobar", 1, "A", "B")
	if expected != actual {
		t.Errorf("kv pair not as expected - expected [%s] actual [%v]", expected, actual)
	}
}

// TestTransformersSortedColumnNames makes sure the column sorting
// works by checking typical alphanumeric and then throw in some
// smily too
func TestTransformersSortedColumnNames(t *testing.T) {
	var colValues = map[string][]interface{}{
		"Z": {"A", "B"},
		"A": {"one"},
		"b": {"1"},
		"1": {"-"},
		"-": {"what?"},
		"?": {"foobar"},
		"😀": {"smile"},
	}
	var expected = []string{"-", "1", "?", "A", "Z", "b", "😀"}
	var actual = transformers.SortedColumnNames(colValues)

	if len(actual) != len(expected) {
		t.Errorf("len mismatch!")
	}

	for i, v := range actual {
		var ex = expected[i]
		if v != ex {
			t.Errorf("order failed - expected [%s] actual [%v]", ex, v)
		}
	}

}

// TestTransformersColumnValuesList checks that the source map
// is converted into a sorted and merged set of values correctly
func TestTransformersColumnValuesList(t *testing.T) {
	var colValues = map[string][]interface{}{
		"Z": {"A", "B"},
		"A": {"one"},
		"b": {"1"},
		"1": {"-"},
		"-": {"what?"},
		"?": {"foobar"},
		"😀": {"smile"},
	}
	var expected = [][]string{
		{"-:what?^"},
		{"1:-^"},
		{"?:foobar^"},
		{"A:one^"},
		{"Z:A^", "Z:B^"},
		{"b:1^"},
		{"😀:smile^"},
	}
	var actual = transformers.ColumnValuesList(colValues)

	if len(actual) != len(expected) {
		t.Errorf("len mismatch!")
	}

	for i, actualI := range actual {
		var exI = expected[i]
		for x, val := range actualI {
			if exI[x] != val {
				t.Errorf("col values failed - expected [%s] actual [%v]", exI[x], val)
			}
		}

	}

}
