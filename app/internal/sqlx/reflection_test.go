package sqlx

import (
	"reflect"
	"testing"
)

type testReflectionObj struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type testReflectionNested struct {
	ID      int                `json:"id"`
	JoinVal testReflectionObj  `json:"joinval"`
	JoinPtr *testReflectionObj `json:"joinptr"`
}

func TestSQLxReflectionValStruct(t *testing.T) {

	// should be a passed by value struct
	if !byValueStruct(reflect.TypeOf(testReflectionObj{})) {
		t.Error("unexpected error - should be flagged as a struct by value")
	}
	// should not be true, its passed by reference / ptr
	if byValueStruct(reflect.TypeOf(&testReflectionObj{})) {
		t.Error("unexpected error - should be flagged as a ptr")
	}
	// check random other type
	if byValueStruct(reflect.TypeOf(0)) {
		t.Error("unexpected error - int should be seen as a struct")
	}

}

func TestSQLxReflectionPtrStruct(t *testing.T) {
	// struct passed by ptr should be true
	if !byPtrStruct(reflect.TypeOf(&testReflectionObj{})) {
		t.Error("unexpected error - should be flagged as a ptr")
	}
	// value should not be viewed as a ptr
	if byPtrStruct(reflect.TypeOf(testReflectionObj{})) {
		t.Error("unexpected error - struct by value should not be flagged as ptr")
	}
	// check random other type
	if byPtrStruct(reflect.TypeOf(0)) {
		t.Error("unexpected error - int should be seen as a struct")
	}

}

func TestSQLxReflectionSrc(t *testing.T) {

	r := newReflection(&testReflectionObj{})
	if r.Src == nil {
		t.Errorf("unexpected error - reflection src is empty")
	}

}

func TestSQLxReflectionIsStruct(t *testing.T) {
	var r *reflection

	r = newReflection(&testReflectionObj{})
	if !r.IsStruct {
		t.Error("unexpected error - should be a struct")
	}
	r = newReflection(testReflectionObj{})
	if !r.IsStruct {
		t.Error("unexpected error - should be a struct")
	}

	r = newReflection(0)
	if r.IsStruct {
		t.Error("unexpected error - should not be a struct")
	}

}

func TestSQLxReflectionFields(t *testing.T) {
	var r *reflection

	r = newReflection(&testReflectionNested{})
	if len(r.Fields()) != 3 {
		t.Errorf("unexpected number of fields")
	}

	r = newReflection(testReflectionNested{})
	if len(r.Fields()) != 3 {
		t.Errorf("unexpected number of fields")
	}

	r = newReflection(0)
	if len(r.Fields()) > 0 {
		t.Errorf("unexpected fields found")
	}

}

func TestSQLxReflectionRecursiveFields(t *testing.T) {
	var all []reflect.StructField
	var r *reflection

	r = newReflection(&testReflectionNested{})
	all = r.RecursiveFields(nil)
	// the nested struct has 3 fields, two of which are
	// structs that have 2 fields each - so there should
	// be 7 fields
	if len(all) != 7 {
		t.Error("unexpected error - incorrect number of fields returned")
	}

	r = newReflection(testReflectionNested{})
	all = r.RecursiveFields(nil)
	if len(all) != 7 {
		t.Error("unexpected error - incorrect number of fields returned")
	}

	r = newReflection(0)
	all = r.RecursiveFields(nil)
	if len(all) > 0 {
		t.Error("unexpected error - int should not return any fields")
	}

}
