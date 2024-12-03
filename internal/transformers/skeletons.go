package transformers

import "github.com/ministryofjustice/opg-reports/pkg/consts"

// DateTableSkeleton creates a series of skeleton rows from the known
// column values and date range to form a map of table rows.
//
// `columnValues` is from the api result `.column_values` and contains
// all the possible column values, but not the dates.
//
// `dateRange` is from the api result `.date_range` and contains
// all the dates that will be used in a row.
//
// # Example
//
// Input:
//
//	columnValues := map[string][]string {
//		"environment": []string {"dev", "prod"},
//		"service": []string {"ec2", "ecs"}
//	}
//	dateRange := []string {"2024-01", "2024-02"}
//
// Output:
//
//	map[string]map[string]interface{}{
//		"environment:dev^service:ec2" : map[string]interface{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:dev^service:ecs" : map[string]interface{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ec2" : map[string]interface{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ecs" : map[string]interface{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//	}
func DateTableSkeleton(columnValues map[string][]interface{}, dateRange []string) (skel map[string]map[string]interface{}) {
	var columnValuesAsList [][]string
	var columnValuePermutations []string
	skel = map[string]map[string]interface{}{}

	columnValuesAsList = ColumnValuesList(columnValues)

	if len(columnValuesAsList) == 1 {
		columnValuePermutations = columnValuesAsList[0]
	} else if len(columnValuesAsList) > 1 {
		columnValuePermutations = Permutations(columnValuesAsList...)
	}

	for _, permutation := range columnValuePermutations {
		// insert the unique line into the skeleton
		if _, ok := skel[permutation]; !ok {
			skel[permutation] = map[string]interface{}{}
		}
		// insert the columns into the row with empty strings
		for column := range columnValues {
			skel[permutation][column] = ""
		}

		// loop over the date and set an empty value for it
		for _, date := range dateRange {
			skel[permutation][date] = consts.DefaultFloatString
		}

	}
	return
}
