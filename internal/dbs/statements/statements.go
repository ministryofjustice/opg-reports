package statements

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// Named is a interface to use with structs for named parameters
// within a sql statement like `SELECT * from table WHERE id = :id`
// and the `id` would map a value on the struct
type Named interface{}

type Stmt string

// Names returns all the `:field` values within the stmt
// so they can be checked against the parameters etc
func (self Stmt) Names() (needs []string) {
	needs = Names(self)
	return
}

// Create is a wrapper around string for "CREATE" sql statements
type Create Stmt

// Select is a wrapper around string for "SELECT" sql statements
type Select Stmt

// Insert is a wrapper around string for "INSERT" sql statements
type Insert Stmt

// Names returns all the `:field` values within the stmt
// so they can be checked against the parameters etc
func Names(stmt Stmt) (needs []string) {
	var (
		query   string         = string(stmt)
		pattern string         = `(?m)(:[\w-]+)`
		prefix  string         = ":"
		re      *regexp.Regexp = regexp.MustCompile(pattern)
	)
	needs = []string{}
	for _, match := range re.FindAllString(string(query), -1) {
		needs = append(needs, strings.TrimPrefix(match, prefix))
	}
	return
}

// Validate checks that the named fields needed in stmt are present in
// the parameters passed.
func Validate(stmt Stmt, parameters Named) (err error) {
	var (
		missingFields string
		missing       []string               = []string{}
		mapped        map[string]interface{} = map[string]interface{}{}
		needs         []string               = Names(stmt)
	)
	// grab a map of the struct
	if mapped, err = structs.ToMap(parameters); err != nil {
		return
	}

	for _, field := range needs {
		if _, ok := mapped[field]; !ok {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		missingFields = strings.Join(missing, ",")
		err = fmt.Errorf("statement validation failed; missing the following named parameters: [%s]", missingFields)
	}

	return
}
