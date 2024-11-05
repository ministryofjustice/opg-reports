package uptimefront_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimefront"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeio"
)

// check the function map
var _ navigation.ResponseTransformer = uptimefront.TransformResult

var dateRanges = []string{"2024-01", "2024-02", "2024-03", "2024-04"}

var colValues = map[string][]interface{}{
	"unit": {"A", "B", "C"},
}

var standardSampleData = []*uptime.Uptime{
	{Unit: "A", Average: 55.01, Date: "2024-01"},
	{Unit: "A", Average: 99.01, Date: "2024-02"},
	{Unit: "B", Average: 15.01, Date: "2024-02"},
	{Unit: "A", Average: 5.01, Date: "2024-03"},
	{Unit: "B", Average: 72.01, Date: "2024-04"},
	{Unit: "C", Average: 100.0, Date: "2024-04"},
}
var expected = map[string]map[string]interface{}{
	"unit:A^": {
		"unit":    "A",
		"2024-01": "55.0100",
		"2024-02": "99.0100",
		"2024-03": "5.0100",
		"2024-04": "0.0000",
	},
	"unit:B^": {
		"unit":    "B",
		"2024-01": "0.0000",
		"2024-02": "15.0100",
		"2024-03": "0.0000",
		"2024-04": "72.0100",
	},
	"unit:C^": {
		"unit":    "C",
		"2024-01": "0.0000",
		"2024-02": "0.0000",
		"2024-03": "0.0000",
		"2024-04": "100.0000",
	},
}

// TestUptimeApiTransformersResultStandard checks that a preset
// series of data that mimics api info for uptime will come out
// in the expected way by calling the transformation directly
func TestUptimeApiTransformersResultStandard(t *testing.T) {

	bdy := &uptimeio.UptimeBody{
		Type:         "unit",
		ColumnOrder:  []string{"unit"},
		DateRange:    dateRanges,
		ColumnValues: colValues,
		Result:       standardSampleData,
	}

	actual := uptimefront.TransformResult(bdy).(*uptimeio.UptimeBody)

	for key, actualRow := range actual.TableRows {
		var expectedRow = expected[key]

		for field, actualValue := range actualRow {
			var expectedValue = expectedRow[field]

			if expectedValue.(string) != actualValue.(string) {
				t.Errorf("error with table data [%s] expected field [%s] does not match - [%v]==[%v]", key, field, expectedValue, actualValue)
			}

		}

	}
}
