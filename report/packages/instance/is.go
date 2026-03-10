package instance

import "reflect"

// IsPtr returns true if T is a pointer
func IsPtr[T any](item T) bool {
	return reflect.ValueOf(item).Kind() == reflect.Ptr
}
