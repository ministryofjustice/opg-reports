package api

import (
	"fmt"
	"strings"
)

func IsQueryParameterSelectable(value *string) bool {
	return (value != nil && *value != "")
}
func IsQueryParameterWhereable(value *string) bool {
	return (value != nil && *value != "" && strings.ToLower(*value) != "true")
}
func IsQueryParameterGroupable(value *string) bool {
	return (value != nil && strings.ToLower(*value) == "true")
}
func IsQueryParameterOrderable(value *string) bool {
	return (value != nil && strings.ToLower(*value) == "true")
}

type Field struct {
	Key   string
	Value *string

	SelectAs  string
	WhereAs   string
	GroupByAs string
	OrderByAs string
}

func Selects(fields ...*Field) (str string) {
	var eol string = ",\n"
	for _, field := range fields {
		var set = (field.Value == nil || IsQueryParameterSelectable(field.Value))
		if stmt := field.SelectAs; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)
	return
}

func Wheres(fields ...*Field) (str string) {
	var eol string = " AND \n"
	for _, field := range fields {
		var set = (field.Value == nil || IsQueryParameterWhereable(field.Value))
		if stmt := field.WhereAs; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

func Groups(fields ...*Field) (str string) {
	var eol string = ",\n"
	for _, field := range fields {
		var set = (field.Value == nil || IsQueryParameterGroupable(field.Value))
		if stmt := field.GroupByAs; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

func Orders(fields ...*Field) (str string) {
	var eol string = ",\n"
	for _, field := range fields {
		var set = (field.Value == nil || IsQueryParameterOrderable(field.Value))
		if stmt := field.OrderByAs; stmt != "" && set {
			str += stmt + eol
		}
	}
	str = strings.Trim(str, eol)

	return
}

func BuildGroupSelect(table string, joins string, fields ...*Field) (sql string) {
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
		Selects(fields...),
		table,
		joins,
		Wheres(fields...),
		Groups(fields...),
		Orders(fields...),
	)

	return
}
