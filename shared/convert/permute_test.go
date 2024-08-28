package convert_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/shared/convert"
)

func TestSharedConvertPermuteStings(t *testing.T) {
	test := [][]string{
		[]string{"account", "unit"},
		[]string{"prod"},
	}

	res := convert.PermuteStrings(test...)
	if len(res) != 2 {
		t.Errorf("failed generating possibles")
	}

}
