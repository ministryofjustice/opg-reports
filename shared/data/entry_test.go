package data

import (
	"testing"
)

type testEntry struct {
	Id string `json:"id"`
}

func (i *testEntry) UID() string {
	return i.Id
}

func (i *testEntry) Valid() bool {
	return true
}

func TestSharedDataEntryComparable(t *testing.T) {
	t1 := &testEntry{Id: "t1"}
	t2 := &testEntry{Id: "t2"}
	t3 := &testEntry{Id: "t1"}

	if *t1 == *t2 {
		t.Errorf("comparison matched when shouldn't")
	}

	if *t3 != *t1 {
		t.Errorf("comparison failed when shouldn't")
	}

}

func TestSharedDataMapConversion(t *testing.T) {
	te := &testEntry{Id: "001"}

	m, _ := ToMap(te)
	if m["id"] != te.Id {
		t.Errorf("error converting to map")
	}

	p, _ := FromMap[*testEntry](m)
	if p.Id != te.Id {
		t.Errorf("error converting back from map")
	}

	content := `{"id": "100"}`
	i, _ := FromJson[*testEntry]([]byte(content))
	if i.Id != "100" {
		t.Errorf("error converting from json")
	}

}
