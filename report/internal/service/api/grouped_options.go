package api

import (
	"opg-reports/report/internal/utils"
)

type GroupableApiOptions interface{}

// GetGroupedByColumns is used to show how the uptime was grouped
// together. Typically this is used by api recievers to generate
// table headers and so on
//
// Order does matter
func GetGroupedByColumns(groupable GroupableApiOptions) (groups []string) {
	groups = []string{}
	mapped := map[string]string{}
	utils.Convert(groupable, &mapped)

	for k, v := range mapped {
		if v == "true" {
			groups = append(groups, k)
		}
	}
	return
}
