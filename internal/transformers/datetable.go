package transformers

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/pkg/nums"
)

// DateWideTable interface is used for tables whose
// date values are headers rather than column values
type dateWideTable interface {
	DateWideDateValue() string
	DateWideValue() string
}

// recordToDateRow takes a DateTable struct and adds its data into an existing table row.
//
// By using the list of columns it generates a key-value string to identify which table row
// this record relates to.
//
// After adding the cost records column data it will then add the cost values to the date
// column. This looks for and sets a column on the row matching the .Date property on the cost.Cost
// and sets the value of that to be .Cost. If a value is already represent, it "adds" to it.
//
// For ease, it returns the key-value used so this can be tracked
func recordToDateRow[T dateWideTable](item T, columns []string, existingData map[string]map[string]interface{}) (key string, err error) {
	var (
		ok          bool
		v           interface{}
		date        string                 = item.DateWideDateValue()
		value       string                 = item.DateWideValue()
		asMap       map[string]interface{} = map[string]interface{}{}
		existingRow map[string]interface{} = map[string]interface{}{}
	)

	if err = structs.Convert(item, asMap); err != nil {
		return
	}
	// this is generated id this cost item would use in the possible list
	key = RowKV(columns, asMap)
	// look for the existing data within the dataset
	// - if cant find it, error
	existingRow, ok = existingData[key]
	if !ok {
		err = fmt.Errorf("failed to find existing data with key [%s]", key)
		return
	}

	// now we try and set values on the existingRow from this cost by using the columns
	for _, field := range columns {
		existingRow[field] = asMap[field]
	}
	// now we set the date cost value
	// if there is a value present, and it is not the default value
	// then Add those together
	if v, ok = existingRow[date]; ok && v.(string) != defaultFloat {
		existingRow[date] = nums.Add(existingRow[date], value)
	} else {
		existingRow[date] = value
	}

	return
}

// ResultsToDateRows converts a list of data (generally from the api) into a series of table rows
// that can then be used for rendering.
//
// It uses the columnValues & dateRange to generate a skeleton set or rows that cover all
// possible combinations of the columns and fills these with empty values ("" for column,
// "0.0000" for the date fields).
//
// Any of the possible rows that do not have actual data in the result are then removed from
// the end result.
//
// # Example
//
//	Inputs:
//		apiData []*costs.Cost{
//			{Unit: "A", Environment: "development", Service: "ecs", Date: "2024-01", Cost: "-1.01"},
//			{Unit: "A", Environment: "development", Service: "ecs", Date: "2024-02", Cost: "3.01"},
//			{Unit: "A", Environment: "development", Service: "ec2", Date: "2024-01", Cost: "3.51"},
//			{Unit: "B", Environment: "development", Service: "ecs", Date: "2024-01", Cost: "10.0"},
//			{Unit: "B", Environment: "development", Service: "ec2", Date: "2024-01", Cost: "-4.72"},
//		},
//		columnValues map[string][]string{
//			"unit":        {"A", "B", "C"},
//			"environment": {"development", "pre-production", "production"},
//			"service":     {"ecs", "ec2", "rds"},
//		},
//		dateRange []string{"2024-01", "2024-02", "2024-03"},
//
// Output:
//
//	map[string]map[string]interface{}{
//		"environment:development^service:ecs^unit:A^": map[string]interface{}{
//			"environment": "development",
//			"unit":        "A",
//			"service":     "ecs",
//			"2024-01":     "-1.0100",
//			"2024-02":     "3.0100",
//			"2024-03":     "0.0000",
//		},
//		"environment:development^service:ecs^unit:B^": map[string]interface{}{
//			"environment": "development",
//			"unit":        "B",
//			"service":     "ecs",
//			"2024-01":     "10.0000",
//			"2024-02":     "0.0000",
//			"2024-03":     "0.0000",
//		},
//		"environment:development^service:ec2^unit:A^": map[string]interface{}{
//			"environment": "development",
//			"unit":        "A",
//			"service":     "ec2",
//			"2024-01":     "3.5100",
//			"2024-02":     "0.0000",
//			"2024-03":     "0.0000",
//		},
//		"environment:development^service:ec2^unit:B^": map[string]interface{}{
//			"environment": "development",
//			"unit":        "B",
//			"service":     "ec2",
//			"2024-01":     "-4.7200",
//			"2024-02":     "0.0000",
//			"2024-03":     "0.0000",
//		},
//	}
func ResultsToDateRows[T dateWideTable](apiData []T, columnValues map[string][]interface{}, dateRange []string) (dataAsMap map[string]map[string]interface{}, err error) {
	// columns is sorted column names only - this is to ensure 'key' order is a match
	var columns []string = SortedColumnNames(columnValues)
	// found tracks which 'key' has real data and inserted in to the data map
	// so anything that is not in this list can be removed - as it
	// will not have and values
	var found []string = []string{}

	// generate a skel of table data
	dataAsMap = DateTableSkeleton(columnValues, dateRange)

	for _, item := range apiData {
		rowKey, e := recordToDateRow(item, columns, dataAsMap)
		if e != nil {
			slog.Error("[transformers] recordToDateRow failed", slog.String("err", e.Error()))
			return
		}
		// insert to the list of done rows
		if !slices.Contains(found, rowKey) {
			found = append(found, rowKey)
		}

	}

	// remove any row that has not been marked as 'done' - these are empty combinations
	for key := range dataAsMap {
		if !slices.Contains(found, key) {
			delete(dataAsMap, key)
		}
	}

	return

}
