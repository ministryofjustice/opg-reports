package sqlx

import (
	"reflect"
)

// Reflection used to help identify if item is a
// struct or not via reflection and also provides
// helper methods to fetch all attached fields
// which are then used in validation / parsing
// of sql statements
type Reflection struct {
	Src       any
	K         reflect.Kind
	V         reflect.Value
	T         reflect.Type
	IsStruct  bool
	valStruct bool
	ptrStruct bool
}

// func (self *Reflection) ElemValue() (el reflect.Value) {
// 	if self.valStruct {
// 		el = reflect.ValueOf(&self.Src).Elem()
// 	} else if self.ptrStruct {
// 		el = reflect.ValueOf(self.Src).Elem()
// 	}
// 	return
// }

// Fields returns all visible fields for this struct
//
// Works out sorce based on struct type
func (self *Reflection) Fields() (fields []reflect.StructField) {
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
func (self *Reflection) RecursiveFields(fields []reflect.StructField) (all []reflect.StructField) {
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
		if ByValueStruct(t) {
			sub = self.RecursiveFields(reflect.VisibleFields(t))
		} else if ByPtrStruct(t) {
			sub = self.RecursiveFields(reflect.VisibleFields(t.Elem()))
		}
		if len(sub) > 0 {
			all = append(all, sub...)
		}
	}

	return
}

// NewReflection creates a reflection struct based on the
// src passed in.
//
// Types, values and struct type (pointer / value) are
// calculated at this point
func NewReflection(src any) (r *Reflection) {
	var (
		val reflect.Value = reflect.ValueOf(src)
		typ reflect.Type  = reflect.TypeOf(src)
	)

	r = &Reflection{
		Src:       src,
		V:         val,
		T:         typ,
		valStruct: ByValueStruct(typ),
		ptrStruct: ByPtrStruct(typ),
	}
	r.IsStruct = (r.valStruct || r.ptrStruct)
	return
}

// ByValueStruct checks for structs passed by value
func ByValueStruct(t reflect.Type) bool {
	return (t.Kind() == reflect.Struct)
}

// ByPtrStruct decides if the type is a struct pointer
func ByPtrStruct(t reflect.Type) bool {
	return (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}
