package data

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/shared/files"
)

var (
	ErrStoreItemNotFound error = errors.New("store item not found via index") // Flag that the request index is not present in the store
)

// IStore is the common interface for all stores so we can always change the underlying
// data source (eg to database / redis etc) easily
type IStore[T IEntry] interface {
	// Add includes the item into the stores dataset, returns an error if it cant
	Add(item T) error
	// Get looks for the primary index for the item within the dataset. If it cant
	// be found, it returns an error
	Get(uid string) (T, error)
	// Filter will return items that match *ANY* of the filters, acting as *OR* logic
	// Chain multiple Filter calls together for AND logic
	Filter(filters ...IStoreFilter[T]) IStore[T]
	// Group merges parts of data into chunks
	Group(group IStoreGrouper[T]) (stores map[string]IStore[T])
	// All returns all of the items within the store
	All() map[string]T
	// List returns all things as a slice
	List() []T
	// Length returns the number of items within the store
	Length() int
}

// IStoreFilter is used to enforce a signature on functions that are used to filter the
// data store. Each IStoreFilter is called against and item, those that return true
// are added to the result
type IStoreFilter[T IEntry] func(item T) bool

// IStoreGrouper enforces a signature for functions used to group the data store items.
// These functions should return a string used as index for the data store - something
// like YYYY-MM from a timestamp field etc
type IStoreGrouper[T IEntry] func(item T) string

// IStoreIdxer are used by ToIdxF to generate a field & value string that is then used
// as an index - typically for grouping the data together
type IStoreIdxer[T IEntry] func(item T) (string, string)

// Store is memory based store that operates from the items map.
// Impliments: IStore
type Store[T IEntry] struct {
	items map[string]T
}

// Add uses the item.UID() as the key for this item and inserts
// item into the data store.
// This will overwrite items that are already present without error
func (s *Store[T]) Add(item T) error {
	s.items[item.UID()] = item
	slog.Debug("[data/store] Add", slog.String("UID", item.UID()))
	return nil
}

// Get returns the item from the iternal map whose index matches
// the passed value.
// If it cant be found an error is returned
func (s *Store[T]) Get(uid string) (i T, err error) {
	if value, ok := s.items[uid]; ok {
		i = value
	} else {
		err = ErrStoreItemNotFound
	}
	slog.Debug("[data/store] Get", slog.String("UID", uid), slog.String("err", fmt.Sprintf("%v", err)))
	return
}

// All returns all items from the store
func (s *Store[T]) All() map[string]T {
	slog.Debug("[data/store] All")
	return s.items
}

// List returns all items from the store as a slice
func (s *Store[T]) List() (l []T) {
	l = []T{}
	for _, item := range s.All() {
		l = append(l, item)
	}
	slog.Debug("[data/store] list")
	return
}

// Length returns the count of items in the store
func (s *Store[T]) Length() int {
	slog.Debug("[data/store] length")
	return len(s.items)
}

// Group iterates over all items in the data store, runs the groupF func for each
// and then creates a new store for each key returned from the group func and Adds
// each relevant item to it
//
// Used to group data into chunks, like YYYY-MM out of a timestamp field
func (s *Store[T]) Group(groupF IStoreGrouper[T]) (stores map[string]IStore[T]) {
	stores = map[string]IStore[T]{}

	for _, item := range s.List() {
		key := groupF(item)
		if _, ok := stores[key]; !ok {
			stores[key] = NewStore[T]()
		}
		stores[key].Add(item)
	}
	slog.Debug("[data/store] group")
	return
}

// Filter compares each item in the data store against each filter. If the item
// matches *ANY* filter it will be included in the result, functioning as *OR*
// style logic
//
// To allow *AND* filters this returns a new Store[T] with matches, allowing
// Filter calls to be chained together
func (s *Store[T]) Filter(filters ...IStoreFilter[T]) (store IStore[T]) {
	found := map[string]T{}

	for key, item := range s.All() {
		match := false
		// check all filters, break on first match for speed
		for _, filterFunc := range filters {
			if filterFunc(item) == true {
				match = true
				break
			}
		}
		// add to found if its true
		if match {
			found[key] = item
		}
	}

	if len(found) > 0 {
		store = NewStoreFromMap(found)
	} else {
		store = NewStore[T]()
	}
	slog.Debug("[data/store] filter")
	return
}

// NewStore returns an empty store
func NewStore[T IEntry]() *Store[T] {
	return &Store[T]{items: map[string]T{}}
}

// NewStoreFromMap creates the store with preset items in form of a map
func NewStoreFromMap[T IEntry](items map[string]T) *Store[T] {
	return &Store[T]{items: items}
}

// NewStoreFromList create a store and adds each item in the list to the store
// by calling .Add()
func NewStoreFromList[T IEntry](items []T) *Store[T] {
	s := NewStore[T]()
	for _, i := range items {
		s.Add(i)
	}
	return s
}

// NewStoreFromFS loads all data files from a filesystem into the store.
// It unmarshals each file into a slice of T
// Note: errors with either loading the file or unmarshaling are
// ignored
func NewStoreFromFS[T IEntry, F files.IWriteFS](fS F) *Store[T] {
	allFiles := files.All(fS, true)
	allItems := []T{}

	for _, f := range allFiles {
		if content, err := files.ReadFile(fS, f); err == nil {
			if items, err := FromJsonList[T](content); err == nil {
				allItems = append(allItems, items...)
			}
		}
	}

	return NewStoreFromList(allItems)
}
