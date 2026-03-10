package instance

import (
	"reflect"
)

// Of will create a new instance of T
func Of[T any]() (result T) {
	var from T
	var t reflect.Type
	var instance reflect.Value

	// create result T in different way depending if its a pointer
	if IsPtr(from) {
		t = reflect.TypeOf(from).Elem()
		instance = reflect.New(t)
		result = instance.Interface().(T)
	} else {
		t = reflect.TypeOf(from)
		instance = reflect.New(t)
		result = instance.Elem().Interface().(T)
	}

	return
}

func Copy[T any](item T) (result T) {
	result = Of[T]()
	return
}
