package utils

import (
	"testing"
)

func TestTransformersPermutations(t *testing.T) {

	var checkN = []string{"1", "2", "3"}
	var checkA = []string{"A", "B"}
	var check = [][]string{checkN, checkA}
	var expected = len(checkN) * len(checkA)

	p := Permutations(check...)
	actual := len(p)
	if actual != expected {
		t.Errorf("permutation length mismtach")
	}
}
