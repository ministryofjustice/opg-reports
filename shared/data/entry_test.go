package data

import (
	"opg-reports/shared/logger"
	"testing"
	"time"
)

var nowT = time.Now().UTC()

type testEntry struct {
	Id string `json:"id"`
}

func (i *testEntry) UID() string {
	return i.Id
}

func (i *testEntry) Valid() bool {
	return true
}
func (i *testEntry) TS() time.Time {
	return nowT
}

type testEntryExt struct {
	Id       string `json:"id"`
	Tag      string `json:"tag"`
	Category string `json:"category"`
}

func (i *testEntryExt) UID() string {
	return i.Id
}
func (i *testEntryExt) TS() time.Time {
	return nowT
}
func (i *testEntryExt) Valid() bool {
	return true
}

func TestSharedDataEntryComparable(t *testing.T) {
	logger.LogSetup()
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
	logger.LogSetup()
	indent = false

	te := &testEntry{Id: "001"}

	m, _ := ToMap(te)
	if m["id"] != te.Id {
		t.Errorf("error converting to map")
	}

	p, _ := FromMap[*testEntry](m)
	if p.Id != te.Id {
		t.Errorf("error converting back from map")
	}

	content := `{"id":"100"}`
	i, _ := FromJson[*testEntry]([]byte(content))
	if i.Id != "100" {
		t.Errorf("error converting from json")
	}

	b, _ := ToJson(i)
	if string(b) != content {
		t.Errorf("error converting to json: (%s)=(%s)", string(b), content)
	}

	list := []*testEntry{
		te, {Id: "002"},
	}
	by, err := ToJsonList(list)
	if err != nil {
		t.Errorf("unexpected error")
	}
	asStr := `[{"id":"001"},{"id":"002"}]`
	if string(by) != asStr {
		t.Errorf("error, value doesnt match")
	}

	j, _ := FromJsonList[*testEntry]([]byte(asStr))
	if len(j) != len(list) {
		t.Errorf("lists dont match")
	}

}
