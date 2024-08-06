// package Row handles a series of table cells and bundles them in a logical
// row structure akin to a spreadsheet or html table.
// To help, the cells within the row are categorised to mark if they
// are a header (like the first cell in a row), supplementary (such
// as an extra column tracking row rotals) or not
package row

import "opg-reports/shared/server/resp/cell"

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
