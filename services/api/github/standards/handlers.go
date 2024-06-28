package standards

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server/response"
	"sort"
)

func (a *Api[V, F]) List(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store

	errs := resp.GetErrors()
	if len(errs) == 0 {
		activeOnly := func(item *std.Repository) bool {
			return item.Archived == false
		}

		// get everything in range
		onlyActive := store.Filter(activeOnly)
		rows := []*response.Row[*response.Cell]{}

		list := onlyActive.List()
		sort.Slice(list, func(i, j int) bool {
			return list[i].FullName < list[j].FullName
		})

		for _, item := range list {
			row := std.ToRow(item)
			rows = append(rows, row)
			// Add data times to resp
			resp.AddTimestamp(item.TS())
		}
		result := response.NewData(rows...)
		resp.SetResult(result)

	}
	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)
}
