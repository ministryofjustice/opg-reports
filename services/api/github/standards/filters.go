package standards

import (
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"strconv"
	"strings"
)

// EndpointFilters generates a map of functions that can be run against the data store based on the parameters passed.
//
// Archived filter is as simple match against the boolean
//   - `?archived=true`
//   - `?archived=false`
//
// Teams filter checks if any of the passed teams match any of the set values. For OR logic, pass multiple `team` parameters
// for AND logic pass mutiple teams in one team parameter:
//   - `?team=<TEAM-A>&team=<OR-TEAM-B>`
//   - `?team=<TEAM-A>,<AND-TEAM-B>&team=<OR-TEAM-C>`
func EndpointFilters(parameters map[string][]string) (funcs map[string]data.IStoreFilterFunc[*std.Repository]) {

	// -- setup for the archived check
	// By default, we only want non-archived repos
	//  - Use the last value in the parameters
	archivedValue := false
	if archValues, ok := parameters["archived"]; ok {
		if s, err := strconv.ParseBool(archValues[len(archValues)-1]); err == nil {
			archivedValue = s
		}
	}
	var archivedFilter = func(item *std.Repository) bool {
		return item.Archived == archivedValue
	}

	// -- setup for team filter
	teams := []string{}
	if teamValues, ok := parameters["team"]; ok {
		teams = teamValues
	}
	// If no teams are passed, then everything should be returned
	// otherwise, only those that match team
	var teamFilter = func(item *std.Repository) (found bool) {
		if len(teams) > 0 {
			found = false
			merged := strings.ToLower(strings.Join(item.Teams, ","))
			for _, s := range teams {
				if strings.Contains(merged, strings.ToLower(s)) {
					found = true
					break
				}
			}
		} else {
			found = true
		}
		return
	}

	funcs = map[string]data.IStoreFilterFunc[*std.Repository]{
		"archived": archivedFilter,
		"team":     teamFilter,
	}

	return

}
