package costsfront_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/pkg/transformers"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsfront"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsio"
)

// check the function map
var _ navigation.ResponseTransformer = costsfront.TransformResult

var dateRanges = []string{"2024-01", "2024-02", "2024-03"}
var colValues = map[string][]interface{}{
	"unit":        {"A", "B", "C"},
	"environment": {"development", "pre-production", "production"},
	"service":     {"ecs", "ec2", "rds"},
}
var standardSampleData = []*costs.Cost{
	{Unit: "A", Environment: "development", Service: "ecs", Date: "2024-01", Cost: "-1.01"},
	{Unit: "A", Environment: "development", Service: "ecs", Date: "2024-02", Cost: "3.01"},
	{Unit: "A", Environment: "development", Service: "ec2", Date: "2024-01", Cost: "3.51"},
	{Unit: "B", Environment: "development", Service: "ecs", Date: "2024-01", Cost: "10.0"},
	{Unit: "B", Environment: "development", Service: "ec2", Date: "2024-01", Cost: "-4.72"},
}
var expected = map[string]map[string]interface{}{
	"environment:development^service:ecs^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ecs",
		"2024-01":     "-1.01",
		"2024-02":     "3.01",
		"2024-03":     "0.0000",
	},
	"environment:development^service:ecs^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ecs",
		"2024-01":     "10.0",
		"2024-02":     "0.0000",
		"2024-03":     "0.0000",
	},
	"environment:development^service:ec2^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ec2",
		"2024-01":     "3.51",
		"2024-02":     "0.0000",
		"2024-03":     "0.0000",
	},
	"environment:development^service:ec2^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ec2",
		"2024-01":     "-4.72",
		"2024-02":     "0.0000",
		"2024-03":     "0.0000",
	},
}

// TestCostApiTransformersDataRowsStandard checks that a preset
// series of data that mimics api info for costs will come out
// in the expected setup via ResultsToRows
func TestCostApiTransformersDataRowsStandard(t *testing.T) {

	actual, err := transformers.ResultsToRows(standardSampleData, colValues, dateRanges)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	for key, actualRow := range actual {
		var expectedRow = expected[key]

		for field, actualValue := range actualRow {
			var expectedValue = expectedRow[field]

			if expectedValue.(string) != actualValue.(string) {
				t.Errorf("error with table data [%s] expected field [%s] does not match - expected [%v] != actual [%v]", key, field, expectedValue, actualValue)
			}

		}

	}

}

// TestCostApiTransformersResultStandard checks that a preset
// series of data that mimics api info for costs will come out
// in the expected way by calling the handler that is attached
// to navigation data.
func TestCostApiTransformersResultStandard(t *testing.T) {
	bdy := &costsio.CostsStandardBody{
		Type:         "unit-environment",
		ColumnOrder:  []string{"unit", "environment"},
		DateRange:    dateRanges,
		ColumnValues: colValues,
		Result:       standardSampleData,
	}

	actual := costsfront.TransformResult(bdy).(*costsio.CostsStandardBody)
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
