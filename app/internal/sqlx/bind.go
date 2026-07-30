package sqlx

import (
	"context"
	"errors"
	"fmt"
	"opg-reports/app/internal/convert"
	"reflect"
	"regexp"
	"strings"
)

var (
	ErrModelNotStruct   error = errors.New("filterModel passed to Bind is not a struct.") // raised when a non-struct is passed to `Bind`
	ErrModelMissingTags error = errors.New("filterModel is missing json tags on fields.") // raised when the struct passed to `Bind` has fields without json tags
)

var (
	ErrBindingNoKey error = errors.New("error when binding sql statement to filterModel, a parameter (`:x`) was found that has no coresponding value on filterModel.") // raised when the sql passed to `Bind` contains a placeholder that does not exist on the filterModel.
)

// Bind finds the values of json tagged fields on the filterModel struct and
// uses those to replace matching placeholders within the `sqlStmt` ready to
// be used in a `db.ExecContext` or `db.QueryContext` call.
//
// The intention is to allow our local code to use prepared sql statements
// that contain placeholders that map to values on a struct which can then
// be passed into Bind to then generate sql - avoiding complex logic within
// the main loop and ant ORM usage.
//
// # Example usage
//
//	type filterByMonth struct {
//	    Months []string `json:"months"`
//	}
//	var filter = &filterByMonth{
//	    Months: []string{"2026-05", "2026-06"}
//	}
//	var selectStmt string = "SELECT * FROM costs WHERE month IN(:months);"
//	sql, args, err := Bind(ctx, selectStmt, filter)
//	rows, err := db.QueryContext(ctx, sql, args...)
//
// Result in `sql` being returned as:
//
//	"SELECT * FROM costs WHERE month IN(?,?);"
//
// And `args` will then contain the values required in the correct order:
//
//	[]interface{"2026-05", "2026-06"}
//
// # Validation / errors
//
// Bind will validate that the `filterModel` passed is either a struct or
// a point to a struct.
//
// The `filterModel` fields (and any embedded struct fields) will be
// checked to ensure they have a json tag (`json:"$name"`) set. This is
// used for determining the names and values of the sql placeholders
// (`:$name`).
//
// Validation uses reflection to check all visible fields, but the
// value replacement utilises json marshaling.
func Bind(ctx context.Context, sqlStmt string, filterModel any) (sql string, args []interface{}, err error) {
	var (
		reflected   *reflection              // used to inspect the struct for its type and field data
		mbSql       *modelBoundSql           // used to recursively update the sql and ordered args form the `sqlStmt` and `filterModel`
		boundValues map[string][]interface{} // the flat map of field names and values
	)
	args = []interface{}{}
	// setup the reflection and the validate the model
	reflected = newReflection(filterModel)
	err = validateModel(reflected)
	if err != nil {
		return
	}
	// get the values from the model into a flat key map
	boundValues, err = bindingValues(reflected)
	if err != nil {
		return
	}

	// create the modelBoundSql instance which will recursively replace
	// field placeholders with ? and provide the ordered arguments
	mbSql = newBoundSql(sqlStmt, boundValues)
	err = mbSql.Parameterise()
	if err != nil {
		return
	}
	// grab the return values
	sql, args = mbSql.Values()

	return
}

// modelBoundSql is used within the Bind call to parse and update the
// sql statement, replacing each `:$x` placeholder with correct number
// of `?` and pushing values into the orderedArgs list
type modelBoundSql struct {
	Statement   string
	OrderedArgs []interface{}
	values      map[string][]interface{}
	re          *regexp.Regexp
}

// Parameterise recursively calls itself and replaces the first matching
// field placeholder with the correct number of `?` and pushes the values
// of that field (found in `.values`) into the `OrderedArgs`.
//
// Returns an error if the field placeholder (`:$name`) it finds does
// not exist in the `.values` slice.
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
	// match.. so lets grab the field placeholder ...
	key = stmt[matches[0]:matches[1]]
	// if the key is not in the bound value data then return an error
	if _, ok := self.values[key]; !ok {
		err = errors.Join(
			ErrBindingNoKey,
			fmt.Errorf("the extra parameter found in the sql statement was [%s]", key))
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

// Values is shorthand to fetch the end result sql and ordered arguments
func (self *modelBoundSql) Values() (sql string, args []interface{}) {
	return self.Statement, self.OrderedArgs
}

// replaceStmtSegment replaces the `:$x` string segment found at loc[0]..loc[1]
// with the correct number of `?` that should be used for the sql prepared
// statement to work
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

// newBoundSql generates
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
func validateModel(ref *reflection) (err error) {

	if !ref.IsStruct {
		err = ErrModelNotStruct
		return
	}

	if !allFieldsHaveJSONTags(ref) {
		err = ErrModelMissingTags
		return
	}
	return
}

// bindingValues is used to create a lookup of key/values to allow Bind to find the
// replacement values quickly, so the paramBindings use `:jsonTagName` as keys.
//
// bindings makes every value a slice for consistency within the sql parsing and
// handling or ordered arguments.
//
// # Example
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
// Results in:
//
//	{
//		":address": [ "address test line 1" ],
//		":email": [ "test@example.com", "test2@example.com" ],
//		":id": [ 1 ]
//	}
func bindingValues(ref *reflection) (paramBindingSlices map[string][]interface{}, err error) {
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
