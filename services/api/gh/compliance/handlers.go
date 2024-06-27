package compliance

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/gh/comp"
	"opg-reports/shared/server/response"
)

func (a *Api[V, F]) List(w http.ResponseWriter, r *http.Request) {
	resp := response.NewResponse()
	resp.Start()
	store := a.store

	errs := resp.GetErrors()
	if len(errs) == 0 {
		activeOnly := func(item *comp.Compliance) bool {
			return item.Archived == false
		}
		// get everything in range
		onlyActive := store.Filter(activeOnly)
		rows := []*response.Row[*response.Cell]{}

		for _, item := range onlyActive.List() {
			row := comp.ToRow(item)
			rows = append(rows, row)
		}
		result := response.NewData(rows...)

		resp.SetResult(result)
	}
	resp.End()
	content, _ := json.Marshal(resp)
	a.Write(w, resp.GetStatus(), content)
}
