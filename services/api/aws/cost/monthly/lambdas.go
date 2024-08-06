package monthly

import (
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"strings"
)

// Filters
var excludeTax = func(item *cost.Cost) bool {
	return strings.ToLower(item.Service) != taxServiceName
}

// Helpers used within grouping
var unit = func(i *cost.Cost) (string, string) {
	return "account_unit", strings.ToLower(i.AccountUnit)
}
var account_id = func(i *cost.Cost) (string, string) {
	return "account_id", i.AccountId
}
var account_env = func(i *cost.Cost) (string, string) {
	return "account_environment", strings.ToLower(i.AccountEnvironment)
}
var service = func(i *cost.Cost) (string, string) {
	return "service", strings.ToLower(i.Service)
}

// Group by month
var byUnit = func(item *cost.Cost) string {
	return data.ToIdxF(item, unit)
}
var byUnitEnv = func(item *cost.Cost) string {
	return data.ToIdxF(item, unit, account_env)
}
var byAccountService = func(item *cost.Cost) string {
	return data.ToIdxF(item, account_id, unit, account_env, service)
}
