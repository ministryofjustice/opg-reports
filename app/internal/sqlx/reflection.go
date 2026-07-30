package sqlx

import (
	"reflect"
)

// reflection used to help identify if item is a
// struct or not via reflection and also provides
// helper methods to fetch all attached fields
// which are then used in validation / parsing
// of sql statements
type reflection struct {
	Src       any
	K         reflect.Kind
	V         reflect.Value
	T         reflect.Type
	IsStruct  bool
	valStruct bool
	ptrStruct bool
}

// Fields returns all visible fields for this struct
//
// Works out sorce based on struct type
func (self *reflection) Fields() (fields []reflect.StructField) {
	fields = []reflect.StructField{}
	if self.valStruct {
		fields = reflect.VisibleFields(self.T)
	} else if self.ptrStruct {
		fields = reflect.VisibleFields(self.T.Elem())
	}
	return
}

// RecursiveFields fetches all visible fields on this struct and
// if any field itself a struct it will recurse and fetch the
// the fields from there as well
//
// Passing `nil` as the initial fields param will cause the func
// to use the current `T` fields
func (self *reflection) RecursiveFields(fields []reflect.StructField) (all []reflect.StructField) {
	if fields == nil {
		fields = self.Fields()
	}
	all = []reflect.StructField{}

	// check each field for being a nested struct
	for _, f := range fields {
		var t = f.Type
		var sub = []reflect.StructField{}
		// add this field
		all = append(all, f)
		// now recursively find all other fields
		if byValueStruct(t) {
			sub = self.RecursiveFields(reflect.VisibleFields(t))
		} else if byPtrStruct(t) {
			sub = self.RecursiveFields(reflect.VisibleFields(t.Elem()))
		}
		if len(sub) > 0 {
			all = append(all, sub...)
		}
	}

	return
}

// newReflection creates a reflection struct based on the
// src passed in.
//
// The reflection struct is used to determine if the filterModel
// passed to `Bind` is actually a struct or not and also to
// query all fields (recursively) on a struct - generally to
// checking tag information about each struct field.
//
// Types, values and struct type (pointer / value) are
// calculated at this point
func newReflection(src any) (r *reflection) {
	var (
		val reflect.Value = reflect.ValueOf(src)
		typ reflect.Type  = reflect.TypeOf(src)
	)

	r = &reflection{
		Src:       src,
		V:         val,
		T:         typ,
		valStruct: byValueStruct(typ),
		ptrStruct: byPtrStruct(typ),
	}
	r.IsStruct = (r.valStruct || r.ptrStruct)
	return
}

// byValueStruct checks for structs passed by value
func byValueStruct(t reflect.Type) bool {
	return (t.Kind() == reflect.Struct)
}

// byPtrStruct decides if the type is a struct pointer
func byPtrStruct(t reflect.Type) bool {
	return (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}
