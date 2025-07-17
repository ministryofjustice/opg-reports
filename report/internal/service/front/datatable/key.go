package datatable

import (
	"fmt"
	"slices"
)

func CombinationKey(item map[string]string, identifiers ...string) (key string) {
	slices.Sort(identifiers)
	identifiers = slices.Compact(identifiers)

	key = ""
	for _, k := range identifiers {
		key += fmt.Sprintf("%s:%s^", k, item[k])
	}
	return
}
