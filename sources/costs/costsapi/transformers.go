package costsapi

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/tmplfuncs"
	"github.com/ministryofjustice/opg-reports/pkg/transformers"
	"github.com/ministryofjustice/opg-reports/sources/costs"
)

// costRecordToRow takes a costs.Cost struct and adds its data into an existing table row.
//
// By using the list of columns it generates a key-value string to identify which table row
// this record relates to.
//
// After adding the cost records column data it will then add the cost values to the date
// column. This looks for and sets a column on the row matching the .Date property on the cost.Cost
// and sets the value of that to be .Cost. If a value is already represent, it "adds" to it.
//
// For ease, it returns the key-value used so this can be tracked
func costRecordToRow(cost *costs.Cost, columns []string, existingData map[string]map[string]interface{}) (key string, err error) {
	var (
		ok          bool
		date        string                 = cost.Date
		value       string                 = cost.Cost
		costAsMap   map[string]interface{} = map[string]interface{}{}
		existingRow map[string]interface{} = map[string]interface{}{}
	)

	if err = convert.Cast(cost, costAsMap); err != nil {
		return
	}
	// this is generated id this cost item would use in the possible list
	key = transformers.RowKV(columns, costAsMap)
	// look for the existing data within the dataset
	// - if cant find it, error
	existingRow, ok = existingData[key]
	if !ok {
		err = fmt.Errorf("failed to find existing data with key [%s]", key)
		return
	}

	// now we try and set values on the existingRow from this cost by using the columns
	for _, field := range columns {
		existingRow[field] = costAsMap[field]
	}
	// now we set the date cost value
	if _, ok = existingRow[date]; ok {
		existingRow[date] = tmplfuncs.Add(existingRow[date], value)
	} else {
		existingRow[date] = value
	}

	return
}

// ResultsToRows converts a list of costs (generally from the api) into a series of table rows
// that can then be used for rendering.
//
// It uses the columnValues & date range to generate a skeleton set or rows that cover all
// possible combinations of the columns and fills these with empty values ("" for column,
// "0.0000" for the date fields).
//
// Any of the possible rows that do not have actual data in the result are then removed from
// the end result.
//
// # Example
//
//	Inputs:
//		apiData []*costs.Cost {
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
//		dateRange []string{"2024-01", "2024-02", "2024-03"}
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
func ResultsToRows(apiData []*costs.Cost, columnValues map[string][]string, dateRange []string) (dataAsMap map[string]map[string]interface{}, err error) {

	var (
		columns  []string = transformers.SortedColumnNames(columnValues)
		rowsDone []string = []string{}
	)
	dataAsMap = transformers.DateTableSkeleton(columnValues, dateRange)

	for _, item := range apiData {
		rowKey, e := costRecordToRow(item, columns, dataAsMap)
		if e != nil {
			slog.Error("[costsapi.ResultsToRows] failed", slog.String("err", e.Error()))
			return
		}
		// insert to the list of done rows
		if !slices.Contains(rowsDone, rowKey) {
			rowsDone = append(rowsDone, rowKey)
		}

	}

	// remove any row that has not been marked as 'done' - these are empty combinations
	for key := range dataAsMap {
		if !slices.Contains(rowsDone, key) {
			delete(dataAsMap, key)
		}
	}

	return

}
