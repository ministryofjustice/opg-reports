package sqlx

import (
	"errors"
	"testing"
)

type bindTestInputs struct {
	SQL   string
	Model any
}

type bindTestOutputs struct {
	SQL  string
	Args []interface{}
	Err  error
}

type bindTestScenario struct {
	In       *bindTestInputs
	Expected *bindTestOutputs
}

var bindTestScenarios = []*bindTestScenario{
	// (success)
	// simple example that will swap out a single value from a flat
	// struct to use - a simple filter like team
	{
		Expected: &bindTestOutputs{
			SQL:  "SELECT * FROM activity_log WHERE username = ?",
			Args: []interface{}{"test-user"},
			Err:  nil,
		},
		In: &bindTestInputs{
			SQL: "SELECT * FROM activity_log WHERE username = :username",
			Model: &struct {
				Username string `json:"username"`
			}{Username: "test-user"},
		},
	},
	// (success)
	// example showing slices and how to use them for an IN style
	// query with other values
	{
		Expected: &bindTestOutputs{
			SQL:  "SELECT * FROM activity_log WHERE username IN(?,?) AND month = ?",
			Args: []interface{}{"test-user-a", "test-user-b", "2026-01"},
			Err:  nil,
		},
		In: &bindTestInputs{
			SQL: "SELECT * FROM activity_log WHERE username IN(:usernames) AND month = :month",
			Model: &struct {
				Usernames []string `json:"usernames"`
				Month     string   `json:"month"`
			}{Usernames: []string{"test-user-a", "test-user-b"}, Month: "2026-01"},
		},
	},
	// (success)
	// example testing repeated usage of the same field in more than one place mixed
	// with other filters to ensure order or args is correct
	{
		Expected: &bindTestOutputs{
			SQL:  "SELECT * FROM activity_log WHERE month = ? AND username IN(?,?) AND month = ?",
			Args: []interface{}{"2026-01", "test-user-a", "test-user-b", "2026-01"},
			Err:  nil,
		},
		In: &bindTestInputs{
			SQL: "SELECT * FROM activity_log WHERE month = :month AND username IN(:usernames) AND month = :month",
			Model: &struct {
				Usernames []string `json:"usernames"`
				Month     string   `json:"month"`
			}{Usernames: []string{"test-user-a", "test-user-b"}, Month: "2026-01"},
		},
	},
	// (error)
	// check that a non-struct model triggers the correct
	// validation errors
	{
		Expected: &bindTestOutputs{
			SQL:  "",              // empty as error will happen first
			Args: []interface{}{}, // empty as error will happen before parsing
			Err:  ErrModelNotStruct,
		},
		In: &bindTestInputs{
			SQL:   "SELECT * FROM users WHERE username = :user",
			Model: true,
		},
	},
	// (error)
	// check that a struct without js tagging on all of its fields will
	// trigger validation errors.
	// Use nested struct without a tag to check recursion of the field
	// processing
	{
		Expected: &bindTestOutputs{
			SQL:  "",              // empty as error will happen first
			Args: []interface{}{}, // empty as error will happen before parsing
			Err:  ErrModelMissingTags,
		},
		In: &bindTestInputs{
			SQL: "SELECT * FROM users WHERE username = :user",
			Model: &struct {
				Email  string `json:"email"`
				Person struct {
					Name string
				} `json:"person"`
			}{Email: "test@example.com"},
		},
	},
	// (error)
	// check error handling of a field that is not on the struct
	// is flagged with correct error type and no data is returned
	{
		Expected: &bindTestOutputs{
			SQL:  "",              // empty as error will happen first
			Args: []interface{}{}, // empty as error will happen before parsing
			Err:  ErrBindingNoKey,
		},
		In: &bindTestInputs{
			SQL: "SELECT * FROM users WHERE username = :user",
			Model: &struct {
				Email string `json:"email"`
			}{Email: "test@example.com"},
		},
	},
}

// TestSQLxBindScenarios runs over a mix of scenarios to check output aligns with
// what is expected of bind
func TestSQLxBindScenarios(t *testing.T) {
	// logx.Set(slog.LevelDebug)

	var ctx = t.Context()

	for i, test := range bindTestScenarios {
		var expected = test.Expected
		var actual = &bindTestOutputs{}

		actual.SQL, actual.Args, actual.Err = Bind(ctx, test.In.SQL, test.In.Model)

		if actual.SQL != expected.SQL {
			t.Errorf("[%d] actual SQL \n[%s]\n did not match expected \n[%s]\n", i, actual.SQL, expected.SQL)
		}

		if !errors.Is(actual.Err, expected.Err) {
			t.Errorf("[%d] actual Err \n[%s]\n did not match expected \n[%s]\n", i, actual.Err, expected.Err)
		}
		if len(actual.Args) != len(expected.Args) {
			t.Errorf("[%d] actual args differs in length from expected values.", i)
		}

		for j, ex := range expected.Args {
			var arg = actual.Args[j]
			if arg != ex {
				t.Errorf("[%d] actual Arg[%d] \n[%s]\n did not match expected \n[%s]\n", i, j, arg, ex)
			}
		}

	}
}

type testFilter struct {
	Team     string   `json:"team"`
	Months   []string `json:"months"`
	Services []string `json:"services"`
}

func TestSQLxParameteriseSimpleSelectFilter(t *testing.T) {
	var stmt = `SELECT * FROM costs WHERE team = :team AND month IN (:months) ;`
	var expectedsql = `SELECT * FROM costs WHERE team = ? AND month IN (?,?) ;`
	var filter = &testFilter{
		Team:   "test-a",
		Months: []string{"2026-06", "2026-05"},
	}
	var expectedArgs = []interface{}{
		"test-a", "2026-06", "2026-05",
	}
	var bound, _ = bindingValues(newReflection(filter))

	bsql := newBoundSql(stmt, bound)
	bsql.Parameterise()

	// check the sql statement
	actualsql, actualArgs := bsql.Values()
	if actualsql != expectedsql {
		t.Errorf("parsed sql does not match expected value.. \nexpected:\n[%s]\nactual:\n[%s]\n", expectedsql, actualsql)
	}
	// check the args
	matched := true
	for i, actual := range actualArgs {
		var expected = expectedArgs[i]
		if actual != expected {
			matched = false
		}
	}
	if !matched {
		t.Errorf("actual arguments were not as expected...")
	}

}

type testComms struct {
	Address string   `json:"address"`
	Emails  []string `json:"email"`
}
type testBindings struct {
	ID        int       `json:"id"`
	privateID int       `json:"-"`
	Comms     testComms `json:"comms"`
}

func TestSQLxBindings(t *testing.T) {
	var r *reflection

	tm := &testBindings{
		ID:        1,
		privateID: 10,
		Comms: testComms{
			Address: "address test line 1",
			Emails: []string{
				"test@example.com",
				"test2@example.com",
			},
		},
	}
	r = newReflection(tm)

	res, err := bindingValues(r)
	if err != nil {
		t.Errorf("unexpected error - [%s]", err.Error())
	}
	// privateId should be ignored
	if len(res) != 3 {
		t.Errorf("incorrect length of bindings")
	}

}
