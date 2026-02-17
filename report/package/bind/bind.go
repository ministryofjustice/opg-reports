package bind

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var ErrBadPlaceHolder = errors.New("found a placeholder syntax (:name) without a value to replace it with.")

var re *regexp.Regexp = regexp.MustCompile(`(?m):[[:alnum:]_-]+`)

// Bind replaces all `:name` elements within a sql statement with '?' and returns list of args to use for within Query / Exec
func Bind(named string, model map[string]interface{}) (statement string, args []interface{}, err error) {
	var (
		modelValues map[string][]interface{}
	)
	args = []interface{}{}
	// get all the values for the mode
	modelValues = map[string][]interface{}{}
	for key, value := range model {
		modelValues[key] = asSlice(value)
	}

	statement = named
	statement, args, err = reduceReplace(statement, modelValues, args)

	return
}

// reduceReplace recursively removes each instance or `:name`
//
// If there is no value, the place holder is fully removed
func reduceReplace(statement string, modelValues map[string][]interface{}, args []interface{}) (string, []interface{}, error) {
	var err error
	var matches = re.FindAllStringIndex(statement, 1)

	if len(matches) >= 1 {
		sub := ""
		idx := matches[0]
		i := idx[0]
		j := idx[1]
		k := statement[i+1 : j] // removes the `:`
		// if we find values, replace them
		if values, ok := modelValues[k]; ok {
			// add the args
			args = append(args, values...)
			// replace the chunk
			for x := 0; x < len(values); x++ {
				sub += "?,"
			}
			sub = strings.TrimSuffix(sub, ",")
			// reform the string without this chunk
			statement = statement[0:i] + sub + statement[j:]
		} else {
			err = errors.Join(ErrBadPlaceHolder, fmt.Errorf("incorrect placeholder: [:%s]", k))
		}
		if err == nil && len(matches) > 1 {
			statement, args, err = reduceReplace(statement, modelValues, args)
		}
	}
	return statement, args, err
}

// asSlice uses reflection to expand val into multiples if its a p[slice etc]
func asSlice[T any](val T) (values []interface{}) {
	var v reflect.Value
	var t reflect.Type
	values = []interface{}{}

	v = reflect.ValueOf(val)
	t = v.Type()
	if t.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface().(T))
		}
	} else {
		values = append(values, val)
	}
	return
}

// generateBindVar returns a string with '?' bind var markers for the value depending on length.
//
// For a slice, it returns length instances of '?', otherwise, just one
func generateBindVar(value interface{}) (s []string) {

	var v reflect.Value
	var t reflect.Type
	s = []string{}

	if value == nil {
		return
	}
	v = reflect.ValueOf(value)
	t = v.Type()

	if t.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			s = append(s, "?")
		}
	} else {
		s = append(s, "?")
	}

	return
}
