package tabulate

import (
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tabulate/rows"
	"testing"
)

func TestTabularTableTabulate(t *testing.T) {
	headers := &headers.Headers{Headers: []*headers.Header{
		{Field: "name", Type: headers.KEY, Default: ""},
		{Field: "team", Type: headers.KEY, Default: ""},
		{Field: "2026-01", Type: headers.DATA, Default: 0.0},
		{Field: "2025-12", Type: headers.DATA, Default: 0.0},
		{Field: "trend", Type: headers.EXTRA, Default: ""},
		{Field: "total", Type: headers.END, Default: 0.0},
	}}
	data := []map[string]interface{}{
		{
			"name": "foo",
			"team": "bar",
			"cost": -1.087,
			"date": "2025-12",
		},
		{
			"name": "foo",
			"team": "bar",
			"cost": 8.107,
			"date": "2026-01",
		},
		{
			"name": "test",
			"team": "bar",
			"cost": 57.10,
			"date": "2025-12",
		},
		{
			"name": "test",
			"team": "bar",
			"cost": -1.1,
			"date": "2026-01",
		},
	}
	res := Tabulate[float64](data, headers, &Options{
		ColumnKey:    "date",
		ValueKey:     "cost",
		SortByColumn: "2026-01",
		RowEndF:      rows.TotalF,
		TableEndF:    TotalF,
	})

	if len(res) != 3 {
		t.Errorf("summary or other row missing")
	}

}
