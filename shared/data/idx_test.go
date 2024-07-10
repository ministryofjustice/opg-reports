package data

import (
	"fmt"
	"testing"
)

func TestSharedDataToIdx(t *testing.T) {
	i := &testEntryExt{Id: "01", Tag: "test-tag", Category: "cat1"}

	idx := ToIdx(i, "tag", "category")

	if idx != fmt.Sprintf("%s^%s.%s^%s.", "tag", "test-tag", "category", "cat1") {
		t.Errorf("idx did not generate correctly")
	}

	i = &testEntryExt{Id: "01", Tag: "test-tag"}
	idx = ToIdx(i, "tag", "category")

	if idx != fmt.Sprintf("%s^%s.%s^%s.", "tag", "test-tag", "category", "-") {
		t.Errorf("idx did not generate correctly")
	}

	idxF := func(i *testEntryExt) (string, string) {
		return "tag", i.Tag
	}
	idxFe := func(i *testEntryExt) (string, string) {
		return "tag", ""
	}

	str := ToIdxF(i, idxF)
	if str != "tag^test-tag." {
		t.Errorf("idxf failed")
	}
	str = ToIdxF(i, idxFe)
	if str != "tag^-." {
		t.Errorf("idxf failed")
	}
}
