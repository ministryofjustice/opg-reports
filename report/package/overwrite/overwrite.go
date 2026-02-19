package overwrite

import "slices"

type keyed interface {
	ID() string
}

// Overwrite rebuilds the list of items to use, starting with the
// overwrites and only adding the existing items back in that
// weren't in the overwrites list.
func Overwrite[T keyed](existing []T, overwrites ...T) (list []T) {
	var added = []string{}

	list = []T{}
	// firstly, add all overwrites into the list by default
	for _, i := range overwrites {
		list = append(list, i)
		added = append(added, i.ID())
	}

	// now we add the existing items that are not already in the list
	for _, item := range existing {
		if !slices.Contains(added, item.ID()) {
			list = append(list, item)
			added = append(added, item.ID())
		}
	}

	return
}
