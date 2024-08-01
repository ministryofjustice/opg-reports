package standards

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/server/response"
	"sort"
)

// List: /github/standards/v1/list/
func (a *Api[V, F, C, R]) List(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()

	errs := resp.GetError()
	if len(errs) == 0 {

		getFilters := a.FiltersForGetParameters(r)
		if len(getFilters) > 0 {
			store = store.Filter(getFilters...)
		}
		rows := []R{}

		list := store.List()
		sort.Slice(list, func(i, j int) bool {
			return list[i].FullName < list[j].FullName
		})

		for _, item := range list {
			row := data.ToRow(item).(R)
			rows = append(rows, row)
			// Add data times to resp
			resp.SetDataAge(item.TS())
		}
		table := response.NewTable(rows...)
		resp.SetData(table)

	}
	a.End(w, r)
}
