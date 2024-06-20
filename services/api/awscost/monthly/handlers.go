package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/dates"
)

// Index: /aws/costs/v1/monthly
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := a.Response()
	res.Start()
	all := a.store.List()
	res.Set(all, http.StatusOK, nil)
	res.End()
	a.Write(w, res)
}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	items := []*cost.Cost{}
	errs := []error{}
	store := a.store
	res := a.Response()
	res.Start()

	// Get the dates
	startDate, err := dates.StringToDate(r.PathValue("start"))
	if err != nil {
		errs = append(errs, err)
	}
	endDate, err := dates.StringToDate(r.PathValue("end"))
	if err != nil {
		errs = append(errs, err)
	}

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

	res.Set(items, http.StatusOK, errs)
	res.End()
	a.Write(w, res)
}
