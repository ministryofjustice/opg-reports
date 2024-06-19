package data

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSharedDataStoreNew(t *testing.T) {
	s := NewStore[*testEntry]()
	if s.Length() != 0 {
		t.Errorf("not empty")
	}
}
func BenchmarkSharedDataStoreNew(b *testing.B) {
	NewStore[*testEntry]()
}

func TestSharedDataStoreNewFromMap(t *testing.T) {
	items := map[string]*testEntry{}
	for i := 0; i < 200; i++ {
		e := &testEntry{Id: fmt.Sprintf("%d", i+1)}
		items[e.Idx()] = e
	}
	s := NewStoreFromMap[*testEntry](items)
	if s.Length() != len(items) {
		t.Errorf("incorrect length")
	}
}

func BenchmarkSharedDataStoreNewFromMap(b *testing.B) {
	items := map[string]*testEntry{}
	for i := 0; i < 2000; i++ {
		e := &testEntry{Id: fmt.Sprintf("%d", i+1)}
		items[e.Idx()] = e
	}
	NewStoreFromMap[*testEntry](items)

}

func TestSharedDataStoreNewFromList(t *testing.T) {
	items := []*testEntry{}
	for i := 0; i < 200; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	s := NewStoreFromList[*testEntry](items)
	if s.Length() != len(items) {
		t.Errorf("incorrect length")
	}
}

func BenchmarkSharedDataStoreNewFromList(b *testing.B) {
	items := []*testEntry{}
	for i := 0; i < 200000; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	NewStoreFromList[*testEntry](items)
}

func TestSharedDataStoreAdd(t *testing.T) {
	store := NewStore[*testEntry]()
	for i := 0; i < 100; i++ {
		store.Add(&testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	if store.Length() != len(store.items) || store.Length() != 100 {
		t.Errorf("incorrect length")
	}
}

func BenchmarkSharedDataStoreAdd(b *testing.B) {
	store := NewStore[*testEntry]()
	for i := 0; i < 20000; i++ {
		store.Add(&testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
}

func TestSharedDataStoreAll(t *testing.T) {
	store := NewStore[*testEntry]()
	for i := 0; i < 100; i++ {
		store.Add(&testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	list := store.All()

	if len(list) != 100 || len(list) != store.Length() {
		t.Errorf("length mismatch")
	}
}

func BenchmarkSharedDataStoreAll(b *testing.B) {
	items := []*testEntry{}
	for i := 0; i < 200000; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	store := NewStoreFromList[*testEntry](items)
	store.All()
}

func TestSharedDataStoreGet(t *testing.T) {
	store := NewStore[*testEntry]()
	for i := 0; i < 1000; i++ {
		store.Add(&testEntry{Id: fmt.Sprintf("%d", i+1)})
	}

	if _, err := store.Get("2000"); err == nil {
		t.Errorf("should have failed to find idx 2000")
	}

	item, err := store.Get("100")
	if err != nil {
		t.Errorf("should have found item")
	}
	if item.Id != "100" {
		t.Errorf("item id should match search for idx")
	}

}

func BenchmarkSharedDataStoreGet(b *testing.B) {
	items := []*testEntry{}
	for i := 0; i < 200000; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	store := NewStoreFromList[*testEntry](items)
	store.Get("1000")

}

func TestSharedDataStoreFilter(t *testing.T) {
	items := []*testEntry{}
	for i := 0; i < 40; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	store := NewStoreFromList[*testEntry](items)

	under20 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i <= 20)
	}
	over20 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i > 20)
	}
	under10 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i < 10)
	}
	under5 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i < 5)
	}

	f1 := store.Filter(under20)

	if f1.Length() != 20 {
		t.Errorf("incorrect length on filter")
	}

	// should error
	f3 := f1.Filter(over20)
	if f3.Length() != 0 {
		t.Errorf("expected an error about no matching items")
	}

	fin := store.Filter(under10).Filter(under5)
	if fin.Length() != 4 {
		t.Errorf("unexpected length of chanined filters")
	}

}

func BenchmarkSharedDataStoreFilter(b *testing.B) {
	items := []*testEntry{}
	for i := 0; i < 200000; i++ {
		items = append(items, &testEntry{Id: fmt.Sprintf("%d", i+1)})
	}
	store := NewStoreFromList[*testEntry](items)

	u1 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i <= 199999)
	}

	u2 := func(item *testEntry) bool {
		i, _ := strconv.Atoi(item.Id)
		return (i <= 99999)
	}

	f1 := store.Filter(u1)

	f1.Filter(u2)

}
