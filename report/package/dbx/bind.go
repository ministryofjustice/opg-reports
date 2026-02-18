package dbx

import (
	"context"
	"errors"
	"fmt"
	"opg-reports/report/package/cntxt"
	"reflect"
	"regexp"
	"strings"
)

var ErrBadPlaceHolder = errors.New("found a placeholder syntax (:name) without a value to replace it with.")
var re *regexp.Regexp = regexp.MustCompile(`(?m):[[:alnum:]_-]+`)

// Bind replaces all `:name` elements within a sql statement with '?' and returns list of args to use for within Query / Exec
func Bind(ctx context.Context, named string, model map[string]interface{}) (statement string, args []interface{}, err error) {
	var modelValues map[string][]interface{}
	var log = cntxt.GetLogger(ctx).With("package", "bind", "func", "Bind")
	args = []interface{}{}
	// get all the values for the mode
	modelValues = map[string][]interface{}{}
	for key, value := range model {
		modelValues[key] = asSlice(value)
	}
	res := &reducer{Statement: named, ModelValues: modelValues, Args: args}
	err = reduceReplace(res)

	if strings.Contains(res.Statement, ":") {
		log.Warn("possibly placeholders left within sql")
	}

	return res.Statement, res.Args, err
}

type reducer struct {
	Statement   string
	ModelValues map[string][]interface{}
	Args        []interface{}
}

func (self *reducer) MatchedKey(matches [][]int) string {
	var i = matches[0][0]
	var j = matches[0][1]
	return self.Statement[i+1 : j]
}
func (self *reducer) ReplaceMatch(matches [][]int, str string) string {
	var i = matches[0][0]
	var j = matches[0][1]
	return self.Statement[0:i] + str + self.Statement[j:]
}

// reduceReplace recursive function to replace each part of the sql statement
func reduceReplace(reduce *reducer) (err error) {
	var matches = re.FindAllStringIndex(reduce.Statement, 1)

	if len(matches) == 0 {
		return
	}
	var (
		sub        = ""
		key        = reduce.MatchedKey(matches)
		values, ok = reduce.ModelValues[key]
	)
	// this field doesnt exist, so throw an error
	if !ok {
		fmt.Println("no data " + key)
		errors.Join(ErrBadPlaceHolder, fmt.Errorf("incorrect placeholder: [:%s]", key))
		return
	}
	// attached the values for the :name
	reduce.Args = append(reduce.Args, values...)
	// add ?
	for i := 0; i < len(values); i++ {
		sub += "?,"
	}
	sub = strings.TrimSuffix(sub, ",")
	// adjust the statement at the match point
	reduce.Statement = reduce.ReplaceMatch(matches, sub)
	// see if need to recurse
	matches = re.FindAllStringIndex(reduce.Statement, 1)
	if len(matches) > 0 {
		err = reduceReplace(reduce)
	}
	return
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
