package must

func Must[T any](v T, err error) (t T) {
	if err != nil {
		var n T
		t = n
	} else {
		t = v
	}
	return
}
