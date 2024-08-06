package standards

import (
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/cell"
	"opg-reports/shared/server/resp/row"
	"sort"
)

var DisplayRow endpoint.DisplayRowFunc[*std.Repository] = func(group string, store data.IStore[*std.Repository], resp *resp.Response) (rows []*row.Row) {
	rows = []*row.Row{}
	// sort alphabetically
	list := store.List()
	sort.Slice(list, func(i, j int) bool {
		return list[i].FullName < list[j].FullName
	})

	for _, item := range list {
		cells := []*cell.Cell{}
		mapped, _ := data.ToMap(item)
		for k, v := range mapped {
			cells = append(cells, cell.New(k, v, false, false))
		}
		r := row.New(cells...)
		rows = append(rows, r)
	}
	return
}
