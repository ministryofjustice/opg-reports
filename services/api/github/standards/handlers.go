package standards

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server/response"
	"sort"
)

func (a *Api[V, F, C, R]) List(w http.ResponseWriter, r *http.Request) {
	a.Start(w, r)
	resp := a.GetResponse()
	store := a.Store()

	errs := resp.GetError()
	if len(errs) == 0 {
		activeOnly := func(item *std.Repository) bool {
			return item.Archived == false
		}

		// get everything in range
		onlyActive := store.Filter(activeOnly)
		rows := []R{}

		list := onlyActive.List()
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
