package transformers_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/transformers"
)

type tabTest struct {
	ID          int    `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`                                 // ID is a generated primary key
	Ts          string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."`                                                            // TS is timestamp when the record was created
	Unit        string `json:"unit,omitempty" db:"unit" faker:"oneof: A, B, C" doc:"The name of the unit / team that owns this account."`                               // Unit is the team that owns this account, passed directly
	Environment string `json:"environment,omitempty" db:"environment" faker:"oneof: production, pre-production, development" doc:"Environment type."`                   // Environment is passed along to show if this is production, development etc account
	Region      string `json:"region,omitempty" db:"region" faker:"oneof: NoRegion, eu-west-1, eu-west-2, us-east-2" doc:"Region this cost was generated within."`      // From the cost data, this is the region the service cost aws generated in
	Service     string `json:"service,omitempty" db:"service" faker:"oneof: Tax, ecs, ec2, s3, sqs, waf, ses, rds" doc:"Name of the service that generated this cost."` // The AWS service name
	Date        string `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."`                                                        // The data the cost was incurred - provided from the cost explorer result
	Cost        string `json:"cost,omitempty" db:"cost" faker:"float_string" doc:"Cost value."`                                                                         // The actual cost value as a string - without an currency, but is USD by default}
}

// UID
// Record interface
func (self *tabTest) UID() string {
	return fmt.Sprintf("%s-%d", "costs", self.ID)
}

// TDate
// Transformable interface
func (self *tabTest) TDate() string {
	return self.Date
}

// TValue
// Transformable interface
func (self *tabTest) TValue() string {
	return self.Cost
}

var dateRanges = []string{"2024-01", "2024-02", "2024-03"}
var colValues = map[string][]interface{}{
	"unit":        {"A", "B", "C"},
	"environment": {"development", "pre-production", "production"},
	"service":     {"ecs", "ec2", "rds"},
}
var standardSampleData = []*tabTest{
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

// TestTransformersResultsToRows checks that a preset
// series of data that mimics api info for costs will come out
// in the expected setup via ResultsToRows
func TestTransformersResultsToRows(t *testing.T) {

	actual, err := transformers.ResultsToRows(standardSampleData, colValues, dateRanges)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	for key, actualRow := range actual {
		var expectedRow = expected[key]

		for field, actualValue := range actualRow {
			var expectedValue = expectedRow[field]

			if expectedValue.(string) != actualValue.(string) {
				t.Errorf("error with table data [%s] expected field [%s] does not match - [%v]==[%v]", key, field, expectedValue, actualValue)
			}

		}

	}

}
