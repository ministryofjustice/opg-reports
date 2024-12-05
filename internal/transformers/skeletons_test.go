package transformers_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/transformers"
)

// TestTransformersDataTableSkeleton checks that with
// preset data the skeleton returns contains the correct columns
// and date values for each row
func TestTransformersDateTableSkeleton(t *testing.T) {
	var dateRanges = []string{"2024-01", "2024-02", "2024-03"}
	var colValues = map[string][]interface{}{
		"unit":        {"A", "B", "C"},
		"environment": {"development", "pre-production", "production"},
		"service":     {"ecs", "ec2", "rds"},
	}
	var expectedLen = 1
	for _, v := range colValues {
		expectedLen = expectedLen * len(v)
	}
	res := transformers.DateTableSkeleton(colValues, dateRanges)

	if expectedLen != len(res) {
		t.Error("skeleton is missing data")
	}
	// check each skeleton for the columns and dates
	for _, row := range res {
		// check columns have been set
		for field := range colValues {
			if _, ok := row[field]; !ok {
				t.Errorf("column [%s] was not set!", field)
			}
		}
		// check the dates have been set
		for _, date := range dateRanges {
			if _, ok := row[date]; !ok {
				t.Errorf("date [%s] has not been set", date)
			}
		}
	}

}

func TestTransformersTableSkeleton(t *testing.T) {
	var colValues = map[string][]interface{}{
		"unit":        {"A", "B"},
		"environment": {"development", "production"},
		"service":     {"ecs", "ec2"},
		"date":        {"2024-01"},
	}
	var expectedLen = 1
	for _, v := range colValues {
		expectedLen = expectedLen * len(v)
	}
	res := transformers.TableSkeleton(colValues)

	if expectedLen != len(res) {
		t.Error("skeleton is missing data")
	}
	// check each skeleton for the columns and dates
	for _, row := range res {
		// check columns have been set
		for field := range colValues {
			if _, ok := row[field]; !ok {
				t.Errorf("column [%s] was not set!", field)
			}
		}
	}
}
