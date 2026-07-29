package sqlx

import (
	"opg-reports/app/internal/fmtx"
	"testing"
)

type testFilter struct {
	Team     string   `json:"team"`
	Months   []string `json:"months"`
	Services []string `json:"services"`
}

func TestSQLxParameteriseFilter(t *testing.T) {
	var stmt = `
SELECT
	SUM(cost)
FROM cost
WHERE
	team = :team AND
	month IN (:months)
;`

	var filter = &testFilter{
		Team:   "test-a",
		Months: []string{"2026-06", "2026-05"},
	}
	var bound, _ = bindingValues(NewReflection(filter))

	bsql := newBoundSql(stmt, bound)
	bsql.Parameterise()
	fmtx.Printj(bsql)
	t.FailNow()
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
	var r *Reflection

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
	r = NewReflection(tm)

	res, err := bindingValues(r)
	if err != nil {
		t.Errorf("unexpected error - [%s]", err.Error())
	}
	// privateId should be ignored
	if len(res) != 3 {
		t.Errorf("incorrect length of bindings")
	}

}
