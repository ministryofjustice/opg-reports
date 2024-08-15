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

func First[T any](v []T) (t T) {
	if len(v) > 0 {
		t = v[0]
	}
	return
}

func FirstOrDefault[T any](v []T, defaultValue T) (t T) {
	if len(v) > 0 {
		t = v[0]
	} else {
		t = defaultValue
	}
	return
}
