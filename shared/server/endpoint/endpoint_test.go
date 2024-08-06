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
	"strconv"
	"strings"
	"testing"
)

func TestSharedServerEndpointFull(t *testing.T) {
	logger.LogSetup()

	mux := testhelpers.Mux()
	min, max, _ := testhelpers.Dates()
	create := 15

	foos := []*testhelpers.TestIEntry{}
	others := []*testhelpers.TestIEntry{}
	items := []*testhelpers.TestIEntry{}
	for i := 0; i < create; i++ {
		item := &testhelpers.TestIEntry{
			Id:       fake.IntAsStr(1000, 9999),
			Tags:     []string{"default", fake.Choice([]string{"foo", "bar"})},
			Date:     fake.Date(min, max),
			Category: fake.Choice([]string{"main", "other"}),
			Status:   true,
		}
		if slices.Contains(item.Tags, "foo") {
			foos = append(foos, item)
		}
		if item.Category == "other" {
			others = append(others, item)
		}

		items = append(items, item)
	}
	store := data.NewStoreFromList[*testhelpers.TestIEntry](items)
	// -- data filtering and group

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
	rowF := func(store data.IStore[*testhelpers.TestIEntry], resp *resp.Response) []*row.Row {
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
		return []*row.Row{r}
	}

	// Setup the endpoint
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		qp := NewQueryable([]string{
			"active",
			"tag",
		})
		params := qp.Parse(r)
		// filter by status
		activeF := func(i *testhelpers.TestIEntry) bool {
			status := true
			if v, ok := params["active"]; ok {
				check, _ := strconv.ParseBool(v[0])
				status = (i.Status == check)
			}
			return status
		}
		// filter by tag
		tagF := func(i *testhelpers.TestIEntry) bool {
			res := false
			if tags, ok := params["tag"]; ok {
				merged := strings.ToLower(strings.Join(i.Tags, ","))
				for _, s := range tags {
					if strings.Contains(merged, strings.ToLower(s)) {
						res = true
						break
					}
				}
			}
			return res
		}

		filters := map[string]data.IStoreFilterFunc[*testhelpers.TestIEntry]{
			"active": activeF,
			"tag":    tagF,
		}
		resp := resp.New()
		data := NewEndpointData[*testhelpers.TestIEntry](store, group, filters)
		display := NewEndpointDisplay[*testhelpers.TestIEntry](headF, rowF, nil)

		ep := New[*testhelpers.TestIEntry]("test", resp, data, display, params)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

	// should be empty
	w, r := testhelpers.WRGet("/test/?active=false")
	mux.ServeHTTP(w, r)
	_, b := resp.Stringify(w.Result())
	re := resp.New()
	resp.FromJson(b, re)

	l := len(re.Result.Body)
	if l != 0 {
		t.Errorf("unexpected results returned")
	}

	w, r = testhelpers.WRGet("/test/?tag=foo")
	mux.ServeHTTP(w, r)
	str, b := resp.Stringify(w.Result())
	resp.FromJson(b, re)

	total := 0
	for _, r := range re.Result.Body {
		total += int(r.Supplementary[0].Value.(float64))
	}

	if total != len(foos) {
		t.Errorf("tag filter failed: [%v] [%v]", total, len(foos))
		fmt.Println(str)
	}
}
