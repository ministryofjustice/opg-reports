package transformers

import (
	"fmt"
	"testing"
)

type tabTest struct {
	ID          int    `json:"id,omitempty" db:"id" faker:"-" doc:"Database primary key."`                                                                              // ID is a generated primary key
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

// DateWideDateValue returns the value of the date field
// Interfaces:
//   - transformers.dateWideTable
func (self *tabTest) DateWideDateValue() string {
	return self.Date
}

func (self *tabTest) DateDeepDateColumn() string {
	return "date"
}

// DateWideValue returns the value to use in the date column
// Interfaces:
//   - transformers.dateWideTable
func (self *tabTest) DateWideValue() string {
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
var expectedWide = map[string]map[string]interface{}{
	"environment:development^service:ecs^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ecs",
		"2024-01":     "-1.01",
		"2024-02":     "3.01",
		"2024-03":     defaultFloat,
	},
	"environment:development^service:ecs^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ecs",
		"2024-01":     "10.0",
		"2024-02":     defaultFloat,
		"2024-03":     defaultFloat,
	},
	"environment:development^service:ec2^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ec2",
		"2024-01":     "3.51",
		"2024-02":     defaultFloat,
		"2024-03":     defaultFloat,
	},
	"environment:development^service:ec2^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ec2",
		"2024-01":     "-4.72",
		"2024-02":     defaultFloat,
		"2024-03":     defaultFloat,
	},
}

var expectedDeep = map[string]map[string]interface{}{
	"date:2024-01^environment:development^service:ecs^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ecs",
		"date":        "2024-01",
		"cost":        "-1.01",
	},
	"date:2024-02^environment:development^service:ecs^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ecs",
		"date":        "2024-02",
		"cost":        "3.01",
	},
	"date:2024-01^environment:development^service:ec2^unit:A^": {
		"environment": "development",
		"unit":        "A",
		"service":     "ec2",
		"date":        "2024-01",
		"cost":        "3.51",
	},
	"date:2024-01^environment:development^service:ecs^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ecs",
		"date":        "2024-01",
		"cost":        "10.0",
	},
	"date:2024-01^environment:development^service:ec2^unit:B^": {
		"environment": "development",
		"unit":        "B",
		"service":     "ec2",
		"date":        "2024-01",
		"cost":        "10.0",
	},
}

// TestTransformersResultsToDateRows checks that a preset
// series of data that mimics api info for costs will come out
// in the expected setup via DateResultsToRows
func TestTransformersResultsToDateRows(t *testing.T) {

	actual, err := ResultsToDateRows(standardSampleData, colValues, dateRanges)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	for key, actualRow := range actual {
		var expectedRow = expectedWide[key]

		for field, actualValue := range actualRow {
			var expectedValue = expectedRow[field]

			if expectedValue.(string) != actualValue.(string) {
				t.Errorf("error with table data [%s] expected field [%s] does not match - [%v]==[%v]", key, field, expectedValue, actualValue)
			}

		}

	}

}

// TestTransformersResultsToDateRows
func TestTransformersResultsToDeepRows(t *testing.T) {

	actual, err := ResultsToDeepRows(standardSampleData, colValues, dateRanges)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	for key, actualRow := range actual {
		var expectedRow = expectedDeep[key]

		for field, actualValue := range actualRow {
			var expectedValue = expectedRow[field]
			if expectedValue.(string) != actualValue.(string) {
				t.Errorf("error with table data [%s] expected field [%s] does not match - [%v]==[%v]", key, field, expectedValue, actualValue)
			}
		}

	}

}
