package standards

import (
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"strconv"
)

func (a *Api[V, F, C, R]) FiltersForGetParameters(r *http.Request) (filters []data.IStoreFilterFunc[*std.Repository]) {
	slog.Info("filtering based on get parameters")
	filters = []data.IStoreFilterFunc[*std.Repository]{}
	filterValues := a.GetParameters(a.AllowedGetParameters(), r)
	slog.Info("filtering based on get parameters",
		slog.String("filterValues", fmt.Sprintf("%+v", filterValues)),
	)
	// if archived is set, then check if the item status matches the
	// value we are looking for
	if values, ok := filterValues["archived"]; ok {

		status := false
		if s, err := strconv.ParseBool(values[len(values)-1]); err == nil {
			status = s
		}
		slog.Info("found archived filter", slog.Bool("status", status))
		filters = append(filters, func(item *std.Repository) bool {
			return item.Archived == status
		})
	}

	return
}
