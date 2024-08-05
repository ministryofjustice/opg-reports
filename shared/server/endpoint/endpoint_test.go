package endpoint

import (
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"opg-reports/shared/logger"
	"opg-reports/shared/server"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/resp/cell"
	"opg-reports/shared/server/resp/row"
	"slices"
	"testing"
)

// func TestRegister(mux *http.ServeMux) {
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		st := data.NewStore[*testEntry]()
// 		epData := NewEndpointData[*testEntry](st, nil, nil)
// 		epDisplay := NewEndpointDisplay(nil)
// 		a.ep := New[*testEntry]("/home", epData, epDisplay)
// 		server.Middleware(a.Index, server.LoggingMW, server.SecurityHeadersMW)
// 	})
// }

func TestSharedServerEndpointFull(t *testing.T) {
	logger.LogSetup()

	mux := testhelpers.Mux()
	min, max, _ := testhelpers.Dates()
	create := 10

	foos := []*testhelpers.TestIEntry{}
	others := []*testhelpers.TestIEntry{}
	active := []*testhelpers.TestIEntry{}
	items := []*testhelpers.TestIEntry{}
	for i := 0; i < create; i++ {
		st := fake.Choice([]bool{true, false})
		if i == 0 {
			st = true
		}
		item := &testhelpers.TestIEntry{
			Id:       fake.IntAsStr(1000, 9999),
			Tags:     []string{"default", fake.Choice([]string{"foo", "bar"})},
			Date:     fake.Date(min, max),
			Category: fake.Choice([]string{"main", "other"}),
			Status:   st,
		}
		if slices.Contains(item.Tags, "foo") {
			foos = append(foos, item)
		}
		if item.Category == "other" {
			others = append(others, item)
		}
		if item.Status {
			active = append(active, item)
		}
		items = append(items, item)
	}
	store := data.NewStoreFromList[*testhelpers.TestIEntry](items)
	// -- data filtering and group
	filters := map[string]data.IStoreFilterFunc[*testhelpers.TestIEntry]{
		"active": func(i *testhelpers.TestIEntry) bool {
			return i.Status == true
		},
	}
	group := func(item *testhelpers.TestIEntry) string {
		return item.Category
	}
	// -- display
	months := dates.Strings(dates.Months(min, max), dates.FormatYM)
	// func() *row.Row
	headF := func() *row.Row {
		r := row.New()
		r.Add(cell.New("Category", "Category", true, false))
		for _, m := range months {
			r.Add(cell.New(m, m, true, false))
		}
		r.Add(cell.New("Total", "Total", true, true))
		return r
	}
	// func(store data.IStore[V]) *row.Row
	rowF := func(store data.IStore[*testhelpers.TestIEntry], resp *resp.Response) *row.Row {
		totalCount := 0
		list := store.List()
		first := list[0]

		r := row.New()
		r.Add(
			cell.New(first.Category, first.Category, true, false),
		)
		for _, m := range months {
			inM := func(item *testhelpers.TestIEntry) bool {
				return dates.InMonth(item.Date.Format(dates.FormatYMD), []string{m})
			}
			values := store.Filter(inM)
			count := values.Length()
			r.Add(cell.New(m, count, false, false))

			totalCount += count
		}
		for _, i := range store.List() {
			resp.AddDataAge(i.TS())
		}
		r.Add(cell.New("Total", totalCount, false, true))
		return r
	}

	// Setup the endpoint
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		resp := resp.New()
		data := NewEndpointData[*testhelpers.TestIEntry](store, group, filters)
		display := NewEndpointDisplay[*testhelpers.TestIEntry](headF, rowF, nil)
		ep := New[*testhelpers.TestIEntry]("test", resp, data, display)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

	w, r := testhelpers.WRGet("/test/")
	mux.ServeHTTP(w, r)

	str, _ := resp.Stringify(w.Result())
	fmt.Println(str)

}
