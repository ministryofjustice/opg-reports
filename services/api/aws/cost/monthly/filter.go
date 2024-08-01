package monthly

import (
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"strings"
)

// AllowedGetParameters allows this data to be filtered by
// - teams
func (a *Api[V, F, C, R]) AllowedGetParameters() []string {
	return []string{
		"unit",
		"environment",
	}
}

// FiltersForGetParameters generates a series of filters to pass into the data store .Filter method based on GET parameters.
func (a *Api[V, F, C, R]) FiltersForGetParameters(r *http.Request) (filters []data.IStoreFilterFunc[*cost.Cost]) {
	slog.Info("filtering based on get parameters")
	// check what get parameters are found that are allowed
	filters = []data.IStoreFilterFunc[*cost.Cost]{}
	filterValues := a.GetParameters(a.AllowedGetParameters(), r)
	slog.Info("filtering based on get parameters",
		slog.String("filterValues", fmt.Sprintf("%+v", filterValues)),
	)

	if len(filterValues) > 0 {
		resp := a.GetResponse()
		resp.SetMetadata("filters", filterValues)
	}

	// if unit is passed, check all of then against the account unit details
	if units, ok := filterValues["unit"]; ok {
		slog.Info("adding filter for unit check", slog.String("wanted one of", fmt.Sprintf("%+v", units)))

		filters = append(filters, func(item *cost.Cost) bool {
			check := strings.ToLower(item.AccountUnit)
			found := false
			for _, u := range units {
				if check == strings.ToLower(u) {
					found = true
					break
				}
			}
			return found
		})
	}

	// if environment is passed, check all of then against the account environemt
	if environments, ok := filterValues["environment"]; ok {
		slog.Info("adding filter for environment check", slog.String("wanted one of", fmt.Sprintf("%+v", environments)))

		filters = append(filters, func(item *cost.Cost) bool {
			check := strings.ToLower(item.AccountEnvironment)
			found := false
			for _, e := range environments {
				if check == strings.ToLower(e) {
					found = true
					break
				}
			}
			return found
		})
	}

	return
}
