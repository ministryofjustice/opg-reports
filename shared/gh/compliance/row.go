package compliance

import (
	"opg-reports/shared/data"
	"opg-reports/shared/server/response"
)

func ToRow(cmp *Compliance) (row *response.Row[*response.Cell]) {
	mapped, _ := data.ToMap(cmp)
	cells := []*response.Cell{}

	for k, v := range mapped {
		cells = append(cells, response.NewCell(k, v))
	}
	row = response.NewRow(cells...)
	return
}

func FromRow(row *response.Row[*response.Cell]) (cmp *Compliance) {
	mapped := map[string]string{}

	for _, c := range row.GetCells() {
		mapped[c.Name] = c.Value
	}

	cmp, _ = data.FromMap[*Compliance](mapped)
	return
}
