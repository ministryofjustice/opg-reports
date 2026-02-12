// Package query is used to generate a sql statment that contains named vars (`:x“) and compatible
// with the dbselect package
//
// Mainly used in combination with an incoming request (converts to a `map[string]interface{}`) and
// map of request files to sql segments to convert incoming api parameters into a usable sql
// statment. Input vars are not directly used (assumes bind vars etc)
package query

import (
	"fmt"
	"strings"
)

// base is the structure of a select
const baseStmt string = `
{select}
FROM {from}
{joins}
{where}
{group_by}
{having}
{order_by}
;`

// Type represents the selects of the sql statment this will generate
type Type string

const (
	SELECT  Type = "select"
	WHERE   Type = "where"
	GROUPBY Type = "group_by"
	HAVING  Type = "having"
	ORDERBY Type = "order_by"
)

// QuerySegment is used to build a part of the overall sql statment
// and depending on the `Type` will be used in the sql clause differently.
//
// For example, in WHERE clause segements are joined using `AND` but in
// GROUP BY they are joined via `,`
type QuerySegment struct {
	Type Type
	Stmt string
}

// Query is the struct used to generate a string that can be used as a named
// statement (containing bind var syntax of `:x`).
//
// Generally combination of mapped query segments and inpur request is used to create
// each part of the statement.. see the test for an example.
//
// Notes:
//
//	`From` attribute is used as `FROM {from}` replacement.
//	`Joins` are merged together to create the set of join statements after the `FROM` - as these can be various types this struct does not add the JOIN notation at the start like ti does for FROM.
type Query struct {
	From     string
	Joins    []string
	Segments map[string][]*QuerySegment
}

// FromRequest builds a sql statement from the reuest, mapping each input key against the `Segments` to generate
// each clause.
//
// See test for examples.
func (self *Query) FromRequest(request map[string]interface{}) (query string, blocks map[Type][]string) {
	var qs = newQueryStr(baseStmt, self)

	for key, value := range request {
		for _, segment := range self.Segments[key] {
			qs.Add(segment, value)
		}
	}
	blocks = qs.Blocks()
	query = qs.String()
	return
}

type queryStr struct {
	base string
	q    *Query
	strs map[Type][]string
}

// Add inserts query segments into the correct category.
//
// If the qs.Stmt contains a `:` then its presumed to be a filter (`where x = :value`) so
// it will only be added if the `val` is set to a real value (not "true"). This allows
// filters to be handled on where / having etc
func (self *queryStr) Add(qs *QuerySegment, val interface{}) {
	var (
		hasColon bool   = strings.Contains(qs.Stmt, ":")
		add      bool   = false
		v        string = strings.ToLower(val.(string))
	)
	// if there is no colon, its not a filter statement, so add directly
	// otherwise only add if it has a filterable value
	if !hasColon || (hasColon && v != "" && v != "true") {
		add = true
	}
	if add {
		// set the default
		if _, ok := self.strs[qs.Type]; !ok {
			self.strs[qs.Type] = []string{}
		}
		self.strs[qs.Type] = append(self.strs[qs.Type], qs.Stmt)
	}

}

// Blocks returns the internal strs set
func (self *queryStr) Blocks() map[Type][]string {
	return self.strs
}

// String will return the generate sql statement from all added blocks
func (self *queryStr) String() (stmt string) {
	stmt = self.base
	// generate the string from the current set of slices
	for k, values := range self.strs {
		var (
			joined string = ""
			eol    string = ",\n" // end of line is normally a , but not for where or having so adjust
			key    string = fmt.Sprintf(`{%s}`, string(k))
			prefix string = strings.ReplaceAll(strings.ToUpper(string(k)), "_", " ") + "\n"
		)
		if k == WHERE || k == HAVING {
			eol = " AND\n"
		}
		if len(values) > 0 {
			joined = prefix + strings.TrimSuffix(strings.Join(values, eol), eol)
		}
		stmt = strings.ReplaceAll(stmt, key, joined)
	}
	// remove empty lines and trailing lines
	stmt = strings.ReplaceAll(stmt, "\n\n", "\n")
	stmt = strings.TrimPrefix(stmt, "\n")
	stmt = strings.TrimSuffix(stmt, "\n")
	// proces t he from & join clauses
	stmt = strings.ReplaceAll(stmt, `{from}`, self.q.From)
	stmt = strings.ReplaceAll(stmt, `{joins}`, strings.Join(self.q.Joins, "\n"))
	return
}

func newQueryStr(sql string, q *Query) *queryStr {
	return &queryStr{
		q:    q,
		base: sql,
		strs: map[Type][]string{},
	}
}
