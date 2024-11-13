package statements_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs/statements"
)

func TestStatementNames(t *testing.T) {
	var checks = map[statements.Stmt][]string{
		"SELECT * FROM test":                                  {},
		"SELECT * FROM test WHERE id = :id":                   {"id"},
		"SELECT * FROM test WHERE id = :id_withWeird-Naming1": {"id_withWeird-Naming1"},
	}

	for check, expected := range checks {
		var actual = check.Names()
		// check the length
		if len(expected) != len(actual) {
			t.Errorf("names result length mismatch - expected [%d] actual [%v]", len(expected), len(actual))
		}
		// check the content
		for _, exp := range expected {
			var found = false
			for _, act := range actual {
				if exp == act {
					found = true
				}
			}
			if !found {
				t.Errorf("failed to find value in results - expected [%s]", exp)
			}
		}

	}
}

type testStmtValidateFixture struct {
	stmt  statements.Stmt
	named statements.Named
	err   error
}
type testNames struct {
	ID string `json:"id"`
}

func TestStatementValidate(t *testing.T) {
	var checks = []testStmtValidateFixture{
		{stmt: "SELECT * from test where id = 1"},
		{stmt: "SELECT * from test where id = :id", named: &testNames{ID: "1"}},
		{stmt: "SELECT * from test where id = :id", err: fmt.Errorf("statement validation failed; missing the following named parameters: [%s]", "id")},
	}

	for _, test := range checks {
		err := statements.Validate(test.stmt, test.named)
		if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("statement validation error:\nexpected\n[%v]\nactual\n[%v]", test.err, err)
		}
		if test.err == nil && err != nil {
			t.Errorf("unexpected error returned [%v]", err.Error())
		}
	}
}
