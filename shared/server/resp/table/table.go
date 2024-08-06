package table

import "opg-reports/shared/server/resp/row"

type Table struct {
	Head *row.Row   `json:"head"`
	Body []*row.Row `json:"body"`
	Foot *row.Row   `json:"foot"`
}

func (t *Table) SetBody(rows ...*row.Row) {
	t.Body = rows
}

func New(body ...*row.Row) *Table {
	t := &Table{
		Head: row.New(),
		Foot: row.New(),
		Body: []*row.Row{},
	}
	t.SetBody(body...)
	return t
}
