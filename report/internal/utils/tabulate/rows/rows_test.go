package rows

import (
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/tabulate/headers"
	"slices"
	"testing"
)

func TestTabularRowKey(t *testing.T) {
	headers := &headers.Headers{Headers: []*headers.Header{
		{Field: "name", Type: headers.KEY, Default: ""},
		{Field: "team", Type: headers.KEY, Default: ""},
		{Field: "2026-01", Type: headers.DATA, Default: 0.0},
		{Field: "2025-12", Type: headers.DATA, Default: 0.0},
		{Field: "trend", Type: headers.EXTRA, Default: ""},
		{Field: "total", Type: headers.END, Default: 0.0},
	}}
	row := map[string]interface{}{
		"name": "foo",
		"team": "bar",
		"cost": -1.087,
		"date": "2025-12",
	}
	key := Key(row, headers)
	if key != "name=foo^team=bar^" {
		t.Errorf("key generation failed: [%s]", key)
	}
}

func TestTabularRowEmpty(t *testing.T) {

	headers := &headers.Headers{Headers: []*headers.Header{
		{Field: "name", Type: headers.KEY, Default: ""},
		{Field: "team", Type: headers.KEY, Default: ""},
		{Field: "2026-01", Type: headers.DATA, Default: 0.0},
		{Field: "2025-12", Type: headers.DATA, Default: 0.0},
		{Field: "trend", Type: headers.EXTRA, Default: ""},
		{Field: "total", Type: headers.END, Default: 0.0},
	}}
	expected := []string{"name", "team", "2025-12", "2026-01", "trend", "total"}
	emp := Empty(headers)
	// make sure we find all columns
	// check all headers are in expected
	for k, _ := range emp {
		if !slices.Contains(expected, k) {
			t.Errorf("unexpected key in empty row: [%s]", k)
		}
	}
	// check all expected rows are in empty
	for _, k := range expected {
		if _, ok := emp[k]; !ok {
			t.Errorf("expected key is missing: [%s]", k)
		}
	}
}

func TestTabularRowPopulate(t *testing.T) {

	headers := &headers.Headers{Headers: []*headers.Header{
		{Field: "name", Type: headers.KEY, Default: ""},
		{Field: "team", Type: headers.KEY, Default: ""},
		{Field: "2026-01", Type: headers.DATA, Default: 0.0},
		{Field: "2025-12", Type: headers.DATA, Default: 0.0},
		{Field: "trend", Type: headers.EXTRA, Default: ""},
		{Field: "total", Type: headers.END, Default: 0.0},
	}}
	dest := Empty(headers)
	val := -1.087
	row := map[string]interface{}{
		"name": "foo",
		"team": "bar",
		"cost": val,
		"date": "2025-12",
	}
	Populate(row, dest, headers, &Options{ColumnKey: "date", ValueKey: "cost"})
	debugger.Dump(dest)
	if dest["name"].(string) != "foo" ||
		dest["team"].(string) != "bar" ||
		dest["2025-12"].(float64) != val {
		t.Errorf("destination data not as expected.")
	}

}
