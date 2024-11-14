package must

// Must is wrapper for dealing with single value returns and error combinations
// from functions that will panic on error
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
