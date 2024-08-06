package endpoint

import (
	"opg-reports/shared/data"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/row"
)

type IDisplay[V data.IEntry] interface {
	Head() *row.Row
	Rows(groupKey string, store data.IStore[V], resp *resp.Response) []*row.Row
	Foot(body []*row.Row) *row.Row
}

type DisplayHeadFunc func() *row.Row
type DisplayFootFunc func(bodyRows []*row.Row) *row.Row
type DisplayRowFunc[V data.IEntry] func(group string, store data.IStore[V], resp *resp.Response) []*row.Row

type Display[V data.IEntry] struct {
	headF DisplayHeadFunc
	footF DisplayFootFunc
	rowF  DisplayRowFunc[V]
}

func (d *Display[V]) Foot(body []*row.Row) (ro *row.Row) {
	ro = row.New()
	if d.footF != nil {
		ro = d.footF(body)
	}
	return
}
func (d *Display[V]) Rows(groupKey string, store data.IStore[V], resp *resp.Response) (rows []*row.Row) {
	rows = []*row.Row{}
	if d.rowF != nil {
		rows = d.rowF(groupKey, store, resp)
	}
	return
}

func (d *Display[V]) Head() (ro *row.Row) {
	ro = row.New()
	if d.headF != nil {
		ro = d.headF()
	}
	return
}

func NewEndpointDisplay[V data.IEntry](
	headF DisplayHeadFunc,
	rowF DisplayRowFunc[V],
	footF DisplayFootFunc) IDisplay[V] {
	return &Display[V]{
		headF: headF,
		rowF:  rowF,
		footF: footF,
	}
}
