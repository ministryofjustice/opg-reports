package monthly

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"time"
)

// Index: /aws/costs/v1/monthly
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiSimpleResponse()
	res.Start()
	res.SetStatus(http.StatusOK)
	res.End()
	content, _ := json.Marshal(res)
	a.Write(w, res.GetStatus(), content)
}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
// Note: if {start} or {end} are "-" it uses current month
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiResponseWithResult[*cost.Cost, map[string][]*cost.Cost]()
	res.Start()

	store := a.store

	startDate, endDate := startAndEndDates(r, res)

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		f := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		store := store.Filter(f)
		// For totals, we group data just the month in YYYY-MM format
		ym := func(i *cost.Cost) (string, string) {
			return "month", dates.Reformat(i.Date, dates.FormatYM)
		}
		byMonth := func(item *cost.Cost) string {
			return data.ToIdxF(item, ym)
		}

		items := map[string][]*cost.Cost{}
		for k, g := range store.Group(byMonth) {
			for _, v := range g.List() {
				items[k] = append(items[k], v)
			}
		}
		res.SetResult(items)

	}

	res.End()
	content, _ := json.Marshal(res)
	a.Write(w, res.GetStatus(), content)
}

func startAndEndDates(r *http.Request, res server.IApiResponseBase) (startDate time.Time, endDate time.Time) {
	var err error
	now := time.Now().UTC().Format(dates.FormatYM)
	// Get the dates
	startDate, err = dates.StringToDateDefault(r.PathValue("start"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	endDate, err = dates.StringToDateDefault(r.PathValue("end"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	return

}
