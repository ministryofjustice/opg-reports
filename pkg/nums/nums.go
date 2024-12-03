package nums

import "github.com/ministryofjustice/opg-reports/internal/strutils"

type nums interface {
	float32 | float64 | int
}
type adders interface {
	nums | string
}

func add[T adders](a T, args ...any) (result T) {
	result = a
	for _, arg := range args {
		// check the T casting works before doing the +
		if val, ok := arg.(T); ok {
			result += val
		}
	}
	return
}

// Add handles "adding" floats, ints and strings being added.
//
// For strings, it will try to treat them as floats first (via
// `addString`) but if that fails due to parsing errors it will
// instead concatenate them (via `add`).
//
// Examples:
//
//	Add(1, 2, 3 ) 	// 6
//	Add("1", "2")	// 3
//	Add(1.0, 2.0)	// 3.0
//	Add("A", "b")	// "Ab"
func Add(a interface{}, args ...interface{}) (result interface{}) {
	switch a.(type) {
	case float64:
		result = add(a.(float64), args...)
	case int:
		result = add(a.(int), args...)
	case string:
		v, err := strutils.Add(a.(string), args...)
		if err != nil {
			result = add(a.(string), args...)
		} else {
			result = v
		}
	default:
		result = ""
	}
	return
}
