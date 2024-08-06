package endpoint

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/row"
	"opg-reports/shared/server/resp/table"
)

type IEndpoint[V data.IEntry] interface {
	Data() IData[V]
	Display() IDisplay[V]
	ProcessRequest(w http.ResponseWriter, r *http.Request)
}
type Endpoint[V data.IEntry] struct {
	endpoint string
	data     IData[V]
	display  IDisplay[V]
	resp     *resp.Response
	params   map[string][]string
}

func (e *Endpoint[V]) Data() IData[V] {
	return e.data
}

func (e *Endpoint[V]) Display() IDisplay[V] {
	return e.display
}

func (e *Endpoint[V]) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	slog.Info("processing endpoint request", slog.String("endpoint", e.endpoint))

	response := e.resp
	response.Start(w, r)

	data := e.Data()
	data.ApplyFilters()

	display := e.Display()
	table := table.New()

	table.Head = display.Head()
	bdy := []*row.Row{}

	for key, g := range data.ApplyGroupBy() {
		lines := display.Rows(key, g, response)
		bdy = append(bdy, lines...)
		// add the timestamp data
		for _, i := range g.List() {
			response.AddDataAge(i.TS())
		}
	}
	table.Body = bdy
	table.Foot = display.Foot(bdy)
	response.Result = table
	response.End(w, r)
}

func New[V data.IEntry](endpoint string, resp *resp.Response, dataCnf IData[V], displayCnf IDisplay[V], params map[string][]string) IEndpoint[V] {
	return &Endpoint[V]{
		endpoint: endpoint,
		resp:     resp,
		data:     dataCnf,
		display:  displayCnf,
		params:   params,
	}
}
