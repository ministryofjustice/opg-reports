package datatable

import "strings"

// SkeletonTable reates a series of skeleton rows from the known
// column values and date range to form a map of table rows.
//
// `keys` is the output from `PossibleCombinationsAsKeys`
// `cells` is normally the date columns
//
// Output:
//
//	map[string]map[string]string{}{
//		"environment:dev^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:dev^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//	}
func SkeletonTable(keys []string, cells []string) (table map[string]map[string]string) {
	table = map[string]map[string]string{}

	for _, key := range keys {
		row := map[string]string{}
		// recreate the column name and value from the formatted key
		// 	- "environment:backup^account:A" => {"environment":"backup", "account":"A"}
		key = strings.TrimSuffix(key, "^")
		for _, columnAndValue := range strings.Split(key, "^") {
			sp := strings.Split(columnAndValue, ":")
			col, val := sp[0], sp[1]
			row[col] = val
		}
		// append the extra cells on to the table
		for _, cell := range cells {
			row[cell] = emptyCell
		}
		table[key+"^"] = row
	}

	return
}
