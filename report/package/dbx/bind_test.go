package dbx

import (
	"database/sql"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/logger"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestBindSimpleSelect(t *testing.T) {
	var err error
	sql := `SELECT date as d, time as :a FROM table WHERE (x = :b OR y = :b )AND (monthA IN (:c) OR monthB IN (:c)) LIMIT 1;`
	data := map[string]interface{}{
		"a": "date",
		"b": 1,
		"c": []string{"2025-01", "2025-02"},
	}

	expected := `SELECT date as d, time as ? FROM table WHERE (x = ? OR y = ? )AND (monthA IN (?,?) OR monthB IN (?,?)) LIMIT 1;`
	expectedArgs := []interface{}{
		"date", 1, 1, "2025-01", "2025-02", "2025-01", "2025-02",
	}
	s, args, err := Bind(t.Context(), sql, data)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	if expected != s {
		t.Errorf("statement does not match expected:\nactual:\n%s\nexpected:\n%s", s, expected)
	}
	for i, exp := range expectedArgs {
		actual := args[i]
		if exp != actual {
			t.Errorf("arg does not match expected:\nactual:\n%s\nexpected:\n%s", s, expected)
		}
	}

}

type testStruct struct {
	Region    string `json:"region,omitempty"`
	Service   string `json:"service,omitempty"`
	Month     string `json:"month,omitempty"`
	Cost      string `json:"cost,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

func TestBindSimpleInsert(t *testing.T) {
	sql := `
INSERT INTO costs (
	region,
	service,
	month,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:month,
	:cost,
	:account_id
) ON CONFLICT (account_id, month, region, service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id
;
`
	s := &testStruct{
		Region:    "NoRegion",
		Service:   "Tax",
		Month:     "2025-01",
		Cost:      "0.01556",
		AccountID: "36700000",
	}
	data := map[string]interface{}{}
	cnv.Convert(s, &data)

	st, _, _ := Bind(t.Context(), sql, data)

	if strings.Count(st, "?") != 5 {
		t.Errorf("failed to bind string")
	}

}

type tmpS struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestBindWorking(t *testing.T) {
	var (
		err    error
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    = t.TempDir()
		driver = "sqlite3"
		dbpath = filepath.Join(dir, "test-binders.db")
	)
	db, err := sql.Open(driver, dbpath)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
		t.FailNow()
	}
	defer db.Close()

	// create a dummy table and write a row to it
	create := `CREATE TABLE IF NOT EXISTS testing(id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;`
	_, err = db.ExecContext(ctx, create)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	insert := `INSERT INTO testing (name) VALUES ('test') RETURNING id;`
	_, err = db.ExecContext(ctx, insert)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	sel := `SELECT id, name FROM testing WHERE name = :name;`
	sel, args, err := Bind(t.Context(), sel, map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	rows, err := db.QueryContext(ctx, sel, args...)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	all := []*tmpS{}
	defer rows.Close()
	for rows.Next() {
		var r = &tmpS{}
		err = rows.Scan(&r.ID, &r.Name)
		all = append(all, r)
	}

	if len(all) != 1 {
		t.Errorf("unexpected number of rows.")
	}
}
