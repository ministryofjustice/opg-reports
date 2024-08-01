package standards

import (
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"strconv"
	"strings"
)

// AllowedGetParameters allows this data to be filtered by
// - archived
// - teams
func (a *Api[V, F, C, R]) AllowedGetParameters() []string {
	return []string{
		"archived",
		"teams",
	}
}

// FiltersForGetParameters generates a series of filters to pass into the data store .Filter method based on GET parameters.
// - Find the allowed get parameters and their values from the current http.Request
// - If `archived` or `teams` are present in the filter values then add filters functions for those
//
// Archived filter is as simple match against the boolean
//   - `?archived=true`
//   - `?archived=false`
//
// Teams filter checks if any of the passed teams match any of the set values. For OR logic, pass multiple `team` parameters
// for AND logic pass mutiple teams in one team parameter:
//   - `?teams=<TEAM-A>&teams=<OR-TEAM-B>`
//   - `?teams=<TEAM-A>,<AND-TEAM-B>&team=<OR-TEAM-C>`
func (a *Api[V, F, C, R]) FiltersForGetParameters(r *http.Request) (filters []data.IStoreFilterFunc[*std.Repository]) {
	slog.Info("filtering based on get parameters")
	// check what get parameters are found that are allowed
	filters = []data.IStoreFilterFunc[*std.Repository]{}
	filterValues := a.GetParameters(a.AllowedGetParameters(), r)
	slog.Info("filtering based on get parameters",
		slog.String("filterValues", fmt.Sprintf("%+v", filterValues)),
	)

	// if archived is set, then check if the item status matches the
	// value we are looking for
	if values, ok := filterValues["archived"]; ok {

		status := false
		// check the last value of the archived get parameter
		if s, err := strconv.ParseBool(values[len(values)-1]); err == nil {
			status = s
		}
		slog.Info("adding filter for archived check", slog.Bool("wanted value", status))
		filters = append(filters, func(item *std.Repository) bool {
			return item.Archived == status
		})
	}

	// if teams is set, then see if any of the teams on the item matches any of the teams passed
	if teams, ok := filterValues["teams"]; ok {
		slog.Info("adding filter for teams check", slog.String("wanted one of", fmt.Sprintf("%+v", teams)))

		filters = append(filters, func(item *std.Repository) bool {
			mergedTeams := strings.ToLower(strings.Join(item.Teams, ","))
			found := false
			for _, s := range teams {
				if strings.Contains(mergedTeams, strings.ToLower(s)) {
					found = true
					break
				}
			}
			return found
		})
	}

	return
}
