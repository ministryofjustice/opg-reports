package sqlx

import (
	"reflect"
)

// allFieldsHaveJSONTags checks all fields via reflection
// and looks for the json tag beign present
//
// Non-struct reflectiosn will always return false
func allFieldsHaveJSONTags(ref *reflection) bool {
	var (
		jsonTag string                = "json"
		fields  []reflect.StructField = ref.RecursiveFields(nil)
	)

	if !ref.IsStruct {
		return false
	}

	for _, f := range fields {
		if _, ok := f.Tag.Lookup(jsonTag); !ok {
			return false
		}
	}

	return true
}
