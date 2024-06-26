package data

import (
	"fmt"
	"log/slog"
	"strings"
)

const endOfField string = "."
const endOfKey string = "^"

// ToIdx generates a string index for grouping that merges the field name and the field value.
// This allows a 1 depth map (map[string][]T) that is grouped by multiple fields
//
//	item := &IEntry{id: "01", "tag": "tOne"}
//	ToIdx(item, "id", "tag")
//	// Output: "id^01.tag^tOne"
func ToIdx[T IEntry](item T, fields ...string) string {
	str := ""

	if mapped, err := ToMap(item); err == nil {
		for _, key := range fields {
			key = strings.ToLower(key)
			var value string
			if v, ok := mapped[key]; !ok || v == "" {
				value = "-"
			} else {
				value = v
			}
			str += fmt.Sprintf("%s%s%s%s", key, endOfKey, value, endOfField)
		}

	}
	slog.Debug("[data/entry] ToIdx", slog.String("UID", item.UID()), slog.String("idx", str))
	return str
}

// ToIdxF generates a striung index for grouping that merges the field name and the field value.
// This allows a 1 depth map (map[string][]T) that is grouped by multiple fields
//
// Operates like [ToIdx], ubt instead of a list of fields it uses a series of functions. By using a
// function we can adjust content of the item, in particular reducing timestamps to just their
// month
func ToIdxF[T IEntry](item T, funcs ...IStoreIdxer[T]) string {
	str := ""
	for _, f := range funcs {
		key, value := f(item)
		key = strings.ToLower(key)

		if value == "" {
			value = "-"
		}
		str += fmt.Sprintf("%s%s%s%s", key, endOfKey, value, endOfField)
	}
	slog.Debug("[data/entry] ToIdx", slog.String("UID", item.UID()), slog.String("idx", str))
	return str
}
