package nums_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/nums"
)

// TestNumsAdd tests the various versions
// of adding values using the tempalte helper
func TestNumsAdd(t *testing.T) {

	// -- add non mixed values

	if nums.Add(1.0, 2.1, 3.9, -1.0) != 6.0 {
		t.Errorf("adding floats failed")
	}
	if nums.Add(1, 2, 3, 1, -1) != 6 {
		t.Errorf("adding int failed")
	}
	if nums.Add("A", "b", "z") != "Abz" {
		t.Errorf("adding strings failed")
	}
	if v := nums.Add("1", "2.1", "z"); v != "3.1" {
		t.Errorf("adding strings failed: [%v]", v)
	}

	// -- test mix
	// the values not matching type of first
	// param are ignored
	if nums.Add(1, 2, 3.0, "t", 3) != 6 {
		t.Errorf("adding int with mix failed")
	}

}
