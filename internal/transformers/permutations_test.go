package transformers_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/transformers"
)

func TestTransformersPermutations(t *testing.T) {

	var checkN = []string{"1", "2", "3"}
	var checkA = []string{"A", "B"}
	var check = [][]string{checkN, checkA}
	var expected = len(checkN) * len(checkA)

	p := transformers.Permutations(check...)
	actual := len(p)
	if actual != expected {
		t.Errorf("permutation length mismtach")
	}
}
