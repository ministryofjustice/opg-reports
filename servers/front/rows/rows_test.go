package rows_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/front/rows"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestServersFrontRowsSkeleton(t *testing.T) {

	s := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	months := dates.Strings(dates.Range(s, e, dates.MONTH), dates.FormatYM)
	intervals := map[string][]string{"interval": months}

	columns := map[string][]interface{}{
		"unit":    {"foo"},
		"env":     {"prod"},
		"account": {"1"},
	}
	skel := rows.Skeleton(columns, intervals)
	if len(skel) != 1 {
		t.Errorf("permuations incorrect")
	}

	columns = map[string][]interface{}{
		"unit":    {"foo", "bar"},
		"env":     {"prod", "dev"},
		"account": {"1", "2", "3"},
	}
	// math to determine combinations
	l := len(columns["unit"]) * len(columns["env"]) * len(columns["account"])

	skel = rows.Skeleton(columns, intervals)
	if len(skel) != l {
		t.Errorf("permuations incorrect")
	}

}

func TestServersFrontRowsRowMap(t *testing.T) {
	s := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	months := dates.Strings(dates.Range(s, e, dates.MONTH), dates.FormatYM)
	intervals := map[string][]string{"interval": months}

	columns := map[string][]interface{}{
		"unit":    {"foo", "bar"},
		"env":     {"prod", "dev"},
		"account": {"1", "2", "3"},
	}

	data := []interface{}{
		map[string]interface{}{
			"unit":     "foo",
			"env":      "prod",
			"account":  "1",
			"interval": "2024-05",
			"cost":     10.5,
		},
		map[string]interface{}{
			"unit":     "foo",
			"env":      "dev",
			"account":  "3",
			"interval": "2024-06",
			"cost":     65,
		},
		map[string]interface{}{
			"unit":     "bar",
			"env":      "dev",
			"account":  "2",
			"interval": "2024-06",
			"cost":     50,
		},
	}

	values := map[string]string{
		"interval": "cost",
	}

	rows := rows.RowMap(data, columns, intervals, values)
	for key, row := range rows {
		fmt.Printf("[%s]  ", key)

		js := convert.String(row)
		fmt.Println(js)

	}

}
