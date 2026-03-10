package handler

import (
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
	"opg-reports/report/packages/types/models"
)

const testSetup string = `
CREATE TABLE IF NOT EXISTS test_model (
	id INTEGER PRIMARY KEY,
	name TEXT,
	month TEXT
) STRICT;

INSERT INTO test_model (name, month) VALUES ('Z', '2025-12');
INSERT INTO test_model (name, month) VALUES ('A', '2026-01');
INSERT INTO test_model (name, month) VALUES ('B', '2026-02');
INSERT INTO test_model (name, month) VALUES ('C', '2026-03');
INSERT INTO test_model (name, month) VALUES ('D', '2026-03');
`

const testSelect string = `
SELECT
	id,
	name,
	month
FROM test_model
WHERE
	month IN(:months)
ORDER BY name ASC;
`

type testFilter struct {
	Months []string `json:"months,omitempty"`
}

type testResult struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Month string `json:"month"`
}

func (self *testResult) Sequence() []any {
	return []any{
		&self.ID,
		&self.Name,
		&self.Month,
	}
}

// Map returns a map of all fields on this struct
func (self *testFilter) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

type testStatement struct {
	sql string
}

func (self *testStatement) SQL(filter interfaces.Filterable) (s string) {
	var model = &models.Filter{}
	convert.Between(filter.Map(), &model)
	s = self.sql
	return
}

// func TestPackagesHandlerAPI(t *testing.T) {
// 	var (
// 		err    error
// 		db     *sql.DB
// 		ctx    = t.Context()
// 		dir    = t.TempDir()
// 		driver = "sqlite3"
// 		dbpath = filepath.Join(dir, "test-handler.db")
// 	)

// 	// setup a small test db
// 	db, err = sql.Open(driver, dbpath)
// 	if err != nil {
// 		t.Errorf("unexpected error: %s", err.Error())
// 		t.FailNow()
// 	}
// 	defer db.Close()

// 	// run db setup
// 	_, err = db.ExecContext(ctx, testSetup)
// 	if err != nil {
// 		t.Errorf("unexpected error: %s", err.Error())
// 		t.FailNow()
// 	}

// 	// SETUP CONFIG
// 	r := &models.APIResponse[*testResult, *testResult, *models.Request, *models.Filter]{
// 		SHA:     "000",
// 		Version: "1.0.0",
// 	}
// 	handler := &APIConfig[*testResult, *testResult, *models.Request, *models.Filter]{
// 		Label:     "test-api-call",
// 		DB:        &args.DB{Driver: driver, DB: dbpath},
// 		Request:   &models.Request{},
// 		Statement: &testStatement{sql: testSelect},
// 		Results:   []*testResult{},
// 		Response:  r,
// 	}

// 	// SETUP MUX
// 	url := `/test/?date_start=2026-01&date_end=2026-02`
// 	req := httptest.NewRequest(http.MethodGet, url, nil)
// 	writer := httptest.NewRecorder()
// 	mux := http.NewServeMux()
// 	RegisterAPI(ctx, mux, handler, `/test/`)
// 	mux.ServeHTTP(writer, req)
// 	// PROCESS RESULT
// 	res := writer.Result()
// 	resp := &models.APIResponse[*testResult, *testResult, *models.Request, *models.Filter]{}

// 	err = convert.Response(res, &resp)
// 	if err != nil {
// 		t.Errorf("unexpected error: %s", err.Error())
// 	}

// 	// test the response data ..
// 	if resp.Version != r.Version || resp.SHA != r.SHA || resp.Label != handler.Label {
// 		t.Errorf("version / sha mismatch in result")
// 	}

// 	if len(resp.Data) != 2 {
// 		t.Errorf("unexpected number of items returned")
// 	}
// 	// check the data, should only contain A & B
// 	for _, item := range resp.Data {
// 		if item.Name != "A" && item.Name != "B" {
// 			t.Errorf("unexpected item returned in data")
// 		}
// 	}

// }
