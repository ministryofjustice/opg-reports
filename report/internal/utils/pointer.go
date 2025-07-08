package utils

// Ptr converts an item to the pointer version of itself.
func Ptr[T any](item T) *T {
	var ptr = &item
	return ptr
}
