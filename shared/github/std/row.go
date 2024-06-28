package std

import (
	"opg-reports/shared/data"
	"opg-reports/shared/server/response"
)

func ToRow(cmp *Repository) (row *response.Row[*response.Cell]) {
	mapped, _ := data.ToMap(cmp)
	cells := []*response.Cell{}

	for k, v := range mapped {
		cells = append(cells, response.NewCell(k, v))
	}
	row = response.NewRow(cells...)
	return
}

func FromRow(row *response.Row[*response.Cell]) (cmp *Repository) {
	mapped := map[string]interface{}{}

	for _, c := range row.GetCells() {
		mapped[c.Name] = c.Value
	}

	cmp, _ = data.FromMap[*Repository](mapped)
	return
}
