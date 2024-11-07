package datastore

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ministryofjustice/opg-reports/pkg/convert"
)

// Exec is a string used that contains a
// sql command such as CREATE TABLE or similar
// that causes a change, but returns no value
type ExecStatement string

// CreateStatement is a subtype of ExecStatement
// specifically for running create operations
type CreateStatement ExecStatement

// InsertStatement is a string used to run INSERT
// operations against the database
type InsertStatement string

// SelectStatement is a string used as enum-esque
// type contraints for sql queries that contain SELECT
// operations
type SelectStatement string

// NamedSelectStatement is a SELECT operation that
// contains :name placeholders
type NamedSelectStatement SelectStatement

// NamedParameters are structs with fields that will be converted into named values
// within statements
type NamedParameters interface{}

// Needs is used in part of the validate check of the named parameters and returns
// the field names the NamedSelectStatement passed in should have
// Uses a regex to find words starting with :
func Needs(query NamedSelectStatement) (needs []string) {
	var namedParamPattern string = `(?m)(:[\w-]+)`
	var prefix string = ":"
	var re = regexp.MustCompile(namedParamPattern)
	for _, match := range re.FindAllString(string(query), -1) {
		needs = append(needs, strings.TrimPrefix(match, prefix))
	}
	return
}

// ValidateParameters checks if the parameters passed meets all the required
// needs for the query being run
func ValidateParameters[P NamedParameters](params P, needs []string) (err error) {
	mapped, err := convert.Map(params)
	if err != nil {
		return
	}

	missing := []string{}
	// check each need if that exists as a key in the map
	for _, need := range needs {
		if _, ok := mapped[need]; !ok {
			missing = append(missing, need)
		}
	}
	// if any field is missing then set error
	if len(missing) > 0 {
		cols := strings.Join(missing, ",")
		err = fmt.Errorf("missing required fields for this query: [%s]", cols)
	}

	return
}
