package data

import (
	"fmt"
	"opg-reports/shared/fake"
	"opg-reports/shared/files"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestSharedDataStoreGroup generates a series of test data that is grouped
// by Tag, Category & then both. These are then compared to the counters of
// each to make sure grouping is working
func TestSharedDataStoreGroup(t *testing.T) {
	items := []*testEntryExt{}
	tags := []string{"t1", "t2", "t3"}
	cats := []string{"c1", "c2"}

	tcounts := map[string]int{}
	ccounts := map[string]int{}
	mcount := map[string]int{}
	l := 100
	for i := 0; i < l; i++ {
		tg := fake.Choice(tags)
		ct := fake.Choice(cats)
		mk := fmt.Sprintf("%s:%s", tg, ct)

		if _, ok := tcounts[tg]; !ok {
			tcounts[tg] = 0
		}
		if _, ok := ccounts[ct]; !ok {
			ccounts[ct] = 0
		}
		if _, ok := mcount[mk]; !ok {
			mcount[mk] = 0
		}
		tcounts[tg]++
		ccounts[ct]++
		mcount[mk]++

		te := &testEntryExt{
			Id:       fmt.Sprintf("%d", i+1),
			Tag:      tg,
			Category: ct,
		}
		items = append(items, te)

	}
	s := NewStoreFromList[*testEntryExt](items)
	if s.Length() != l {
		t.Error("error with length")
	}

	// check the tag grouping matches the random data
	tg := func(item *testEntryExt) string {
		return ToIdx(item, "tag")
	}
	tstores := s.Group(tg)
	if len(tstores) != len(tcounts) {
		t.Errorf("tag grouping falied")
	}
	for tag, count := range tcounts {
		tkey := "tag^" + tag + "."
		if tstores[tkey].Length() != count {
			t.Errorf("counts dont match - expected [%d] actual [%d]", count, tstores[tkey].Length())
		}
	}

	// check the category grouping matches the data
	cg := func(item *testEntryExt) string {
		return ToIdx(item, "category")
	}
	cstores := s.Group(cg)
	if len(cstores) != len(ccounts) {
		t.Errorf("cat grouping falied")
	}
	for cat, count := range ccounts {
		ckey := "category^" + cat + "."
		if cstores[ckey].Length() != count {
			t.Errorf("counts dont match - expected [%d] actual [%d]", count, cstores[ckey].Length())
		}
	}
	mg := func(item *testEntryExt) string {
		return ToIdx(item, "tag", "category")
	}
	mstores := s.Group(mg)
	// check tag -> cat grouping is working
	for mk, count := range mcount {
		sp := strings.Split(mk, ":")
		k := fmt.Sprintf("tag^%s.category^%s.", sp[0], sp[1])
		st, ok := mstores[k]
		if !ok {
			t.Errorf("store does not exist")
		} else if st.Length() != count {
			t.Errorf("counts dont match - expected [%d] actual [%d]", count, st.Length())
		}

	}

}

func TestSharedDataStoreNewFromFS(t *testing.T) {
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "test-data-store-fs-*")
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)

	defer os.RemoveAll(tDir)
	// create 2 files, each with 100 items in for testing loading from a filesystem
	for x := 0; x < 2; x++ {
		fn, _ := os.CreateTemp(tDir, "dummy-*.json")
		defer os.Remove(fn.Name())

		items := []*testEntry{}
		for i := 0; i < 100; i++ {
			items = append(items, &testEntry{Id: fmt.Sprintf("%d000%d", x, i)})
		}
		j, _ := ToJsonList(items)
		files.WriteFile(tDir, fn.Name(), j)
	}

	store := NewStoreFromFS[*testEntry, *files.WriteFS](f)

	if store.Length() != 200 {
		t.Errorf("error loading store")
	}

}

func BenchmarkSharedDataStoreNewFromFS(b *testing.B) {
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "b-data-store-fs-*")
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)

	defer os.RemoveAll(tDir)
	// create 2 files, each with 100 items in for testing loading from a filesystem
	for x := 0; x < 2; x++ {
		fn, _ := os.CreateTemp(tDir, "dummy-*.json")
		defer os.Remove(fn.Name())

		items := []*testEntry{}
		for i := 0; i < 100; i++ {
			items = append(items, &testEntry{Id: fmt.Sprintf("%d000%d", x, i)})
		}
		j, _ := ToJsonList(items)
		files.WriteFile(tDir, fn.Name(), j)
	}

	NewStoreFromFS[*testEntry, *files.WriteFS](f)
}

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
		items[e.UID()] = e
	}
	s := NewStoreFromMap[*testEntry](items)
	if s.Length() != len(items) {
		t.Errorf("incorrect length: [%d] [%d] ", s.Length(), len(items))
	}
}

func BenchmarkSharedDataStoreNewFromMap(b *testing.B) {
	items := map[string]*testEntry{}
	for i := 0; i < 2000; i++ {
		e := &testEntry{Id: fmt.Sprintf("%d", i+1)}
		items[e.UID()] = e
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
	all := store.All()

	if len(all) != 100 || len(all) != store.Length() {
		t.Errorf("length mismatch")
	}
	list := store.List()
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
	store.List()
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
		t.Errorf("unexpected length of chanined filters: %d", fin.Length())
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
