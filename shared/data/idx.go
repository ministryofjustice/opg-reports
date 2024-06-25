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

func ToIdxKV(key string, value string) string {
	return fmt.Sprintf("%s%s%s%s", key, endOfKey, value, endOfField)
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
		if value == "" {
			value = "-"
		}
		str += fmt.Sprintf("%s%s%s%s", key, endOfKey, value, endOfField)
	}
	slog.Debug("[data/entry] ToIdx", slog.String("UID", item.UID()), slog.String("idx", str))
	return str
}

// FromIdx converts a string from ToIdx to a map of field:value
func FromIdx(idx string) (fieldsAndValues map[string]string) {
	fieldsAndValues = map[string]string{}

	for _, fv := range strings.Split(idx, endOfField) {
		chunks := strings.Split(fv, endOfKey)
		if len(chunks) == 2 {
			fieldsAndValues[chunks[0]] = chunks[1]
		}
	}
	slog.Debug("[data/entry] FromIdx", slog.String("idx", idx), slog.String("fieldsAndValues", fmt.Sprintf("%+v", fieldsAndValues)))
	return
}
