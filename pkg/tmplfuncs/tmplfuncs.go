// Package tmplfuncs contains series of funcs to use within
// front end templates
//
// Exposed as a map (`All`)
package tmplfuncs

type adders interface {
	float32 | float64 | int | string
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

// Add handles "adding" floats, ints and strings
// - string is concatenated without spaces
// - any `args` of a (type) different to `a` are ignored
// - if `a` is not float64, int or string, "" is returned
func Add(a interface{}, args ...interface{}) (result interface{}) {
	switch a.(type) {
	case float64:
		result = add(a.(float64), args...)
	case int:
		result = add(a.(int), args...)
	case string:
		result = add(a.(string), args...)
	default:
		result = ""
	}
	return
}

var All map[string]interface{} = map[string]interface{}{
	"add": Add,
}
