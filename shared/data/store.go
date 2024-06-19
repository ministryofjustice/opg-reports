package data

import "errors"

var (
	ErrStoreItemNotFoundViaIdx error = errors.New("store item not found via index") // Flag that the request index is not present in the store
)

// IStore is the common interface for all stores so we can always change the underlying
// data source (eg to database / redis etc) easily
type IStore[T IEntry] interface {
	// Add includes the item into the stores dataset, returns an error if it cant
	Add(item T) error
	// Get looks for the primary index for the item within the dataset. If it cant
	// be found, it returns an error
	Get(idx string) (T, error)
	// Filter will return items that match *ANY* of the filters, acting as *OR* logic
	// Chain multiple Filter calls together for AND logic
	Filter(filters ...IStoreFilter[T]) IStore[T]
	// All returns all of the items within the store
	All() map[string]T
	// Length returns the number of items within the store
	Length() int
}

type IStoreFilter[T IEntry] func(item T) bool

// Store is memory based store that operates from the items map.
// Impliments: IStore
type Store[T IEntry] struct {
	items map[string]T
}

// Add uses the item.Idx() as the key for this item and inserts
// item into the data store.
// This will overwrite items that are already present without error
func (s *Store[T]) Add(item T) error {
	s.items[item.Idx()] = item
	return nil
}

// Get returns the item from the iternal map whose index matches
// the passed value.
// If it cant be found an error is returned
func (s *Store[T]) Get(idx string) (i T, err error) {
	if value, ok := s.items[idx]; ok {
		i = value
	} else {
		err = ErrStoreItemNotFoundViaIdx
	}
	return
}

// List returns all items from the store
func (s *Store[T]) All() map[string]T {
	return s.items
}

// Length returns the count of items in the store
func (s *Store[T]) Length() int {
	return len(s.items)
}

// Filter compares each item in the data store against each filter. If the item
// matches *ANY* filter it will be included in the result, functioning as *OR*
// style logic
//
// To allow *AND* filters this returns a new Store[T] with matches, allowing
// Filter calls to be chained together
func (s *Store[T]) Filter(filters ...IStoreFilter[T]) (store *Store[T]) {
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
