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
		store = store.Filter(f)
		slog.Info("applied filter",
			slog.String("name", name),
			slog.Int("datastore count", store.Length()))
	}
	d.store = store
	return
}

func (d *Data[V]) ApplyGroupBy() (g map[string]data.IStore[V]) {

	// provide a default group of the store
	g = map[string]data.IStore[V]{"all": d.store}
	if d.groupBy != nil {
		g = d.store.Group(d.groupBy)
	}
	slog.Info("applied groupby",
		slog.Int("datastore groups", len(g)))
	return
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
