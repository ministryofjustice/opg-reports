package endpoint

import (
	"opg-reports/shared/data"
)

type IEndpoint[V data.IEntry] interface {
	Data() IData[V]
	Display() IDisplay
}
type Endpoint[V data.IEntry] struct {
	endpoint string
	data     IData[V]
	display  IDisplay
}

func (e *Endpoint[V]) Data() IData[V] {
	return e.data
}

func (e *Endpoint[V]) Display() IDisplay {
	return e.display
}

func New[V data.IEntry](endpoint string, dataCnf IData[V], displayCnf IDisplay) IEndpoint[V] {
	return &Endpoint[V]{
		endpoint: endpoint,
		data:     dataCnf,
		display:  displayCnf,
	}
}
