// Package qb (query builder) is used to generate a sql statment that contains named vars (`:x“) and compatible
// with the dbselect package
//
// Mainly used in combination with an incoming request (converts to a `map[string]interface{}`) and
// map of request files to sql segments to convert incoming api parameters into a usable sql
// statment. Input vars are not directly used (assumes bind vars etc)
package qb

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
	JOIN    Type = "joins"
	GROUPBY Type = "group_by"
	HAVING  Type = "having"
	ORDERBY Type = "order_by"
)

func segmentTypes() []string {
	return []string{
		`select`,
		`where`,
		`joins`,
		`group_by`,
		`having`,
		`order_by`,
	}
}

// Segment is used to build a part of the overall sql statment
// and depending on the `Type` will be used in the sql clause differently.
//
// For example, in WHERE clause segements are joined using `AND` but in
// GROUP BY they are joined via `,`
type Segment struct {
	Type Type
	Stmt string
}

// Builder is the struct used to generate a string that can be used as a named
// statement (containing bind var syntax of `:x`).
//
// Generally combination of mapped query segments and inpur request is used to create
// each part of the statement.. see the test for an example.
//
// Notes:
//
//	`From` attribute is used as `FROM {from}` replacement.
//	`Joins` are merged together to create the set of join statements after the `FROM` - as these can be various types this struct does not add the JOIN notation at the start like ti does for FROM.
type Builder struct {
	From     string
	Segments map[string][]*Segment
}

// FromRequest builds a sql statement from the reuest, mapping each input key against the `Segments` to generate
// each clause.
//
// See test for examples.
func (self *Builder) FromRequest(request map[string]string) (query string, blocks map[Type][]string) {
	var qs = newBuilderStr(baseStmt, self)

	if def, ok := self.Segments["_default"]; ok {
		qs.Add(nil, def...)
	}

	for key, value := range request {
		for _, segment := range self.Segments[key] {
			qs.Add(&value, segment)
		}
	}
	blocks = qs.Blocks()
	query = qs.String()
	return
}

// New
func New(from string, segments map[string][]*Segment) *Builder {
	return &Builder{
		From:     from,
		Segments: segments,
	}
}

type queryStr struct {
	base      string
	q         *Builder
	sqlBlocks map[Type][]string
}

// Add inserts query segments into the correct category.
//
// If the qs.Stmt contains a `:` then its presumed to be a filter (`where x = :value`) so
// it will only be added if the `val` is set to a real value (not "true"). This allows
// filters to be handled on where / having etc
func (self *queryStr) Add(value *string, segments ...*Segment) {
	var lowerValue string = ""
	var add []*Segment = []*Segment{}

	// if val is exactly nil then add everything - used to enable defaults
	// otherwise decide which should be added based on type
	if value == nil {
		add = segments
	} else {
		lowerValue = strings.ToLower(*value)
		for _, segment := range segments {
			var hasColon bool = strings.Contains(segment.Stmt, ":")
			var isFilter bool = (hasColon && lowerValue != "true")

			if lowerValue != "" && (!hasColon || isFilter) {
				add = append(add, segment)
			}
		}
	}
	// add selected segments
	for _, segment := range add {
		if _, ok := self.sqlBlocks[segment.Type]; !ok {
			self.sqlBlocks[segment.Type] = []string{}
		}
		self.sqlBlocks[segment.Type] = append(self.sqlBlocks[segment.Type], segment.Stmt)
	}

}

// Blocks returns the internal sqlBlocks set
func (self *queryStr) Blocks() map[Type][]string {
	return self.sqlBlocks
}

// String will return the generate sql statement from all added blocks
func (self *queryStr) String() (stmt string) {
	stmt = self.base

	for blockType, blocks := range self.sqlBlocks {
		var (
			joined  string = ""
			eol     string = ""
			prefix  string = ""
			lower   string = strings.ToLower(string(blockType))
			upper   string = strings.ToUpper(lower)
			replace string = fmt.Sprintf(`{%s}`, lower)
		)
		// skip if there are no blocks to process
		if len(blocks) == 0 {
			continue
		}

		switch blockType {
		case WHERE, HAVING:
			prefix = strings.ReplaceAll(upper, "_", " ") + "\n"
			eol = " AND\n"
		case JOIN:
			prefix = ""
			eol = "\n"
		default:
			prefix = strings.ReplaceAll(upper, "_", " ") + "\n"
			eol = ",\n"
		}
		joined = strings.Join(blocks, eol)
		joined = prefix + strings.TrimSuffix(joined, eol)
		// update end string
		stmt = strings.ReplaceAll(stmt, replace, joined)
	}
	// remove any defaults
	for _, def := range segmentTypes() {
		stmt = strings.ReplaceAll(stmt, fmt.Sprintf(`{%s}`, def), "")
	}
	stmt = strings.ReplaceAll(stmt, "\n\n", "\n")
	stmt = strings.TrimPrefix(stmt, "\n")
	stmt = strings.TrimSuffix(stmt, "\n")
	// proces the from & join clauses
	stmt = strings.ReplaceAll(stmt, `{from}`, self.q.From)

	return
}

func newBuilderStr(sql string, q *Builder) *queryStr {
	return &queryStr{
		q:         q,
		base:      sql,
		sqlBlocks: map[Type][]string{},
	}
}
