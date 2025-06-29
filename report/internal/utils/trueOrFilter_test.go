package utils_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestTrueOrFilter(t *testing.T) {

	// a string of "true" should be, Selectable, Groupable, & Orderable
	var isTrue string = "true"

	test := utils.TrueOrFilter(isTrue)
	if !test.Selectable() {
		t.Errorf("string == true should be selected as a column")
	}
	if !test.Groupable() {
		t.Errorf("string == true should be groupable")
	}
	if !test.Orderable() {
		t.Errorf("string == true should be orderable")
	}
	if test.Whereable() {
		t.Errorf("string == true should NOT be used for a where")
	}

	var isFilter = "team-name"
	test = utils.TrueOrFilter(isFilter)
	if !test.Selectable() {
		t.Errorf("string == value should be selected as a column")
	}
	if test.Groupable() {
		t.Errorf("string == value should NOT be groupable")
	}
	if test.Orderable() {
		t.Errorf("string == value should NOT be orderable")
	}
	if !test.Whereable() {
		t.Errorf("string == value should be used for a where")
	}

}
