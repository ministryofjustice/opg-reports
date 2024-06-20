package data

import (
	"fmt"
	"strings"
)

const endOfField string = "."
const endOfKey string = "^"

// ToIdx generates a striung index for grouping that merges the field name and the field value.
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

	return
}
