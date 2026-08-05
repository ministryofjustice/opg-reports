package sqlx

import (
	"log/slog"
	"opg-reports/app/internal/logx"
	"reflect"
)

// allFieldsHaveJSONTags checks all fields via reflection
// and looks for the json tag beign present
//
// Non-struct reflectiosn will always return false
func allFieldsHaveJSONTags(ref *reflection) bool {
	var (
		jsonTag string                = "json"                                   // tag to look for on the field
		fields  []reflect.StructField = ref.RecursiveFields(nil)                 // all visible fields
		lg      *slog.Logger          = logx.Default().With("T", ref.T.String()) // grab the default logger
	)
	lg.Debug("validating reflection struct fields have json tags.")
	if !ref.IsStruct {
		lg.Debug("reflection says this is not a struct, returning false.")
		return false
	}

	for _, f := range fields {
		if _, ok := f.Tag.Lookup(jsonTag); !ok {
			lg.Debug("reflection struct field does not have the json tag present.", "fieldName", f.Name)
			return false
		}
	}
	lg.Debug("all fields have json tags present.")
	return true
}
