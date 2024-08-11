// package Row handles a series of table cells and bundles them in a logical
// row structure akin to a spreadsheet or html table.
// To help, the cells within the row are categorised to mark if they
// are a header (like the first cell in a row), supplementary (such
// as an extra column tracking row rotals) or not
package row

import (
	"github.com/ministryofjustice/opg-reports/servers/shared/resp/cell"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type Row struct {
	Headers       []*cell.Cell `json:"headers"`
	Data          []*cell.Cell `json:"data"`
	Supplementary []*cell.Cell `json:"supplementary"`
}

func (r *Row) Add(cells ...*cell.Cell) {
	for _, c := range cells {
		if c.IsSupplementary {
			r.Supplementary = append(r.Supplementary, c)
		} else if c.IsHeader {
			r.Headers = append(r.Headers, c)
		} else {
			r.Data = append(r.Data, c)
		}
	}
	return
}

func (r *Row) All() (all []*cell.Cell) {
	all = []*cell.Cell{}
	for _, h := range r.Headers {
		all = append(all, h)
	}
	for _, d := range r.Data {
		all = append(all, d)
	}
	for _, s := range r.Supplementary {
		all = append(all, s)
	}
	return
}

func ToStruct[T any](r *Row) (item T) {
	mapped := map[string]interface{}{}
	for _, c := range r.All() {
		mapped[c.Name] = c.Value
	}
	item, _ = convert.FromMap[T](mapped)
	return
}

func ToStructs[T any](rows []*Row) (items []T) {
	items = []T{}
	for _, r := range rows {
		item := ToStruct[T](r)
		items = append(items, item)
	}
	return
}

func New(cells ...*cell.Cell) (r *Row) {
	r = &Row{
		Headers:       []*cell.Cell{},
		Data:          []*cell.Cell{},
		Supplementary: []*cell.Cell{},
	}
	if len(cells) > 0 {
		r.Add(cells...)
	}
	return
}
