package sqlx

import (
	"context"
	"errors"
	"fmt"
	"opg-reports/app/internal/convert"
	"opg-reports/app/internal/fmtx"
	"reflect"
	"regexp"
	"strings"
)

var (
	ErrModelNotStruct   error = errors.New("filterModel passed to Bind is not a struct.")
	ErrModelMissingTags error = errors.New("filterModel is missing json tags on fields.")
)

var (
	ErrBindingNoKey error = errors.New("error when binding sql statement to filterModel, a parameter (`:x`) was found that has no coresponding value on filterModel.")
)

// Bind finds `:x` placeholders within sql statement and generates a set of values
// while also updating the sql statement to match the required bound parameters usage.
//
// Validates that the filterModel passed is a struct (or ptr to a struct) and that
// all of its fields (and sub struct fields) have json tagging set. Will
// return an error if not
//
// WARNING: Validation uses reflection to check all fields within the model have json
// tags. However the call to generate binding parameters to replace `:x` values
// uses json conversion.
func Bind(ctx context.Context, statement string, filterModel any) (err error) {
	var (
		reflected   *Reflection
		boundValues map[string][]interface{}
		// generatedStmt string = statement
		// args          []interface{}
	)
	// setup the reflection and the validate the model
	reflected = NewReflection(filterModel)
	err = validateModel(reflected)
	if err != nil {
		return
	}
	// get the values from the model into a flat key map
	boundValues, err = bindingValues(reflected)
	if err != nil {
		return
	}

	fmtx.Printj(boundValues)
	return
}

// modelBoundSql is used within the Bind call to parse and update the
// sql statement, replacing each `:$x` with correct number of `?` and pushing
// values into the orderedArgs list
type modelBoundSql struct {
	Statement   string
	OrderedArgs []interface{}
	values      map[string][]interface{}
	re          *regexp.Regexp
}

func (self *modelBoundSql) Parameterise() (err error) {
	var (
		key     string
		stmt    string = self.Statement
		matches []int  = self.re.FindStringIndex(stmt)
	)
	// no matches, nothing to do
	if len(matches) < 2 {
		return
	}
	// match.. so lets find things...
	key = stmt[matches[0]:matches[1]]
	// if the key is not in the bound value data then return an error
	if _, ok := self.values[key]; !ok {
		err = errors.Join(ErrBindingNoKey, fmt.Errorf("the extra parameter found in the sql statement was [%s]", key))
		return
	}

	// now replace this segment of the string with current values
	stmt = self.replaceStmtSegment(stmt, matches, self.values[key])
	// now add the values to the ordered arguments
	self.OrderedArgs = append(self.OrderedArgs, self.values[key]...)
	// now recurse....
	self.Statement = stmt
	err = self.Parameterise()

	return
}

// replaceStmtSegment replaces the `:$x` string segment found at loc[0]..loc[1] with the correct
// number of `?` that should be used for the sql prepared statement to work
func (self *modelBoundSql) replaceStmtSegment(stmt string, loc []int, values []interface{}) string {
	var (
		i int    = loc[0]
		j int    = loc[1]
		s string = ""
	)
	for i := 0; i < len(values); i++ {
		s += "?,"
	}
	s = strings.TrimSuffix(s, ",")
	return stmt[0:i] + s + stmt[j:]
}

func newBoundSql(statement string, values map[string][]interface{}) *modelBoundSql {
	return &modelBoundSql{
		OrderedArgs: []interface{}{},
		Statement:   statement,
		values:      values,
		re:          regexp.MustCompile(`(?m):[[:alnum:]_-]+`),
	}
}

// validateModel is used by bind to confirm the model is
// a struct and has json tagged attributes.
func validateModel(ref *Reflection) (err error) {

	if !ref.IsStruct {
		err = ErrModelNotStruct
		return
	}

	if !AllFieldsHaveJSONTags(ref) {
		err = ErrModelMissingTags
		return
	}
	return
}

// jsonBindings recursively walks over a map generated from converting a struct
// into a map[string]interface{} and appends key / value combinations to the flat
// valueMap result
//
// Used within bindings to generate a lookup of field names that will be within the
// sql statement
func jsonBindings(jsonified map[string]interface{}, valueMap map[string]interface{}) {

	for k, v := range jsonified {
		// handle empty slices etc
		if v == nil {
			continue
		}
		var ty = reflect.TypeOf(v).Kind()
		var key = fmt.Sprintf(":%s", k)

		if ty == reflect.Map {
			jsonBindings(v.(map[string]interface{}), valueMap)
		} else {
			valueMap[key] = v
		}
	}
}

// bindingValues is used to create a lookup of key/values to allow Bind to find the
// replacement values quickly, so the paramBindings use `:jsonTagName` as keys
//
// bindings makes every value a slice, for consistency within the sql parsing so should
// always look convert:
//
//	struct{
//		ID:        1,
//		Address: "address test line 1",
//		Emails: []string{
//				"test@example.com",
//				"test2@example.com",
//		},
//	}
//
// into:
//
//	{
//		":address": [ "address test line 1" ],
//		":email": [ "test@example.com", "test2@example.com" ],
//		":id": [ 1 ]
//	}
func bindingValues(ref *Reflection) (paramBindingSlices map[string][]interface{}, err error) {
	var jsonified map[string]interface{} = map[string]interface{}{}
	var paramBindings = map[string]interface{}{}
	paramBindingSlices = map[string][]interface{}{}

	err = convert.Convert(ref.Src, &jsonified)
	if err != nil {
		return
	}
	jsonBindings(jsonified, paramBindings)
	// convert into slices if they arent
	for k, v := range paramBindings {
		paramBindingSlices[k] = asSlice(v)
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
