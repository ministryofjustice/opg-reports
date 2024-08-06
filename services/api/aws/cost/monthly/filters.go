package monthly

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/resp"
	"strings"
)

var excludeTax = func(item *cost.Cost) bool {
	return strings.ToLower(item.Service) != taxServiceName
}

func FilterFunctions(parameters map[string][]string, response *resp.Response) (funcs map[string]data.IStoreFilterFunc[*cost.Cost]) {
	// -- get the start & end dates as well as list of all months
	startDate, endDate := startEnd(parameters)
	months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
	slog.Debug("[aws/costs/monthly] FilterFunctions",
		slog.String("start", startDate.String()),
		slog.String("end", endDate.String()),
		slog.String("months", fmt.Sprintf("%+v", months)),
		slog.String("parameters", fmt.Sprintf("%+v", parameters)))

	response.Metadata["StartDate"] = startDate
	response.Metadata["EndDate"] = endDate
	response.Metadata["filters"] = parameters

	// -- month range filters
	var inMonth = func(item *cost.Cost) bool {
		return dates.InMonth(item.Date, months)
	}

	// -- setup for units filter
	units := []string{}
	if unitValues, ok := parameters["unit"]; ok {
		units = unitValues
	}
	// If no units are passed, then everything should be returned
	// otherwise, check those that are an exact match
	var unitFilter = func(item *cost.Cost) (found bool) {
		if len(units) > 0 {
			found = false
			check := strings.ToLower(item.AccountUnit)
			for _, u := range units {
				if check == strings.ToLower(u) {
					found = true
					break
				}
			}
		} else {
			found = true
		}
		return
	}
	// -- setup for environments filter
	envs := []string{}
	if envValues, ok := parameters["environment"]; ok {
		envs = envValues
	}
	var envFilter = func(item *cost.Cost) (found bool) {
		if len(envs) > 0 {
			found = false
			check := strings.ToLower(item.AccountEnvironment)
			for _, i := range envs {
				if check == strings.ToLower(i) {
					found = true
					break
				}
			}
		} else {
			found = true
		}
		return
	}

	funcs = map[string]data.IStoreFilterFunc[*cost.Cost]{
		"inMonth":     inMonth,
		"unit":        unitFilter,
		"environment": envFilter,
	}

	return
}
