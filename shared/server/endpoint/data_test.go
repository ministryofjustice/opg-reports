package endpoint

import (
	"opg-reports/shared/data"
	"opg-reports/shared/fake"
	"opg-reports/shared/logger"
	"slices"
	"testing"
	"time"
)

type testEntry struct {
	Id       string   `json:"id"`
	Tags     []string `json:"tags"`
	Category string   `json:"category"`
}

func (i *testEntry) UID() string {
	return i.Id
}
func (i *testEntry) TS() time.Time {
	return time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
}
func (i *testEntry) Valid() bool {
	return true
}

func TestSharedServerEndpointDataFilters(t *testing.T) {
	logger.LogSetup()

	tagFoo := 5
	categoryBar := 10
	create := 20

	items := []*testEntry{}
	for i := 0; i < create; i++ {
		item := &testEntry{
			Id:       fake.IntAsStr(1000, 9999),
			Tags:     []string{"default"},
			Category: "all",
		}
		if i < tagFoo {
			item.Tags = append(item.Tags, "foo")
		}
		if i < categoryBar {
			item.Category = "bar"
		}
		items = append(items, item)
	}
	store := data.NewStoreFromList[*testEntry](items)

	filters := map[string]data.IStoreFilterFunc[*testEntry]{
		"category-is-bar": func(i *testEntry) bool {
			return i.Category == "bar"
		},
	}

	d := NewEndpointData[*testEntry](store, nil, filters)

	s := d.ApplyFilters()
	li := s.List()
	if len(li) != categoryBar {
		t.Errorf("failed to filter correct number of items for category")
	}

	filters = map[string]data.IStoreFilterFunc[*testEntry]{
		"category-is-bar": func(i *testEntry) bool {
			return i.Category == "bar"
		},
		"tag-contains-foo": func(i *testEntry) bool {
			return slices.Contains(i.Tags, "foo")
		},
	}

	d = NewEndpointData[*testEntry](store, nil, filters)
	s = d.ApplyFilters()
	li = s.List()
	if len(li) != tagFoo {
		t.Errorf("failed to filter correct number of items for tag")
	}
}

func TestSharedServerEndpointDataGrouping(t *testing.T) {
	logger.LogSetup()

	categoryBar := 10
	create := 20

	items := []*testEntry{}
	for i := 0; i < create; i++ {
		item := &testEntry{
			Id:       fake.IntAsStr(1000, 9999),
			Tags:     []string{"default"},
			Category: "all",
		}
		if i < categoryBar {
			item.Category = "bar"
		}
		items = append(items, item)
	}
	store := data.NewStoreFromList[*testEntry](items)
	group := func(item *testEntry) string {
		return item.Category
	}
	d := NewEndpointData[*testEntry](store, group, nil)

	groups := d.ApplyGroupBy()

	if len(groups) != 2 {
		t.Errorf("grouping failed")
	}
}
