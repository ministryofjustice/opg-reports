package daily

import (
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"strings"
)

// Helpers used within grouping
var unit = func(i *uptime.Uptime) (string, string) {
	return "account_unit", strings.ToLower(i.AccountUnit)
}

// Group by unit
var byUnit = func(item *uptime.Uptime) string {
	return data.ToIdxF(item, unit)
}
