package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/dates"
	"time"
)

// Index: /aws/costs/v1/monthly
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := a.Response()
	res.Start()

	res.SetStatus(http.StatusOK)

	res.End()
	a.Write(w, res)
}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
// Note: if {start} or {end} are "-" it uses current month
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	items := []*cost.Cost{}
	store := a.store
	res := a.Response()
	res.Start()
	now := time.Now().UTC().Format(dates.FormatYM)
	// Get the dates
	startDate, err := dates.StringToDateDefault(r.PathValue("start"), "-", now)
	if err != nil {
		res.AddStatusError(http.StatusConflict, err)
	}
	endDate, err := dates.StringToDateDefault(r.PathValue("end"), "-", now)
	if err != nil {
		res.AddStatusError(http.StatusConflict, err)
	}

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		f := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		store := store.Filter(f)

		// set items to be from the store for now
		items = store.List()
	}

	if res.GetStatus() == http.StatusOK {
		res.SetResults(items)
	}
	res.End()
	a.Write(w, res)
}
