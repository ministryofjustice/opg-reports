package endpoint

import (
	"log/slog"
	"opg-reports/shared/data"
)

type IData[V data.IEntry] interface {
	ApplyFilters() data.IStore[V]
	ApplyGroupBy() map[string]data.IStore[V]
}

type Data[V data.IEntry] struct {
	store   data.IStore[V]
	groupBy data.IStoreGrouper[V]
	filters map[string]data.IStoreFilterFunc[V]
}

func (d *Data[V]) ApplyFilters() (store data.IStore[V]) {
	store = d.store
	for name, f := range d.filters {
		slog.Info("applying filter", slog.String("name", name))
		store = store.Filter(f)
	}
	d.store = store
	return
}

func (d *Data[V]) ApplyGroupBy() map[string]data.IStore[V] {
	slog.Info("applying groupby")
	return d.store.Group(d.groupBy)
}

func NewEndpointData[V data.IEntry](
	store data.IStore[V],
	groupBy data.IStoreGrouper[V],
	filters map[string]data.IStoreFilterFunc[V],
) IData[V] {
	return &Data[V]{
		store:   store,
		groupBy: groupBy,
		filters: filters,
	}
}
