package api

import (
	"fmt"
	"strings"
)

// various end of line strings for different part of the statements
const (
	selectEOL  string = ",\n"
	whereEOL   string = " AND \n"
	groupByEOL string = ",\n"
	orderByEOL string = ",\n"
)

// Field is used represent an item with in the SQL statement that we want to generate.
//
// It is used to handle more dynamic SQL generation that handles both fixed values
// such as a SUM or COUNT but also handle the different between including a field in
// the SELECT to that field being part of a WHERE or GROUP clause.
//
// # Examples
//
// A fFeld that should always be present in the end SQL:
//
//	 &Field{
//			Key: 		"count",
//			Select:		"COUNT(*) as counter",
//	 }
//
// A Field that will always be present and used to group data by the month and is
// used in ordering:
//
//	 &Field{
//			Key: 		"date",
//			Select:		"strftime('%Y-%m', date) as date",
//			GroupBy: 	"strftime('%Y-%m', date)",
//			OrderBy:	"strftime('%Y-%m', date) ASC",
//
//	}
//
// A Field that will always be present and used to group data by the month, used
// for ordering and will filter the data based on fixed months
//
//	 &Field{
//			Key: 		"date",
//			Select:		"strftime('%Y-%m', date) as date",
//			Where: 		"(date >= '2025-01' AND date <= '2025-06')",
//			GroupBy: 	"strftime('%Y-%m', date)",
//			OrderBy: 	"strftime('%Y-%m', date) ASC",
//
//	}
//
// See `BuildGroupSelect` func for how these work to generate a complete SQL statement.
type Field struct {
	Key     string  // used to help identify the database / request field this relates to
	Value   *string // allows a specific value to be passed which will decide if the field is included in hte SQL or not.
	Select  string  // contains how the field should be selected - this can include date formatting, counts, aliases etc
	Where   string  // contains how this field should be handled in a WHERE - it can be a multi conditions
	GroupBy string  // contains how this field should used in GROUP BY
	OrderBy string  // contains how this field should be used in ORDER BY
}

// isQueryParameterSelectable used to decide if this value should be part of the
// column select
//
// If the pointer isnt nil and the value is not equilivant to empty, then it should
// be used in the select
func isQueryParameterSelectable(value *string) bool {
	return (value != nil && *value != "")
}

// isQueryParameterWhereable determines if the field should be part of a where
// statement
//
// If the pointers value is set to something other that "true" than it is a where
// query
func isQueryParameterWhereable(value *string) bool {
	return (value != nil && *value != "" && strings.ToLower(*value) != "true")
}

// isQueryParameterGroupable is used to decide if the value means this field
// should be used in the group by section.
//
// If the value is "true" then it should be used in group by
func isQueryParameterGroupable(value *string) bool {
	return (value != nil && strings.ToLower(*value) == "true")
}

// isQueryParameterOrderable is used to decide if the value means the field
// should be part of the order by section.
//
// When the value is "true", it should be part of the order by
func isQueryParameterOrderable(value *string) bool {
	return (value != nil && strings.ToLower(*value) == "true")
}

// generateSelect checks each Field to see if it should be included in the
// SELECT section of the SQL statement and single string containing all
// fields that should be in the SELECT.
//
// It decides if a field should be in the SELECT by checking if the Field.Value
// is `nil` or `isQueryParameterSelectable`. Will only be included when
// there is something in `.Select`
//
// When a field is included, but has a `nil` Value, then this field should
// always be used within the select - and is typically a expression such as
// `COUNT(*) as cnt` that is always required.
func generateSelect(fields ...*Field) (str string) {
	var eol string = selectEOL
	for _, field := range fields {
		var set = (field.Value == nil || isQueryParameterSelectable(field.Value))
		if stmt := field.Select; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)
	return
}

// generateWhere checks each Field to see if it should be included in the
// WHERE section of the SQL statement and single string containing all
// fields that should be in the WHERE is returned - generate from each
// Field.Where
//
// It decides if a field should be in the WHERE by checking if the Field.Value
// is `nil` or `isQueryParameterWhereable`. Will only be included when
// there is something in `.Where`
//
// When a Field as `nil` Value, then this field should always be included.
func generateWhere(fields ...*Field) (str string) {
	var eol string = whereEOL
	for _, field := range fields {
		var set = (field.Value == nil || isQueryParameterWhereable(field.Value))
		if stmt := field.Where; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

// generateGroupBy checks each Field to see if it should be included in the
// GROUP BY section of the SQL statement and single string containing all
// fields that should be in the GROUP BY is returned - generated from each
// Field.GroupBy
//
// It decides if a field should use by checking if the Field.Value is `nil`
// or `isQueryParameterGroupable`. Will only be included when there is
// something in `.GroupBy`
//
// When a Field.Value is `nil`, then this field should always be included.
func generateGroupBy(fields ...*Field) (str string) {
	var eol string = groupByEOL
	for _, field := range fields {
		var set = (field.Value == nil || isQueryParameterGroupable(field.Value))
		if stmt := field.GroupBy; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

// generateOrderBy checks each Field to see if it should be included in the
// ORDER BY section of the SQL statement and single string containing all
// fields that should be in the ORDER BY is returned - generated from each
// Field.OrderBy
//
// It decides if a field should use by checking if the Field.Value is `nil`
// or `isQueryParameterOrderable`. Will only be included when there is
// something in `.OrderBy`
//
// When a Field.Value is `nil`, then this field should always be included.
func generateOrderBy(fields ...*Field) (str string) {
	var eol string = orderByEOL
	for _, field := range fields {
		var set = (field.Value == nil || isQueryParameterOrderable(field.Value))
		if stmt := field.OrderBy; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

// BuildSelectFromFields is used to generate a dynamic SQL statement based on fixed fields
// and values that have come from API requests.
//
// The API allows fields like `account` which can be either empty (so not used),
// `true` (acting as a group by) or a real value like `account-01` (which would filter
// the data).
//
// To simplify the code, SQL is constructed from multiple Fields, each of which can
// either be a always required, or only used depending on the `Value` (from the API
// request).
//
// This allows the SQL to flex to include only the fields and clauses required
// based on the configured Fields
func BuildSelectFromFields(table string, joins string, fields ...*Field) (sql string) {
	var stmt string = `
SELECT
%s
FROM %s
%s
WHERE
%s
GROUP BY
%s
ORDER BY
%s
;`
	sql = fmt.Sprintf(stmt,
		generateSelect(fields...),
		table,
		joins,
		generateWhere(fields...),
		generateGroupBy(fields...),
		generateOrderBy(fields...),
	)

	return
}
